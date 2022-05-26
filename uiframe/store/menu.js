// import Menu_main from "@/config/menu/main";
// import Menu_system from "@/config/menu/system";
// import Mock_menu from "@/config/menu/mock-menu.js";

// let Menu = Menu_main;
// if (localStorage.getItem("Menu_sidebar") && localStorage.getItem("Menu_sidebar") === "system") {
//     Menu = Menu_system;
// }

const session = window.sessionStorage;
const CACHE_KEYS = {
    'main_tabs': 'cache_main_tabs',
    'main_tabs_active_name': 'cache_main_tabs_active_name',
}

const setCache = (key, value) => {
    switch (toString.call(value)) {
        case '[object Object]':
        case '[object Array]':
            value = JSON.stringify(value);
            break;
    }
    session.setItem(CACHE_KEYS[key], value);
}

let cache_main_tabs = [];
try {
    cache_main_tabs = JSON.parse(session.getItem(CACHE_KEYS['main_tabs']));
    if (toString.call(cache_main_tabs) !== "[object Array]") {
        cache_main_tabs = [];
    }
} catch (e) {
    session.removeItem(CACHE_KEYS['main_tabs']);
}

export default {
    namespaced: true,
    state: () => ({
        main_tabs: cache_main_tabs,   // 当前存在的标签页
        main_tabs_active_name: '',    // 当前激活的标签页名字
        menu_setting: [],             // 后台返回的菜单配置，包括顶部和边栏的数据
        menu_sidebar: [],             // 单独边栏菜单的数据
    }),
    getters: {
        getMenu: (state) => state.menu_sidebar,
        getMenuSetting: (state) => state.menu_setting,
    },
    mutations: {
        SET_MAIN_TABS_ACTIVE_NAME: (state, name) => {
            state.main_tabs_active_name = name;
            setCache('main_tabs_active_name', state.main_tabs_active_name)
        },
        // 添加标签页
        MAIN_TABS_ADD: (state, tab) => {
            let item = state.main_tabs.find((item => item.name == tab.name));
            if (!item) {
                state.main_tabs.push(tab);
                setCache('main_tabs', state.main_tabs)
            }
            state.main_tabs_active_name = tab.name;
            setCache('main_tabs_active_name', state.main_tabs_active_name)
        },
        // 更新tab栏目
        MAIN_TABS_UPDATE: (state, tabs) => {
            state.main_tabs = tabs;
            setCache('main_tabs', tabs)
        },
        SET_MAIN_SIDE: (state, type) => {
            state.menu_sidebar = state.menu_setting[Number(type)].children;
        },
        // 这里是给api拿到菜单后注入vuex的commit
        SET_MENUDATA: (state, data = []) => {
            state.menu_sidebar = data[0].children;
            state.menu_setting = data;
        }
    },
};
