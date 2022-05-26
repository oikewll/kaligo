// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/check_token": {
            "post": {
                "tags": [
                    "Auth"
                ],
                "summary": "CheckToken 检查 CSRF Token 是否存在",
                "parameters": [
                    {
                        "type": "string",
                        "description": "CSRF Token",
                        "name": "csrf",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "tags": [
                    "Auth"
                ],
                "summary": "Login 账户登陆",
                "parameters": [
                    {
                        "type": "string",
                        "default": "test",
                        "description": "账号",
                        "name": "username",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "test",
                        "description": "密码",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "default": true,
                        "description": "记住密码",
                        "name": "remember",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "delete": {
                "tags": [
                    "Auth"
                ],
                "summary": "Logout 账户退出",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/auth/token": {
            "get": {
                "tags": [
                    "Auth"
                ],
                "summary": "Token 产生一个 CSRF Token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/init": {
            "get": {
                "tags": [
                    "Home"
                ],
                "summary": "初始化接口",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/permissions": {
            "get": {
                "tags": [
                    "Home"
                ],
                "summary": "权限列表: 权限选择、权限展示",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/sessions": {
            "get": {
                "tags": [
                    "Session"
                ],
                "summary": "Session 信息",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Session 添加简介",
                "tags": [
                    "Session"
                ],
                "summary": "Session 添加",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "tags": [
                    "Session"
                ],
                "summary": "Session 删除",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/sessions/destory": {
            "delete": {
                "tags": [
                    "Session"
                ],
                "summary": "Session 销毁",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/user": {
            "get": {
                "tags": [
                    "User"
                ],
                "summary": "List 分页获取用户信息",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "当前页数",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 20,
                        "description": "每页记录",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.User"
                            }
                        }
                    }
                }
            },
            "post": {
                "tags": [
                    "User"
                ],
                "summary": "Create 添加一条数据",
                "parameters": [
                    {
                        "type": "string",
                        "default": "test",
                        "description": "账号",
                        "name": "username",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "test",
                        "description": "密码",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "test",
                        "description": "姓名",
                        "name": "realname",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "enum": [
                            "1",
                            "2",
                            "3"
                        ],
                        "type": "string",
                        "description": "所属组IDs",
                        "name": "groups",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "test@gmail.com",
                        "description": "邮箱",
                        "name": "emali",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "状态",
                        "name": "status",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                }
            }
        },
        "/user/createform": {
            "get": {
                "tags": [
                    "User"
                ],
                "summary": "CreateForm 用户添加表单",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Form"
                        }
                    }
                }
            }
        },
        "/user/{id}": {
            "get": {
                "tags": [
                    "User"
                ],
                "summary": "Detail 用户信息",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "账号ID",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "put": {
                "tags": [
                    "User"
                ],
                "summary": "Update 更新单条或多条数据",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "账号ID",
                        "name": "id",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "default": "test",
                        "description": "账号",
                        "name": "username",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "test",
                        "description": "密码",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "test",
                        "description": "姓名",
                        "name": "realname",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "enum": [
                            "1",
                            "2",
                            "3"
                        ],
                        "type": "string",
                        "description": "所属组IDs",
                        "name": "groups",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "test@gmail.com",
                        "description": "邮箱",
                        "name": "emali",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "状态",
                        "name": "status",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.User"
                            }
                        }
                    }
                }
            }
        },
        "/user{id}": {
            "delete": {
                "tags": [
                    "User"
                ],
                "summary": "Delete 删除单条或多条数据",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "账号ID",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Component": {
            "type": "object",
            "properties": {
                "field": {
                    "description": "字段",
                    "type": "string"
                },
                "placeholder": {
                    "description": "输入框占位文本",
                    "type": "string"
                },
                "props": {
                    "description": "属性",
                    "$ref": "#/definitions/model.Props"
                },
                "tips": {
                    "description": "格式:bm://open:com.xxx.xxx",
                    "type": "string"
                },
                "title": {
                    "description": "标题",
                    "type": "string"
                },
                "type": {
                    "description": "类型: input、editor、textarea、file、image",
                    "type": "string"
                },
                "validate": {
                    "description": "验证: required、numeric..., 详细看本页最下面",
                    "$ref": "#/definitions/model.Validate"
                },
                "value": {
                    "description": "默认值",
                    "type": "string"
                }
            }
        },
        "model.Form": {
            "type": "object",
            "properties": {
                "components": {
                    "description": "表单组件",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Component"
                    }
                },
                "csrf": {
                    "type": "string"
                },
                "method": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                }
            }
        },
        "model.Props": {
            "type": "object",
            "properties": {
                "type": {
                    "description": "类型: text、password",
                    "type": "string"
                }
            }
        },
        "model.Purview": {
            "description": "用户权限",
            "type": "object"
        },
        "model.User": {
            "description": "User account information",
            "type": "object",
            "required": [
                "username"
            ],
            "properties": {
                "avatar": {
                    "description": "用户头像地址",
                    "type": "string"
                },
                "email": {
                    "description": "邮箱地址",
                    "type": "string"
                },
                "first_login": {
                    "description": "是否首次登录",
                    "type": "boolean"
                },
                "groups": {
                    "description": "用户所属权限组",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "id": {
                    "description": "用户ID",
                    "type": "integer"
                },
                "purviews": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Purview"
                    }
                },
                "realname": {
                    "description": "用户昵称",
                    "type": "string"
                },
                "status": {
                    "description": "状态",
                    "type": "integer"
                },
                "uid": {
                    "description": "UID",
                    "type": "string"
                },
                "username": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        },
        "model.Validate": {
            "type": "object",
            "properties": {
                "enum": {
                    "description": "枚举类型",
                    "type": "string"
                },
                "len": {
                    "description": "字段长度",
                    "type": "integer"
                },
                "max": {
                    "description": "最大长度",
                    "type": "integer"
                },
                "message": {
                    "description": "校验文案",
                    "type": "string"
                },
                "min": {
                    "description": "最小长度",
                    "type": "integer"
                },
                "pattern": {
                    "description": "正则表达式校验",
                    "type": "string"
                },
                "required": {
                    "description": "是否必选",
                    "type": "boolean"
                },
                "type": {
                    "description": "内建校验类型: date、datetimerange、emali、array",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Kaligo Example API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
