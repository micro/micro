import Vue from "vue";

import '@/plugins/tools';
import '@/plugins/vuetify';
import i18n from '@/plugins/I18n';
import {setLan} from '@/plugins/elementui';
import '@/theme/style';

import App from "./App.vue";
import router from "./router/router";
import store from "./store";

import XTools from '@/utils/index';


Vue.config.productionTip = false;


new Vue({
    i18n,
    router,
    store,
    render: h => h(App),
    template: "<App/>"
}).$mount("#app")
    .$on('localeChange', (locale: string) => {
        i18n.locale = locale
        setLan(locale)
        XTools.Utils.setCookie("locale", locale, 30)
    });


