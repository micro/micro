<template>
    <el-container v-loading="pageLoading">
        <el-header>
            <el-card :height="60" :body-style="{ padding: '10px 10px 10px 20px'}">
                <el-row>
                    <el-col :span="4">
                        <el-select v-model="serviceName" :placeholder="$t('base.service')" @change="changeService"
                                   clearable
                                   style="width:90%">
                            <el-option
                                    v-for="item in services"
                                    :key="item.name"
                                    :label="item.name"
                                    :value="item.name">
                            </el-option>
                        </el-select>
                    </el-col>
                    <el-col :span="3">
                        <el-select v-model="serviceNode" :placeholder="$t('base.address')" @change="changeNode"
                                   clearable>
                            <el-option
                                    v-for="(item, index) in nodes"
                                    :key="index"
                                    :label="item.metadata.server_address"
                                    :value="item.metadata.server_address">
                            </el-option>
                        </el-select>
                    </el-col>
                    <el-col :span="3" style="float: right;">
                        <el-button style="float: right;" :disabled="!(this.serviceName && this.serviceNode)"
                                   @click="changeNode">{{$t("base.refresh")}}
                        </el-button>
                    </el-col>
                </el-row>
            </el-card>
        </el-header>

        <el-container>
            <el-aside width="400px">
                <el-card>
                    <div>
                        <span>Info</span>
                        <el-table
                                :data="infoItems"
                                border
                                :show-header="false"
                                style="width: 100%">
                            <el-table-column
                                    width="100">
                                <template slot-scope="scope">
                                    <span class="rowName">{{$t('stats.' + scope.row.name)}}</span>
                                </template>
                            </el-table-column>
                            <el-table-column>
                                <template slot-scope="scope">
                                    <span>{{ infoData[scope.row.key] && scope.row.formatter(infoData[scope.row.key])}}</span>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                </el-card>
                <el-card style="margin-top: 15px;">
                    <div>
                        <span>Requests</span>
                        <el-table
                                :data="requestsItems"
                                border
                                :show-header="false"
                                style="width: 100%">
                            <el-table-column
                                    width="100">
                                <template slot-scope="scope">
                                    <span class="rowName">{{$t('stats.' + scope.row.name)}}</span>
                                </template>
                            </el-table-column>
                            <el-table-column
                                    prop="value">
                                <template slot-scope="scope">
                                    <span class="rowName">{{requestTableData[scope.row.name] && requestTableData[scope.row.name]}}</span>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                </el-card>
            </el-aside>
            <el-main style="padding-top: 0px">
                <el-card>
                    <div>
                        <span style="float: right"> {{lastUpdateTime && ($t('stats.lastUpdated') + lastUpdateTime.toLocaleTimeString())}}</span>
                        <div style="height: 582px">
                            <v-chart :options="linearOptions" :autoresize="true"/>
                        </div>
                    </div>
                </el-card>
            </el-main>
        </el-container>
    </el-container>
</template>

<style scoped>
    .echarts {
        width: 100%;
        height: 100%;
    }

    .rowName {
        font-weight: 400;
        color: #1f2f3d;
    }

    .el-header {
        padding: 0 20px 0 0;
        height: 70px !important;
    }

    .el-card__body {
        padding: 10px 10px 10px 20px !important;
    }


</style>

