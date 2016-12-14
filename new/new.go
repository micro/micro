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
)

type config struct {
	// foo
	Alias string
	// go.micro
	Namespace string
	// api, srv, web
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
}

type file struct {
	Path string
	Tmpl string
}

func write(c config, file, tmpl string) error {
	fn := template.FuncMap{
		"title": strings.Title,
	}

	fmt.Println("creating", file)
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

	fmt.Println("creating service", c.FQDN)

	// write the files
	for _, file := range c.Files {
		f := filepath.Join(c.GoDir, file.Path)
		dir := path.Dir(f)

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Println("creating", dir)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		if err := write(c, f, file.Tmpl); err != nil {
			return err
		}
	}

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

	// attempt to split path if not windows
	if runtime.GOOS != "windows" {
		goPath = strings.Split(goPath, ":")[0]
	}

	goDir := filepath.Join(goPath, "src", path.Clean(dir))

	if len(alias) == 0 {
		// set as last part
		alias = filepath.Base(dir)
	}

	if len(fqdn) == 0 {
		fqdn = strings.Join([]string{namespace, atype, alias}, ".")
	}

	var c config

	switch atype {
	case "srv":
		// create srv config
		c = config{
			Alias:     alias,
			Namespace: namespace,
			Type:      atype,
			FQDN:      fqdn,
			Dir:       dir,
			GoDir:     goDir,
			GoPath:    goPath,
			Files: []file{
				{"main.go", tmpl.MainSRV},
				{"handler/example.go", tmpl.HandlerSRV},
				{"subscriber/example.go", tmpl.SubscriberSRV},
				{"proto/example/example.proto", tmpl.ProtoSRV},
				{"Dockerfile", tmpl.DockerSRV},
				{"README.md", tmpl.Readme},
			},
			Comments: []string{
				"\ndownload protobuf for micro:\n",
				"go get github.com/micro/protobuf/{proto,protoc-gen-go}",
				"\ncompile the proto file example.proto:\n",
				fmt.Sprintf("protoc -I%s \\\n\t--go_out=plugins=micro:%s \\\n\t%s\n",
					goPath+"/src", goPath+"/src", goDir+"/proto/example/example.proto"),
			},
		}
	case "api":
		// create api config
		c = config{
			Alias:     alias,
			Namespace: namespace,
			Type:      atype,
			FQDN:      fqdn,
			Dir:       dir,
			GoDir:     goDir,
			GoPath:    goPath,
			Files: []file{
				{"main.go", tmpl.MainAPI},
				{"client/example.go", tmpl.WrapperAPI},
				{"handler/example.go", tmpl.HandlerAPI},
				{"proto/example/example.proto", tmpl.ProtoAPI},
				{"Dockerfile", tmpl.DockerSRV},
				{"README.md", tmpl.Readme},
			},
			Comments: []string{
				"\ndownload protobuf for micro:\n",
				"go get github.com/micro/protobuf/{proto,protoc-gen-go}",
				"\ncompile the proto file example.proto:\n",
				fmt.Sprintf("protoc -I%s \\\n\t--go_out=plugins=micro:%s \\\n\t%s\n",
					goPath+"/src", goPath+"/src", goDir+"/proto/example/example.proto"),
			},
		}
	case "web":
		// create srv config
		c = config{
			Alias:     alias,
			Namespace: namespace,
			Type:      atype,
			FQDN:      fqdn,
			Dir:       dir,
			GoDir:     goDir,
			GoPath:    goPath,
			Files: []file{
				{"main.go", tmpl.MainWEB},
				{"handler/handler.go", tmpl.HandlerWEB},
				{"html/index.html", tmpl.HTMLWEB},
				{"Dockerfile", tmpl.DockerWEB},
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
			Usage: "Create a new Micro service by specifying a directory path relative to your $GOPATH",
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
					Name:  "fqdn",
					Usage: "FQDN of service e.g com.example.srv.service (defaults to namespace.type.alias)",
				},
				cli.StringFlag{
					Name:  "alias",
					Usage: "Alias is the short name used as part of combined name if specified",
				},
			},
			Action: func(c *cli.Context) {
				run(c)
			},
		},
	}
}
