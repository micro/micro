package stats

var (
	layoutTemplate = `
{{define "layout"}}
<html>
	<head>
		<title>Micro Stats</title>
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">
		<style>
		{{ template "style" . }}
		</style>
	</head>
	<body>
	  <nav class="navbar navbar-inverse">
	    <div class="container">
	      <div class="navbar-header">
		<a class="navbar-brand" href="/">Micro</a>
	      </div>
	    </div>
	  </nav>
          <div class="container">
            <div class="row">
              <span class="pull-right update h6"></span>
	      <div class="col-sm-4">
                <h4>&nbsp;</h4>
                <table class="table table-bordered">
                  <caption>Info</caption>
                  <tbody>
                    <tr>
                      <th>Started</th>
                      <td class="started"></td>
                    </tr>
                    <tr>
                      <th>Uptime</th>
                      <td class="uptime"></td>
                    </tr>
                    <tr>
                      <th>Memory</th>
                      <td class="memory"></td>
                    </tr>
                    <tr>
                      <th>Threads</th>
                      <td class="threads"></td>
                    </tr>
                    <tr>
                      <th>GC</th>
                      <td class="gc"></td>
                    </tr>
                  </tbody>
                </table>

                <table class="table table-bordered">
                  <caption>Requests</caption>
                  <tbody>
                    <tr>
                      <th>Total</th>
                      <td class="total"></td>
                    </tr>
                    <tr>
                      <th>20x</th>
                      <td class="20x"></td>
                    </tr>
                    <tr>
                      <th>40x</th>
                      <td class="40x"></td>
                    </tr>
                    <tr>
                      <th>50x</th>
                      <td class="50x"></td>
                    </tr>
                  </tbody>
                </table>
	      </div>
	      <div class="col-sm-8">
                {{ template "content" . }}
              </div>
            </div>
          </div>
	  <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js"></script>
	  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js" integrity="sha384-0mSbJDEHialfmuBBQP6A4Qrprq5OVfW37PRR3j5ELqxss1yVqOtnepnHVP9aJ7xS" crossorigin="anonymous"></script>
	  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/canvasjs/1.7.0/canvasjs.min.js"></script>
	  {{template "script" . }}
	</body>
</html>
{{end}}
{{ define "style" }}{{end}}
{{ define "script" }}{{end}}
{{ define "title" }}{{end}}
`

	statsTemplate = `
{{define "title"}}Stats{{end}}
{{define "content"}}
  <div id="chart" style="height: 300px; width: 100%;">
{{end}}
{{define "script"}}
<script>
  function loadChart(counters) {
	// dataPoints
	var dataPoints1 = [];
	var dataPoints2 = [];
	var dataPoints3 = [];

	var chart = new CanvasJS.Chart("chart",{
		zoomEnabled: true,
		title: {
			text: "Request Load"		
		},
		toolTip: {
			shared: true
			
		},
		legend: {
			verticalAlign: "top",
			horizontalAlign: "center",
			fontSize: 14,
			fontWeight: "bold",
			fontFamily: "calibri",
			fontColor: "dimGrey"
		},
		axisX: {
			title: "updates every 5 secs"
		},
		axisY:{
			includeZero: false
		}, 
		data: [{ 
			// dataSeries1
			type: "line",
			xValueType: "dateTime",
			showInLegend: true,
			name: "20x",
			dataPoints: dataPoints1
		},
		{				
			// dataSeries2
			type: "line",
			xValueType: "dateTime",
			showInLegend: true,
			name: "40x" ,
			dataPoints: dataPoints2
		},
		{				
			// dataSeries3
			type: "line",
			xValueType: "dateTime",
			showInLegend: true,
			name: "50x" ,
			dataPoints: dataPoints3
		}],
                legend:{
                cursor:"pointer",
                itemclick : function(e) {
                  if (typeof(e.dataSeries.visible) === "undefined" || e.dataSeries.visible) {
                    e.dataSeries.visible = false;
                  }
                  else {
                    e.dataSeries.visible = true;
                  }
                  chart.render();
                }
              }
	});

	var two = 0;
	var four = 0;
	var five = 0;

	for (i = 0; i < counters.length; i++) {
		var time = new Date((counters[i].timestamp + 5) * 1000);
                var counter = counters[i];

                if (counter["status_codes"]["20x"] != undefined) {
		  two = counter["status_codes"]["20x"];
                } else {
                  two = 0;
                }

                if (counter["status_codes"]["50x"] != undefined) {
		  five = counter["status_codes"]["50x"];
                } else {
                  five = 0;
                }

                if (counter["status_codes"]["40x"] != undefined) {
		  four = counter["status_codes"]["40x"];
                } else {
                  four = 0;
                }

		// pushing the new values
		dataPoints1.push({
			x: time.getTime(),
			y: two
		});
		dataPoints2.push({
			x: time.getTime(),
			y: four
		});
		dataPoints3.push({
			x: time.getTime(),
			y: five
		});
	}

	// updating legend text with  updated with y Value 
	chart.options.data[0].legendText = " 20x  " + two;
	chart.options.data[1].legendText = " 40x  " + four; 
	chart.options.data[2].legendText = " 50x  " + five;
	chart.render();
  };


  function loadStats() {
    var req = new XMLHttpRequest();
    req.onreadystatechange = function() {
	if (req.readyState == 4 && req.status == 200) {
	    console.log(req.responseText);

            var data = JSON.parse(req.responseText);
            var started = new Date(data["started"]*1000);
            var uptime = new Date() - started;

            // uptime
            uptime = uptime / 1000;
            if (uptime > 3600) {
              var time = uptime;
	      var hours   = Math.floor(time / 3600);
	      var minutes = Math.floor((time - (hours * 3600)) / 60);
	      var seconds = Math.floor(time - (hours * 3600) - (minutes * 60));

	      if (hours   < 10) {hours   = "0"+hours;}
	      if (minutes < 10) {minutes = "0"+minutes;}
	      if (seconds < 10) {seconds = "0"+seconds;}
	      uptime = hours+':'+minutes+':'+seconds;
            } else {
              uptime = uptime + "s";
            }

            // info
            $('.update').text("Last updated " + (new Date()).toUTCString());
            $('.started').text(started.toUTCString());
            $('.uptime').text(uptime);
            $('.memory').text(data["memory"]);
            $('.threads').text(data["threads"]);
            $('.gc').text(data["gc_pause"]);

            // requests
            var total = 0;
            var tx = 0;
            var fx = 0;
            var fox = 0;

            for (i = 0;  i < data["counters"].length; i++) {
              var counter = data["counters"][i];
              total += counter["total_reqs"];
              if (counter["status_codes"]["20x"] != undefined) {
                tx += counter["status_codes"]["20x"];
              };
              if (counter["status_codes"]["40x"] != undefined) {
                fox += counter["status_codes"]["40x"];
              };
              if (counter["status_codes"]["50x"] != undefined) {
                fx += counter["status_codes"]["50x"];
              };
            };

            $('.total').text(total);
            $('.20x').text(tx);
            $('.40x').text(fox);
            $('.50x').text(fx);

            loadChart(data["counters"]);
	}
    }

    var request = {};
    req.open("GET", window.location.href, true);
    req.setRequestHeader("Content-type","application/json");				
    req.send(JSON.stringify(request));

    setTimeout(function() {
      loadStats();
    }, 5000);
  };

  loadStats();
</script>
{{end}}
`
)
