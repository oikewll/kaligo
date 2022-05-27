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
    getForm: () =>{
        return app.$axios.get("/api/user/createform");
    }
});
