package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	storeproto "github.com/micro/go-micro/v2/store/service/proto"
	"github.com/olekukonko/tablewriter"
)

// Databases is the entrypoint for micro store databases
func Databases(ctx *cli.Context) error {
	client := *cmd.DefaultOptions().Client
	dbReq := client.NewRequest(ctx.String("store"), "Store.Databases", &storeproto.DatabasesRequest{})
	dbRsp := &storeproto.DatabasesResponse{}
	if err := client.Call(context.TODO(), dbReq, dbRsp); err != nil {
		return err
	}
	t := tablewriter.NewWriter(os.Stdout)
	t.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	t.SetCenterSeparator("|")
	t.SetHeader([]string{"Databases"})
	for _, table := range dbRsp.Databases {
		t.Append([]string{table})
	}
	t.SetFooter([]string{fmt.Sprintf("total %d", len(dbRsp.Databases))})
	t.Render()
	return nil
}

// Tables is the entrypoint for micro store tables
func Tables(ctx *cli.Context) error {
	if len(ctx.String("database")) == 0 {
		return errors.New("database flag is required")
	}
	client := *cmd.DefaultOptions().Client
	tReq := client.NewRequest(ctx.String("store"), "Store.Tables", &storeproto.TablesRequest{
		Database: ctx.String("database"),
	})
	tRsp := &storeproto.TablesResponse{}
	if err := client.Call(context.TODO(), tReq, tRsp); err != nil {
		return err
	}
	t := tablewriter.NewWriter(os.Stdout)
	t.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	t.SetCenterSeparator("|")
	t.SetHeader([]string{"Tables"})
	for _, table := range tRsp.Tables {
		t.Append([]string{table})
	}
	t.SetFooter([]string{fmt.Sprintf("total %d", len(tRsp.Tables))})
	t.Render()
	return nil
}
