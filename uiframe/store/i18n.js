let locale = 'en';

if(process.browser){
    locale = localStorage.getItem('LANG') || 'zh-CN';
}

export default {
    namespaced: true,
    state: () => ({
        locales: ['zh-CN', 'en'],
        locale
    }),
    getters: {
        getLang: (state) => state.locale,
    },
    mutations: {
        set_lang(state, locale) {
            if (state.locales.indexOf(locale) !== -1) {
                state.locale = locale
                localStorage.setItem('LANG',locale);
                this.app.i18n.locale = locale;
            }
        },
        add_lang(state, lang){
            state.locales.indexOf(lang) == -1 && state.locales.push(lang);
        }
    }
}