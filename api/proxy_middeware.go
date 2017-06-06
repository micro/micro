package api

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

type proxyClientWrapper struct {
	client.Client
}

func ProxyClientWrap() client.Wrapper {
	return func(c client.Client) client.Client {
		return &proxyClientWrapper{c}
	}
}

func (p *proxyClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if err := p.Client.Call(ctx, req, rsp, opts...); err != nil {
		return err
	}

	b, err := json.Marshal(req.Request())
	if err != nil {
		return err
	}

	// Publish Announce after handler was called via micro api proxy
	// Every single call from the "outside" of microservice environment will publish
	// An action occur message
	if err := p.Client.Publish(ctx, p.Client.NewPublication(
		globals.AnnounceTopic,
		&announce.AnnounceMessage{
			Handler: fmt.Sprintf("%s.%s", req.Service(), req.Method()),
			Data:    string(b),
		},
	)); err != nil {
		return err
	}
	log.Println("AFTER PUBLISH", req.Service(), req.Method())

	return nil
}
