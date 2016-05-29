package new

var (
	srvDockerTemplate = `FROM alpine:3.2
ADD {{.Alias}}-{{.Type}} /{{.Alias}}-{{.Type}}
ENTRYPOINT [ "/{{.Alias}}-{{.Type}}" ]
`

	webDockerTemplate = `FROM alpine:3.2
ADD html /html
ADD {{.Alias}}-{{.Type}} /{{.Alias}}-{{.Type}}
WORKDIR /
ENTRYPOINT [ "/{{.Alias}}-{{.Type}}" ]
`
)
