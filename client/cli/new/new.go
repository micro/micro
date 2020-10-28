// Package new generates micro service templates
package new

import (
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/micro/micro/v3/cmd"
	tmpl "github.com/micro/micro/v3/internal/template"
	"github.com/micro/micro/v3/internal/usage"
	"github.com/urfave/cli/v2"
	"github.com/xlab/treeprint"
)

func protoComments(goDir, alias string) []string {
	return []string{
		"\ndownload protoc zip packages (protoc-$VERSION-$PLATFORM.zip) and install:\n",
		"visit https://github.com/protocolbuffers/protobuf/releases",
		"\ndownload protobuf for micro:\n",
		"go get -u github.com/golang/protobuf/proto",
		"go get -u github.com/golang/protobuf/protoc-gen-go",
		"go get github.com/micro/micro/v3/cmd/protoc-gen-micro",
		"\ncompile the proto file " + alias + ".proto:\n",
		"cd " + alias,
		"make proto\n",
	}
}

type config struct {
	// foo
	Alias string
	// github.com/micro/foo
	Dir string
	// $GOPATH/src/github.com/micro/foo
	GoDir string
	// $GOPATH
	GoPath string
	// UseGoPath
	UseGoPath bool
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
		"title": func(s string) string {
			return strings.ReplaceAll(strings.Title(s), "-", "")
		},
		"dehyphen": func(s string) string {
			return strings.ReplaceAll(s, "-", "")
		},
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
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
	if _, err := os.Stat(c.Dir); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", c.Dir)
	}

	// create usage report
	u := usage.New("new")
	// a single request/service
	u.Metrics.Count["requests"] = uint64(1)
	u.Metrics.Count["services"] = uint64(1)
	// send report
	go usage.Report(u)

	// just wait
	<-time.After(time.Millisecond * 250)

	fmt.Printf("Creating service %s\n\n", c.Alias)

	t := treeprint.New()

	// write the files
	for _, file := range c.Files {
		f := filepath.Join(c.Dir, file.Path)
		dir := filepath.Dir(f)

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		addFileToTree(t, file.Path)
		if err := write(c, f, file.Tmpl); err != nil {
			return err
		}
	}

	// print tree
	fmt.Println(t.String())

	for _, comment := range c.Comments {
		fmt.Println(comment)
	}

	// just wait
	<-time.After(time.Millisecond * 250)

	return nil
}

func addFileToTree(root treeprint.Tree, file string) {
	split := strings.Split(file, "/")
	curr := root
	for i := 0; i < len(split)-1; i++ {
		n := curr.FindByValue(split[i])
		if n != nil {
			curr = n
		} else {
			curr = curr.AddBranch(split[i])
		}
	}
	if curr.FindByValue(split[len(split)-1]) == nil {
		curr.AddNode(split[len(split)-1])
	}
}

func Run(ctx *cli.Context) error {
	dir := ctx.Args().First()
	if len(dir) == 0 {
		fmt.Println("specify service name")
		return nil
	}

	// check if the path is absolute, we don't want this
	// we want to a relative path so we can install in GOPATH
	if path.IsAbs(dir) {
		fmt.Println("require relative path as service will be installed in GOPATH")
		return nil
	}

	var goPath string
	var goDir string

	goPath = build.Default.GOPATH

	// don't know GOPATH, runaway....
	if len(goPath) == 0 {
		fmt.Println("unknown GOPATH")
		return nil
	}

	// attempt to split path if not windows
	if runtime.GOOS == "windows" {
		goPath = strings.Split(goPath, ";")[0]
	} else {
		goPath = strings.Split(goPath, ":")[0]
	}
	goDir = filepath.Join(goPath, "src", path.Clean(dir))

	c := config{
		Alias:     dir,
		Comments:  protoComments(goDir, dir),
		Dir:       dir,
		GoDir:     goDir,
		GoPath:    goPath,
		UseGoPath: false,
		Files: []file{
			{"micro.mu", tmpl.Service},
			{"main.go", tmpl.MainSRV},
			{"generate.go", tmpl.GenerateFile},
			{"handler/" + dir + ".go", tmpl.HandlerSRV},
			{"proto/" + dir + ".proto", tmpl.ProtoSRV},
			{"Dockerfile", tmpl.DockerSRV},
			{"Makefile", tmpl.Makefile},
			{"README.md", tmpl.Readme},
			{".gitignore", tmpl.GitIgnore},
		},
	}

	// set gomodule
	if os.Getenv("GO111MODULE") != "off" {
		c.Files = append(c.Files, file{"go.mod", tmpl.Module})
	}

	// create the files
	return create(c)
}

func init() {
	cmd.Register(&cli.Command{
		Name:        "new",
		Usage:       "Create a service template",
		Description: `'micro new' scaffolds a new service skeleton. Example: 'micro new helloworld && cd helloworld'`,
		Action:      Run,
	})
}
