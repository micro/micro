package cli

import (
	"context"
	"fmt"

	"github.com/micro/cli/v2"
	storeproto "github.com/micro/go-micro/v2/store/service/proto"
	"github.com/micro/micro/v2/client/cli/namespace"
	"github.com/micro/micro/v2/client/cli/util"
	inclient "github.com/micro/micro/v2/internal/client"
)

// databases is the entrypoint for micro store databases
func databases(ctx *cli.Context) error {
	client, err := inclient.New(ctx)
	if err != nil {
		return err
	}
	dbReq := client.NewRequest(ctx.String("store"), "Store.Databases", &storeproto.DatabasesRequest{})
	dbRsp := &storeproto.DatabasesResponse{}
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

	client, err := inclient.New(ctx)
	if err != nil {
		return err
	}
	tReq := client.NewRequest(ctx.String("store"), "Store.Tables", &storeproto.TablesRequest{
		Database: ns,
	})
	tRsp := &storeproto.TablesResponse{}
	if err := client.Call(context.TODO(), tReq, tRsp); err != nil {
		return err
	}
	for _, table := range tRsp.Tables {
		fmt.Println(table)
	}
	return nil
}
