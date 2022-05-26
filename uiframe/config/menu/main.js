/**
 * 主菜单路由配置
 * @return {JSON}
 */

export default [
    {
        id: "0-0",
        name: "dashboard",
        path: "/",
        meta: {
            title: "menu.dashboard",
            icon: "el-icon-odometer",
        },
        hidden: false,
    },
    {
        id: "0-1",
        name: "analysis",
        path: "/analysis",
        meta: {
            title: "menu.analysis-results",
            icon: "el-icon-data-analysis",
        },
        hidden: false,
    },
    {
        id: "0-3",
        name: "keywords",
        path: "/keywords",
        meta: {
            title: "menu.keywords",
            icon: "el-icon-mobile",
        },
        hidden: false,
    },
    {
        id: "0-2",
        name: "results",
        path: "/results",
        meta: {
            title: "menu.hit-results",
            icon: "el-icon-place",
        },
        hidden: false,
    },
    {
        id: "0-4",
        name: "filter",
        path: "/filter",
        meta: {
            title: "menu.filter-domain",
            icon: "el-icon-coin",
        },
        hidden: false,
    },
    {
        id: "0-5",
        name: "acquisition",
        path: "/acquisition",
        meta: {
            title: "menu.acquisition-settings",
            icon: "el-icon-truck",
        },
        hidden: true,
    },
];
