import Vue from 'vue';
import Vuex from 'vuex';
import state from './state';

import call from './modules/call';
import registry from './modules/registry';


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
        call,
        registry,
    },
});
