package service

import (
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/debug/trace"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/store"

	gcli "github.com/micro/go-micro/v2/client"
	memTracer "github.com/micro/go-micro/v2/debug/trace/memory"
	gsrv "github.com/micro/go-micro/v2/server/grpc"
	memStore "github.com/micro/go-micro/v2/store/memory"
)

func init() {
	// set defaults
	client.DefaultClient = gcli.NewClient()
	server.DefaultServer = gsrv.NewServer()
	store.DefaultStore = memStore.NewStore()
	trace.DefaultTracer = memTracer.NewTracer()
}
