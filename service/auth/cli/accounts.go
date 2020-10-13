package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	pb "github.com/micro/micro/v3/proto/auth"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/urfave/cli/v2"
)

func listAccounts(ctx *cli.Context) error {
	cli := pb.NewAccountsService("auth", client.DefaultClient)

	env, err := util.GetEnv(ctx)
	if err != nil {
		return err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return fmt.Errorf("Error getting namespace: %v", err)
	}

	rsp, err := cli.List(context.DefaultContext, &pb.ListAccountsRequest{
		Options: &pb.Options{Namespace: ns},
	}, client.WithAuthToken())
	if err != nil {
		return fmt.Errorf("Error listing accounts: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	defer w.Flush()

	fmt.Fprintln(w, strings.Join([]string{"ID", "Name", "Scopes", "Metadata"}, "\t\t"))
	for _, r := range rsp.Accounts {
		var metadata string
		for k, v := range r.Metadata {
			metadata = fmt.Sprintf("%v%v=%v ", metadata, k, v)
		}
		scopes := strings.Join(r.Scopes, ", ")

		if len(metadata) == 0 {
			metadata = "n/a"
		}
		if len(scopes) == 0 {
			scopes = "n/a"
		}

		fmt.Fprintln(w, strings.Join([]string{r.Id, r.Name, scopes, metadata}, "\t\t"))
	}

	return nil
}

func createAccount(ctx *cli.Context) error {
	if ctx.Args().Len() == 0 {
		return fmt.Errorf("Missing argument: ID")
	}

	env, err := util.GetEnv(ctx)
	if err != nil {
		return err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return fmt.Errorf("Error getting namespace: %v", err)
	}
	if len(ctx.String("namespace")) > 0 {
		ns = ctx.String("namespace")
	}

	options := []auth.GenerateOption{auth.WithIssuer(ns)}
	if len(ctx.StringSlice("scopes")) > 0 {
		options = append(options, auth.WithScopes(ctx.StringSlice("scopes")...))
	}
	if len(ctx.String("secret")) > 0 {
		options = append(options, auth.WithSecret(ctx.String("secret")))
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
	cli := pb.NewAccountsService("auth", client.DefaultClient)

	env, err := util.GetEnv(ctx)
	if err != nil {
		return err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return fmt.Errorf("Error getting namespace: %v", err)
	}

	_, err = cli.Delete(context.DefaultContext, &pb.DeleteAccountRequest{
		Id:      ctx.Args().First(),
		Options: &pb.Options{Namespace: ns},
	}, client.WithAuthToken())
	if err != nil {
		return fmt.Errorf("Error deleting account: %v", err)
	}

	return nil
}
