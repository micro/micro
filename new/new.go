package new

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/micro/cli"
)

type config struct {
	Name   string
	Dir    string
	Main   string
	Docker string
}

func create(c config) error {
	// check if dir exists
	if _, err := os.Stat(c.Dir); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", c.Dir)
	}

	fmt.Println("creating service", c.Name)

	// create all required dirs
	for _, d := range []string{c.Dir, c.Dir + "/handler", c.Dir + "/proto/example"} {
		fmt.Println("creating", d)
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	// write main.go

	fmt.Println("creating", c.Dir+"/main.go")
	f, err := os.Create(c.Dir + "/main.go")
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.New("main").Parse(c.Main)
	if err != nil {
		return err
	}

	if err := t.Execute(f, c); err != nil {
		return err
	}

	// write README

	// write Dockerfile

	return nil
}

func run(ctx *cli.Context) {
	namespace := ctx.String("namespace")
	name := ctx.String("name")
	atype := ctx.String("type")
	dir := ctx.Args().First()

	if len(dir) == 0 {
		fmt.Println("specify service name")
		return
	}

	if len(namespace) == 0 {
		fmt.Println("namespace not defined")
		return
	}

	if len(atype) == 0 {
		fmt.Println("type not defined")
		return
	}

	// check if the path is absolute, we don't want this
	// we want to a relative path so we can install in GOPATH
	if path.IsAbs(dir) {
		fmt.Println("require relative path as service will be installed in GOPATH")
		return
	}

	goPath := os.Getenv("GOPATH")

	// don't know GOPATH, runaway....
	if len(goPath) == 0 {
		fmt.Println("unknown GOPATH")
		return
	}

	p := strings.Split(goPath, ":")[0]
	dir = path.Join(p, "src", path.Clean(dir))
	parts := strings.Split(dir, "/")

	// if name not specified create it from namespace.type.dir
	if len(name) == 0 {
		name = strings.Join([]string{namespace, atype, parts[len(parts)-1]}, ".")
	}

	var c config

	switch atype {
	case "srv":
		// create srv config
		c = config{
			Name:   name,
			Dir:    dir,
			Main:   srvMainTemplate,
			Docker: srvDockerTemplate,
		}
	default:
		fmt.Println("Unknown type", atype)
		return
	}

	if err := create(c); err != nil {
		fmt.Println(err)
		return
	}

	// create proto/example/example.proto
	// create proto/example/example.pb.go
	// create handler/example.go
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:  "new",
			Usage: "Create a new micro service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace",
					Usage: "Namespace for the service e.g com.example",
					Value: "go.micro",
				},
				cli.StringFlag{
					Name:  "type",
					Usage: "Type of service e.g api, srv, web",
					Value: "srv",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "Name of service e.g com.example.srv.service (defaults to namespace.type.[args[len(args)-1])",
				},
			},
			Action: func(c *cli.Context) {
				run(c)
			},
		},
	}
}
