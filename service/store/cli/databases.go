package cli

import (
	"fmt"

	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	pb "github.com/micro/micro/v3/proto/store"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/urfave/cli/v2"
)

// databases is the entrypoint for micro store databases
func databases(ctx *cli.Context) error {
	dbReq := client.NewRequest(ctx.String("store"), "Store.Databases", &pb.DatabasesRequest{})
	dbRsp := &pb.DatabasesResponse{}
	if err := client.DefaultClient.Call(context.DefaultContext, dbReq, dbRsp, client.WithAuthToken()); err != nil {
		return err
	}
	for _, db := range dbRsp.Databases {
		fmt.Println(db)
	}
	return nil
}

// tables is the entrypoint for micro store tables
func tables(ctx *cli.Context) error {
	env, err := util.GetEnv(ctx)
	if err != nil {
		return err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return err
	}

	tReq := client.NewRequest(ctx.String("store"), "Store.Tables", &pb.TablesRequest{
		Database: ns,
	})
	tRsp := &pb.TablesResponse{}
	if err := client.DefaultClient.Call(context.DefaultContext, tReq, tRsp, client.WithAuthToken()); err != nil {
		return err
	}
	for _, table := range tRsp.Tables {
		fmt.Println(table)
	}
	return nil
}
