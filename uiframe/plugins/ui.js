import Vue from 'vue'
import VueI18n from 'vue-i18n'
import customZhCn from '~/lang/zh-CN'
import customEnUs from '~/lang/en'

import Element from 'element-ui'
import zhCnLocale from 'element-ui/lib/locale/lang/zh-CN'
import enUsLocale from 'element-ui/lib/locale/lang/en'
import '~/assets/less/element-variables.scss'

export default ({ app, store }) => {
    Vue.use(VueI18n);
    // Set i18n instance on app
    // This way we can use it in middleware and pages asyncData/fetch
    const messages = {
        'zh-CN': Object.assign({}, zhCnLocale, customZhCn),
        'en': Object.assign({}, enUsLocale, customEnUs)
    }
    const i18n = new VueI18n({
        locale: store.state.i18n.locale,
        fallbackLocale: 'zh-CN',
        messages
    });

    Vue.use(Element, {
        i18n: (key, value) => i18n.t(key, value)
    })

    app.i18n = i18n;
}
