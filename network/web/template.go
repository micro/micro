package web

var (
	templateFile = `<!DOCTYPE html>
<html>
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <link href="https://fonts.googleapis.com/css?family=Source+Code+Pro&display=swap" rel="stylesheet">
  <style>
    html {
      font-family: 'Source Code Pro', monospace;
    }
    .nav {
      margin-bottom: 10px;
      padding: 0;
    }
    .nav li {
      display: inline-block; 
    }
    table td {
      padding-right: 5px;
    }
  </style>
</head>
<body>
  <div id="head">
    <ul class="nav">
      <li><a href="/network/graph">Graph</a></li>
      <li><a href="/network/nodes">Nodes</a></li>
      <li><a href="/network/routes">Routes</a></li>
      <li><a href="/network/services">Services</a></li>
    </ul>
  </div>
  <div id="content">
    {{.}}
  </div>
</body>
</html>
`

	graphTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <link href="https://fonts.googleapis.com/css?family=Source+Code+Pro&display=swap" rel="stylesheet">
  <script src="https://cdn.jsdelivr.net/npm/chart.js@2.8.0"></script>
  <style>
    html {
      font-family: 'Source Code Pro', monospace;
    }
    .nav {
      margin-bottom: 10px;
      padding: 0;
    }
    .nav li {
      display: inline-block; 
    }
    table td {
      padding-right: 5px;
    }
    #graph {
      margin-top: 25px;
    }
  </style>
</head>
<body>
  <div id="head">
    <ul class="nav">
      <li><a href="/network/graph">Graph</a></li>
      <li><a href="/network/nodes">Nodes</a></li>
      <li><a href="/network/routes">Routes</a></li>
      <li><a href="/network/services">Services</a></li>
    </ul>
  </div>
  <div id="content">
    <canvas id="graph"></canvas>
  </div>
  <script type="text/javascript">
	Chart.defaults.scale.gridLines.display = false;
	Chart.defaults.scale.ticks.display = false;
	var ctx = document.getElementById('graph').getContext('2d');
	var chart = new Chart(ctx, {
	    type: 'radar',
	    data: {
		labels: ['{{ Join .Nodes "', '" }}'],
		datasets: [
		{{ range $label, $values := .Data }}
		{
		    label: '{{$label}}', 
		    borderColor: 'rgb({{ Color }})',
		    borderWidth: 1,
		    fill: false,
		    lineTension: 0.3,
		    data: [{{ Join $values ", " }}]
		}{{if $label}},{{end}}{{end}}
		]
	    },
	    options: {
		legend: { display: false }
	    }
	});
  </script>
</body>
</html>
`
)
