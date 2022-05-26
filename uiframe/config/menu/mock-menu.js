/**
 * 主菜单路由配置
 * @return {JSON}
 */

export default [
    {
        "id": 1,
        "path": "",
        "show": false,
        "top": false,
        "reload": false,
        "method": "",
        "name": "常用",
        "icon": "fa fa-th-list",
        "children": [
            {
                "id": 2,
                "path": "",
                "show": false,
                "top": false,
                "reload": false,
                "method": "",
                "name": "内容管理",
                "icon": "fa fa-file",
                "children": [
                    {
                        "id": 3,
                        "path": "/content",
                        "show": true,
                        "top": true,
                        "reload": true,
                        "method": "GET",
                        "name": "内容列表",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 4,
                        "path": "/content",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "POST",
                        "name": "内容添加",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 5,
                        "path": "/content/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "PUT",
                        "name": "内容修改",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 6,
                        "path": "/content/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "内容删除",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 7,
                        "path": "/content/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "内容详情",
                        "icon": "",
                        "children": null
                    }
                ]
            },
            {
                "id": 8,
                "path": "",
                "show": false,
                "top": false,
                "reload": false,
                "method": "",
                "name": "会员管理",
                "icon": "fa fa-users",
                "children": [
                    {
                        "id": 9,
                        "path": "/member",
                        "show": true,
                        "top": false,
                        "reload": true,
                        "method": "GET",
                        "name": "会员列表",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 10,
                        "path": "/member",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "POST",
                        "name": "会员添加",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 11,
                        "path": "/member/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "PUT",
                        "name": "会员修改",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 12,
                        "path": "/member/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "会员删除",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 13,
                        "path": "/member/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "会员详情",
                        "icon": "",
                        "children": null
                    }
                ]
            }
        ]
    },
    {
        "id": 14,
        "path": "",
        "show": false,
        "top": false,
        "reload": false,
        "method": "",
        "name": "系统",
        "icon": "fa fa-gear",
        "children": [
            {
                "id": 15,
                "path": "",
                "show": false,
                "top": false,
                "reload": false,
                "method": "",
                "name": "个人中心",
                "icon": "fa fa-user",
                "children": [
                    {
                        "id": 16,
                        "path": "/admin/password",
                        "show": true,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "登陆密码",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 17,
                        "path": "/admin/purview",
                        "show": true,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "我的权限",
                        "icon": "",
                        "children": null
                    }
                ]
            },
            {
                "id": 18,
                "path": "",
                "show": false,
                "top": false,
                "reload": false,
                "method": "",
                "name": "用户管理",
                "icon": "fa fa-users",
                "children": [
                    {
                        "id": 19,
                        "path": "/admin_group",
                        "show": true,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "用户组管理",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 20,
                        "path": "/admin_group",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "POST",
                        "name": "用户组添加",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 21,
                        "path": "/admin_group/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "PUT",
                        "name": "用户组修改",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 22,
                        "path": "/admin_group/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "用户组删除",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 23,
                        "path": "/admin_group/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "用户组详情",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 24,
                        "path": "/admin",
                        "show": true,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "用户管理",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 25,
                        "path": "/admin",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "POST",
                        "name": "用户添加",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 26,
                        "path": "/admin/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "PUT",
                        "name": "用户修改",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 27,
                        "path": "/admin/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "用户删除",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 28,
                        "path": "/admin/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "用户详情",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 29,
                        "path": "/admin_log/operation_log",
                        "show": true,
                        "top": false,
                        "reload": true,
                        "method": "GET",
                        "name": "操作日志",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 30,
                        "path": "/admin_log/operation_log/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "操作日志删除",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 31,
                        "path": "/admin_log/login_log",
                        "show": true,
                        "top": false,
                        "reload": true,
                        "method": "GET",
                        "name": "登录日志",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 32,
                        "path": "/admin_log/login_log/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "登录日志删除",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 33,
                        "path": "/admin_log/login_log/clear",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "清空三个月前登录日志",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 34,
                        "path": "/admin/enable/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "PUT",
                        "name": "用户激活",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 35,
                        "path": "/admin/disable/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "用户禁用",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 36,
                        "path": "/admin/reset_mfa",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "重设MFA",
                        "icon": "",
                        "children": null
                    }
                ]
            },
            {
                "id": 37,
                "path": "",
                "show": false,
                "top": false,
                "reload": false,
                "method": "",
                "name": "会话管理",
                "icon": "fa fa-rocket",
                "children": [
                    {
                        "id": 38,
                        "path": "/session/online",
                        "show": true,
                        "top": false,
                        "reload": true,
                        "method": "GET",
                        "name": "在线会话",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 39,
                        "path": "/session/history",
                        "show": true,
                        "top": false,
                        "reload": true,
                        "method": "GET",
                        "name": "历史会话",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 40,
                        "path": "/session/terminate",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "终止会话",
                        "icon": "",
                        "children": null
                    }
                ]
            },
            {
                "id": 41,
                "path": "",
                "show": false,
                "top": false,
                "reload": false,
                "method": "",
                "name": "系统管理",
                "icon": "fa fa-wrench",
                "children": [
                    {
                        "id": 42,
                        "path": "/config",
                        "show": true,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "配置管理",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 43,
                        "path": "/config",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "POST",
                        "name": "配置添加",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 44,
                        "path": "/config/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "PUT",
                        "name": "配置修改",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 45,
                        "path": "/config/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "配置删除",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 46,
                        "path": "/config/batch_edit",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "PUT",
                        "name": "配置批量修改",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 47,
                        "path": "/spam",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "防刷管理",
                        "icon": "",
                        "children": null
                    }
                ]
            },
            {
                "id": 48,
                "path": "",
                "show": false,
                "top": false,
                "reload": false,
                "method": "",
                "name": "缓存管理",
                "icon": "fa fa-cloud",
                "children": [
                    {
                        "id": 49,
                        "path": "/cache",
                        "show": true,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "缓存管理",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 50,
                        "path": "/cache",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "POST",
                        "name": "缓存添加",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 51,
                        "path": "/cache/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "PUT",
                        "name": "缓存修改",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 52,
                        "path": "/cache/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "缓存删除",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 53,
                        "path": "/cache/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "缓存详情",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 54,
                        "path": "/cache/clear",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "缓存清理",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 55,
                        "path": "/redis/keys",
                        "show": true,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "Redis键值管理",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 56,
                        "path": "/redis/server_info",
                        "show": true,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "Redis服务器信息",
                        "icon": "",
                        "children": null
                    }
                ]
            },
            {
                "id": 57,
                "path": "",
                "show": false,
                "top": false,
                "reload": false,
                "method": "",
                "name": "资源管理",
                "icon": "fa fa-file",
                "children": [
                    {
                        "id": 58,
                        "path": "/filemanage",
                        "show": true,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "文件管理",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 59,
                        "path": "/filemanage",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "POST",
                        "name": "文件新增",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 60,
                        "path": "/filemanage/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "PUT",
                        "name": "文件修改",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 61,
                        "path": "/filemanage/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "文件删除",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 62,
                        "path": "/upload/upload",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "POST",
                        "name": "文件上传",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 63,
                        "path": "/upload/upload_chunked",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "POST",
                        "name": "文件分块上传",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 64,
                        "path": "/upload/upload_html5",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "POST",
                        "name": "图片拖拉上传",
                        "icon": "",
                        "children": null
                    }
                ]
            },
            {
                "id": 65,
                "path": "",
                "show": false,
                "top": false,
                "reload": false,
                "method": "",
                "name": "作业中心",
                "icon": "fa fa-coffee",
                "children": [
                    {
                        "id": 66,
                        "path": "/crond",
                        "show": true,
                        "top": false,
                        "reload": false,
                        "method": "GET",
                        "name": "任务列表",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 67,
                        "path": "/crond",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "POST",
                        "name": "任务新增",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 68,
                        "path": "/crond/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "PUT",
                        "name": "任务修改",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 69,
                        "path": "/crond/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "任务删除",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 70,
                        "path": "/crond/enable/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "PUT",
                        "name": "任务激活",
                        "icon": "",
                        "children": null
                    },
                    {
                        "id": 71,
                        "path": "/crond/disable/:id",
                        "show": false,
                        "top": false,
                        "reload": false,
                        "method": "DELETE",
                        "name": "任务禁用",
                        "icon": "",
                        "children": null
                    }
                ]
            }
        ]
    }
];
