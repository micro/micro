import {MutationTree, ActionTree} from 'vuex';
import * as TYPES from '../../mutation-types';

import {RegistryState} from './types';
import {getRegistries} from '@/api/registry';

const namespaced: boolean = true;

const state: RegistryState = {
    registries: [],
    pageLoading: false,
}

const mutations: MutationTree<any> = {
    [TYPES.SET_REGISTRY_SERVICES](state: RegistryState, {registries}): void {
        state.registries = registries
        state.pageLoading = false
    },

    [TYPES.SET_REGISTRY_TABLE_LOADING](state: RegistryState, loading: boolean): void {
        state.pageLoading = loading
    },
};

const actions: ActionTree<any, any> = {

    async getRegistries({commit}) {

        commit(TYPES.SET_REGISTRY_TABLE_LOADING, true);

        const res: Ajax.AjaxResponse = await getRegistries();
        commit(TYPES.SET_REGISTRY_SERVICES, {
            registries: res.data
        });
    },

    async getRegistry({commit}) {

        commit(TYPES.SET_REGISTRY_TABLE_LOADING, true);

        const res: Ajax.AjaxResponse = await getRegistry();
        commit(TYPES.SET_REGISTRY_SERVICES, {
            registries: res.data
        });
    },
};

export default {
    namespaced,
    state,
    mutations,
    actions,
};
