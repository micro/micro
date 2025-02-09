package html

var (
	LoginTemplate = `
{{define "title"}}Login{{end}}
{{define "heading"}}{{end}}
{{define "style" }}{{end}}
{{define "content"}}
	<div id="error">{{ .error }}</div>
	<div id="login"></div>
{{end}}
{{define "script"}}<script>renderLogin();</script>{{end}}
`

	LayoutTemplate = `
{{define "layout"}}
<html>
	<head>
		<title>{{ template "title" . }} | Micro</title>
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
	        <link rel="stylesheet" href="/assets/mu.css">
		<style>
		{{ template "style" . }}
		</style>
		{{ template "head" . }}
	</head>
	<body>
	  <div id="header">
            <a id="logo" href="/">Micro</a>
            {{if .Token}}
	    <ul id="menu">
	        <li><a href="/">Home</a></li>
	        <li><a href="/services">Services</a></li>
	        <li><a href="/logout">Logout</a></li>
	    </ul>
	   {{end}}
          </div>
          <div id="container">
              <div id="heading">{{ template "heading" . }}</div>
              <div id="content">{{ template "content" . }}</div>
          </div>
	  <div id="footer"></div>
	  <script src="/assets/mu.js"></script>
	  <script src="/assets/jquery.min.js"></script>
	  <script type="text/javascript">
		function toggle(e) {
		      var ev = window.event? event : e
		      if (ev.keyCode == 80 && ev.ctrlKey && ev.shiftKey) {
			var el = document.getElementById("dev");
			if (el.style.display == "none") {
			  el.style.display = "block";
			} else {
			  el.style.display = "none";
			}
		    }
		}

		document.onkeydown = toggle;

		// set the api url
		setAPI({{.ApiURL}});
	  </script>
	  {{template "script" . }}
	</body>
</html>
{{end}}
{{ define "style" }}{{end}}
{{ define "head" }}{{end}}
{{ define "script" }}{{end}}
{{ define "title" }}Web{{end}}
{{ define "heading" }}<h3>&nbsp;</h3>{{end}}
`
	IndexTemplate = `
{{define "title"}}Home{{end}}
{{define "heading"}}{{end}}
{{define "style" }}{{end}}
{{define "content"}}
<div id="services"></div>
{{end}}
{{define "script"}}
<script type="text/javascript">main()</script>
{{end}}
`
	NotFoundTemplate = `
{{define "title"}}404: Not Found{{end}}
{{define "heading"}}<h3>404: Not Found</h3>{{end}}
{{define "content"}}<p>The requested page could not be found</p>{{end}}`
)
