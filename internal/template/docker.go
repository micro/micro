package template

var (
	DockerSRV = `FROM alpine
ADD {{.Alias}} /{{.Alias}}
ENTRYPOINT [ "/{{.Alias}}" ]
`
)
