package template

var (
	Readme = `# {{title .Alias}} Service

This is the {{title .Alias}} service

Generated with

` + "```" +
		`
{{.Command}}
` + "```" + `

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: {{.FQDN}}
- Type: {{.Type}}
- Alias: {{.Alias}}

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend consul.

` + "```" +
		`
# install consul
brew install consul

# run consul
consul agent -dev
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
./{{.Alias}}-{{.Type}}
` + "```" + `

Build a docker image
` + "```" +
		`
make docker
` + "```"

	ReadmeFNC = `# {{title .Alias}} Function

This is the {{title .Alias}} function

Generated with

` + "```" +
		`
{{.Command}}
` + "```" + `

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: {{.FQDN}}
- Type: {{.Type}}
- Alias: {{.Alias}}

## Dependencies

Micro functions depend on service discovery. The default is consul.

` + "```" +
		`
# install consul
brew install consul

# run consul
consul agent -dev
` + "```" + `

## Usage

A Makefile is included for convenience

Build the binary

` + "```" +
		`
make build
` + "```" + `

Run the function once
` + "```" +
		`
./{{.Alias}}-{{.Type}}
` + "```" + `

Build a docker image
` + "```" +
		`
make docker
` + "```"
)
