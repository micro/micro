<template>
    <div>
        <el-button v-show="!showFlag" icon="el-icon-back"
                   style="float:right; margin-top: -51px; border: 0; margin-right: 30px;" @click="showFlag=true"></el-button>
        <el-carousel :height="'630px'"
                     :autoplay="false"
                     :arrow="serviceDetail.length < 2 ?'never':'always'"
        >

            <el-carousel-item
                    v-for="(item, i) in serviceDetail"
                    :key="i"
            >
                <el-card v-show="showFlag">
                    <el-form label-width="80px">
                        <el-form-item :label="$t('base.version')">
                            <span>{{item.version}}</span>
                        </el-form-item>
                        <el-form-item :label="$t('base.metadata')">
                            <span>{{item.metadata}}</span>
                        </el-form-item>
                        <el-form-item :label="$t('base.endpoints')">
                            <el-button v-for="(ep, idx) in item.endpoints" :key="idx"
                                       @click="showCall(ep, item)"
                                       type="text">
                                {{ep.name}}
                            </el-button>
                        </el-form-item>
                        <el-form-item :label="$t('base.nodes')">
                            <el-col :span="6" style="float: right;">
                                <el-input v-model="search" :placeholder="$t('base.search')"/>
                            </el-col>
                        </el-form-item>
                    </el-form>
                    <el-table
                            border
                            style="margin-top: 10px;"
                            max-height="400"
                            :data="item.nodes.filter(searchFilter)">
                        <el-table-column
                                type="index"
                                width="50">
                        </el-table-column>
                        <el-table-column
                                :label="$t('base.serviceId')"
                                align="center"
                                prop="id">
                        </el-table-column>
                        <el-table-column
                                :label="$t('base.address')"
                                align="center"
                                prop="address">
                        </el-table-column>
                        <el-table-column
                                :label="$t('base.port')"
                                align="center"
                                prop="port">
                        </el-table-column>
                        <el-table-column
                                :label="$t('base.metadata')"
                                align="left"
                                header-align="center"
                                prop="metadata">
                            <template slot-scope="scope">
                                <span>{{scope.row.metadata}}</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-card>

                <el-card v-show="!showFlag">
                    <call :key="callId" shadow="never" :callData="callData"></call>
                </el-card>
            </el-carousel-item>
        </el-carousel>
    </div>
</template>

<style scoped>
    .el-form-item {
        margin-bottom: 0;
    }

    .el-container {
        margin-right: 0px;
    }
</style>

<script lang="ts">
    import {Component, Vue, Prop} from "vue-property-decorator";
    import {Service, Node, Endpoint} from "@/store/basic/types";

    import Call from "../call/Call.vue"

    @Component({
        components: {"call": Call}
    })
    export default class ServiceDetail extends Vue {

        @Prop()
        private serviceDetail?: Service[];

        private showFlag = true;

        private dialog: boolean = false;

        private callId: number = 0;

        private search: string = "";

        private callData = {
            specialModel: true,
            serviceName: "",
            endpoint: "",
            endpoints: []
        };

        created() {
        }

        mounted() {

        }

        formatEndpoint(endpoints: any) {

            let endpointsStr = JSON.stringify(endpoints)

            if (endpointsStr.length > 50) {
                endpointsStr = endpointsStr.substr(0, 50) + '...'
            }
            return endpointsStr
        }

        searchFilter(n: Node) {
            return !this.search
                || n.id.toLowerCase().includes(this.search.toLowerCase())
                || n.address.includes(this.search.toLowerCase())
                || n.port.toString().includes(this.search.toLowerCase())
                || JSON.stringify(n.metadata).includes(this.search.toLowerCase())
        }

        showCall(ep: Endpoint, service: Service) {
            this.callId = new Date().getTime();
            this.showFlag = false;
            this.callData.specialModel = true;
            this.callData.serviceName = service.name;
            this.callData.endpoint = ep.name;
            this.callData.endpoints = service.endpoints;
        }
    }
</script>
