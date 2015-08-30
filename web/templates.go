package web

var (
	indexTemplate = `
<html>
	<body>
		<h1>Micro Web</h1>
		<a href="registry">Registry</a>
	</body>
</html>
`

	registryTemplate = `
{{define "T"}}
<html>
	<body>
		<h1>Micro Registry</h1>
		<ul>
			{{range .}}
			<li>{{.Name}}</li>
			{{end}}
		</ul>
	</body>
</html>
{{end}}
`
)
