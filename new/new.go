// Package new generates micro service templates
package new

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/micro/cli"
	tmpl "github.com/micro/micro/internal/template"
	"github.com/xlab/treeprint"
)

type config struct {
	// foo
	Alias string
	// micro new example -type
	Command string
	// go.micro
	Namespace string
	// api, srv, web, fnc
	Type string
	// go.micro.srv.foo
	FQDN string
	// github.com/micro/foo
	Dir string
	// $GOPATH/src/github.com/micro/foo
	GoDir string
	// $GOPATH
	GoPath string
	// Files
	Files []file
	// Comments
	Comments []string
	// Plugins registry=etcd:broker=nats
	Plugins []string
}

type file struct {
	Path string
	Tmpl string
}

func write(c config, file, tmpl string) error {
	fn := template.FuncMap{
		"title": strings.Title,
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.New("f").Funcs(fn).Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(f, c)
}

func create(c config) error {
	// check if dir exists
	if _, err := os.Stat(c.GoDir); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", c.GoDir)
	}

	fmt.Printf("Creating service %s in %s\n\n", c.FQDN, c.GoDir)

	t := treeprint.New()

	nodes := map[string]treeprint.Tree{}
	nodes[c.GoDir] = t

	// write the files
	for _, file := range c.Files {
		f := filepath.Join(c.GoDir, file.Path)
		dir := filepath.Dir(f)

		b, ok := nodes[dir]
		if !ok {
			d, _ := filepath.Rel(c.GoDir, dir)
			b = t.AddBranch(d)
			nodes[dir] = b
		}

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		p := filepath.Base(f)

		b.AddNode(p)
		if err := write(c, f, file.Tmpl); err != nil {
			return err
		}
	}

	// print tree
	fmt.Println(t.String())

	for _, comment := range c.Comments {
		fmt.Println(comment)
	}

	return nil
}

