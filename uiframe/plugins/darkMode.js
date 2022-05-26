
const getMode = (isDarkMode = false) => {
    return isDarkMode ? "dark" : "light";
}

export default ({ store }) => {
    if (!window.matchMedia || !window.matchMedia("(prefers-color-scheme: dark)")){
        return console.warn(`浏览器不支持深色模式媒体查询`)
    }
    try {
        const theme = getMode(window.matchMedia("(prefers-color-scheme: dark)").matches);
        setTimeout(() => {
            store.commit("app/UPDATE_THEME", theme);
        }, 1)

        window
            .matchMedia("(prefers-color-scheme: dark)")
            .addEventListener("change", (e) => {
                store.commit("app/UPDATE_THEME", getMode(e.matches));
            });
    } catch(e) {
        // nothing to do with window.matchMedia('').addEventListener error in old safari version
    }
};
