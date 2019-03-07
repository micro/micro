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
                :items="registries"
                hide-actions
                :loading="loading"
                :search="search"
        >
            <template v-slot:items="props">
                <td>{{ props.item.name }}</td>
                <td class="text-xs-center">{{ props.item.version || '-' }}</td>
                <td class="text-xs-center">{{ props.item.metadata || '-' }}</td>
                <td class="text-xs-center" border>{{ props.item.endpoints || '-' }}</td>
                <td class="text-xs-center">{{ props.item.nodes || '-' }}</td>
                <td class="justify-center layout px-0">
                    <v-btn flat icon color="teal" @click="showDetail">
                        <v-icon>list</v-icon>
                    </v-btn>
                </td>
            </template>
            <v-alert v-slot:no-results :value="true" color="error" icon="warning">
                Your search for "{{ search }}" found no results.
            </v-alert>
        </v-data-table>

        <v-dialog width="60%" v-model="serviceDetailDialog">
            <service-detail :me="currentService"></service-detail>
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

        private currentService?: Service = null;

        @State((state: state) => state.registry.registries)
        registries?: Service[];

        @State((state: state) => state.registry.pageLoading)
        loading?: boolean;

        @Action('getRegistries', {namespace})
        getRegistries: any;

        created() {
            this.getRegistries(this.search);
        }

        mounted() {

        }

        headers() {
            return [
                {
                    text: this.$t("table.registry.name"),
                    sortable: false,
                    value: 'name'
                },
                {text: this.$t("table.registry.version"), sortable: false, value: 'version'},
                {text: this.$t("table.registry.metadata"), sortable: false, value: 'metadata'},
                {text: this.$t("table.registry.endpoints"), sortable: false, value: 'endpoints'},
                {text: this.$t("table.registry.nodes"), sortable: false, value: 'nodes'},
                {text: this.$t("table.operation"), sortable: false}
            ]
        }

        showDetail(item) {
            this.getRegistries(item.name);
            this.serviceDetailDialog = true
        }
    }
</script>
