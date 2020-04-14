package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/emicklei/proto"
	registry "github.com/micro/go-micro/v2/registry/service/proto"
	"github.com/spf13/cobra"
)

var (
	path           string
	goImportPrefix string
	skip           = []string{"node_modules"}
)

var docsCmd = &cobra.Command{
	Use:   "doc-gen",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		generate()
	},
}

func init() {
	docsCmd.Flags().StringVarP(&path, "path", "p", "", "Source directory to read from")
	docsCmd.Flags().StringVarP(&goImportPrefix, "go-import-prefix", "g", "github.com/micro/services", "For go client examples, the path must be known ie \"github.com/micro/service\"")
	rootCmd.AddCommand(docsCmd)
}

func generate() {
	protoPaths := getProtoPaths(path)
	for serviceName, protos := range protoPaths {
		err := generateMarkdowns(serviceName, protos)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error processing service \"%v\": %v", serviceName, err))
		}
	}
}

func getProtoPaths(p string) map[string][]string {
	protoPaths := map[string][]string{}
	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && strings.HasPrefix(".", info.Name()) || inSkip(info.Name()) {
			return filepath.SkipDir
		}
		if !strings.HasSuffix(info.Name(), ".proto") {
			return nil
		}

		serviceName := strings.Split(strings.Replace(path, p+"/", "", 1), "/")[0]
		protoPs, ok := protoPaths[serviceName]
		if !ok {
			protoPaths[serviceName] = []string{path}
		} else {
			protoPaths[serviceName] = append(protoPs, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return protoPaths
}

func generateMarkdowns(serviceName string, protoPaths []string) error {
	reader, _ := os.Open(protoPaths[0])
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	services := map[string]*proto.Service{}
	messages := map[string]*proto.Message{}
	rpcs := map[string]*proto.RPC{}
	proto.Walk(definition,
		proto.WithService(func(s *proto.Service) {
			services[s.Name] = s
		}),
		proto.WithMessage(func(s *proto.Message) {
			messages[s.Name] = s
		}),
		proto.WithRPC(func(s *proto.RPC) {
			rpcs[s.Name] = s
		}))

	getNormalisedFields := func(messgeName string) []*proto.NormalField {
		msg, ok := messages[messgeName]
		if !ok {
			return []*proto.NormalField{}
		}
		ret := []*proto.NormalField{}
		for _, element := range msg.Elements {
			if nf, ok := element.(*proto.NormalField); ok {
				ret = append(ret, nf)
			}
		}
		return ret
	}

	var getMessageNames func(fields []*proto.NormalField) []string
	getMessageNames = func(fields []*proto.NormalField) []string {
		ret := []string{}
		for _, field := range fields {
			// this means it's a custom message
			if field.Type != strings.ToLower(field.Type) {
				ret = append(ret, field.Type)
				fields := getNormalisedFields(field.Type)
				ret = append(ret, getMessageNames(fields)...)
			}
		}
		return ret
	}
	tmpl, err := template.New(serviceName).Funcs(template.FuncMap{
		"toJSON":          toJSON,
		"getNormalFields": getNormalisedFields,
		"commentLines": func(comment *proto.Comment) string {
			if comment == nil {
				return ""
			}
			return strings.Join(comment.Lines, "\n")
		},
		"packageNamify": func(s string) string {
			return strings.Replace(s, "-", "_", -1)
		},
		"parentService": func(s *proto.RPC) *proto.Service {
			return s.Parent.(*proto.Service)
		},
		"messagesUsedByReqRsp": func(s *proto.RPC) []string {
			reqFields := getNormalisedFields(s.RequestType)
			rspFields := getNormalisedFields(s.ReturnsType)
			return unique(append(getMessageNames(reqFields), getMessageNames(rspFields)...))
		},
	}).Parse(serviceTemplate)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("./" + serviceName + ".md")
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, map[string]interface{}{
		"serviceName":    serviceName,
		"services":       services,
		"messages":       messages,
		"rpcs":           rpcs,
		"goImportPrefix": goImportPrefix,
	})
	if err != nil {
		return err
	}
	return nil
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func inSkip(dirname string) bool {
	for _, v := range skip {
		if dirname == v {
			return true
		}
	}
	return false
}

func typeToDefaultValue(typ string) string {
	switch typ {
	case "string":
		return `""`
	case "int":
	case "int32":
	case "int64":
		return "0"
	case "bool":
		return "false"
	}
	return "{}"
}

// similarly incredibly ugly implementation can be found at https://github.com/micro/platform/blob/master/web/app/src/app/endpoint-list/endpoint-list.component.ts#L68
func toJSON(input *registry.Value, level int) string {
	if input == nil {
		return "no data"
	}
	indent := strings.Repeat("  ", level)
	const fieldSeparator = ",\n"
	if len(input.GetValues()) == 0 {
		if level == 1 {
			return "{}"
		}
		return indent + "\"" + input.GetName() + "\"" + ": " + typeToDefaultValue(input.GetType())
	}
	lines := []string{}
	for _, value := range input.Values {
		lines = append(lines, toJSON(value, level+1))
	}
	if level == 1 {
		return indent + "{\n" + strings.Join(lines, fieldSeparator) + "\n" + indent + "}"
	}
	return indent + "\"" + input.GetName() + "\"" + ": {\n" + strings.Join(lines, fieldSeparator) + "\n" + indent + "}"
}
