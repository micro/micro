<template>

    <el-carousel :height="'600px'"
                 :autoplay="false"
                 :arrow="serviceDetail.length < 2 ?'never':'always'"
    >
        <el-carousel-item
                v-for="(item,i) in serviceDetail"
                :key="i"
        >
            <el-card>
                <el-form label-width="80px">
                    <el-form-item :label="$t('base.version')">
                        <span>{{item.version}}</span>
                    </el-form-item>
                    <el-form-item :label="$t('base.metadata')">
                        <span>{{item.metadata}}</span>
                    </el-form-item>
                    <el-form-item :label="$t('base.endpoints')">
                        <span>{{formatEndpoint(item.endpoints)}}</span>
                        <el-popover
                                placement="right"
                                width="400"
                                trigger="click">
                            <el-input
                                    type="textarea"
                                    :autosize="{ minRows: 2, maxRows: 16 }"
                                    :value="JSON.stringify(item.endpoints, null, 2)">
                            </el-input>
                            <el-button size="small" type="text" slot="reference">more</el-button>
                        </el-popover>
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

        </el-carousel-item>

    </el-carousel>

</template>

<style scoped>
    .el-form-item {
        margin-bottom: 0;
    }
</style>

<script lang="ts">
    import {Component, Vue, Prop} from "vue-property-decorator";
    import {Service, Node} from "@/store/basic/types";

    // @ts-ignore
    import VuePerfectScrollbar from "vue-perfect-scrollbar";

    @Component({
        components: {VuePerfectScrollbar}
    })
    export default class ServiceDetail extends Vue {

        @Prop()
        private serviceDetail?: Service[];

        private dialog: boolean = false;

        private search: string = "";

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
    }
</script>
