package web

var (
	layoutTemplate = `
{{define "layout"}}
<html>
	<head>
		<title>Micro Web</title>
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">
		<style>
		{{ template "style" . }}
		</style>
		{{ template "head" . }}
	</head>
	<body>
	  <nav class="navbar navbar-inverse">
	    <div class="container">
              <div class="navbar-header">
                <button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#navBar">
                  <span class="icon-bar"></span>
                  <span class="icon-bar"></span>
                  <span class="icon-bar"></span> 
                </button>
                <a class="navbar-brand" href="/">Micro</a>
              </div>
              <div class="collapse navbar-collapse" id="navBar">
	        <ul class="nav navbar-nav navbar-right">
	          <li><a href="cli">CLI</a></li>
	          <li><a href="registry">Registry</a></li>
	          <li><a href="call">Call</a></li>
	          {{if .StatsURL}}<li><a href="{{.StatsURL}}" class="navbar-link">Stats</a></li>{{end}}
	        </ul>
              </div>
	    </div>
	  </nav>
          <div class="container">
            <div class="row">
	      <div class="col-sm-12">
                {{ template "heading" . }}
                {{ template "content" . }}
              </div>
            </div>
          </div>
	  <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.4/jquery.min.js"></script>
	  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js" integrity="sha384-0mSbJDEHialfmuBBQP6A4Qrprq5OVfW37PRR3j5ELqxss1yVqOtnepnHVP9aJ7xS" crossorigin="anonymous"></script>
	  {{template "script" . }}
	</body>
</html>
{{end}}
{{ define "style" }}{{end}}
{{ define "head" }}{{end}}
{{ define "script" }}{{end}}
{{ define "title" }}{{end}}
{{ define "heading" }}<h3>&nbsp;</h3>{{end}}
`

	indexTemplate = `
{{define "heading"}}<h4><input class="form-control input-lg search" type=text placeholder="Search"/></h4>{{end}}
{{define "title"}}Web{{end}}
{{define "content"}}
	{{if .Results.HasWebServices}}
		<div>
			{{range .Results.WebServices}}
			<a href="/{{.}}" data-filter={{.}} class="btn btn-default btn-lg" style="margin: 5px 3px 5px 3px;">{{.}}</a>
			{{end}}
		</div>
	{{else}}
		<div class="alert alert-info" role="alert">
			<strong>No web services found</strong>
		</div>
	{{end}}
{{end}}
{{define "script"}}
<script type="text/javascript">
jQuery(function($, undefined) {
	var refs = $('a[data-filter]');
	$('.search').on('keyup', function() {
		var val = $.trim(this.value);
		refs.hide();
		refs.filter(function() {
			return $(this).data('filter').search(val) >= 0
		}).show();
	});
});
</script>
{{end}}
`
	callTemplate = `
{{define "title"}}Call{{end}}
{{define "style"}}
	pre {
		word-wrap: break-word;
	}
{{end}}
{{define "content"}}
<div class="row">
  <div class="panel panel-default">
    <div class="panel-body">
	<div class="col-sm-5">
		<form id="call-form" onsubmit="return call();">
			<div class="form-group">
				<label for="service">Service</label>
				<ul class="list-group">
					<select class="form-control" type=text name=service id=service> 
					<option disabled selected> -- select a service -- </option>
					{{range $key, $value := .Results}}
					<option class = "list-group-item" value="{{$key}}">{{$key}}</option>
					{{end}}
					</select>
				</ul>
			</div>
			<div class="form-group">
				<label for="method">Method</label>
				<ul class="list-group">
					<select class="form-control" type=text name=method id=method>
					<option disabled selected> -- select a method -- </option>
					</select>
				</ul>
			</div>
			<div class="form-group">
				<label for="othermethod">Other Method</label>
				<ul class="list-group">
					<input class="form-control" type=text name=othermethod id=othermethod disabled placeholder="Method"/>
				</ul>
			</div>
			<div class="form-group">
				<label for="request">Request</label>
				<textarea class="form-control" name=request id=request rows=8>{}</textarea>
			</div>
			<div class="form-group">
				<button class="btn btn-default">Execute</button>
			</div>
		</form>
	</div>
	<div class="col-sm-7">
		<p><b>Response</b></p>
		<pre id="response" style="min-height: 405px;">{}</pre>
	</div>
    </div>
  </div>
</div>
{{end}}
{{define "script"}}
	<script>
		$(document).ready(function(){
			//Function executes on change of first select option field 
			$("#service").change(function(){
				var select = $("#service option:selected").val();
				$("#othermethod").attr("disabled", true);
				$('#othermethod').val('');
				$("#method").empty();
				$("#method").append("<option disabled selected> -- select a method -- </option>");
				var s_map = {};
				{{ range $service, $methods := .Results }}
				var m_list = [];
				{{range $index, $element := $methods}}
				m_list[{{$index}}] = {{$element.Name}}
				{{end}}
				s_map[{{$service}}] = m_list
				{{ end }}
				if (select in s_map) {
				var serviceMethods = s_map[select]
				var len = serviceMethods.length;
					for(var i = 0; i < len; i++) {
						$("#method").append("<option value=\""+serviceMethods[i]+"\">"+serviceMethods[i]+"</option>");	
					}
				}
				$("#method").append("<option value=\"other\"> - Other</option>");
			});
			//Function executes on change of second select option field 
			$("#method").change(function(){
				var select = $("#method option:selected").val();
				if (select == "other") {
					$("#othermethod").attr("disabled", false);
				} else {
					$("#othermethod").attr("disabled", true);
					$('#othermethod').val('');
				}

			});
		});
	</script>
	<script>
		function call() {
			var req = new XMLHttpRequest()
			req.onreadystatechange = function() {
				if (req.readyState == 4 && req.status == 200) {
					document.getElementById("response").innerText = JSON.stringify(JSON.parse(req.responseText), null, 2);
				} else if (req.responseText.slice(0, 1) == "{") {
					document.getElementById("response").innerText = JSON.stringify(JSON.parse(req.responseText), null, 2);
				} else if (req.responseText.length > 0) {
					document.getElementById("response").innerText = req.responseText;
				} else {
					document.getElementById("response").innerText = "Request error " + req.status;
				}
				console.log(req.responseText);
			}
			var method = document.forms[0].elements["method"].value
			if (!($('#othermethod').prop('disabled'))) {
				method = document.forms[0].elements["othermethod"].value
			}
			var request = {
				"service": document.forms[0].elements["service"].value,
				"method": method,
				"request": JSON.parse(document.forms[0].elements["request"].value)
			}
			req.open("POST", "/rpc", true);
			req.setRequestHeader("Content-type","application/json");				
			req.send(JSON.stringify(request));

			return false;
		};	
	</script>
{{end}}
`
	registryTemplate = `
{{define "heading"}}<h4><input class="form-control input-lg search" type=text placeholder="Search"/></h4>{{end}}
{{define "title"}}Registry{{end}}
{{define "content"}}
	<div>
		{{range .Results}}
		<a href="registry?service={{.Name}}" data-filter={{.Name}} class="btn btn-default btn-lg" style="margin: 5px 3px 5px 3px;">{{.Name}}</a>
		{{end}}
	</div>
{{end}}
{{define "script"}}
<script type="text/javascript">
jQuery(function($, undefined) {
	var refs = $('a[data-filter]');
	$('.search').on('keyup', function() {
		var val = $.trim(this.value);
		refs.hide();
		refs.filter(function() {
			return $(this).data('filter').search(val) >= 0
		}).show();
	});
});
</script>
{{end}}
`

	serviceTemplate = `
{{define "title"}}Service{{end}}
{{define "heading"}}<h3>{{with $svc := index .Results 0}}{{$svc.Name}}{{end}}</h3>{{end}}
{{define "content"}}
	<hr>
	<h4>Nodes</h4>
	{{range .Results}}
	<h5>Version {{.Version}}</h5>
	<table class="table table-bordered table-striped">
		<thead>
			<th>Id</th>
			<th>Address</th>
			<th>Port</th>
			<th>Metadata</th>
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
	{{end}}
	{{with $svc := index .Results 0}}
	{{if $svc.Endpoints}}
	<h4>Endpoints</h4>
	<hr/>
	{{end}}
	{{range $svc.Endpoints}}
		<h4>{{.Name}}</h4>
		<table class="table table-bordered">
			<tbody>
				<tr>
					<th class="col-sm-2" scope="row">Metadata</th>
					<td>{{ range $key, $value := .Metadata }}{{$key}}={{$value}} {{end}}</td>
				</tr>
				<tr>
					<th class="col-sm-2" scope="row">Request</th>
					<td><pre>{{format .Request}}</pre></td>
				</tr>
				<tr>
					<th class="col-sm-2" scope="row">Response</th>
					<td><pre>{{format .Response}}</pre></td>
				</tr>
			</tbody>
		</table>
	{{end}}
	{{end}}
{{end}}

`

	cliTemplate = `
{{define "head"}}
<link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/jquery.terminal/2.0.2/css/jquery.terminal.min.css">
{{end}}
{{define "title"}}CLI{{end}}
{{define "content"}}
<div id="shell"></div>
{{end}}
{{define "script"}}
<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery.terminal/2.0.2/js/jquery.terminal.min.js"></script>
<script type="text/javascript">
jQuery(function($, undefined) {
    $('#shell').terminal(function(command, term) {
        if (command == '') {
            term.echo('');
	    return;
        }

	var help = "COMMANDS:\n" +
	"    call       Call a service method using rpc\n" +
	"    health      Query the health of a service\n" +
	"    list        List items in registry\n" +
	"    get         Get item from registry\n";
        try {
	    args = command.split(" ");
	    switch (args[0]) {
	    case "help":
		term.echo(help);
		break;
	    case "list":
		if (args.length == 1 || args[1] != "services") {
		    term.echo("COMMANDS:\n    services    List services in registry\n");
		    return;
		}
		$.ajax({
		  dataType: "json",
		  contentType: "application/json",
		  url: "registry",
		  data: {},
		  success: function(data) {
		    var services = [];
		    for (i = 0; i < data.services.length; i++) {
			services.push(data.services[i].name);
		    }
		    term.echo(services.join("\n"));
		  },
		});
		break;
	    case "get":
		if (args.length < 3 || args[1] != "service") {
		    term.echo("COMMANDS:\n    service    Get service from registry\n");
		    return;
		}

		$.ajax({
		  dataType: "json",
		  contentType: "application/json",
		  url: "registry?service="+args[2],
		  data: {},
		  success: function(data) {
		    if (data.services.length == 0) {
			return
		    }

		    term.echo("service\t"+args[2]);
		    term.echo(" ");

		    var eps = {};

		    for (i = 0; i < data.services.length; i++) {
			var service = data.services[i];
			term.echo("\nversion "+service.version);
			term.echo(" ");
			term.echo("Id\tAddress\t\Port\tMetadata\n");
			for (j = 0; j < service.nodes.length; j++) {
			    var node = service.nodes[j];
			    var metadata = [];
			    $.each(node.metadata, function(key, val) {
				metadata.push(key+"="+val);
			    });
			    term.echo(node.id + "\t" + node.address + "\t" + node.port + "\t" + metadata.join(","));
			}
			term.echo(" ");

			for (k = 0; k < service.endpoints.length; k++) {
			    if (eps[service.endpoints[k].name] == undefined) {
				eps[service.endpoints[k].name] = service.endpoints[k];
			    }
			}
		    }


		    $.each(eps, function(key, ep) {
			term.echo("Endpoint: "+key);
			var metadata = [];
			$.each(ep.metadata, function(mkey, val) {
			    metadata.push(mkey+"="+val);
			});
			term.echo("Metadata: "+metadata.join(","));
		
			// TODO: add request-response endpoints	
		    })
		  },
		});

		break;
	    case "health":
		if (args.length < 2) {
		    term.echo("USAGE:\n    health [service]");
		    return;
		}

		$.ajax({
		  dataType: "json",
		  contentType: "application/json",
		  url: "registry?service="+args[1],
		  data: {},
		  success: function(data) {
			
		    term.echo("service\t"+args[1]);
		    term.echo(" ");

		    for (i = 0; i < data.services.length; i++) {
			var service = data.services[i];

			term.echo("\nversion "+service.version);
			term.echo(" ");
			term.echo("Id\tAddress:Port\tMetadata\n");

			for (j = 0; j < service.nodes.length; j++) {
			    var node = service.nodes[j];

			    $.ajax({
				  method: "POST",
				  dataType: "json",
				  contentType: "application/json",
				  url: "rpc",
				  data: JSON.stringify({
					"service": service.name,
					"method": "Debug.Health",
					"request": {},
					"address": node.address + ":" + node.port,
				  }),
				  success: function(data) {
			    		term.echo(node.id + "\t" + node.address + ":" + node.port + "\t" + data.status);
				  },
				  error: function(xhr) {
			    		term.echo(node.id + "\t" + node.address + ":" + node.port + "\t" + xhr.status);
				  },
			    });

			}

			term.echo(" ");
		    }
		  },
		});


		break;
	    case "call":
		if (args.length < 3) {
		    term.echo("USAGE:\n    call [service] [method] [request]");
		    return;
		}

		var request = "{}"

		if (args.length > 3) {
			request = args.slice(3).join(" ");
		}		

		$.ajax({
		  method: "POST",
		  dataType: "json",
		  contentType: "application/json",
		  url: "rpc",
		  data: JSON.stringify({"service": args[1], "method": args[2], "request": request}),
		  success: function(data) {
		    term.echo(JSON.stringify(data, null, 2));
		  },
		});
		
		break;
	    default:
		term.echo(command +": command not found");
		term.echo(help);
	    }
        } catch(e) {
	    term.error(new String(e));
        }
    }, {
        greetings: '',
        name: 'micro_cli',
        height: 400,
        prompt: 'micro:~$ '});
});
</script>
{{end}}
`
)
