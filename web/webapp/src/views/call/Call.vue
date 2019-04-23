<template>
    <el-container>
        <el-row type="flex" justify="space-around" style="width: 100%">
            <el-col :span="12">
                <el-card :body-style="{ padding: '10px 10px 10px 20px'}">
                    <el-form label-position="left" label-width="120px">
                        <!-- <el-form-item :label="$t('base.service')">
                             <el-select
                                     v-model="service"
                                     filterable
                                     clearable
                                     :placeholder="$t('base.service')"
                                     @change="changeService"
                             >
                                 <el-option
                                         v-for="(item, index) in services"
                                         :key="index"
                                         :label="item.name"
                                         :value="item">
                                 </el-option>
                             </el-select>
                         </el-form-item>-->
                        <el-form-item :label="$t('base.service')">
                            <el-select v-model="serviceName"
                                       filterable
                                       clearable
                                       :disabled="callData && callData.specialModel"
                                       :placeholder="$t('base.service')"
                                       @change="changeService"
                            >
                                <el-option
                                        v-for="(item, index) in services"
                                        :key="index"
                                        :label="item.name"
                                        :value="item.name">
                                </el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item :label="$t('base.endpoint')">
                            <el-select
                                    v-model="endpoint"
                                    filterable
                                    :disabled="callData && callData.specialModel"
                                    clearable
                                    :placeholder="$t('base.endpoint')"
                                    @change="changeEndpoint"
                            >
                                <el-option
                                        v-for="(item, index) in currentEndpoints"
                                        :key="index"
                                        :value="item.name">
                                </el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item :label="$t('base.otherEndpoint')">
                            <el-input :disabled="endpoint != 'other'"
                                      v-model="otherEndpoint"
                                      :placeholder="$t('rpc.inputOtherEndpoint')"
                            ></el-input>
                        </el-form-item>
                        <el-form-item :label="$t('rpc.request')">
                            <div style="float: right">
                                <el-col :span="12">
                                    <el-button
                                            size="small"
                                            @click="formatRequestJSON"
                                    >
                                        <span> {{$t('rpc.formatJSON')}}</span>
                                    </el-button>
                                </el-col>
                                <el-col :span="12">
                                    <el-button
                                            :disabled="!(serviceName && (endpoint && endpoint != 'other' || endpoint == 'other' && otherEndpoint))"
                                            size="small"
                                            @click="postRequest"
                                    >
                                        <span> {{$t('rpc.postRequest')}}</span>
                                    </el-button>
                                </el-col>
                            </div>
                        </el-form-item>
                        <div id="jsonRequestEditor" style="height: 300px" class="json-editor">
                        </div>
                    </el-form>
                </el-card>
            </el-col>
            <el-col :span="12" style="margin-left: 20px">
                <el-card :body-style="{ padding: '10px 10px 10px 20px'}">

                    <el-row>
                        <el-col :span="12"><span class="title font-weight-light">{{$t('rpc.result')}}</span></el-col>
                        <el-col :span="12">
                            <el-button
                                    size="small"
                                    style="float:right"
                                    @click="copyResult"
                            >
                                <span> {{$t('rpc.copy')}}</span>
                            </el-button>
                        </el-col>
                    </el-row>
                    <v-card-text>
                        <div id="jsonResponseEditor" style="height: 484px" class="json-editor">

                        </div>
                    </v-card-text>
                </el-card>
            </el-col>
        </el-row>
    </el-container>
</template>
<style scoped>
    .el-container {
        margin-right: 20px;
    }
</style>

<script lang="ts">
    import {Component, Vue, Watch, Prop} from "vue-property-decorator";
    import {State, Action,} from 'vuex-class';

    import {Endpoint, Service, Value} from "@/store/basic/types";

    // @ts-ignore
    import JSONEditor from "jsoneditor"
    import "jsoneditor/dist/jsoneditor.css";

    const namespace: string = 'call';

    @Component({components: {}})
    export default class Call extends Vue {

        @Prop()
        private callData: any;

        private currentEndpoints: any = null;

        private serviceName?: string = "";

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

        @State(state => state.call.xError)
        error?: any;


        @Watch("requestResult")
        resultChange(rr: any) {
            this.rspJSONEditor.set(rr)
            // this.rspJSONEditor.expandAll();
        }

        created() {
        }

        mounted() {
            this.renderJSONEditor();
            this.getServiceDetails()
            if (this.callData && this.callData.specialModel) {
                this.serviceName = this.callData.serviceName;
                this.endpoint = this.callData.endpoint;
                this.renderRequestJSON(this.endpoint, this.callData.endpoints)
            }
        }

        postRequest() {

            let endpoint = this.endpoint;
            if (endpoint == 'other') {
                endpoint = this.otherEndpoint;
            }

            let postData = {
                endpoint: endpoint,
                request: JSON.stringify(this.reqJSONEditor.get()),
                service: this.serviceName
            }

            this.postServiceRequest(postData);
        }

        changeService(serviceName: string) {

            this.endpoint = null;
            this.otherEndpoint = null;
            this.currentEndpoints = [];
            this.reqJSONEditor.set({});

            this.services.forEach((s: Service, i: number) => {
                if (s.name != serviceName) {
                    return
                }
                if (s.endpoints) {
                    this.currentEndpoints = s.endpoints
                } else {
                    this.currentEndpoints = []
                }
            })

            let hasOther = false
            this.currentEndpoints.forEach((item: Endpoint) => {
                if (item.name == "other") {
                    hasOther = true;
                }
            });

            if (!hasOther) {
                this.currentEndpoints.push({name: 'other', value: -1})
            }
        }

        changeEndpoint(endpoint: string) {
            this.endpoint = endpoint
            this.reqJSONEditor.set({});
            // change request json

            this.renderRequestJSON(endpoint, this.currentEndpoints)

        }

        renderRequestJSON(endpointName: string, endpoints: Endpoint[]) {
            endpoints && endpoints.forEach((item: Endpoint) => {
                if (item.name == endpointName) {
                    let json = {}
                    item.request.values.forEach((value: Value) => {

                        let defaultValue;
                        switch (value.type) {
                            case 'int64':
                            case 'int32':
                            case 'int16':
                            case 'int':
                                defaultValue = 0
                                break
                            case 'string':
                            default:
                                defaultValue = ''
                                break;
                        }
                        json[value.name] = defaultValue
                    })

                    this.reqJSONEditor.set(json)
                    return
                }
            })
        }

        copyResult() {
            let that = this
            // @ts-ignore
            this.$xools.copyTxt(JSON.stringify(this.rspJSONEditor.get(), null, 2),
                function (success: boolean) {
                    that.$message({
                        // @ts-ignore
                        message: that.$t('rpc.copySuccess'),
                        type: 'success'
                    });
                })
        }

        renderJSONEditor() {
            let containerReq = document.getElementById("jsonRequestEditor");
            this.reqJSONEditor = new JSONEditor(containerReq, {mode: 'code', mainMenuBar: false});

            let json = {};
            this.reqJSONEditor.set(json)

            let containerRsp = document.getElementById("jsonResponseEditor");
            this.rspJSONEditor = new JSONEditor(containerRsp, {mode: 'code', mainMenuBar: false});
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
