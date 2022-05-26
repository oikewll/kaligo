/**
 * 通用接口，挂载到common的namespace
 */

export default (app) => ({
    init: () => {
        return app.$axios.get("/api/init");
    },
});
