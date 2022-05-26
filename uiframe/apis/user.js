/**
 * 用户权限接口
 */

export default (app) => ({
    init: () => {
        return app.$axios.get("/api/init");
    },
});
