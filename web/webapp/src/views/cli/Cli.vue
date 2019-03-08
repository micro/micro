<template>
    <v-card
            class="mx-auto"
            max-width="90%"
            :height="height"
    >
        <v-card-title class="pt-0 pb-0">
            <v-spacer></v-spacer>
            <v-btn icon @click="handleFullScreen">
                <v-icon>fullscreen</v-icon>
            </v-btn>
        </v-card-title>
        <v-card-text class="pt-0">
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

        }

        mounted() {
            this.height = (window.innerHeight * 0.8) + 'px'
            this.renderTerminal()
        }

        renderTerminal() {

            let that = this;
            // @ts-ignore
            if (!window.jQuery.terminal) {
                setTimeout(this.renderTerminal, 500)
                return
            }
            // @ts-ignore
            jQuery(function ($, undefined) {
                $('#shell').terminal(function (command: any, term: any) {
                    if (command == '') {
                        term.echo('');
                        return;
                    }

                    var help = "COMMANDS:\n" +
                        "    exit        exit fullscreen\n" +
                        "    call        Call a service endpoint using rpc\n" +
                        "    health      Query the health of a service\n" +
                        "    list        List items in registry\n" +
                        "    get         Get item from registry\n";
                    try {
                        let args = command.split(" ");
                        switch (args[0]) {
                            case "exit":
                                that.$xools.toggleFullScreen('shell')
                                break;
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
                                    success: function (data: any) {
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
                                    success: function (data: any) {
                                        if (data.services.length == 0) {
                                            return
                                        }

                                        term.echo("service\t" + args[2]);
                                        term.echo(" ");

                                        let eps: any = {};

                                        for (let i = 0; i < data.services.length; i++) {
                                            var service = data.services[i];
                                            term.echo("\nversion " + service.version);
                                            term.echo(" ");
                                            term.echo("Id\tAddress\tPort\tMetadata\n");
                                            for (let j = 0; j < service.nodes.length; j++) {
                                                let node = service.nodes[j];
                                                //@ts-ignore
                                                let metadata = [];
                                                $.each(node.metadata, function (key: any, val: any) {
                                                    metadata.push(key + "=" + val);
                                                });

                                                // @ts-ignore
                                                term.echo(node.id + "\t" + node.address + "\t" + node.port + "\t" + metadata.join(","));
                                            }
                                            term.echo(" ");

                                            for (let k = 0; k < service.endpoints.length; k++) {
                                                if (eps[service.endpoints[k].name] == undefined) {
                                                    eps[service.endpoints[k].name] = service.endpoints[k];
                                                }
                                            }
                                        }


                                        $.each(eps, function (key: any, ep: any) {
                                            term.echo("Endpoint: " + key);
                                            // @ts-ignore
                                            let metadata = [];
                                            $.each(ep.metadata, function (mkey: any, val: any) {
                                                metadata.push(mkey + "=" + val);
                                            });
                                            // @ts-ignore
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
                                    success: function (data: any) {

                                        term.echo("service\t" + args[1]);
                                        term.echo(" ");

                                        for (let i = 0; i < data.services.length; i++) {
                                            var service = data.services[i];

                                            term.echo("\nversion " + service.version);
                                            term.echo(" ");
                                            term.echo("Id\tAddress:Port\tMetadata\n");

                                            for (let j = 0; j < service.nodes.length; j++) {
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
                                                    success: function (data: any) {
                                                        term.echo(node.id + "\t" + node.address + ":" + node.port + "\t" + data.status);
                                                    },
                                                    error: function (xhr: any) {
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
                                    success: function (data: any) {
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
                    height: 500,
                    prompt: 'micro:~$ '
                });
            });
        }


        handleFullScreen() {
            this.$xools.toggleFullScreen('shell')
        }
    }
</script>