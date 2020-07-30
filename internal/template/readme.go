package template

var (
	Readme = `# {{title .Alias}} Service

This is the {{title .Alias}} service

Generated with

` + "```" +
		`
micro new {{.Alias}}
` + "```" + `

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- Alias: {{.Alias}}

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

` + "```" +
		`
# install etcd
brew install etcd

# run etcd
etcd
` + "```" + `

## Usage

A Makefile is included for convenience

Build the binary

` + "```" +
		`
make build
` + "```" + `

Run the service
` + "```" +
		`
./{{.Alias}}
` + "```" + `

Build a docker image
` + "```" +
		`
make docker
` + "```"
)
