package cli

import (
	"context"
	"fmt"

	"github.com/micro/cli/v2"
	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/service/client"
	pb "github.com/micro/micro/v3/service/store/proto"
)

// databases is the entrypoint for micro store databases
func databases(ctx *cli.Context) error {
	dbReq := client.NewRequest(ctx.String("store"), "Store.Databases", &pb.DatabasesRequest{})
	dbRsp := &pb.DatabasesResponse{}
	if err := client.Call(context.TODO(), dbReq, dbRsp); err != nil {
		return err
	}
	for _, db := range dbRsp.Databases {
		fmt.Println(db)
	}
	return nil
}

// tables is the entrypoint for micro store tables
func tables(ctx *cli.Context) error {
	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}

	tReq := client.NewRequest(ctx.String("store"), "Store.Tables", &pb.TablesRequest{
		Database: ns,
	})
	tRsp := &pb.TablesResponse{}
	if err := client.Call(context.TODO(), tReq, tRsp); err != nil {
		return err
	}
	for _, table := range tRsp.Tables {
		fmt.Println(table)
	}
	return nil
}
