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
                                                v-model="service"
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

                                        :hint="$t('rpc.inputJSONFormatString')"
                                ></v-textarea>
                            </v-card-text>
                            <v-card-actions>
                                <v-btn
                                        :disabled="!requestJSON"
                                        flat
                                        @click="formatRequestJSON"
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
                    <v-card-text>
                        <v-textarea
                                v-model="resultString"
                                height="500"

                                box
                                :label="$t('rpc.result')"
                        ></v-textarea>
                    </v-card-text>
                </v-card>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script lang="ts">
    import {Component, Vue, Watch} from "vue-property-decorator";
    import {State, Action} from 'vuex-class';

    import state from '@/store/state';
    import {Endpoint, Service} from "@/store/basic/types";


    const namespace: string = 'call';

    @Component({components: {}})
    export default class Call extends Vue {

        private currentEndpoints = [];

        private service: Service = new Service();
        private endpoint: string = "";
        private otherEndpoint: string = "";

        private requestJSON: string = "{}";

        private resultString: string = ""

        @Action('getServiceDetails', {namespace})
        getServiceDetails: any;

        @Action('postServiceRequest', {namespace})
        postServiceRequest: any;

        @State((state: state) => state.call.services)
        services?: Service[];

        @State((state: state) => state.call.requestResult)
        requestResult?: object;

        @State((state: state) => state.call.requestLoading)
        requestLoading?: boolean;


        @Watch("requestResult")
        resultChange(rr) {
            try {
                this.resultString = JSON.stringify(rr, null, 2)
            } catch (e) {
                this.resultString = rr
            }
        }

        created() {

        }

        mounted() {
            this.getServiceDetails()
        }

        postRequest() {

            let endpoint = this.endpoint;
            if (endpoint == 'other') {
                endpoint = this.otherEndpoint;
            }

            let postData = {
                endpoint: endpoint,
                request: JSON.parse(this.requestJSON),
                service: this.service.name
            }

            this.postServiceRequest(postData);
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

        formatRequestJSON() {
            this.requestJSON = this.formatJSON(this.requestJSON)
        }

        formatJSON(str: string) {
            return JSON.stringify(JSON.parse(str), null, 2);
        }
    }
</script>


<style>
    .selectLayout .v-card__text {
        padding-bottom: 0;
        padding-top: 0;
    }
</style>