<script lang="ts">

    import MVue from "@/basic/MVue";
    import {Component, Watch} from "vue-property-decorator";
    import {State, Action} from 'vuex-class';

    // @ts-ignore
    import ECharts from 'vue-echarts'
    import 'echarts/lib/chart/line'
    import 'echarts/lib/component/polar'
    import 'echarts/theme/macarons'
    import {Service, Node, Error} from "@/store/basic/types";
    import {Stats} from '@/store/modules/stats/types';

    const namespace: string = 'apiStats';


    @Component({
        components: {
            'v-chart': ECharts
        }
    })
    export default class Statistics extends MVue {

        private serviceName: string = ''
        private serviceNode: string = ''

        private currentInterval: any;
        private lastUpdateTime: Date = null;

        private nodes: Service[] = [];

        private infoData = {
            'started': 0,
            'memory': 0,
            'threads': 0,
            'gc_pause': 0,
        }

        private infoItems = [
            {
                name: "started",
                key: "started",
                formatter: (date: number) => {
                    return new Date(date * 1000).toLocaleString()
                },
                value: "",
            },
            {
                name: "uptime",
                key: "started",
                value: "",
                formatter: (date: number) => {
                    // @ts-ignore
                    return this.$xools.secondsToHHMMSS(((new Date() - date * 1000) / 1000).toFixed(0))
                },
            },
            {
                name: "memory",
                key: "memory",
                value: "",
                formatter: (memory: number) => {
                    return memory
                },
            },

            {
                name: "threads",
                key: "threads",
                value: "",
                formatter: (threads: number) => {
                    return threads
                },
            },

            {
                name: "gc",
                key: "gc_pause",
                value: "",
                formatter: (gc: number) => {
                    return gc
                },
            },
        ];
        private requestsItems = [
            {
                name: "total",
                value: "",
            },
            {
                name: "20x",
                value: "",
            },
            {
                name: "30x",
                value: "",
            },
            {
                name: "40x",
                value: "",
            },
            {
                name: "50x",
                value: "",
            },
        ]

        private requestTableData: any = {
            "total": 0,
            "20x": 0,
            "30x": 0,
            "40x": 0,
            "50x": 0,
        }

        private linearOptions = {
            title: {},
            tooltip: {
                trigger: 'axis'
            },
            color: ['#1E9FAC', '#ED7C30', '#C74344', '#7F6083'],
            legend: {
                data: ['20x', '30x', '40x', '50x'],
                x: 0,
            },
            grid: {
                left: '3%',
                right: '4%',
                bottom: '3%',
                containLabel: true
            },
            toolbox: {
                feature: {}
            },
            xAxis: {
                type: 'category',
                boundaryGap: false,
                data: [],
            },
            yAxis: {
                type: 'value'
            },
            series: [
                {
                    name: '20x',
                    type: 'line',
                    data: []
                },
                {
                    name: '30x',
                    type: 'line',
                    data: []
                },
                {
                    name: '40x',
                    type: 'line',
                    data: []
                },
                {
                    name: '50x',
                    type: 'line',
                    data: []
                }
            ]
        }

        @State(state => state.apiStats.loaded)
        loaded ?: boolean;

        @State(state => state.apiStats.pageLoading)
        pageLoading?: boolean;

        @State(state => state.apiStats.services)
        services?: Service[];

        @State(state => state.apiStats.currentNodeStats)
        currentNodeStats?: Stats;

        @State(state => state.apiStats.xError)
        xError?: string;

        @Action('getAPIGatewayServices', {namespace})
        getAPIGatewayServices: any;

        @Action('getStats', {namespace})
        getStats: any;


        created() {
            if (!this.loaded) {
                this.getAPIGatewayServices()
            }
        }

        mounted() {

        }

        changeService(name: string) {
            this.clean()

            if (name) {
                for (let i = 0; i < this.services.length; i++) {
                    if (this.services[i].name == name) {
                        this.nodes = this.services[i].nodes
                        break;
                    }
                }
            }
        }

        changeNode() {
            this.clean()

            if (!this.serviceNode) {
                return
            }

            let go = () => {
                this.getStats({name: this.serviceName, address: this.serviceNode})
            }

            go()

            this.currentInterval = setInterval(go, 5000)
        }

        totalRequestTableData() {
            for (let key in this.requestTableData) {
                this.requestTableData['total'] += this.requestTableData[key]
            }
        }

        clean() {
            clearInterval(this.currentInterval)
            this.cleanLinerData();
            this.cleanRequestsTableData()
            this.cleanInfoTableData()
        }

        cleanLinerData() {
            this.linearOptions.xAxis.data = []
            this.linearOptions.series[0].data = []
            this.linearOptions.series[1].data = []
            this.linearOptions.series[2].data = []
            this.linearOptions.series[3].data = []
        }

        cleanInfoTableData() {
            this.infoData.started = null
            this.infoData.memory = null
            this.infoData.threads = null
            this.infoData.gc_pause = null
        }

        cleanRequestsTableData() {
            for (let key in this.requestTableData) {
                this.requestTableData[key] = 0
            }
        }

        @Watch("xError")
        catchError(xError: Error) {

            if (xError) {
                clearInterval(this.currentInterval)
                // @ts-ignore
                this.$message.error('Oops, ' + xError.detail || xError);
            }
        }

        @Watch("currentNodeStats")
        asyncData(data: Stats) {

            this.infoData.started = data.started
            this.infoData.memory = data.memory
            this.infoData.threads = data.threads
            this.infoData.gc_pause = data.gc_pause

            if (data.counters) {

                this.cleanLinerData()

                for (let key in this.requestTableData) {
                    this.requestTableData[key] = 0
                }

                for (let i = 0; i < data.counters.length; i++) {

                    let item = data.counters[i];

                    // the X Asix
                    this.linearOptions.xAxis.data.push(new Date(item.timestamp * 1000).toLocaleTimeString())

                    for (let idx in this.linearOptions.legend.data) {

                        let _ = this.linearOptions.legend.data[idx]
                        // series
                        // @ts-ignore
                        let _v = item.status_codes[_] ? item.status_codes[_] : 0
                        this.linearOptions.series[idx].data.push(_v)
                        this.requestTableData[_] += _v
                    }
                }

                this.totalRequestTableData()
            }

            this.lastUpdateTime = new Date()
        }

    }

</script>
