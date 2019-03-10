import Vue from "vue";

import '@/plugins/tools';
import '@/plugins/vuetify';

import "font-awesome/css/font-awesome.css";
import "./theme/default.styl";
import "vuetify/dist/vuetify.min.css";

import App from "./App.vue";
import router from "./router/router";
import store from "./store";
import I18nMsg from '@/i18n/I18n'
import VueI18n from 'vue-i18n'
import XTools from '@/utils/index';

Vue.config.productionTip = false;


Vue.use(VueI18n)

const i18n = new VueI18n({
    locale: getDefaultLan(), // set locale
    messages: I18nMsg// set locale messages,
})

function getDefaultLan() {
    let locale = XTools.Utils.getCookieValue('locale')
    return locale ? locale : 'en';
}

new Vue({
    i18n,
    router,
    store,
    render: h => h(App),
    template: "<App/>"
}).$mount("#app")
    .$on('localeChange', (locale: string) => {
        i18n.locale = locale
        XTools.Utils.setCookie("locale", locale, 30)
    });


