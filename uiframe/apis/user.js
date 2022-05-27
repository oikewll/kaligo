/**
 * 用户权限接口
 */

export default (app) => ({
    init: () => {
        return app.$axios.get("/api/init");
    },
    /**
     * params 有数据的就是post，否则就是get请求数据
     */
    permission: (params = []) => {
        if (params.length === 0){
            return app.$axios.get("/api/permissions")
        }
        return app.$axios.post("/api/permissions", {
            purview: params
        })
    },
    addByForm: (params) =>{
        return app.$axios.post("/api/user", params)
    },
    getList: ({ page = 1, size = 20 } = {})=>{
        return app.$axios.get('/api/user', { params: { page, size } })
    },
    getForm: (id = '') =>{
        if (!!id){
            return app.$axios.get(`/api/user/updateform/${id}`);
        }
        return app.$axios.get("/api/user/createform");
    },
    // 根据id获取用户信息
    get: (id) => {
        return app.$axios.get(`/api/user/${id}`);
    },
    put: (id, params) => {
        return app.$axios.put(`/api/user/${id}`, params);
    },
    delete: (id) => {
        return app.$axios.delete(`/api/user/${id}`);
    },
});
