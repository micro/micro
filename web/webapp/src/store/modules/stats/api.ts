import {MutationTree, ActionTree} from 'vuex';
import * as TYPES from '../../mutation-types';

import {Error} from '@/store/basic/types'
import {APIStatsState, Stats} from './types';

import {getAPIGatewayServices} from '@/api/registry';
import {getAPIStats} from '@/api/stats';

const namespaced: boolean = true;

const state: APIStatsState = {
    loaded: false,
    services: [],
    currentNodeStats: new Stats(),
    xError: null,
    pageLoading: false
}

const mutations: MutationTree<any> = {

    [TYPES.SET_API_STATS_DATA_LOADING](state: APIStatsState, loading): void {
        state.pageLoading = loading
    },

    [TYPES.SET_API_STATS_SERVICES](state: APIStatsState, services): void {
        state.services = services
        state.pageLoading = false
    },

    [TYPES.SET_API_STATS_NODE_STATS](state: APIStatsState, stats: Stats): void {
        state.currentNodeStats = stats
        state.pageLoading = false
        state.loaded = true
    },

    [TYPES.SET_API_STATS_DATA_ERROR](state: APIStatsState, error: Error): void {
        state.xError = error
        state.pageLoading = false
    },
};

const actions: ActionTree<any, any> = {

    async getAPIGatewayServices({commit}) {
        commit(TYPES.SET_API_STATS_DATA_LOADING, true);
        const res: Ajax.AjaxResponse = await getAPIGatewayServices();
        commit(TYPES.SET_API_STATS_SERVICES, res.data);
    },

    async getStats({commit}, {name, address}) {

        commit(TYPES.SET_API_STATS_DATA_LOADING, true);
        const res: Ajax.AjaxResponse = await getAPIStats(name, address);
        // @ts-ignore
        if (res.counters) {
            commit(TYPES.SET_API_STATS_NODE_STATS, res);
        } else {
            commit(TYPES.SET_API_STATS_DATA_ERROR, res);
        }
    },
};


export default {
    namespaced,
    state,
    mutations,
    actions,
};
