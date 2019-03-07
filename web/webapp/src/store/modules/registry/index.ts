import {MutationTree, ActionTree} from 'vuex';
import * as TYPES from '../../mutation-types';

import {RegistryState, Service} from './types';
import {getServices, getService} from '@/api/registry';

const namespaced: boolean = true;

const state: RegistryState = {
    services: [],
    pageLoading: false,
    serviceDetailLoading: false,
    serviceDetail: [],
}

const mutations: MutationTree<any> = {
    [TYPES.SET_REGISTRY_SERVICES](state: RegistryState, {services}): void {
        state.services = services
        state.pageLoading = false
    },

    [TYPES.SET_REGISTRY_TABLE_LOADING](state: RegistryState, loading: boolean): void {
        state.pageLoading = loading
    },

    [TYPES.SET_REGISTRY_SERVICE_DETAIL](state: RegistryState, serviceDetail): void {
        state.serviceDetail = serviceDetail
        state.serviceDetailLoading = false
    },

    [TYPES.SET_REGISTRY_SERVICE_DETAIL_LOADING](state: RegistryState, loading: boolean): void {
        state.serviceDetailLoading = loading
    },
};

const actions: ActionTree<any, any> = {

    async getServices({commit}) {

        commit(TYPES.SET_REGISTRY_TABLE_LOADING, true);

        const res: Ajax.AjaxResponse = await getServices();
        commit(TYPES.SET_REGISTRY_SERVICES, {
            services: res.data
        });
    },

    async getService({commit}, name: string) {

        commit(TYPES.SET_REGISTRY_SERVICE_DETAIL_LOADING, true);

        const res: Ajax.AjaxResponse = await getService(name);
        commit(TYPES.SET_REGISTRY_SERVICE_DETAIL, res.data);
    },
};

export default {
    namespaced,
    state,
    mutations,
    actions,
};
