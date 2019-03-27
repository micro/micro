<template>

    <el-container v-loading="loading">
        <el-header>
            <el-card :height="60" :body-style="{ padding: '10px 10px 10px 20px'}">
                <el-row>
                    <el-col :span="4">
                        <el-input v-model="search" :placeholder="$t('base.search')"/>
                    </el-col>
                    <el-col :span="3" style="float: right;">
                        <el-button style="float: right;" @click="getWebServices">{{$t("base.refresh")}}
                        </el-button>
                    </el-col>
                </el-row>
            </el-card>
        </el-header>

        <el-container>
            <el-table
                    v-loading="loading"
                    :empty-text="$t('base.noDataText')"
                    border
                    :data="webServices.filter(searchFilter)"
            >
                <el-table-column
                        :label="$t('base.serviceName')"
                        align="center"
                        prop="name">
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


    import {Service} from "@/store/basic/types";

    const namespace: string = 'registry';

    @Component({
        components: {}
    })
    export default class RegistryPage extends Vue {

        private search: string = '';

        @State(state => state.registry.webServices)
        webServices?: Service[];

        @State(state => state.registry.pageLoading)
        loading?: boolean;

        @Action('getWebServices', {namespace})
        getWebServices: any;


        created() {
            this.getWebServices();
        }

        mounted() {

        }

        searchFilter(s: Service) {
            return !this.search
                || s.name.toLowerCase().includes(this.search.toLowerCase())
        }

        showDetail(item: Service) {
            window.open("/proxy/" + item.name)
        }
    }
</script>
