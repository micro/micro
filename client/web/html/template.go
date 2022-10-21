package html

var (
	LoginTemplate = `
{{define "title"}}Login{{end}}
{{define "heading"}}{{end}}
{{define "style" }}{{end}}
{{define "content"}}
	<div class="error">{{ .error }}</div>
	<div id="login" class='inner'></div>
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
	  <nav class="navbar">
	    <div class="container">
              <div class="navbar-header">
                <button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#navBar">
                  <span class="icon-bar"></span>
                  <span class="icon-bar"></span>
                  <span class="icon-bar"></span> 
                </button>
                <a class="navbar-brand logo" href="/">Micro</a>
              </div>
              <div class="collapse navbar-collapse" id="navBar">
	        <ul class="nav navbar-nav navbar-right" id="dev">
		  {{if .Token}}
	          <li><a href="/logout">Logout</a></li>
		  {{end}}
	        </ul>
              </div>
	    </div>
	  </nav>
          <div class="container">
            <div class="row">
	      <div class="col-sm-12">
                <div id="heading">{{ template "heading" . }}</div>
                <div id="content">{{ template "content" . }}</div>
              </div>
            </div>
          </div>
	  <script src="/assets/mu.js"></script>
	  {{template "script" . }}
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
		setURL("{{.ApiURL}}");
	  </script>
	</body>
</html>
{{end}}
{{ define "style" }}
.service { border-radius: 100px; }
{{end}}
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
