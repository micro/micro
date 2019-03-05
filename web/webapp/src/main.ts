import Vue from "vue";
import App from "./App.vue";
import Vuetify from "vuetify";
import router from "./router/router";
import store from "./store/store";
import "font-awesome/css/font-awesome.css";
import "./theme/default.styl";
import VeeValidate from "vee-validate";
import VueI18n from 'vue-i18n'
import I18nMsg from '@/i18n/I18n'

import Util from '@/utils'

let Truncate = require("lodash.truncate");
import "vuetify/dist/vuetify.min.css";

Vue.config.productionTip = false;

Vue.filter("truncate", Truncate);
Vue.use(VeeValidate, {fieldsBagName: "formFields"});
Vue.use(Vuetify, {

    options: {
        themeVariations: ["primary", "secondary", "accent"],
        extra: {
            mainToolbar: {
                color: "primary"
            },
            sideToolbar: {},
            sideNav: "primary",
            mainNav: "primary lighten-1",
            bodyBg: ""
        }
    }
});
Vue.use(VueI18n)

const i18n = new VueI18n({
    locale: getDefaultLan(), // set locale
    messages: I18nMsg// set locale messages,
})

new Vue({
    i18n,
    router,
    store,
    render: h => h(App),
    template: "<App/>"
}).$mount("#app")
    .$on('localeChange', (locale: string) => {
        i18n.locale = locale
        Util.setCookie("locale", locale, 30)
    });


function getDefaultLan() {
    let locale = Util.getCookieValue('locale')
    return locale ? locale : 'en';
}
