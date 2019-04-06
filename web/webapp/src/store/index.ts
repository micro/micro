import Vue from 'vue';
import Vuex from 'vuex';
import state from './state';

import call from './modules/call';
import registry from './modules/registry';
import servicesStats from './modules/stats/services';
import apiStats from './modules/stats/api';


// index.js or main.js
import 'vuetify/dist/vuetify.min.css';

import Vuetify from 'vuetify';

Vue.use(Vuetify);
Vue.use(Vuex);

export default new Vuex.Store({
    state,
    mutations: {},
    actions: {
        init: () => {
        },
    },
    modules: {
        apiStats,
        call,
        registry,
        servicesStats
    },
});
