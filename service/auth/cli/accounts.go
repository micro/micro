package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/micro/cli/v2"
	goauth "github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/service/auth"
	pb "github.com/micro/micro/v3/service/auth/proto"
)

func listAccounts(ctx *cli.Context) error {
	accountClient := pb.NewAccountsService("go.micro.auth")

	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return fmt.Errorf("Error getting namespace: %v", err)
	}

	rsp, err := accountClient.List(context.TODO(), &pb.ListAccountsRequest{
		Options: &pb.Options{Namespace: ns},
	}, client.WithServiceToken())
	if err != nil {
		return fmt.Errorf("Error listing accounts: %v", err)
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

	return nil
}

func createAccount(ctx *cli.Context) error {
	if ctx.Args().Len() == 0 {
		return fmt.Errorf("Missing argument: ID")
	}

	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return fmt.Errorf("Error getting namespace: %v", err)
	}

	options := []goauth.GenerateOption{goauth.WithIssuer(ns)}
	if len(ctx.StringSlice("scopes")) > 0 {
		options = append(options, goauth.WithScopes(ctx.StringSlice("scopes")...))
	}
	if len(ctx.String("secret")) > 0 {
		options = append(options, goauth.WithSecret(ctx.String("secret")))
	}
	acc, err := auth.Generate(ctx.Args().First(), options...)
	if err != nil {
		return fmt.Errorf("Error creating account: %v", err)
	}

	json, _ := json.Marshal(acc)
	fmt.Printf("Account created: %v\n", string(json))
	return nil
}

func deleteAccount(ctx *cli.Context) error {
	if ctx.Args().Len() == 0 {
		return fmt.Errorf("Missing argument: ID")
	}
	accountClient := pb.NewAccountsService("go.micro.auth")

	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return fmt.Errorf("Error getting namespace: %v", err)
	}

	_, err = accountClient.Delete(context.TODO(), &pb.DeleteAccountRequest{
		Id:      ctx.Args().First(),
		Options: &pb.Options{Namespace: ns},
	}, client.WithServiceToken())
	if err != nil {
		return fmt.Errorf("Error deleting account: %v", err)
	}

	return nil
}
