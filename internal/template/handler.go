package template

var (
	HandlerFNC = `package handler

import (
	"context"

	{{.Alias}} "{{.Dir}}/proto/{{.Alias}}"
)

type {{title .Alias}} struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *{{title .Alias}}) Call(ctx context.Context, req *{{.Alias}}.Request, rsp *{{.Alias}}.Response) error {
	rsp.Msg = "Hello " + req.Name
	return nil
}
`

	HandlerSRV = `package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	{{.Alias}} "{{.Dir}}/proto/{{.Alias}}"
)

type {{title .Alias}} struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *{{title .Alias}}) Call(ctx context.Context, req *{{.Alias}}.Request, rsp *{{.Alias}}.Response) error {
	log.Info("Received {{title .Alias}}.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *{{title .Alias}}) Stream(ctx context.Context, req *{{.Alias}}.StreamingRequest, stream {{.Alias}}.{{title .Alias}}_StreamStream) error {
	log.Infof("Received {{title .Alias}}.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&{{.Alias}}.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *{{title .Alias}}) PingPong(ctx context.Context, stream {{.Alias}}.{{title .Alias}}_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&{{.Alias}}.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
`

	SubscriberFNC = `package subscriber

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	{{.Alias}} "{{.Dir}}/proto/{{.Alias}}"
)

type {{title .Alias}} struct{}

func (e *{{title .Alias}}) Handle(ctx context.Context, msg *{{.Alias}}.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}
`

	SubscriberSRV = `package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	{{.Alias}} "{{.Dir}}/proto/{{.Alias}}"
)

type {{title .Alias}} struct{}

func (e *{{title .Alias}}) Handle(ctx context.Context, msg *{{.Alias}}.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *{{.Alias}}.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
`

	HandlerAPI = `package handler

import (
	"context"
	"encoding/json"
	log "github.com/micro/go-micro/v2/logger"

	"{{.Dir}}/client"
	"github.com/micro/go-micro/v2/errors"
	api "github.com/micro/go-micro/v2/api/proto"
	{{.Alias}} "path/to/service/proto/{{.Alias}}"
)

type {{title .Alias}} struct{}

func extractValue(pair *api.Pair) string {
	if pair == nil {
		return ""
	}
	if len(pair.Values) == 0 {
		return ""
	}
	return pair.Values[0]
}

// {{title .Alias}}.Call is called by the API as /{{.Alias}}/call with post body {"name": "foo"}
func (e *{{title .Alias}}) Call(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Info("Received {{title .Alias}}.Call request")

	// extract the client from the context
	{{.Alias}}Client, ok := client.{{title .Alias}}FromContext(ctx)
	if !ok {
		return errors.InternalServerError("{{.FQDN}}.{{.Alias}}.call", "{{.Alias}} client not found")
	}

	// make request
	response, err := {{.Alias}}Client.Call(ctx, &{{.Alias}}.Request{
		Name: extractValue(req.Post["name"]),
	})
	if err != nil {
		return errors.InternalServerError("{{.FQDN}}.{{.Alias}}.call", err.Error())
	}

	b, _ := json.Marshal(response)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}
`

	HandlerWEB = `package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/micro/go-micro/v2/client"
	{{.Alias}} "path/to/service/proto/{{.Alias}}"
)

func {{title .Alias}}Call(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	{{.Alias}}Client := {{.Alias}}.New{{title .Alias}}Service("{{.Namespace}}.srv.{{.Alias}}", client.DefaultClient)
	rsp, err := {{.Alias}}Client.Call(context.TODO(), &{{.Alias}}.Request{
		Name: request["name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
`
)
