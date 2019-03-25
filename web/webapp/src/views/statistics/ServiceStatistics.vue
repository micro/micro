<template>
    <el-container>
        <el-header>
            <el-card :height="60" :body-style="{ padding: '10px 10px 10px 20px'}">
                <el-row>
                    <el-col :span="8">
                        <el-select v-model="serviceName" :placeholder="$t('base.service')" @change="changeService"
                                   style="width:90%">
                            <el-option
                                    v-for="item in services"
                                    :key="item.name"
                                    :label="item.name"
                                    :value="item.name">
                            </el-option>
                        </el-select>
                    </el-col>
                    <el-col :span="12">
                        <el-select v-model="serviceNode" :placeholder="$t('base.address')" @change="changeNode">
                            <el-option
                                    v-for="(item, index) in currentNodes"
                                    :key="index"
                                    :label="item.address + ':' + item.port"
                                    :value="item.address + ':' + item.port">
                            </el-option>
                        </el-select>
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
                                    <span class="rowName">{{scope.row.name}}</span>
                                </template>
                            </el-table-column>
                            <el-table-column>
                                <template slot-scope="scope">
                                    <span>{{ nodeStats[scope.row.key] && scope.row.formatter(nodeStats[scope.row.key])}}</span>
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
                                    <span class="rowName">{{scope.row.name}}</span>
                                </template>
                            </el-table-column>
                            <el-table-column
                                    prop="value">
                            </el-table-column>
                        </el-table>
                    </div>
                </el-card>
            </el-aside>
            <el-main style="padding-top: 0px">
                <el-card>
                    <div>
                        <span style="float: right">Last updated Tue, 19 Mar 2019 15:40:42 GMT</span>
                        <div style="height: 533px">
                            <v-chart :options="polar" :autoresize="true"/>
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

    import {Component, Vue, Watch} from "vue-property-decorator";
    import {State, Action} from 'vuex-class';

    // @ts-ignore
    import ECharts from 'vue-echarts'
    import 'echarts/lib/chart/line'
    import 'echarts/lib/component/polar'
    import 'echarts/theme/macarons'
    import {Service, Node} from "@/store/basic/types";
    import {Stats} from '@/store/modules/stats/types';

    const namespace: string = 'servicesStats';


    @Component({
        components: {
            'v-chart': ECharts
        }
    })
    export default class Statistics extends Vue {

        private serviceName: string = ''
        private serviceNode: string = ''

        private currentInterval: number;

        @State(state => state.servicesStats.services)
        services?: Service[];

        @State(state => state.servicesStats.currentNodes)
        currentNodes?: Node[];

        @State(state => state.servicesStats.nodeStats)
        nodeStats?: Stats;

        @State(state => state.servicesStats.xError)
        xError?: string;

        @Action('getServices', {namespace})
        getServices: any;

        @Action('getNodes', {namespace})
        getNodes: any;

        @Action('getStats', {namespace})
        getStats: any;


        private infoItems = [
            {
                name: "Started",
                key: "started",
                formatter: (date: number) => {
                    return new Date(date * 1000).toUTCString()
                },
            },
            {
                name: "Uptime",
                key: "uptime",
                value: "1740.501s",
                formatter: (uptime: number) => {
                    return this.$xools.secondsToHHMMSS(uptime)
                },
            },
            {
                name: "Memory",
                key: "memory",
                value: "1.96mb",
                formatter: (memory: number) => {
                    return (memory / (1024 * 1024)).toFixed(2) + 'mb'
                },
            },

            {
                name: "Threads",
                key: "threads",
                value: "14",
                formatter: (threads: number) => {
                    return threads
                },
            },

            {
                name: "GC",
                key: "gc",
                value: "2.043ms",
                formatter: (gc: number) => {
                    return (gc / (1000 * 1000)).toFixed(3) + 'ms'
                },
            },
        ];
        private requestsItems = [
            {
                name: "Total",
                value: "22",
            },
            {
                name: "20x",
                value: "22",
            },
            {
                name: "40x",
                value: "0",
            },

            {
                name: "50x",
                value: "0",
            },
        ]

        private polar = {
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
                data: ['00:15 AM', '00:15 AM', '00:15 AM', '00:15 AM', '00:15 AM', '00:15 AM', '00:15 AM'],
            },
            yAxis: {
                type: 'value'
            },
            series: [
                {
                    name: '20x',
                    type: 'line',
                    stack: '总量',
                    data: [120, 132, 101, 134, 90, 230, 210]
                },
                {
                    name: '30x',
                    type: 'line',
                    stack: '总量',
                    data: [220, 182, 191, 234, 290, 330, 310]
                },
                {
                    name: '40x',
                    type: 'line',
                    stack: '总量',
                    data: [150, 232, 201, 154, 190, 330, 410]
                },
                {
                    name: '50x',
                    type: 'line',
                    stack: '总量',
                    data: [320, 332, 301, 334, 390, 330, 320]
                }
            ]
        }


        created() {
            this.getServices()
        }

        mounted() {

        }

        changeService(name: string) {
            this.getNodes(name)
        }

        changeNode(address: string) {
            clearInterval(this.currentInterval)
            let go = () => {
                this.getStats({name: this.serviceName, address: this.serviceNode})
            }

            go()

            if (address) {
                this.currentInterval = setInterval(go, 5000)
            }
        }

        @Watch("xError")
        catchError(xError: string) {

            if (xError) {
                clearInterval(this.currentInterval)
                // @ts-ignore
                this.$message.error('Oops, ' + xError.error);
            }
        }

    }

</script>
