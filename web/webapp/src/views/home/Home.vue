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
                :items="webServices"
                hide-actions
                :loading="loading"
                :no-data-text="$t('base.noDataText')"
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
                        <v-icon>open_in_browser</v-icon>
                    </v-btn>
                </td>
            </template>
        </v-data-table>

    </v-card>
</template>

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
            window.open("/proxy/" + item.name)
        }
    }
</script>
