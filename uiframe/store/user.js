export default {
    namespaced: true,
    state: () => ({
        userInfo: (() => {
            const MOCKUSER = {
                username: 'mockuser',
            }
            const userInfoJson = localStorage.getItem("userInfo");
            return userInfoJson ? JSON.parse(userInfoJson) : MOCKUSER;
        })(),
    }),
    getters: {
        getToken: (state) => state.user.userInfo.token,
    },
    mutations: {
        setUserInfo(state, payload) {
            state.userInfo = payload;
            localStorage.setItem("userInfo", JSON.stringify(payload));
        },
        clearUserInfo(state) {
            state.userInfo = {};
            localStorage.removeItem("userInfo");
        },
    },
};
