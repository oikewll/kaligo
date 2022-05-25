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
        "/todo": {
            "get": {
                "tags": [
                    "todo"
                ],
                "summary": "List 获取所有 Todo",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "当前页数, 1开始",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "当前页数, 默认20",
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
                                "$ref": "#/definitions/model.Todo"
                            }
                        }
                    }
                }
            },
            "put": {
                "tags": [
                    "todo"
                ],
                "summary": "Update 更新单条或多条数据",
                "parameters": [
                    {
                        "type": "string",
                        "name": "date",
                        "in": "formData"
                    },
                    {
                        "type": "boolean",
                        "name": "done",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "name": "id",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "name": "title",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Todo"
                            }
                        }
                    }
                }
            },
            "post": {
                "tags": [
                    "todo"
                ],
                "summary": "Create 添加一条数据",
                "parameters": [
                    {
                        "type": "string",
                        "name": "date",
                        "in": "formData"
                    },
                    {
                        "type": "boolean",
                        "name": "done",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "name": "id",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "name": "title",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Todo"
                        }
                    }
                }
            },
            "delete": {
                "tags": [
                    "todo"
                ],
                "summary": "Delete 删除单条或多条数据",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Todo ID",
                        "name": "id",
                        "in": "query"
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
        },
        "/todo/{id}": {
            "get": {
                "tags": [
                    "todo"
                ],
                "summary": "Detail 获取单条数据详情",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Todo ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Todo"
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
                        "description": "当前页数, 1开始",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "当前页数, 默认20",
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
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "账号",
                        "name": "username",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "密码",
                        "name": "password",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "姓名",
                        "name": "realname",
                        "in": "query",
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
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "test@gmail.com",
                        "description": "邮箱",
                        "name": "emali",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "状态",
                        "name": "status",
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
                        "description": "账号",
                        "name": "username",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "密码",
                        "name": "password",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "姓名",
                        "name": "realname",
                        "in": "query",
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
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "test@gmail.com",
                        "description": "邮箱",
                        "name": "emali",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "状态",
                        "name": "status",
                        "in": "query"
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
            },
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
                        "in": "query"
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
        },
        "/user/login": {
            "post": {
                "tags": [
                    "User"
                ],
                "summary": "Login 账户登陆",
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
        "/user/logout": {
            "delete": {
                "tags": [
                    "User"
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
        "/user/{id}": {
            "get": {
                "tags": [
                    "User"
                ],
                "summary": "Detail 用户信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
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
        }
    },
    "definitions": {
        "model.Purview": {
            "type": "object"
        },
        "model.Todo": {
            "type": "object",
            "properties": {
                "date": {
                    "type": "string"
                },
                "done": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "model.User": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "creatorId": {
                    "type": "integer"
                },
                "deletedAt": {
                    "type": "string"
                },
                "deletorId": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "groups": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "isFirstLogin": {
                    "type": "boolean"
                },
                "password": {
                    "type": "string"
                },
                "purviews": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Purview"
                    }
                },
                "realname": {
                    "type": "string"
                },
                "sessionExpire": {
                    "type": "string"
                },
                "sessionID": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "uid": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "updatorId": {
                    "type": "integer"
                },
                "username": {
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
