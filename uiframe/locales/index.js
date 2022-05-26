import Vue from "vue";
import VueI18n from "vue-i18n";
import Cookies from "js-cookie";
import customZhCn from "./lang/zh";
import customEnUs from "./lang/en";
import customKM from "./lang/km";

Vue.use(VueI18n);

const messages = {
    en: {
        ...customEnUs
    },
    "zh-cn": {
        ...customZhCn
    },
    km: {
        ...customKM
    }
};

export function getLocale() {
    const chooseLanguage = Cookies.get("language");
    if (chooseLanguage) return chooseLanguage;

    // if has not choose language
    // const language = (
    //     navigator.language.split('-')[0] || navigator.browserLanguage
    // ).toLowerCase()
    // const locales = Object.keys(messages)
    // for (const locale of locales) {
    //     if (language.includes(locale)) {
    //         return locale
    //     }
    // }
    return "zh-cn";
}

export function setLocale(locale) {
    i18n.locale = locale;
    Cookies.set("language", locale);
}

const i18n = new VueI18n({
    locale: getLocale(),
    fallbackLocale: "zh-cn",
    messages
});

Cookies.set("language", i18n.locale);

export default i18n;
