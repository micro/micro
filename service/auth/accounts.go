package auth

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/micro/v2/internal/client"
)

func listAccounts(ctx *cli.Context) {
	client := accountsFromContext(ctx)

	rsp, err := client.List(context.TODO(), &pb.ListAccountsRequest{})
	if err != nil {
		fmt.Printf("Error listing accounts: %v\n", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	defer w.Flush()

	fmt.Fprintln(w, strings.Join([]string{"ID", "Roles", "Metadata"}, "\t"))
	for _, r := range rsp.Accounts {
		var metadata string
		for k, v := range r.Metadata {
			metadata = fmt.Sprintf("%v %v=%v ", metadata, k, v)
		}
		roles := strings.Join(r.Roles, ", ")
		fmt.Fprintln(w, strings.Join([]string{r.Id, roles, metadata}, "\t"))
	}
}

func createAccount(ctx *cli.Context) {
	var options []auth.GenerateOption
	if len(ctx.StringSlice("roles")) > 0 {
		options = append(options, auth.WithRoles(ctx.StringSlice("roles")...))
	}
	if len(ctx.String("secret")) > 0 {
		options = append(options, auth.WithSecret(ctx.String("secret")))
	}

	_, err := authFromContext(ctx).Generate(ctx.String("id"), options...)
	if err != nil {
		fmt.Printf("Error creating account: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Account created")
}

func accountsFromContext(ctx *cli.Context) pb.AccountsService {
	return pb.NewAccountsService("go.micro.auth", client.New())
}
