<template>
    <v-card>
        <v-card-title>
            Services
            <v-spacer></v-spacer>
            <v-text-field
                    v-model="search"
                    append-icon="search"
                    label="Search"
                    single-line
                    hide-details
            ></v-text-field>
        </v-card-title>
        <v-data-table
                :headers="headers()"
                :items="services"
                hide-actions
                :loading="loading"
                :search="search"
        >
            <template v-slot:items="props">
                <td>{{ props.item.name }}</td>
                <!-- <td class="text-xs-center">{{ props.item.version || '-' }}</td>
                  <td class="text-xs-center">{{ props.item.metadata || '-' }}</td>
                  <td class="text-xs-center" border>{{ props.item.endpoints || '-' }}</td>
                  <td class="text-xs-center">{{ props.item.nodes || '-' }}</td>-->
                <td class="justify-center layout px-0">
                    <v-btn flat icon color="teal" @click="showDetail(props.item)">
                        <v-icon>list</v-icon>
                    </v-btn>
                </td>
            </template>
        </v-data-table>

        <v-dialog width="70%" scrollable v-model="serviceDetailDialog">
            <service-detail :serviceDetail="serviceDetail"></service-detail>
        </v-dialog>

    </v-card>
</template>

<script lang="ts">
    import {Component, Vue} from "vue-property-decorator";
    import {State, Action} from 'vuex-class';

    import state from '@/store/state';
    import {Service} from "@/store/modules/registry/types";
    import ServiceDetail from "@/views/common/ServiceDetail"

    const namespace: string = 'registry';

    @Component({
        components: {ServiceDetail}
    })
    export default class RegistryPage extends Vue {

        private search: string = '';

        private serviceDetailDialog: boolean = false;

        @State((state: state) => state.registry.services)
        services?: Service[];

        @State((state: state) => state.registry.serviceDetail)
        serviceDetail?: Service[];

        @State((state: state) => state.registry.pageLoading)
        loading?: boolean;

        @Action('getServices', {namespace})
        getServices: any;

        @Action('getService', {namespace})
        getService: any;

        created() {
            this.getServices(this.search);
        }

        mounted() {

        }

        headers() {
            return [
                {
                    text: this.$t("base.serviceName"),
                    sortable: false,
                    value: 'name'
                },
                /*   {text: this.$t("base.version"), sortable: false, align: 'center', value: 'version'},
               {text: this.$t("base.metadata"), sortable: false, value: 'metadata'},
                  {text: this.$t("base.endpoints"), sortable: false, value: 'endpoints'},
                  {text: this.$t("base.nodes"), sortable: false, value: 'nodes'},*/
                {text: this.$t("table.operation"), align: 'center', sortable: false}
            ]
        }

        showDetail(item: Service) {
            this.getService(item.name);
            this.serviceDetailDialog = true
        }
    }
</script>
