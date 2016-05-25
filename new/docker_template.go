package new

var (
	srvDockerTemplate = `
FROM alpine:3.2
ADD {{.Name}}-{{.Type}} /{{.Name}}-{{.Type}}
ENTRYPOINT [ "/{{.Name}}-{{.Type}}" ]
`
)
