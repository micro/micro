package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/micro/v2/client/cli/namespace"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/internal/client"
)

func listAccounts(ctx *cli.Context) {
	client := accountsFromContext(ctx)

	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		fmt.Printf("Error getting namespace: %v\n", err)
		os.Exit(1)
	}

	rsp, err := client.List(context.TODO(), &pb.ListAccountsRequest{
		Options: &pb.Options{Namespace: ns},
	})
	if err != nil {
		fmt.Printf("Error listing accounts: %v\n", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	defer w.Flush()

	fmt.Fprintln(w, strings.Join([]string{"ID", "Scopes", "Metadata"}, "\t\t"))
	for _, r := range rsp.Accounts {
		var metadata string
		for k, v := range r.Metadata {
			metadata = fmt.Sprintf("%v %v=%v ", metadata, k, v)
		}
		scopes := strings.Join(r.Scopes, ", ")

		if len(metadata) == 0 {
			metadata = "n/a"
		}
		if len(scopes) == 0 {
			scopes = "n/a"
		}

		fmt.Fprintln(w, strings.Join([]string{r.Id, scopes, metadata}, "\t\t"))
	}
}

func createAccount(ctx *cli.Context) {
	if ctx.Args().Len() == 0 {
		fmt.Println("Missing argument: ID")
		return
	}

	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		fmt.Printf("Error getting namespace: %v\n", err)
		os.Exit(1)
	}

	options := []auth.GenerateOption{auth.WithIssuer(ns)}
	if len(ctx.StringSlice("scopes")) > 0 {
		options = append(options, auth.WithScopes(ctx.StringSlice("scopes")...))
	}
	if len(ctx.String("secret")) > 0 {
		options = append(options, auth.WithSecret(ctx.String("secret")))
	}

	acc, err := authFromContext(ctx).Generate(ctx.Args().First(), options...)
	if err != nil {
		fmt.Printf("Error creating account: %v\n", err)
		os.Exit(1)
	}

	json, _ := json.Marshal(acc)
	fmt.Printf("Account created: %v\n", string(json))
}

func deleteAccount(ctx *cli.Context) {
	if ctx.Args().Len() == 0 {
		fmt.Println("Missing argument: ID")
		return
	}
	client := accountsFromContext(ctx)

	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		fmt.Printf("Error getting namespace: %v\n", err)
		os.Exit(1)
	}

	_, err = client.Delete(context.TODO(), &pb.DeleteAccountRequest{
		Id:      ctx.Args().First(),
		Options: &pb.Options{Namespace: ns},
	})
	if err != nil {
		fmt.Printf("Error deleting account: %v\n", err)
		os.Exit(1)
	}
}

func accountsFromContext(ctx *cli.Context) pb.AccountsService {
	return pb.NewAccountsService("go.micro.auth", client.New(ctx))
}
