import {MutationTree, ActionTree} from 'vuex';
import * as TYPES from '../../mutation-types';

import {ServicesStatsState, Stats} from './types';

import {getServices, getService, getMicroServices} from '@/api/registry';
import {getStats} from '@/api/stats';

const namespaced: boolean = true;

const state: ServicesStatsState = {
    services: [],
    nodesStatsMap: new Map(),
    pageLoading: false,
    cardLoading: new Map(),
    cardLoadingChanged: '',
    xError: '',
}

const mutations: MutationTree<any> = {

    [TYPES.SET_SERVICES_STATS_PAGE_LOADING](state: ServicesStatsState, loading: boolean): void {
        state.pageLoading = loading
    },

    [TYPES.SET_SERVICES_STATS_DATA_LOADING](state: ServicesStatsState, {address, loading}: { address: string, loading: boolean }): void {
        state.cardLoading.set(address, loading)
    },

    [TYPES.SET_SERVICES_STATS_SERVICES](state: ServicesStatsState, {services}): void {
        state.services = services
        state.pageLoading = false
    },

    [TYPES.SET_SERVICES_STATS_NODE_STATS](state: ServicesStatsState, {address, stats}: { address: string, stats: Stats }): void {
        state.nodesStatsMap.set(address, stats)
        state.cardLoading.set(address, false)
        state.cardLoadingChanged = new Date().toJSON()
    },

    [TYPES.SET_SERVICES_STATS_DATA_ERROR](state: ServicesStatsState, error: string): void {
        state.xError = error
    },
};

const actions: ActionTree<any, any> = {

    async getServices({commit}) {

        commit(TYPES.SET_SERVICES_STATS_PAGE_LOADING, true);
        const res: Ajax.AjaxResponse = await getMicroServices();
        commit(TYPES.SET_SERVICES_STATS_SERVICES, {
            services: res.data
        });
    },

    async getStats({commit, dispatch}, {name, address}) {

        commit(TYPES.SET_SERVICES_STATS_DATA_LOADING, {address: address, loading: true});
        // await new Promise(resolve => setTimeout(resolve, 2000));
        const res: Ajax.AjaxResponse = await getStats(name, address);
        if (res.success) {
            commit(TYPES.SET_SERVICES_STATS_NODE_STATS, {address: address, stats: res.data});
        } else {
            commit(TYPES.SET_SERVICES_STATS_DATA_ERROR, res);
            commit(TYPES.SET_SERVICES_STATS_DATA_LOADING, {address: address, loading: false});
        }
    },
};

export default {
    namespaced,
    state,
    mutations,
    actions,
};
