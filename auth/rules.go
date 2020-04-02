package auth

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/micro/cli/v2"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/client"
)

func listRules(ctx *cli.Context) {
	client := rulesFromContext(ctx)

	rsp, err := client.List(context.TODO(), &pb.ListRequest{})
	if err != nil {
		fmt.Printf("Error listing rules: %v\n", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	defer w.Flush()

	fmt.Fprintln(w, strings.Join([]string{"Role", "Access", "Resource"}, "\t"))
	for _, r := range rsp.Rules {
		res := strings.Join([]string{r.Resource.Namespace, r.Resource.Type, r.Resource.Name, r.Resource.Endpoint}, ":")
		fmt.Fprintln(w, strings.Join([]string{r.Role, r.Access.String(), res}, "\t"))
	}
}

func createRule(ctx *cli.Context) {
	client := rulesFromContext(ctx)
	r := constructRule(ctx)

	_, err := client.Create(context.TODO(), &pb.CreateRequest{
		Role:     r.Role,
		Access:   r.Access,
		Resource: r.Resource,
	})
	if err != nil {
		fmt.Printf("Error creating rule: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Rule created")
}

func deleteRule(ctx *cli.Context) {
	client := rulesFromContext(ctx)
	r := constructRule(ctx)

	_, err := client.Delete(context.TODO(), &pb.DeleteRequest{
		Role:     r.Role,
		Access:   r.Access,
		Resource: r.Resource,
	})
	if err != nil {
		fmt.Printf("Error creating rule: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Rule deleted")
}

func constructRule(ctx *cli.Context) *pb.Rule {
	var access pb.Access
	switch ctx.String("access") {
	case "granted":
		access = pb.Access_GRANTED
	case "denied":
		access = pb.Access_DENIED
	default:
		fmt.Printf("Invalid access: %v, must be granted or denied\n", ctx.String("access"))
		os.Exit(1)
	}

	resComps := strings.Split(ctx.String("resource"), ":")
	if len(resComps) != 4 {
		fmt.Println("Invalid resource, must be in the format namespace:type:name:endpoint")
	}

	return &pb.Rule{
		Access: access,
		Role:   ctx.String("role"),
		Resource: &pb.Resource{
			Namespace: resComps[0],
			Type:      resComps[1],
			Name:      resComps[2],
			Endpoint:  resComps[3],
		},
	}
}

func rulesFromContext(ctx *cli.Context) pb.RulesService {
	if ctx.Bool("platform") {
		os.Setenv("MICRO_PROXY", "service")
		os.Setenv("MICRO_PROXY_ADDRESS", "proxy.micro.mu:443")
	}

	return pb.NewRulesService("go.micro.auth", client.DefaultClient)
}
