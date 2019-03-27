import {MutationTree, ActionTree} from 'vuex';
import * as TYPES from '../../mutation-types';

import {CallState} from './types';
import {getServiceDetails, postServiceRequest} from '@/api/call';


const namespaced: boolean = true;

const state: CallState = {
    requestLoading: false,
    services: [],
    requestResult: {},
    xError: null
}

const mutations: MutationTree<any> = {
    [TYPES.SET_CALL_SERVICES](state: CallState, {services}): void {
        state.services = services
        state.requestLoading = false
    },

    [TYPES.SET_CALL_LOADING](state: CallState, loading: boolean): void {
        state.requestLoading = loading
    },

    [TYPES.SET_CALL_RESULT](state: CallState, result: object): void {
        state.requestResult = result
    },

    [TYPES.SET_CALL_ERROR](state: CallState, error: object): void {
        state.xError = error
    },

};

const actions: ActionTree<any, any> = {

    async getServiceDetails({commit}) {

        commit(TYPES.SET_CALL_LOADING, true);

        const res: Ajax.AjaxResponse = await getServiceDetails();
        commit(TYPES.SET_CALL_SERVICES, {
            services: res.data
        });
    },

    async postServiceRequest({commit}, req: object) {

        commit(TYPES.SET_CALL_LOADING, true);
        const res: Ajax.AjaxResponse = await postServiceRequest(req);

        if (res.success) {
            if (res.data.body) {
                res.data.body = JSON.parse(res.data.body)
            }

            commit(TYPES.SET_CALL_RESULT, res);
        } else {
            commit(TYPES.SET_CALL_RESULT, res);
        }
    },
};

export default {
    namespaced,
    state,
    mutations,
    actions,
};
