<template>
    <v-container fluid grid-list-md>
        <v-layout row wrap>
            <v-flex d-flex xs12 sm6 md6>
                <v-layout row wrap>
                    <v-flex d-flex>
                        <v-layout row wrap selectLayout>
                            <v-flex
                                    d-flex
                                    xs12
                            >
                                <v-card>
                                    <v-card-text>
                                        <v-select
                                                :model="service"
                                                :items="services"
                                                item-text="name"
                                                item-value="endpoints"
                                                return-object
                                                :label="$t('base.service')"
                                                @change="changeService"
                                        ></v-select>
                                    </v-card-text>
                                </v-card>
                            </v-flex>
                            <v-flex
                                    d-flex
                                    xs12
                            >
                                <v-card>
                                    <v-card-text>
                                        <v-select
                                                :model="endpoint"
                                                :items="currentEndpoints"
                                                item-text="name"
                                                return-object
                                                :label="$t('base.endpoint')"
                                                @change="changeEndpoint"
                                        >
                                        </v-select>
                                    </v-card-text>
                                </v-card>
                            </v-flex>
                            <v-flex
                                    d-flex
                                    xs12
                            >
                                <v-card>
                                    <v-card-text>
                                        <v-text-field
                                                :disabled="endpoint != 'other'"
                                                v-model="otherEndpoint"
                                                :label="$t('base.otherEndpoint')"
                                                :placeholder="$t('rpc.inputOtherEndpoint')"
                                        ></v-text-field>
                                    </v-card-text>
                                </v-card>
                            </v-flex>
                        </v-layout>
                    </v-flex>
                    <v-flex d-flex>
                        <v-card>
                            <v-card-text>
                                <v-textarea
                                        v-model="requestJSON"
                                        height="300"
                                        box
                                        placeholder="{}"
                                        :label="$t('rpc.request')"
                                        value=""
                                        :hint="$t('rpc.inputJSONFormatString')"
                                ></v-textarea>
                            </v-card-text>
                            <v-card-actions>
                                <v-btn
                                        :disabled="!requestJSON"
                                        flat
                                        @click="formatJSON"
                                >
                                    {{$t('rpc.formatJSON')}}
                                </v-btn>
                                <v-spacer></v-spacer>
                                <v-btn
                                        @click="postRequest"
                                >
                                    {{$t('rpc.postRequest')}}
                                </v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-flex>
                </v-layout>
            </v-flex>
            <v-flex d-flex xs12 sm6 md6>
                <v-card color="lighten-1" dark>
                    <v-card-text placeholder="{}"></v-card-text>
                </v-card>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script lang="ts">
    import {Component, Vue} from "vue-property-decorator";
    import {State, Action} from 'vuex-class';

    import state from '@/store/state';
    import {Endpoint, Service} from "@/store/basic/types";

    import VueJsonPretty from 'vue-json-pretty'

    const namespace: string = 'call';

    @Component({components: {VueJsonPretty}})
    export default class Call extends Vue {

        private lorem = `Lorem ipsum dolor sit amet, mel at clita quando. Te sit oratio vituperatoribus, nam ad ipsum
        posidonium mediocritatem, explicari dissentiunt cu mea. Repudiare disputationi vim in, mollis iriure nec cu, alienum argumentum ius ad. Pri eu justo aeque torquatos.`;


        private items = ['Foo', 'Bar', 'Fizz', 'Buzz'];

        private currentEndpoints = [];

        private service: string = "";
        private endpoint: string = "";
        private otherEndpoint: string = "";

        private requestJSON: string = "";

        @Action('getServiceDetails', {namespace})
        getServiceDetails: any;

        @State((state: state) => state.call.services)
        services?: Service[];

        @State((state: state) => state.registry.pageLoading)
        loading?: boolean;

        created() {

        }

        mounted() {
            this.getServiceDetails()
        }

        postRequest() {

            let postData = {
                endpoint: "VideoService.AddVideo",
                request: {},
                service: "video-service"
            }
        }

        changeService(service: Service) {
            if (service.endpoints) {
                this.currentEndpoints = service.endpoints
            } else {
                this.currentEndpoints = []
            }
            this.currentEndpoints.push({name: 'other', value: -1})
        }

        changeEndpoint(endpoint: Endpoint) {
            this.endpoint = endpoint.name
        }

        formatJSON() {
            this.requestJSON = JSON.stringify(JSON.parse(this.requestJSON), null, 2);
        }
    }
</script>


<style>
    .selectLayout .v-card__text {
        padding-bottom: 0;
        padding-top: 0;
    }
</style>