func run(ctx *cli.Context) {
	namespace := ctx.String("namespace")
	alias := ctx.String("alias")
	fqdn := ctx.String("fqdn")
	atype := ctx.String("type")
	dir := ctx.Args().First()
	useGoPath := ctx.Bool("gopath")
	var plugins []string

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

	// set the command
	command := fmt.Sprintf("micro new %s", dir)
	if len(namespace) > 0 {
		command += " --namespace=" + namespace
	}
	if len(alias) > 0 {
		command += " --alias=" + alias
	}
	if len(fqdn) > 0 {
		command += " --fqdn=" + fqdn
	}
	if len(atype) > 0 {
		command += " --type=" + atype
	}
	if plugins := ctx.StringSlice("plugin"); len(plugins) > 0 {
		command += " --plugin=" + strings.Join(plugins, ":")
	}

	// check if the path is absolute, we don't want this
	// we want to a relative path so we can install in GOPATH
	if path.IsAbs(dir) {
		fmt.Println("require relative path as service will be installed in GOPATH")
		return
	}

	var goPath string
	var goDir string

	// only set gopath if told to use it
	if useGoPath {
		goPath = os.Getenv("GOPATH")

		// don't know GOPATH, runaway....
		if len(goPath) == 0 {
			fmt.Println("unknown GOPATH")
			return
		}

		// attempt to split path if not windows
		if runtime.GOOS == "windows" {
			goPath = strings.Split(goPath, ";")[0]
		} else {
			goPath = strings.Split(goPath, ":")[0]
		}
		goDir = filepath.Join(goPath, "src", path.Clean(dir))
	} else {
		goDir = path.Clean(dir)
	}

	if len(alias) == 0 {
		// set as last part
		alias = filepath.Base(dir)
	}

	if len(fqdn) == 0 {
		fqdn = strings.Join([]string{namespace, atype, alias}, ".")
	}

	for _, plugin := range ctx.StringSlice("plugin") {
		// registry=etcd:broker=nats
		for _, p := range strings.Split(plugin, ":") {
			// registry=etcd
			parts := strings.Split(p, "=")
			if len(parts) < 2 {
				continue
			}
			plugins = append(plugins, path.Join(parts...))
		}
	}

	var c config

	switch atype {
	case "fnc":
		// create srv config
		c = config{
			Alias:     alias,
			Command:   command,
			Namespace: namespace,
			Type:      atype,
			FQDN:      fqdn,
			Dir:       dir,
			GoDir:     goDir,
			GoPath:    goPath,
			Plugins:   plugins,
			Files: []file{
				{"main.go", tmpl.MainFNC},
				{"plugin.go", tmpl.Plugin},
				{"handler/example.go", tmpl.HandlerFNC},
				{"subscriber/example.go", tmpl.SubscriberFNC},
				{"proto/example/example.proto", tmpl.ProtoFNC},
				{"Dockerfile", tmpl.DockerFNC},
				{"Makefile", tmpl.Makefile},
				{"README.md", tmpl.ReadmeFNC},
			},
			Comments: []string{
				"\ndownload protobuf for micro:\n",
				"brew install protobuf",
				"go get -u github.com/golang/protobuf/{proto,protoc-gen-go}",
				"go get -u github.com/micro/protoc-gen-micro",
				"\ncompile the proto file example.proto:\n",
				"cd " + goDir,
				"protoc --proto_path=. --go_out=. --micro_out=. proto/example/example.proto\n",
			},
		}
	case "srv":
		// create srv config
		c = config{
			Alias:     alias,
			Command:   command,
			Namespace: namespace,
			Type:      atype,
			FQDN:      fqdn,
			Dir:       dir,
			GoDir:     goDir,
			GoPath:    goPath,
			Plugins:   plugins,
			Files: []file{
				{"main.go", tmpl.MainSRV},
				{"plugin.go", tmpl.Plugin},
				{"handler/example.go", tmpl.HandlerSRV},
				{"subscriber/example.go", tmpl.SubscriberSRV},
				{"proto/example/example.proto", tmpl.ProtoSRV},
				{"Dockerfile", tmpl.DockerSRV},
				{"Makefile", tmpl.Makefile},
				{"README.md", tmpl.Readme},
			},
			Comments: []string{
				"\ndownload protobuf for micro:\n",
				"brew install protobuf",
				"go get -u github.com/golang/protobuf/{proto,protoc-gen-go}",
				"go get -u github.com/micro/protoc-gen-micro",
				"\ncompile the proto file example.proto:\n",
				"cd " + goDir,
				"protoc --proto_path=. --go_out=. --micro_out=. proto/example/example.proto\n",
			},
		}
	case "api":
		// create api config
		c = config{
			Alias:     alias,
			Command:   command,
			Namespace: namespace,
			Type:      atype,
			FQDN:      fqdn,
			Dir:       dir,
			GoDir:     goDir,
			GoPath:    goPath,
			Plugins:   plugins,
			Files: []file{
				{"main.go", tmpl.MainAPI},
				{"plugin.go", tmpl.Plugin},
				{"client/example.go", tmpl.WrapperAPI},
				{"handler/example.go", tmpl.HandlerAPI},
				{"proto/example/example.proto", tmpl.ProtoAPI},
				{"Makefile", tmpl.Makefile},
				{"Dockerfile", tmpl.DockerSRV},
				{"README.md", tmpl.Readme},
			},
			Comments: []string{
				"\ndownload protobuf for micro:\n",
				"brew install protobuf",
				"go get -u github.com/golang/protobuf/{proto,protoc-gen-go}",
				"go get -u github.com/micro/protoc-gen-micro",
				"\ncompile the proto file example.proto:\n",
				"cd " + goDir,
				"protoc --proto_path=. --go_out=. --micro_out=. proto/example/example.proto\n",
			},
		}
	case "web":
		// create srv config
		c = config{
			Alias:     alias,
			Command:   command,
			Namespace: namespace,
			Type:      atype,
			FQDN:      fqdn,
			Dir:       dir,
			GoDir:     goDir,
			GoPath:    goPath,
			Plugins:   plugins,
			Files: []file{
				{"main.go", tmpl.MainWEB},
				{"plugin.go", tmpl.Plugin},
				{"handler/handler.go", tmpl.HandlerWEB},
				{"html/index.html", tmpl.HTMLWEB},
				{"Dockerfile", tmpl.DockerWEB},
				{"Makefile", tmpl.Makefile},
				{"README.md", tmpl.Readme},
			},
			Comments: []string{},
		}
	default:
		fmt.Println("Unknown type", atype)
		return
	}

	if err := create(c); err != nil {
		fmt.Println(err)
		return
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:  "new",
			Usage: "Create a new micro service by specifying a directory path relative to your $GOPATH",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace",
					Usage: "Namespace for the service e.g com.example",
					Value: "go.micro",
				},
				cli.StringFlag{
					Name:  "type",
					Usage: "Type of service e.g api, fnc, srv, web",
					Value: "srv",
				},
				cli.StringFlag{
					Name:  "fqdn",
					Usage: "FQDN of service e.g com.example.srv.service (defaults to namespace.type.alias)",
				},
				cli.StringFlag{
					Name:  "alias",
					Usage: "Alias is the short name used as part of combined name if specified",
				},
				cli.StringSliceFlag{
					Name:  "plugin",
					Usage: "Specify plugins e.g --plugin=registry=etcd:broker=nats or use flag multiple times",
				},
				cli.BoolTFlag{
					Name:  "gopath",
					Usage: "Create the service in the gopath. Defaults to true.",
				},
			},
			Action: func(c *cli.Context) {
				run(c)
			},
		},
	}
}
