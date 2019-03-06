<template>
    <v-card
            class="mx-auto"
            max-width="90%"
    >
        <v-card-text class="headline font-weight-bold">
            <div id="shell"></div>
        </v-card-text>
    </v-card>
</template>

<script lang="ts">
    import {Component, Vue} from "vue-property-decorator";

    @Component({
        components: {}
    })

    export default class Cli extends Vue {

        private height: string = '500px';

        created() {
            this.loadTerminalJsCss()
        }

        mounted() {
            this.height = (window.innerHeight * 0.8) + 'px'
            this.renderTerminal()
        }

        loadTerminalJsCss() {

            let jQTerminalCss = document.createElement('link')
            jQTerminalCss.setAttribute("rel", "stylesheet")
            jQTerminalCss.setAttribute("href", "https://cdnjs.cloudflare.com/ajax/libs/jquery.terminal/2.0.2/css/jquery.terminal.min.css")
            document.head.appendChild(jQTerminalCss)

            let jQTerminalScript = document.createElement('script')
            jQTerminalScript.setAttribute("src", "https://cdnjs.cloudflare.com/ajax/libs/jquery.terminal/2.0.2/js/jquery.terminal.min.js")
            document.head.appendChild(jQTerminalScript)
        }

        renderTerminal() {

            if (!window.jQuery.terminal) {
                setTimeout(this.renderTerminal, 500)
                return
            }

            jQuery(function ($, undefined) {
                $('#shell').terminal(function (command, term) {
                    if (command == '') {
                        term.echo('');
                        return;
                    }

                    var help = "COMMANDS:\n" +
                        "    call       Call a service endpoint using rpc\n" +
                        "    health      Query the health of a service\n" +
                        "    list        List items in registry\n" +
                        "    get         Get item from registry\n";
                    try {
                        let args = command.split(" ");
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
                                    success: function (data) {
                                        let services = [];
                                        for (let i = 0; i < data.services.length; i++) {
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
                                    url: "registry?service=" + args[2],
                                    data: {},
                                    success: function (data) {
                                        if (data.services.length == 0) {
                                            return
                                        }

                                        term.echo("service\t" + args[2]);
                                        term.echo(" ");

                                        let eps = {};

                                        for (let i = 0; i < data.services.length; i++) {
                                            var service = data.services[i];
                                            term.echo("\nversion " + service.version);
                                            term.echo(" ");
                                            term.echo("Id\tAddress\tPort\tMetadata\n");
                                            for (let j = 0; j < service.nodes.length; j++) {
                                                var node = service.nodes[j];
                                                var metadata = [];
                                                $.each(node.metadata, function (key, val) {
                                                    metadata.push(key + "=" + val);
                                                });
                                                term.echo(node.id + "\t" + node.address + "\t" + node.port + "\t" + metadata.join(","));
                                            }
                                            term.echo(" ");

                                            for (let k = 0; k < service.endpoints.length; k++) {
                                                if (eps[service.endpoints[k].name] == undefined) {
                                                    eps[service.endpoints[k].name] = service.endpoints[k];
                                                }
                                            }
                                        }


                                        $.each(eps, function (key, ep) {
                                            term.echo("Endpoint: " + key);
                                            var metadata = [];
                                            $.each(ep.metadata, function (mkey, val) {
                                                metadata.push(mkey + "=" + val);
                                            });
                                            term.echo("Metadata: " + metadata.join(","));

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
                                    url: "registry?service=" + args[1],
                                    data: {},
                                    success: function (data) {

                                        term.echo("service\t" + args[1]);
                                        term.echo(" ");

                                        for (i = 0; i < data.services.length; i++) {
                                            var service = data.services[i];

                                            term.echo("\nversion " + service.version);
                                            term.echo(" ");
                                            term.echo("Id\tAddress:Port\tMetadata\n");

                                            for (j = 0; j < service.nodes.length; j++) {
                                                var node = service.nodes[j];

                                                $.ajax({
                                                    endpoint: "POST",
                                                    dataType: "json",
                                                    contentType: "application/json",
                                                    url: "rpc",
                                                    data: JSON.stringify({
                                                        "service": service.name,
                                                        "endpoint": "Debug.Health",
                                                        "request": {},
                                                        "address": node.address + ":" + node.port,
                                                    }),
                                                    success: function (data) {
                                                        term.echo(node.id + "\t" + node.address + ":" + node.port + "\t" + data.status);
                                                    },
                                                    error: function (xhr) {
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
                                    term.echo("USAGE:\n    call [service] [endpoint] [request]");
                                    return;
                                }

                                var request = "{}"

                                if (args.length > 3) {
                                    request = args.slice(3).join(" ");
                                }

                                $.ajax({
                                    endpoint: "POST",
                                    dataType: "json",
                                    contentType: "application/json",
                                    url: "rpc",
                                    data: JSON.stringify({"service": args[1], "endpoint": args[2], "request": request}),
                                    success: function (data) {
                                        term.echo(JSON.stringify(data, null, 2));
                                    },
                                });

                                break;
                            default:
                                term.echo(command + ": command not found");
                                term.echo(help);
                        }
                    } catch (e) {
                        term.error(new String(e));
                    }
                }, {
                    greetings: '',
                    name: 'micro_cli',
                    height: 400,
                    prompt: 'micro:~$ '
                });
            });
        }
    }
</script>