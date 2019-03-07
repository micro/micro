<template>

    <v-carousel :hide-controls="serviceDetail.length < 2"
                hide-delimiters
                :interval="10000000000"
    >
        <v-carousel-item
                v-for="(item,i) in serviceDetail"
                :key="i"
        >
            <v-card min-height="100%">
                <v-card-title
                        class="headline blue lighten-1"
                >
                    {{item.name}}
                </v-card-title>

                <v-card-text>
                    <v-list>
                        <v-list-tile
                        >
                            <v-list-tile-action>
                                <span color="blue">{{$t("base.version")}}:</span>
                            </v-list-tile-action>

                            <v-list-tile-content>
                                <v-list-tile-title v-text="item.version"></v-list-tile-title>
                            </v-list-tile-content>
                        </v-list-tile>
                        <v-list-tile
                        >
                            <v-list-tile-action>
                                <span color="blue">{{$t("base.metadata")}}:</span>
                            </v-list-tile-action>

                            <v-list-tile-content>
                                <v-list-tile-title v-text="item.metadata"></v-list-tile-title>
                            </v-list-tile-content>
                        </v-list-tile>
                        <v-list-tile
                        >
                            <v-list-tile-action>
                                <span color="blue">{{$t("base.endpoints")}}:</span>
                            </v-list-tile-action>

                            <v-list-tile-content>
                                <v-list-tile-title v-text="item.endpoints"></v-list-tile-title>
                            </v-list-tile-content>
                        </v-list-tile>
                    </v-list>

                    <v-divider></v-divider>

                    <v-list
                            subheader
                            dense expand
                            two-line
                    >
                        <v-subheader>{{$t("base.nodes")}}({{item.nodes && item.nodes.length || 0}})

                            <v-spacer></v-spacer>
                            <v-text-field
                                    v-model="search"
                                    append-icon="search"
                                    label="Search"
                                    single-line
                                    hide-details
                            ></v-text-field>
                        </v-subheader>

                        <v-layout column style="height: 50vh">
                            <v-flex md6 style="overflow: scroll">
                                <v-data-table
                                        :headers="headers()"
                                        :items="item.nodes"
                                        hide-actions
                                        :search="search"
                                >
                                    <template v-slot:items="props" style="overflow-y: scroll">
                                        <td>{{ props.item.id }}</td>
                                        <td class="text-xs-center">{{ props.item.address || '-' }}</td>
                                        <td class="text-xs-center">{{ props.item.port || '-' }}</td>
                                        <td class="text-xs-center" border>{{ props.item.metadata || '-' }}</td>
                                    </template>
                                </v-data-table>
                            </v-flex>
                        </v-layout>
                    </v-list>

                </v-card-text>
            </v-card>
        </v-carousel-item>
    </v-carousel>
</template>


<script lang="ts">
    import {Component, Vue, Prop} from "vue-property-decorator";
    import {Service} from "@/store/modules/registry/types";

    // @ts-ignore
    import VuePerfectScrollbar from "vue-perfect-scrollbar";

    @Component({
        components: {VuePerfectScrollbar}
    })
    export default class ServiceDetail extends Vue {

        @Prop()
        private serviceDetail: Service[];

        private dialog: boolean = false;

        private search: string = "";


        private serviceDetail: Service;


        mounted() {

        }

        headers() {
            return [
                {
                    text: this.$t("base.serviceId"),
                    sortable: false,
                    value: 'id'
                },
                {text: this.$t("base.address"), sortable: false, align: 'center', value: 'address'},
                {text: this.$t("base.port"), sortable: false, align: 'center', value: 'port'},
                {text: this.$t("base.metadata"), sortable: false, align: 'center', value: 'metadata'}
            ]
        }
    }
</script>