package api

import (
	"encoding/json"
	"fmt"
	"go-micro.dev/v5/broker"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/config"
	"go-micro.dev/v5/errors"
	"go-micro.dev/v5/registry"
	"go-micro.dev/v5/store"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"tailscale.com/tsnet"
)

func normalize(v string) string {
	return strings.Title(v)
}

func init() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// assuming we're just going to parse headers
		if r.URL.Path == "/" {
			service := r.Header.Get("Micro-Service")
			endpoint := r.Header.Get("Micro-Endpoint")
			request, _ := io.ReadAll(r.Body)
			if len(request) == 0 {
				request = []byte(`{}`)
			}

			// defaulting to json
			w.Header().Set("Content-Type", "application/json")

			if len(service) == 0 || len(endpoint) == 0 {
				err := errors.New("api.error", "missing service/endpoint", 400)
				w.Header().Set("Micro-Error", err.Error())
				http.Error(w, err.Error(), 400)
				return
			}

			req := client.NewRequest(service, endpoint, &bytes.Frame{Data: request})
			var rsp bytes.Frame
			err := client.Call(r.Context(), req, &rsp)
			if err != nil {
				gerr := errors.New("api.error", err.Error(), 500)
				w.Header().Set("Micro-Error", gerr.Error())
				http.Error(w, gerr.Error(), 500)
			}

			// write the response
			w.Write(rsp.Data)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/store/") {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/store/write":
				if r.Method != "POST" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"error":"method not allowed"}`))
					return
				}
				var req struct {
					Key   string `json:"key"`
					Value string `json:"value"`
					Table string `json:"table,omitempty"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"invalid request body"}`))
					return
				}
				rec := &store.Record{Key: req.Key, Value: []byte(req.Value)}
				var opts []store.WriteOption
				if req.Table != "" {
					opts = append(opts, tableWriteOption(req.Table))
				}
				if err := store.DefaultStore.Write(rec, opts...); err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"` + err.Error() + `"}`))
					return
				}
				w.Write([]byte(`{"result":"ok"}`))
				return
			case "/store/read":
				if r.Method != "GET" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"error":"method not allowed"}`))
					return
				}
				key := r.URL.Query().Get("key")
				table := r.URL.Query().Get("table")
				if key == "" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"missing key"}`))
					return
				}
				var opts []store.ReadOption
				if table != "" {
					opts = append(opts, tableReadOption(table))
				}
				recs, err := store.DefaultStore.Read(key, opts...)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"error":"` + err.Error() + `"}`))
					return
				}
				if len(recs) == 0 {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(`{"error":"not found"}`))
					return
				}
				resp := map[string]interface{}{"key": recs[0].Key, "value": string(recs[0].Value)}
				if recs[0].Expiry > 0 {
					resp["expiry"] = recs[0].Expiry
				}
				b, _ := json.Marshal(resp)
				w.Write(b)
				return
			case "/store/delete":
				if r.Method != "DELETE" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"error":"method not allowed"}`))
					return
				}
				var req struct {
					Key   string `json:"key"`
					Table string `json:"table,omitempty"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"invalid request body"}`))
					return
				}
				var opts []store.DeleteOption
				if req.Table != "" {
					opts = append(opts, tableDeleteOption(req.Table))
				}
				if err := store.DefaultStore.Delete(req.Key, opts...); err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"` + err.Error() + `"}`))
					return
				}
				w.Write([]byte(`{"result":"ok"}`))
				return
			case "/store/list":
				if r.Method != "GET" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"error":"method not allowed"}`))
					return
				}
				prefix := r.URL.Query().Get("prefix")
				table := r.URL.Query().Get("table")
				var opts []store.ReadOption
				if table != "" {
					opts = append(opts, tableReadOption(table))
				}
				if prefix != "" {
					opts = append(opts, prefixReadOption())
				}
				recs, err := store.DefaultStore.Read(prefix, opts...)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"` + err.Error() + `"}`))
					return
				}
				keys := make([]string, 0, len(recs))
				for _, rec := range recs {
					keys = append(keys, rec.Key)
				}
				b, _ := json.Marshal(map[string]interface{}{ "keys": keys })
				w.Write(b)
				return
			}
		}

		// --- Broker API ---
		if strings.HasPrefix(r.URL.Path, "/broker/") {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/broker/publish":
				if r.Method != "POST" {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"error":"method not allowed"}`))
					return
				}
				var req struct {
					Topic   string `json:"topic"`
					Message string `json:"message"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"invalid request body"}`))
					return
				}
				if req.Topic == "" || req.Message == "" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"missing topic or message"}`))
					return
				}
				if err := broker.DefaultBroker.Publish(req.Topic, &broker.Message{Body: []byte(req.Message)}); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"` + err.Error() + `"}`))
					return
				}
				w.Write([]byte(`{"result":"ok"}`))
				return
			case "/broker/subscribe":
				if r.Method != "POST" {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"error":"method not allowed"}`))
					return
				}
				var req struct {
					Topic string `json:"topic"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"invalid request body"}`))
					return
				}
				if req.Topic == "" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"missing topic"}`))
					return
				}
				ch := make(chan *broker.Message, 1)
				_, err := broker.DefaultBroker.Subscribe(req.Topic, func(p broker.Event) error {
					ch <- p.Message()
					return nil
				})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"` + err.Error() + `"}`))
					return
				}
				select {
				case msg := <-ch:
					w.Write([]byte(`{"topic":"` + req.Topic + `","message":"` + string(msg.Body) + `"}`))
				case <-r.Context().Done():
					w.WriteHeader(http.StatusRequestTimeout)
					w.Write([]byte(`{"error":"timeout"}`))
				}
				return
			}
		}

		// --- Config API ---
		if strings.HasPrefix(r.URL.Path, "/config/") {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/config/get":
				if r.Method != "GET" {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"error":"method not allowed"}`))
					return
				}
				key := r.URL.Query().Get("key")
				if key == "" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"missing key"}`))
					return
				}
				val, err := config.DefaultConfig.Get(key)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"error":"` + err.Error() + `"}`))
					return
				}
				b, _ := json.Marshal(map[string]interface{}{ "key": key, "value": val.String() })
				w.Write(b)
				return
			case "/config/set":
				if r.Method != "POST" {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"error":"method not allowed"}`))
					return
				}
				var req struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"invalid request body"}`))
					return
				}
				if req.Key == "" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"missing key"}`))
					return
				}
				if err := config.DefaultConfig.Set(req.Key, req.Value); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"` + err.Error() + `"}`))
					return
				}
				w.Write([]byte(`{"result":"ok"}`))
				return
			}
		}

		// --- Registry API ---
		if strings.HasPrefix(r.URL.Path, "/registry/") {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/registry/list":
				if r.Method != "GET" {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"error":"method not allowed"}`))
					return
				}
				services, err := registry.DefaultRegistry.ListServices()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"` + err.Error() + `"}`))
					return
				}
				b, _ := json.Marshal(map[string]interface{}{ "services": services })
				w.Write(b)
				return
			case "/registry/get":
				if r.Method != "GET" {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"error":"method not allowed"}`))
					return
				}
				name := r.URL.Query().Get("name")
				if name == "" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"missing name"}`))
					return
				}
				svc, err := registry.DefaultRegistry.GetService(name)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"error":"` + err.Error() + `"}`))
					return
				}
				b, _ := json.Marshal(map[string]interface{}{ "service": svc })
				w.Write(b)
				return
			case "/registry/register":
				if r.Method != "POST" {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"error":"method not allowed"}`))
					return
				}
				var req registry.Service
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"invalid request body"}`))
					return
				}
				if req.Name == "" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"missing name"}`))
					return
				}
				if err := registry.DefaultRegistry.Register(&req); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"` + err.Error() + `"}`))
					return
				}
				w.Write([]byte(`{"result":"ok"}`))
				return
			case "/registry/deregister":
				if r.Method != "POST" {
					w.WriteHeader(http.StatusMethodNotAllowed)
					w.Write([]byte(`{"error":"method not allowed"}`))
					return
				}
				var req registry.Service
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"invalid request body"}`))
					return
				}
				if req.Name == "" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"missing name"}`))
					return
				}
				if err := registry.DefaultRegistry.Deregister(&req); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"` + err.Error() + `"}`))
					return
				}
				w.Write([]byte(`{"result":"ok"}`))
				return
			}
		}

		parts := strings.Split(r.URL.Path, "/")

		if len(parts) < 3 {
			return
		}

		service := parts[1]
		endpoint := parts[2]

		if len(parts) == 4 {
			endpoint = normalize(endpoint) + "." + normalize(parts[3])
		} else {
			endpoint = normalize(service) + "." + normalize(endpoint)
		}

		request, _ := io.ReadAll(r.Body)
		if len(request) == 0 {
			request = []byte(`{}`)

			if r.Method == "GET" {
				req := map[string]interface{}{}
				r.ParseForm()
				for k, v := range r.Form {
					req[k] = strings.Join(v, ",")
				}
				if len(req) > 0 {
					request, _ = json.Marshal(req)
				}
			}
		}

		req := client.NewRequest(service, endpoint, &bytes.Frame{Data: request})
		var rsp bytes.Frame
		err := client.Call(r.Context(), req, &rsp)
		if err != nil {
			gerr := errors.New("api.error", err.Error(), 500)
			w.Header().Set("Micro-Error", gerr.Error())
			http.Error(w, gerr.Error(), 500)
		}

		// write the response
		w.Write(rsp.Data)

	})

	cmd.Register(&cli.Command{
		Name:  "api",
		Usage: "Run the micro api on port :8080",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "network", Value: "", Usage: "Set the network e.g --network=tailscale requires TS_AUTHKEY"},
		},
		Action: func(c *cli.Context) error {
			var network string
			var key string

			if c.IsSet("network") {
				network = c.Value("network").(string)
			}

			if network == "tailscale" {
				// check for TS_AUTHKEY
				key = os.Getenv("TS_AUTHKEY")
				if len(key) == 0 {
					return fmt.Errorf("missing TS_AUTHKEY")
				}

				srv := new(tsnet.Server)
				srv.AuthKey = key
				srv.Hostname = "micro"

				ln, err := srv.Listen("tcp", ":8080")
				if err != nil {
					return err
				}

				return http.Serve(ln, h)
			}

			return http.ListenAndServe(":8080", h)
		},
	})
}

// Helper for store.Table as WriteOption, ReadOption, DeleteOption
func tableWriteOption(table string) store.WriteOption {
	return store.WriteOption(store.Table(table))
}
func tableReadOption(table string) store.ReadOption {
	return store.ReadOption(store.Table(table))
}
func tableDeleteOption(table string) store.DeleteOption {
	return store.DeleteOption(store.Table(table))
}
// Helper for store.Prefix as ReadOption
func prefixReadOption() store.ReadOption {
	return store.ReadOption(store.Prefix())
}
