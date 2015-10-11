package web

var (
	indexTemplate = `
<html>
	<head>
		<title>Micro Web</title>
		<style>
			html {
				font-family: helvetica;
			}
		</style>
	</head>
	<body>
		<h1>Micro Web</h1>
		<h3><a href="registry">Registry</a></h3>
	</body>
</html>
`

	registryTemplate = `
{{define "T"}}
<html>
	<head>
		<title>Micro Web</title>
		<style>
			html {
				font-family: helvetica;
			}
		</style>
	</head>
	<body>
		<h1>Micro Web</h1>
		<h3>Registry</h3>
		<ul>
			{{range .}}
			<li><a href="registry?service={{.Name}}">{{.Name}}</a></li>
			{{end}}
		</ul>
	</body>
</html>
{{end}}
`

	serviceTemplate = `
{{define "T"}}
<html>
	<head>
		<title>Micro Web</title>
		<style>
			html {
				font-family: helvetica;
			}
		</style>
	</head>
	<body>
		<h1>Micro Web</h1>
		<h3>Service {{.Name}}</h3>
		<h4>Nodes</h4>
		<table>
			<thead>
				<td>Id</td>
				<td>Address</td>
				<td>Port</td>
				<td>Metadata</td>
			<thead>
			<tbody>
				{{range .Nodes}}
				<tr>
					<td>{{.Id}}</td>
					<td>{{.Address}}</td>
					<td>{{.Port}}</td>
					<td>{{ range $key, $value := .Metadata }}{{$key}}={{$value}} {{end}}</td>
				</tr>
				{{end}}
			</tbody>
		</table>
		<h4>Endpoints</h4>
		{{range .Endpoints}}
			Name: {{.Name}}</br>
			Metadata: {{ range $key, $value := .Metadata }}{{$key}}={{$value}} {{end}}</br>
			Request:</br>
			<pre>{{format .Request}}</pre>
			Response:</br>
			<pre>{{format .Response}}</pre>
		{{end}}
	</body>
</html>
{{end}}

`
)
