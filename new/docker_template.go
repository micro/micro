package new

var (
	srvDockerTemplate = `FROM alpine:3.2
ADD {{.Alias}}-{{.Type}} /{{.Alias}}-{{.Type}}
ENTRYPOINT [ "/{{.Alias}}-{{.Type}}" ]
`
)
