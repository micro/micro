<template>

    <el-container v-loading="loading">
        <el-header>
            <el-card :height="60" :body-style="{ padding: '10px 10px 10px 20px'}">
                <el-row>
                    <el-col :span="4">
                        <el-input v-model="search" :placeholder="$t('base.search')"/>
                    </el-col>
                    <el-col :span="3" style="float: right;">
                        <el-button style="float: right;" @click="getServices">{{$t("base.refresh")}}
                        </el-button>
                    </el-col>
                </el-row>
            </el-card>
        </el-header>

        <el-container>
            <el-table
                    border
                    :empty-text="$t('base.noDataText')"
                    :data="services.filter(searchFilter)">
                <el-table-column
                        :label="$t('base.serviceName')"
                        align="center"
                        prop="name">
                </el-table-column>
                <el-table-column
                        :label="$t('base.nodes')"
                        align="left"
                        header-align="center"
                        prop="nodes">
                    <template slot-scope="scope">
                        {{
                        parseNodes(scope.row.nodes)
                        }}
                    </template>
                </el-table-column>
                <el-table-column
                        :label="$t('table.operation')"
                        align="center">
                    <template slot-scope="scope">
                        <el-button
                                type="text"
                                size="mini"
                                @click="showDetail(scope.row)">Detail
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </el-container>

        <el-dialog width="70%" :title="detailTitle" :visible.sync="serviceDetailDialog">
            <service-detail :serviceDetail="serviceDetail"></service-detail>
        </el-dialog>
    </el-container>
</template>

<style scoped>

    .el-container .el-container {
        margin-right: 20px;
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


    import {Service, Node} from "@/store/basic/types";
    import {mergeAddressAndPort} from '@/store/basic/funcs'
    import ServiceDetail from "@/views/common/ServiceDetail.vue"

    const namespace: string = 'registry';

    @Component({
        components: {ServiceDetail}
    })
    export default class RegistryPage extends Vue {

        private search: string = '';

        private serviceDetailDialog: boolean = false;

        private detailTitle = ""

        @State(state => state.registry.services)
        services?: Service[];

        @State(state => state.registry.serviceDetail)
        serviceDetail?: Service[];

        @State(state => state.registry.pageLoading)
        loading?: boolean;

        @Action('getServices', {namespace})
        getServices: any;

        @Action('getService', {namespace})
        getService: any;

        created() {
            this.getServices();
        }

        mounted() {

        }

        searchFilter(s: Service) {
            return !this.search
                || s.name.toLowerCase().includes(this.search.toLowerCase())
                || this.parseNodes(s.nodes).includes(this.search.toLowerCase())
        }

        parseNodes(nodes: Node[]) {

            if (nodes) {
                let nodesStr: any[] = [];
                nodes.forEach(node => {
                    nodesStr.push(mergeAddressAndPort(node.address, node.port))
                })
                return nodesStr.join(", ") + ' | ' + nodesStr.length
            }
            return '';
        }


        showDetail(item: Service) {
            this.detailTitle = item.name;
            this.getService(item.name);
            this.serviceDetailDialog = true
        }
    }
</script>
