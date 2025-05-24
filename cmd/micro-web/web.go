package web

import (
	"context"
	"encoding/json"
	"fmt"
	//"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/registry"
	"go-micro.dev/v5/store"
	"tailscale.com/tsnet"
)

var sidebarTemplate = `
<div id="sidebar">
  <a href="/">Services</a>
  <a href="/store">Store</a>
  <a href="/broker">Broker</a>
  <a href="/config">Config</a>
  <a href="/registry">Registry</a>
</div>
`

var htmlTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width" />
    <title>Micro Web</title>
    <style>
      html, body {
        font-family: Arial;
        font-size: 16px;
        margin: 0;
        padding: 0;
      }
      #head {
        margin: 25px;
      }
      #head a { color: black; text-decoration: none; }
      .container {
         padding: 25px;
         max-width: 1400px;
         margin: 0 auto;
         display: flex;
      }
      #sidebar {
        min-width: 180px;
        max-width: 220px;
        margin-right: 40px;
        display: flex;
        flex-direction: column;
        gap: 10px;
        background: #f7f7f7;
        border-radius: 8px;
        padding: 20px 10px;
        height: fit-content;
      }
      #sidebar a {
        color: #222;
        text-decoration: none;
        font-weight: bold;
        padding: 8px 12px;
        border-radius: 5px;
        transition: background 0.2s;
      }
      #sidebar a:hover {
        background: #e0e0e0;
      }
      #main-content {
        flex: 1;
      }
      a { color: black; text-decoration: none; font-weight: bold; margin-bottom: 10px;}
      pre { background: #f5f5f5; border-radius: 5px; padding: 10px; overflow: scroll;}
      input, button { border-radius: 5px; padding: 10px; display: block; margin-bottom: 5px; }
      button:hover { cursor: pointer; }
    </style>
  </head>
  <body>
     <div id="head">
       <h1><a href="/">Micro</a></h1>
     </div>
     <div class="container">
       %s
       <div id="main-content">%s</div>
     </div>
  </body>
</html>
`

var serviceTemplate = `
<h2>%s</h2>
<div>%s</div>
<pre>%s</pre>
`

var endpointTemplate = `
<h2>%s</h2>
<form action=%s method=POST>%s</form>
`
var responseTemplate = `<h2>Response</h2><div>%s</div>`

func normalize(v string) string {
	return strings.Title(v)
}

func render(w http.ResponseWriter, v string) error {
	html := fmt.Sprintf(htmlTemplate, sidebarTemplate, v)
	_, err := w.Write([]byte(html))
	return err
}

func rpcCall(service, endpoint string, request []byte) ([]byte, error) {
	data := []byte(`{}`)

	if len(request) > 0 {
		data = request
	}

	req := client.NewRequest(service, endpoint, &bytes.Frame{Data: data})
	var rsp bytes.Frame
	err := client.Call(context.TODO(), req, &rsp)
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}

func init() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// list serfvices
		if r.URL.Path == "/" {
			services, _ := registry.ListServices()

			var html string

			for _, service := range services {
				html += fmt.Sprintf(`<p><a href="/%s">%s</a></p>`, url.QueryEscape(service.Name), service.Name)
			}

			html = fmt.Sprintf(htmlTemplate, html)
			w.Write([]byte(html))

			return
		}

		// got more e.g /helloworld
		parts := strings.Split(r.URL.Path, "/")

		if len(parts) < 2 {
			return
		}

		service := parts[1]

		// get service
		s, err := registry.GetService(service)
		if err != nil {
			return
		}

		// no service
		if len(s) == 0 {
			return
		}

		// service definition for /helloworld
		if len(parts) < 3 {
			var endpoints string
			for _, ep := range s[0].Endpoints {
				parts := strings.Split(ep.Name, ".")
				uri := fmt.Sprintf("/%s/%s/%s", service, parts[0], parts[1])
				endpoints += fmt.Sprintf(`<div><a href="%s">%s</a>`, uri, ep.Name)
			}

			b, _ := json.MarshalIndent(s[0], "", "    ")

			serviceHTML := fmt.Sprintf(serviceTemplate, service, endpoints, string(b))
			render(w, serviceHTML)
			return
		}

		endpoint := parts[2]

		if len(parts) == 4 {
			endpoint = normalize(endpoint) + "." + normalize(parts[3])
		} else {
			endpoint = normalize(service) + "." + normalize(endpoint)
		}

		// get the endpoint
		var ep *registry.Endpoint
		for _, eps := range s[0].Endpoints {
			if eps.Name == endpoint {
				ep = eps
				break
			}
		}

		// no endpoint match
		if ep == nil {
			return
		}

		// render form
		if r.Method == "GET" {
			var inputs string
			if ep != nil {
				for _, input := range ep.Request.Values {
					inputs += fmt.Sprintf(`<input id=%s name=%s placeholder=%s>`, input.Name, input.Name, input.Name)
				}
			} else {
				inputs = `<textarea></textarea>`
			}

			inputs += `<button>Submit</button>`

			formHTML := fmt.Sprintf(endpointTemplate, service, r.URL.Path, inputs)
			render(w, formHTML)
			return
		}

		if r.Method == "POST" {
			r.ParseForm()

			request := map[string]interface{}{}

			for k, v := range r.Form {
				request[k] = strings.Join(v, ",")
			}

			b, _ := json.Marshal(request)

			rsp, err := rpcCall(service, endpoint, b)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var response map[string]interface{}
			json.Unmarshal(rsp, &response)

			var output string
			for k, v := range response {
				output += fmt.Sprintf(`<div>%s: %s</div>`, k, v)
			}

			pretty, _ := json.MarshalIndent(response, "", "    ")
			output += fmt.Sprintf(`<pre>%s</pre>`, string(pretty))
			render(w, fmt.Sprintf(responseTemplate, output))
			return
		}

		// Store Admin UI
		if r.URL.Path == "/store" && r.Method == "GET" {
			storeHTML := `
			<h2>Store Admin</h2>
			<form action="/store/write" method="POST">
			  <h3>Write</h3>
			  <input name="key" placeholder="Key">
			  <input name="value" placeholder="Value">
			  <input name="table" placeholder="Table (optional)">
			  <button>Write</button>
			</form>
			<form action="/store/read" method="POST">
			  <h3>Read</h3>
			  <input name="key" placeholder="Key">
			  <input name="table" placeholder="Table (optional)">
			  <button>Read</button>
			</form>
			<form action="/store/delete" method="POST">
			  <h3>Delete</h3>
			  <input name="key" placeholder="Key">
			  <input name="table" placeholder="Table (optional)">
			  <button>Delete</button>
			</form>
			<form action="/store/list" method="POST">
			  <h3>List (by prefix)</h3>
			  <input name="prefix" placeholder="Prefix">
			  <input name="table" placeholder="Table (optional)">
			  <button>List</button>
			</form>
			`
			render(w, storeHTML)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/store/") && r.Method == "POST" {
			var output string
			var backLink = `<p><a href="/store">Back to Store Admin</a></p>`
			switch r.URL.Path {
			case "/store/write":
				key := r.FormValue("key")
				value := r.FormValue("value")
				table := r.FormValue("table")
				rec := &store.Record{Key: key, Value: []byte(value)}
				var opts []store.WriteOption
				if table != "" {
					opts = append(opts, store.Table(table))
				}
				err := store.DefaultStore.Write(rec, opts...)
				if err != nil {
					output = "<div>Error: " + err.Error() + "</div>" + backLink
				} else {
					output = "<div>Write OK</div>" + backLink
				}
			case "/store/read":
				key := r.FormValue("key")
				table := r.FormValue("table")
				var opts []store.ReadOption
				if table != "" {
					opts = append(opts, store.Table(table))
				}
				recs, err := store.DefaultStore.Read(key, opts...)
				if err != nil {
					output = "<div>Error: " + err.Error() + "</div>" + backLink
				} else if len(recs) == 0 {
					output = "<div>No record found</div>" + backLink
				} else {
					pretty, _ := json.MarshalIndent(recs[0], "", "    ")
					output = "<pre>" + string(pretty) + "</pre>" + backLink
				}
			case "/store/delete":
				key := r.FormValue("key")
				table := r.FormValue("table")
				var opts []store.DeleteOption
				if table != "" {
					opts = append(opts, store.Table(table))
				}
				err := store.DefaultStore.Delete(key, opts...)
				if err != nil {
					output = "<div>Error: " + err.Error() + "</div>" + backLink
				} else {
					output = "<div>Delete OK</div>" + backLink
				}
			case "/store/list":
				prefix := r.FormValue("prefix")
				table := r.FormValue("table")
				var opts []store.ReadOption
				if table != "" {
					opts = append(opts, store.Table(table))
				}
				if prefix != "" {
					opts = append(opts, store.Prefix())
				}
				recs, err := store.DefaultStore.Read(prefix, opts...)
				if err != nil {
					output = "<div>Error: " + err.Error() + "</div>" + backLink
				} else {
					pretty, _ := json.MarshalIndent(recs, "", "    ")
					output = "<pre>" + string(pretty) + "</pre>" + backLink
				}
			default:
				w.WriteHeader(404)
				return
			}
			render(w, output)
			return
		}

		// Broker Admin UI
		if r.URL.Path == "/broker" && r.Method == "GET" {
			brokerHTML := `
			<h2>Broker Admin</h2>
			<form action="/broker/publish" method="POST">
			  <h3>Publish</h3>
			  <input name="topic" placeholder="Topic">
			  <input name="message" placeholder="Message">
			  <button>Publish</button>
			</form>
			<form action="/broker/subscribe" method="POST">
			  <h3>Subscribe (one message)</h3>
			  <input name="topic" placeholder="Topic">
			  <button>Subscribe</button>
			</form>
			`
			render(w, brokerHTML)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/broker/") && r.Method == "POST" {
			var output string
			var backLink = `<p><a href="/broker">Back to Broker Admin</a></p>`
			switch r.URL.Path {
			case "/broker/publish":
				topic := r.FormValue("topic")
				msg := r.FormValue("message")
				if topic == "" || msg == "" {
					output = "<div>Error: missing topic or message</div>" + backLink
				} else if err := broker.DefaultBroker.Publish(topic, &broker.Message{Body: []byte(msg)}); err != nil {
					output = "<div>Error: " + err.Error() + "</div>" + backLink
				} else {
					output = "<div>Publish OK</div>" + backLink
				}
			case "/broker/subscribe":
				topic := r.FormValue("topic")
				if topic == "" {
					output = "<div>Error: missing topic</div>" + backLink
				} else {
					ch := make(chan *broker.Message, 1)
					_, err := broker.DefaultBroker.Subscribe(topic, func(p broker.Event) error {
						ch <- p.Message()
						return nil
					})
					if err != nil {
						output = "<div>Error: " + err.Error() + "</div>" + backLink
					} else {
						select {
						case msg := <-ch:
							output = "<div>Received: " + string(msg.Body) + "</div>" + backLink
						case <-r.Context().Done():
							output = "<div>Timeout</div>" + backLink
						}
					}
				}
			default:
				w.WriteHeader(404)
				return
			}
			render(w, output)
			return
		}

		// Config Admin UI
		if r.URL.Path == "/config" && r.Method == "GET" {
			configHTML := `
			<h2>Config Admin</h2>
			<form action="/config/get" method="POST">
			  <h3>Get</h3>
			  <input name="key" placeholder="Key">
			  <button>Get</button>
			</form>
			<form action="/config/set" method="POST">
			  <h3>Set</h3>
			  <input name="key" placeholder="Key">
			  <input name="value" placeholder="Value">
			  <button>Set</button>
			</form>
			<form action="/config/delete" method="POST">
			  <h3>Delete</h3>
			  <input name="key" placeholder="Key">
			  <button>Delete</button>
			</form>
			<form action="/config/list" method="POST">
			  <h3>List</h3>
			  <button>List</button>
			</form>
			`
			render(w, configHTML)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/config/") && r.Method == "POST" {
			var output string
			var backLink = `<p><a href="/config">Back to Config Admin</a></p>`
			switch r.URL.Path {
			case "/config/get":
				key := r.FormValue("key")
				if key == "" {
					output = "<div>Error: missing key</div>" + backLink
				} else if val, err := config.DefaultConfig.Get(key); err != nil {
					output = "<div>Error: " + err.Error() + "</div>" + backLink
				} else {
					output = "<div>Value: " + val.String() + "</div>" + backLink
				}
			case "/config/set":
				key := r.FormValue("key")
				value := r.FormValue("value")
				if key == "" {
					output = "<div>Error: missing key</div>" + backLink
				} else if err := config.DefaultConfig.Set(key, value); err != nil {
					output = "<div>Error: " + err.Error() + "</div>" + backLink
				} else {
					output = "<div>Set OK</div>" + backLink
				}
			case "/config/delete":
				key := r.FormValue("key")
				if key == "" {
					output = "<div>Error: missing key</div>" + backLink
				} else if err := config.DefaultConfig.Delete(key); err != nil {
					output = "<div>Error: " + err.Error() + "</div>" + backLink
				} else {
					output = "<div>Delete OK</div>" + backLink
				}
			case "/config/list":
				vals, err := config.DefaultConfig.List()
				if err != nil {
					output = "<div>Error: " + err.Error() + "</div>" + backLink
				} else {
					pretty, _ := json.MarshalIndent(vals, "", "    ")
					output = "<pre>" + string(pretty) + "</pre>" + backLink
				}
			default:
				w.WriteHeader(404)
				return
			}
			render(w, output)
			return
		}

		// Registry Admin UI
		if r.URL.Path == "/registry" && r.Method == "GET" {
			registryHTML := `
			<h2>Registry Admin</h2>
			<form action="/registry/list" method="POST">
			  <h3>List Services</h3>
			  <button>List</button>
			</form>
			<form action="/registry/get" method="POST">
			  <h3>Get Service</h3>
			  <input name="name" placeholder="Service Name">
			  <button>Get</button>
			</form>
			<form action="/registry/register" method="POST">
			  <h3>Register Service</h3>
			  <input name="name" placeholder="Service Name">
			  <input name="version" placeholder="Version">
			  <input name="node_id" placeholder="Node ID">
			  <input name="address" placeholder="Node Address">
			  <button>Register</button>
			</form>
			<form action="/registry/deregister" method="POST">
			  <h3>Deregister Service</h3>
			  <input name="name" placeholder="Service Name">
			  <input name="node_id" placeholder="Node ID">
			  <button>Deregister</button>
			</form>
			`
			render(w, registryHTML)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/registry/") && r.Method == "POST" {
			var output string
			var backLink = `<p><a href="/registry">Back to Registry Admin</a></p>`
			switch r.URL.Path {
			case "/registry/list":
				services, err := registry.DefaultRegistry.ListServices()
				if err != nil {
					output = "<div>Error: " + err.Error() + "</div>" + backLink
				} else {
					pretty, _ := json.MarshalIndent(services, "", "    ")
					output = "<pre>" + string(pretty) + "</pre>" + backLink
				}
			case "/registry/get":
				name := r.FormValue("name")
				if name == "" {
					output = "<div>Error: missing name</div>" + backLink
				} else if svc, err := registry.DefaultRegistry.GetService(name); err != nil {
					output = "<div>Error: " + err.Error() + "</div>" + backLink
				} else {
					pretty, _ := json.MarshalIndent(svc, "", "    ")
					output = "<pre>" + string(pretty) + "</pre>" + backLink
				}
			case "/registry/register":
				name := r.FormValue("name")
				version := r.FormValue("version")
				nodeID := r.FormValue("node_id")
				address := r.FormValue("address")
				if name == "" || nodeID == "" || address == "" {
					output = "<div>Error: missing required fields</div>" + backLink
				} else {
					svc := &registry.Service{
						Name:    name,
						Version: version,
						Nodes: []*registry.Node{{
							Id:      nodeID,
							Address: address,
						}},
					}
					if err := registry.DefaultRegistry.Register(svc); err != nil {
						output = "<div>Error: " + err.Error() + "</div>" + backLink
					} else {
						output = "<div>Register OK</div>" + backLink
					}
				}
			case "/registry/deregister":
				name := r.FormValue("name")
				nodeID := r.FormValue("node_id")
				if name == "" || nodeID == "" {
					output = "<div>Error: missing required fields</div>" + backLink
				} else {
					svc := &registry.Service{
						Name: name,
						Nodes: []*registry.Node{{Id: nodeID}},
					}
					if err := registry.DefaultRegistry.Deregister(svc); err != nil {
						output = "<div>Error: " + err.Error() + "</div>" + backLink
					} else {
						output = "<div>Deregister OK</div>" + backLink
					}
				}
			default:
				w.WriteHeader(404)
				return
			}
			render(w, output)
			return
		}
	})

	cmd.Register(&cli.Command{
		Name:  "web",
		Usage: "Launch the web app on port :8082",
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

				ln, err := srv.Listen("tcp", ":8082")
				if err != nil {
					return err
				}

				return http.Serve(ln, h)
			}

			return http.ListenAndServe(":8082", h)
		},
	})
}
