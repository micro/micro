package template

var (
	Readme = `# {{title .Alias}} {{title .Type}}

This is the {{title .Alias}} service with fqdn {{.FQDN}}.

## Getting Started

### Prerequisites

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
` + "```" +
		`
$ consul agent -dev -advertise=127.0.0.1
` + "```" +
		`

### Run Service

` + "```" +
		`
$ go run main.go
` + "```" +
		`

### Building a container

If you would like to build the docker container do the following
` + "```" +
		`
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o {{.Alias}}-{{.Type}} ./main.go
docker build -t {{.Alias}}-{{.Type}} .

` + "```"

	ReadmeFNC = `# {{title .Alias}} {{title .Type}}

This is the {{title .Alias}} function with fqdn {{.FQDN}}.

## Getting Started

### Service Discovery

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
` + "```" +
		`
$ consul agent -dev
` + "```" +
		`
### Micro Toolkit

Install Micro

` + "```" +
		`
go get github.com/micro/micro
` + "```" +
		`

### Run Function

` + "```" +
		`
$ micro run -r {{.Dir}}
` + "```" +
		`

### Building a container

If you would like to build the docker container do the following
` + "```" +
		`
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o {{.Alias}}-{{.Type}} ./main.go
docker build -t {{.Alias}}-{{.Type}} .

` + "```"
)
