import Utils from "@/utils";

export default {
    namespaced: true,
    state: () => ({
        theme: 'light',
        sidebar: {
            opened: (() => {
                const status = localStorage.getItem("sidebarStatus");
                return Utils.isMobile() ? false : status == false ? false : true;
            })(),
            withoutAnimation: false,
        },
        device: "desktop",
        size: localStorage.getItem("size") || "medium",
        msgCenter: false,
        lockScreen: !!localStorage.getItem("lockScreen"),
    }),
    getters: {
        getSidebar: (state) => state.sidebar,
        getLockScreen: (state) => state.lockScreen,
        getTheme: (state) => state.theme,
    },
    mutations: {
        UPDATE_THEME: (state, payload) => state.theme = payload,
        UPDATE_SIDEBAR(state){
            state.sidebar.opened = !state.sidebar.opened;
            state.sidebar.withoutAnimation = false;
            if (state.sidebar.opened) {
                localStorage.setItem("sidebarStatus", 1);
            } else {
                localStorage.setItem("sidebarStatus", 0);
            }
        },
        UPDATE_MSGCENTER(state) {
            state.msgCenter = !state.msgCenter;
        },
        UPDATE_LOCKSTATE(state) {
            if (!state.lockScreen){
                localStorage.setItem("lockScreen", 1);
            }else{
                localStorage.removeItem("lockScreen");
            }
            state.lockScreen = !state.lockScreen;
        },
        // 通用更新app的state
        UPDATE_STATUS(state, obj) {
            state[obj.key] = obj.val;
        }
    },
    action: {},
};
