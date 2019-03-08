<template>
    <v-container fluid grid-list-xl>
        <v-layout row wrap>
            <v-flex d-flex xs12 sm6 md5>
                <v-layout row wrap>
                    <v-flex d-flex>
                        <v-layout row wrap>
                            <v-flex
                                    d-flex
                                    xs12
                                    pt-0
                                    pb-2
                            >
                                <v-card>
                                    <v-card-text class="pt-0 pb-0">
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
                                    pt-0
                                    pb-2
                            >
                                <v-card>
                                    <v-card-text class="pt-0 pb-0">
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
                                    pt-0
                                    pb-2
                            >
                                <v-card>
                                    <v-card-text class="pt-0 pb-0">
                                        <v-text-field
                                                :disabled="endpoint != 'other'"
                                                v-model="otherEndpoint"
                                                :label="$t('base.otherEndpoint')"
                                                :placeholder="$t('rpc.inputOtherEndpoint')"
                                        ></v-text-field>
                                    </v-card-text>
                                </v-card>
                            </v-flex>
                            <v-flex d-flex codeFlex pt-0 pb-2>
                                <v-card>
                                    <v-card-title>
                                        <span class="title font-weight-light">{{$t('rpc.request')}}</span>
                                    </v-card-title>
                                    <v-card-text>
                                        <div id="jsonRequestEditor" style="height: 300px" class="json-editor">

                                        </div>
                                    </v-card-text>
                                </v-card>
                            </v-flex>
                        </v-layout>
                    </v-flex>
                </v-layout>
            </v-flex>
            <v-flex md1>

                <v-btn
                        small
                        @click="formatRequestJSON"
                >
                    <span> {{$t('rpc.formatJSON')}}</span>
                </v-btn>
                <v-btn
                        small
                        @click="postRequest"
                >
                    <span>{{$t('rpc.postRequest')}}</span>
                </v-btn>

            </v-flex>
            <v-flex d-flex xs12 sm6 md6 pt-0>
                <v-card color="lighten-1">
                    <v-card-text>
                        <v-card-title>
                            <span class="title font-weight-light">{{$t('rpc.result')}}</span>

                            <v-spacer/>

                            <v-btn
                                    small
                                    @click="copyResult"
                            >
                                <span> {{$t('rpc.copy')}}</span>
                            </v-btn>
                        </v-card-title>
                        <v-card-text>
                            <v-alert
                                    :value="copySuccess"
                                    type="success"
                                    transition="scale-transition"
                            >
                                {{$t('rpc.copySuccess')}}
                            </v-alert>
                            <div id="jsonResponseEditor" style="height: 500px" class="json-editor">

                            </div>
                        </v-card-text>
                    </v-card-text>
                </v-card>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script lang="ts">
    import {Component, Vue, Watch} from "vue-property-decorator";
    import {State, Action} from 'vuex-class';

    import {Endpoint, Service} from "@/store/basic/types";

    // @ts-ignore
    import JSONEditor from "jsoneditor"
    import "jsoneditor/dist/jsoneditor.css";

    const namespace: string = 'call';

    @Component({components: {}})
    export default class Call extends Vue {

        private currentEndpoints: any = null;

        private copySuccess: boolean = false

        private service: Service = new Service();
        private endpoint: string = "";
        private otherEndpoint: string = "";

        private reqJSONEditor?: JSONEditor;

        private rspJSONEditor?: JSONEditor;

        @Action('getServiceDetails', {namespace})
        getServiceDetails: any;

        @Action('postServiceRequest', {namespace})
        postServiceRequest: any;

        @State(state => state.call.services)
        services?: Service[];

        @State(state => state.call.requestResult)
        requestResult?: object;

        @State(state => state.call.requestLoading)
        requestLoading?: boolean;


        @Watch("requestResult")
        resultChange(rr: any) {
            this.rspJSONEditor.set(rr)
            this.rspJSONEditor.expandAll();
        }

        created() {

        }

        mounted() {
            this.renderJSONEditor();
            this.getServiceDetails()
        }

        postRequest() {

            let endpoint = this.endpoint;
            if (endpoint == 'other') {
                endpoint = this.otherEndpoint;
            }

            let postData = {
                endpoint: endpoint,
                request: JSON.stringify(this.reqJSONEditor.get()),
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

        copyResult() {
            let that = this
            // @ts-ignore
            this.copySuccess = this.$xools.copyTxt(JSON.stringify(this.rspJSONEditor.get(), null, 2),
                function (success: boolean) {
                    that.copySuccess = success
                    setTimeout(() => {
                        that.copySuccess = false
                    }, 2000)
                })
        }

        renderJSONEditor() {
            let containerReq = document.getElementById("jsonRequestEditor");
            this.reqJSONEditor = new JSONEditor(containerReq, {mode: 'code', mainMenuBar: false});

            let json = {};
            this.reqJSONEditor.set(json)

            let containerRsp = document.getElementById("jsonResponseEditor");
            this.rspJSONEditor = new JSONEditor(containerRsp, {mode: 'tree', search: true});
        }

        formatRequestJSON() {
            this.reqJSONEditor.format()
        }
    }
</script>


<style>

    .codeFlex .v-card__text {
        padding-bottom: 30px;
    }
</style>