<template>
    <el-container>
        <el-header>
            <el-card :height="60" :body-style="{ padding: '10px 10px 10px 20px'}">
                <el-select v-model="serviceName" multiple :placeholder="$t('base.service')">
                    <el-option
                            v-for="item in options"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value">
                    </el-option>
                </el-select>
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
                                    prop="name"

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

    import {Component, Vue} from "vue-property-decorator";
    import {State, Action} from 'vuex-class';

    // @ts-ignore
    import ECharts from 'vue-echarts'
    import 'echarts/lib/chart/line'
    import 'echarts/lib/component/polar'
    import 'echarts/theme/macarons'


    @Component({
        components: {
            'v-chart': ECharts
        }
    })
    export default class RegistryPage extends Vue {

        private serviceName: string = ''
        private options2 = [{
            value: '选项1',
            label: '黄金糕'
        }, {
            value: '选项2',
            label: '双皮奶',
            disabled: true
        }, {
            value: '选项3',
            label: '蚵仔煎'
        }, {
            value: '选项4',
            label: '龙须面'
        }, {
            value: '选项5',
            label: '北京烤鸭'
        }]

        private value2 = ''


        private infoItems = [
            {
                name: "Started",
                value: "Tue, 19 Mar 2019 15:19:02 GMT",
            },
            {
                name: "Uptime",
                value: "1740.501s",
            },
            {
                name: "Memory",
                value: "1.96mb",
            },

            {
                name: "Threads",
                value: "14",
            },

            {
                name: "GC",
                value: "2.043ms",
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
            autoresize: true,
            title: {},
            tooltip: {
                trigger: 'axis'
            },
            legend: {
                data: ['邮件营销', '联盟广告', '视频广告', '直接访问', '搜索引擎'],
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
                data: ['00:15 AM', '00:15 AM', '00:15 AM', '00:15 AM', '00:15 AM', '00:15 AM', '00:15 AM']
            },
            yAxis: {
                type: 'value'
            },
            series: [
                {
                    name: '邮件营销',
                    type: 'line',
                    stack: '总量',
                    data: [120, 132, 101, 134, 90, 230, 210]
                },
                {
                    name: '联盟广告',
                    type: 'line',
                    stack: '总量',
                    data: [220, 182, 191, 234, 290, 330, 310]
                },
                {
                    name: '视频广告',
                    type: 'line',
                    stack: '总量',
                    data: [150, 232, 201, 154, 190, 330, 410]
                },
                {
                    name: '直接访问',
                    type: 'line',
                    stack: '总量',
                    data: [320, 332, 301, 334, 390, 330, 320]
                },
                {
                    name: '搜索引擎',
                    type: 'line',
                    stack: '总量',
                    data: [820, 932, 901, 934, 1290, 1330, 1320]
                }
            ]
        }

        mounted() {
        }

    }

</script>
