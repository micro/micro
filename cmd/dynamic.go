package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/cmd"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/micro/v2/client/cli/namespace"
	"github.com/micro/micro/v2/client/cli/util"
	inclient "github.com/micro/micro/v2/internal/client"
)

// lookupService queries the service for a service with the given alias. If
// no services are found for a given alias, the registry will return nil and
// the error will also be nil. An error is only returned if there was an issue
// listing from the registry.
func lookupService(ctx *cli.Context) (*registry.Service, error) {
	// get the namespace to query the services from
	dom, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return nil, err
	}

	// lookup from the registry in the current namespace
	reg := *cmd.DefaultCmd.Options().Registry
	srvs, err := reg.ListServices(registry.ListDomain(dom))
	if err != nil {
		return nil, err
	}

	// filter to services with the correct suffix
	for _, s := range srvs {
		if strings.HasSuffix(s.Name, "."+ctx.Args().First()) {
			srvs, err = reg.GetService(s.Name, registry.GetDomain(dom))
			if err == nil && len(srvs) > 0 {
				return srvs[0], nil
			}
		}
	}

	// check for the service in the default namespace also
	if dom == registry.DefaultDomain {
		return nil, nil
	}
	srvs, err = reg.ListServices()
	if err != nil {
		return nil, err
	}

	// filter to services with the correct suffix
	for _, s := range srvs {
		if strings.HasSuffix(s.Name, "."+ctx.Args().First()) {
			srvs, err = reg.GetService(s.Name, registry.GetDomain(dom))
			if err == nil && len(srvs) > 0 {
				return srvs[0], nil
			}
		}
	}

	// no service was found
	return nil, nil
}

// formatServiceUsage returns a string containing the service usage.
func formatServiceUsage(srv *registry.Service, alias string) string {
	commands := make([]string, len(srv.Endpoints))
	for i, e := range srv.Endpoints {
		// map "Helloworld.Call" to "helloworld.call"
		name := strings.ToLower(e.Name)

		// remove the prefix if it is the service name, e.g. rather than
		// "micro run helloworld helloworld call", it would be
		// "micro run helloworld call".
		name = strings.TrimPrefix(name, alias+".")

		// instead of "micro run helloworld foo.bar", the command should
		// be "micro run helloworld foo bar".
		commands[i] = strings.Replace(name, ".", " ", 1)
	}

	// sort the command names alphabetically
	sort.Strings(commands)

	result := fmt.Sprintf("NAME:\n\t%v\n\n", srv.Name)
	result += fmt.Sprintf("VERSION:\n\t%v\n\n", srv.Version)
	result += fmt.Sprintf("USAGE:\n\tmicro %v [flags] [command]\n\n", alias)
	result += fmt.Sprintf("COMMANDS:\n\t%v\n\n", strings.Join(commands, "\n\t"))
	return result
}

// callService will call a service using the arguments and flags provided
// in the context. It will print the result or error to stdout. If there
// was an error performing the call, it will be returned.
func callService(srv *registry.Service, ctx *cli.Context) error {
	// parse the flags and args
	args, flags, err := splitCmdArgs(ctx)
	if err != nil {
		return err
	}

	// construct the endpoint
	endpoint, err := constructEndpoint(args)
	if err != nil {
		return err
	}

	// ensure the endpoint exists on the service
	var ep *registry.Endpoint
	for _, e := range srv.Endpoints {
		if e.Name == endpoint {
			ep = e
			break
		}
	}
	if ep == nil {
		return fmt.Errorf("Endpoint %v not found for service %v", endpoint, srv.Name)
	}

	// parse the flags
	body, err := flagsToRequest(flags, ep.Request)
	if err != nil {
		return err
	}

	// construct and execute the request using the json content type
	cli := inclient.New(ctx)
	req := cli.NewRequest(srv.Name, endpoint, body, client.WithContentType("application/json"))
	var rsp json.RawMessage
	if err := cli.Call(ctx.Context, req, &rsp); err != nil {
		return err
	}

	// format the response
	var out bytes.Buffer
	defer out.Reset()
	if err := json.Indent(&out, rsp, "", "\t"); err != nil {
		return err
	}
	out.WriteTo(os.Stdout)

	return nil
}

// splitCmdArgs takes a cli context and parses out the args and flags, for
// example "micro helloworld --name=foo call apple" would result in "call",
// "apple" as args and {"name":"foo"} as the flags.
func splitCmdArgs(ctx *cli.Context) ([]string, map[string]string, error) {
	args := []string{}
	flags := map[string]string{}

	for _, a := range ctx.Args().Slice() {
		if !strings.HasPrefix(a, "--") {
			args = append(args, a)
			continue
		}

		// comps would be "foo", "bar" for "--foo=bar"
		comps := strings.Split(strings.TrimPrefix(a, "--"), "=")
		if len(comps) != 2 {
			return nil, nil, fmt.Errorf("Invalid flag: %v. Expected format: --foo=bar", a)
		}
		flags[comps[0]] = comps[1]
	}

	return args, flags, nil
}

// constructEndpoint takes a slice of args and converts it into a valid endpoint
// such as Helloworld.Call or Foo.Bar, it will return an error if an invalid number
// of arguments were provided
func constructEndpoint(args []string) (string, error) {
	var epComps []string
	switch len(args) {
	case 2:
		epComps = args
	case 3:
		epComps = args[1:3]
	default:
		return "", fmt.Errorf("Incorrect number of arguments")
	}

	// transform the endpoint components, e.g ["helloworld", "call"] to the
	// endpoint name: "Helloworld.Call".
	return fmt.Sprintf("%v.%v", strings.Title(epComps[0]), strings.Title(epComps[1])), nil
}

// flagsToRequeest parses a set of flags, e.g {name:"Foo", "options_surname","Bar"} and
// converts it into a request body. If the key is not a valid object in the request, an
// error will be returned.
func flagsToRequest(flags map[string]string, req *registry.Value) (map[string]interface{}, error) {
	result := map[string]interface{}{}

loop:
	for key, value := range flags {
		for _, attr := range req.Values {
			// matches at a top level
			if attr.Name == key {
				result[key] = value
				continue loop
			}

			// check for matches at the second level
			if !strings.HasPrefix(key, attr.Name+"_") {
				continue
			}
			for _, attr2 := range attr.Values {
				if attr.Name+"_"+attr2.Name != key {
					continue
				}

				if _, ok := result[attr.Name]; !ok {
					result[attr.Name] = map[string]string{}
				} else if _, ok := result[attr.Name].(map[string]string); !ok {
					return nil, fmt.Errorf("Error parsing request, duplicate key: %v", key)
				}
				result[attr.Name].(map[string]string)[attr2.Name] = value
				continue loop
			}
		}

		return nil, fmt.Errorf("Unknown flag: %v", key)
	}

	return result, nil
}
