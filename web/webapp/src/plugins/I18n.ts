import Vue from 'vue';
import VueI18n from 'vue-i18n';
import I18nMsg from '@/i18n/I18n';
import XTools from "@/utils";

Vue.use(VueI18n)


const i18n = new VueI18n({
    locale: XTools.Utils.getDefaultLan(), // set locale
    messages: I18nMsg// set locale messages,
})


export default i18n
