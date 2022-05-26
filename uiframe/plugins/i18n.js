export default ({ app, store }, inject) => {
    const messages = {
        "zh-CN": { ...require("~/lang/zh-CN.js") },
        "en": { ...require("~/lang/en.js") },
    };

    // 暴露出来
    const lang = (item) => {
        const [langClass, langName] = item.split(/\./i);
        const langText = messages[store.state.i18n.locale].default[langClass][langName];
        return langText;
    };

    inject("lang", lang);
};
