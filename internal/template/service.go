package template

var (
	// Service template is the Micro .mu definition of a service
	Service = `service {{lower .Alias}}
`
)
