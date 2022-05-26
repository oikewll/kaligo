/**
 * 系统路由配置
 * @return {JSON}
 */

export default [
    {
        id: "1-0",
        name: "user_manage",
        path: "/user-manage",
        meta: {
            title: "menu.user-manage",
            icon: "el-icon-user"
        },
        hidden: false,
    },
    {
        id: "1-1",
        name: "session_manage",
        path: "/session-manage",
        meta: {
            title: "menu.session-manage",
            icon: "el-icon-chat-round"
        },
        hidden: false,
    },
    {
        id: "1-2",
        name: "system_manage",
        path: "/system-manage",
        meta: {
            title: "menu.system-manage",
            icon: "el-icon-document-copy"
        },
        hidden: false,
    },
    {
        id: "1-3",
        name: "cache_manage",
        path: "/cache-manage",
        meta: {
            title: "menu.cache-manage",
            icon: "el-icon-coin"
        },
        hidden: false,
    },
    {
        id: "1-4",
        name: "resources_manage",
        path: "/resources-manage",
        meta: {
            title: "menu.resources-manage",
            icon: "el-icon-notebook-2"
        },
        hidden: false,
        children: [
            {
                id: "1-4-1",
                name: "images_manage",
                path: "/images-manage",
                meta: {
                    title: "menu.images-manage",
                    icon: "el-icon-aim"
                },
                hidden: false,
            },
            {
                id: "1-4-2",
                name: "videos_manage",
                path: "/videos-manage",
                meta: {
                    title: "menu.videos-manage",
                    icon: "el-icon-orange"
                },
                hidden: false,
            },
        ],
    },
    {
        id: "1-5",
        name: "work_manage",
        path: "/work-manage",
        meta: {
            title: "menu.work-manage",
            icon: "el-icon-set-up"
        },
        hidden: false,
    },
]