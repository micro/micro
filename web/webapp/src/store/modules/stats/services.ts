import {MutationTree, ActionTree} from 'vuex';
import * as TYPES from '../../mutation-types';

import {ServicesStatsState, Stats, mergeNodes} from './types';

import {getServices, getService} from '@/api/registry';
import {getStats} from '@/api/stats';

const namespaced: boolean = true;

const state: ServicesStatsState = {
    services: [],
    currentNodes: [],
    nodeStats: new Stats(),
    xError: '',
    pageLoading: false
}

const mutations: MutationTree<any> = {

    [TYPES.SET_SERVICES_STATS_DATA_LOADING](state: ServicesStatsState, loading): void {
        state.pageLoading = loading
    },

    [TYPES.SET_SERVICES_STATS_SERVICES](state: ServicesStatsState, {services}): void {
        state.services = services
        state.pageLoading = false
    },

    [TYPES.SET_SERVICES_STATS_SERVICE_DETAIL](state: ServicesStatsState, {name, detail}): void {
        state.currentNodes = mergeNodes(detail)
        state.pageLoading = false
    },

    [TYPES.SET_SERVICES_STATS_NODE_STATS](state: ServicesStatsState, stats: Stats): void {
        state.nodeStats = stats
        state.pageLoading = false
    },

    [TYPES.SET_SERVICES_STATS_DATA_ERROR](state: ServicesStatsState, error: string): void {
        state.xError = error
        state.pageLoading = false
    },
};

const actions: ActionTree<any, any> = {

    async getServices({commit}) {

        commit(TYPES.SET_SERVICES_STATS_DATA_LOADING, true);
        const res: Ajax.AjaxResponse = await getServices();
        commit(TYPES.SET_SERVICES_STATS_SERVICES, {
            services: res.data
        });
    },

    async getNodes({commit}, name: string) {

        commit(TYPES.SET_SERVICES_STATS_DATA_LOADING, true);
        const res: Ajax.AjaxResponse = await getService(name);
        commit(TYPES.SET_SERVICES_STATS_SERVICE_DETAIL, {name: name, detail: res.data});
    },

    async getStats({commit}, {name, address}) {

        commit(TYPES.SET_SERVICES_STATS_DATA_LOADING, true);
        const res: Ajax.AjaxResponse = await getStats(name, address);
        if (res.success) {
            commit(TYPES.SET_SERVICES_STATS_NODE_STATS, res.data);
        } else {
            commit(TYPES.SET_SERVICES_STATS_DATA_ERROR, res);
        }
    },


};

export default {
    namespaced,
    state,
    mutations,
    actions,
};
