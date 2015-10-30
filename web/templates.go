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
		<h3><a href="query">Query</a></h3>
	</body>
</html>
`
	queryTemplate = `
{{define "T"}}
<html>
	<head>
		<title>Micro Web</title>
		<style>
			html {
				font-family: helvetica;
			}
			div {
				width: 45%;
				display: inline-block;
				height: 100%;
			}
			input, textarea {
				width: 100%;
			}
			#pre {
				float: right;
			}
			pre {
				word-wrap: break-word;
				padding: 10px;
				border: 1px solid;
			}
		</style>
	</head>
	<body>
		<h1>Micro Web</h1>
		<h3>Query</h3>
		<div>
			<form id="query-form" onsubmit="return query();">
				<p><input type=text name=service id=service placeholder=service /></p>
				<p><input type=text name=method id=method placeholder=method /></p>
				<p><textarea name=request id=request rows=30></textarea></p>
				<p><button>Go!</button></p>
			</form>
		</div>
		<div id="pre">
			<pre id="response"></pre>
		</div>
		<script>
			function query() {
				var req = new XMLHttpRequest()
				req.onreadystatechange = function() {
					if (req.readyState == 4 && req.status == 200) {
						document.getElementById("response").innerText = JSON.stringify(JSON.parse(req.responseText), null, 2);
						console.log(req.responseText);
					}
				}
				var request = {
					"service": document.forms[0].elements["service"].value,
					"method": document.forms[0].elements["method"].value,
					"request": JSON.parse(document.forms[0].elements["request"].value)
				}
				req.open("POST", "/rpc", true);
				req.setRequestHeader("Content-type","application/json");				
				req.send(JSON.stringify(request));

				return false;
			};	
		</script>
	</body>
</html>
{{end}}
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
