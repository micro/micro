import Vue from 'vue';
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import locale from 'element-ui/lib/locale'
// @ts-ignore
import enLocale from 'element-ui/lib/locale/lang/en'
import cnLocale from 'element-ui/lib/locale/lang/zh-CN'
import XTools from '@/utils/index';

if (XTools.Utils.getDefaultLan() == 'en') {
    // configure language
    locale.use(enLocale)
}

Vue.use(ElementUI);


export function setLan(lan: string) {

    if (lan == 'en') {
        locale.use(enLocale)
        return
    }

    debugger

    if (lan == 'cn') {
        locale.use(cnLocale)
        return
    }
}
