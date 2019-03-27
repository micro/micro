<template>
    <div id="serviceStatsDiv" style="padding-right: 20px">
        <el-row v-for="(s, i) in services" :key="i" style="padding-top: 10px;">
            <el-card class="box-card">
                <div slot="header" class="clearfix">
                    <span>{{s.name}}</span>

                    <el-button size="small" type="text" style="float: right; padding: 3px 0;"
                               @click="refreshService(s.name)">
                        {{$t("base.refresh")}}
                    </el-button>
                </div>

                <el-col :span="6" v-for="(n, j) in s.nodes" :key="j" style="padding: 0 5px 10px 0;">
                    <el-card :id="mergeAddressAndPort(n.address, n.port)+'Card'">
                        <div>
                            <span>{{mergeAddressAndPort(n.address, n.port)}}</span>

                            <el-button size="small" type="text" style="float: right; padding: 3px 0;"
                                       @click="getStats({
                                     name: s.name,
                                     address: mergeAddressAndPort(n.address, n.port)
                                 })">
                                {{$t("base.refresh")}}
                            </el-button>
                            <el-table
                                    :data="infoItems"
                                    border
                                    :empty-text="$t('base.noDataText')"
                                    :show-header="false"
                                    style="width: 100%">
                                <el-table-column width="100">
                                    <template slot-scope="scope">
                                        <span class="rowName">{{scope.row.name}}</span>
                                    </template>
                                </el-table-column>
                                <el-table-column>
                                    <template slot-scope="scope">
                                        <span>{{cardLoadingChanged && parseNodeStats(n, scope.row.key, scope.row.formatter)}}</span>
                                    </template>
                                </el-table-column>
                            </el-table>
                        </div>
                    </el-card>
                </el-col>
            </el-card>
        </el-row>
    </div>
</template>

<style scoped>
</style>

<script lang="ts">

    import {Component, Vue, Watch} from "vue-property-decorator";
    import {State, Action} from 'vuex-class';

    import 'echarts/lib/chart/line'
    import 'echarts/theme/macarons'
    import {Stats} from '@/store/modules/stats/types';
    import {Loading} from 'element-ui';

    import {Service, Node} from "@/store/basic/types";
    import {mergeAddressAndPort} from '@/store/basic/funcs'

    const namespace: string = 'servicesStats';

    @Component({
        components: {}
    })
    export default class Statistics extends Vue {

        private currentInterval: number;

        private loadingInstance: any;

        private cardLoadingInstance = new Map<string, any>();

        private mergeAddressAndPort = mergeAddressAndPort;

        @State(state => state.servicesStats.services)
        services?: Service[];

        @State(state => state.servicesStats.cardLoading)
        cardLoading?: Map<string, boolean>;

        @State(state => state.servicesStats.cardLoadingChanged)
        cardLoadingChanged?: boolean;

        @State(state => state.servicesStats.nodesStatsMap)
        nodesStatsMap?: Map<string, Stats>;


        @State(state => state.servicesStats.pageLoading)
        pageLoading?: boolean;

        @State(state => state.servicesStats.xError)
        xError?: string;

        @Action('getServices', {namespace})
        getServices: any;

        @Action('getStats', {namespace})
        getStats: any;

        private infoItems = [
            {
                name: "Started",
                key: "started",
                formatter: (date: number) => {
                    return new Date(date * 1000).toLocaleString()
                },
            },
            {
                name: "Uptime",
                key: "uptime",
                formatter: (uptime: number) => {
                    // @ts-ignore
                    return this.$xools.secondsToHHMMSS(uptime)
                },
            },
            {
                name: "Memory",
                key: "memory",
                formatter: (memory: number) => {
                    return (memory / (1024 * 1024)).toFixed(2) + 'mb'
                },
            },

            {
                name: "Threads",
                key: "threads",
                formatter: (threads: number) => {
                    return threads
                },
            },

            {
                name: "GC",
                key: "gc",
                formatter: (gc: number) => {
                    return (gc / (1000 * 1000)).toFixed(3) + 'ms'
                },
            },
        ];

        created() {
            this.getServices()
        }

        mounted() {
            this.loadingInstance = Loading.service({target: document.getElementById("serviceStatsDiv")})
        }

        refreshService(name: string) {

            for (let i = 0; i < this.services.length; i++) {

                if (name != this.services[i].name) {
                    continue
                }

                let nodes = this.services[i].nodes
                for (let j = 0; j < nodes.length; j++) {
                    this.getStats({
                        name: this.services[i].name,
                        address: this.mergeAddressAndPort(nodes[j].address, nodes[j].port)
                    })
                }
            }
        }

        refreshNode() {

        }

        loadNodesStats() {
            for (let i = 0; i < this.services.length; i++) {
                let nodes = this.services[i].nodes
                for (let j = 0; j < nodes.length; j++) {
                    this.getStats({
                        name: this.services[i].name,
                        address: this.mergeAddressAndPort(nodes[j].address, nodes[j].port)
                    })
                }
            }
        }

        parseNodeStats(n: Node, key: string, formatter: Function) {

            let address = this.mergeAddressAndPort(n.address, n.port)
            let stats = this.nodesStatsMap.get(address)

            // @ts-ignore
            if (stats && stats[key]) {
                // @ts-ignore
                return formatter(stats[key])
            }
        }

        @Watch("xError")
        catchErrorHandler(xError: any) {
            if (xError) {
                clearInterval(this.currentInterval)
                // @ts-ignore
                this.$message.error('Oops, ' + xError.error);
            }
        }

        @Watch("pageLoading")
        loadingHandler(ld: boolean) {
            ld ? this.loadingInstance.close() : null
        }

        @Watch("services")
        servicesLoadedHandler(services: any) {
            this.loadNodesStats();
        }

        /*
            @Watch("cardLoadingChanged")
            cardLoadingChangedHandler() {


            this.cardLoading.forEach((v, k) => {
                    let instance = this.cardLoadingInstance.get(k);
                    if (instance == null) {
                        instance = Loading.service({target: document.getElementById(k + 'Card')})
                        this.cardLoadingInstance.set(k, instance)
                    }

                    if (!v) {
                        instance.close();
                    }
                })

        }*/

        /*
                @Watch("cardLoading")
                cardLoadingHandler(cardLoading: Map<string, boolean>) {
                    cardLoading.forEach((v, k) => {
                        let instance = this.cardLoadingInstance.get(k);
                        if (instance == null) {
                            instance = Loading.service({target: document.getElementById(k + 'Card')})
                            this.cardLoadingInstance.set(k, instance)
                        }

                        if (!v) {
                            instance.close();
                        }
                    })
                }*/


    }

</script>
