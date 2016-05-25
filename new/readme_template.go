package new

var (
	readmeTemplate = `
# {{.Name}} Service

This is the {{.Name}} service. It's of type {{.Type}} with namespace {{.Namespace}}

### Prerequisites

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
` + "```" +
		`
$ consul agent -dev -advertise=127.0.0.1
` + "```" +
		`
Run Service
` + "```" +
		`
$ go run main.go
` + "```"
)
