package controller

import (
    "errors"
    "examples/model"

    "github.com/owner888/kaligo"
)

type Todo struct {
    kaligo.Controller
}

// List 分页获取数据列表
func (c *Todo) List() {
    data, _, err := model.Todo{}.List()
    result(c.Context, data, err)
}

// Detail 获取单条数据详情
func (c *Todo) Detail() {
    id := c.QueryValue("id")
    if len(id) == 0 {
        result(c.Context, nil, errors.New("id is required"))
    }
    data, err := model.Todo{}.Detail(id)
    result(c.Context, data, err)
}

// Create 添加一条数据
func (c *Todo) Create() {
    var todo model.Todo
    err := c.JsonBodyValue(&todo)
    if err != nil {
        result(c.Context, nil, err)
    }
    data, err := model.Todo{}.Create(todo)
    result(c.Context, data, err)
}

// Update 更新单条或多条数据
func (c *Todo) Update() {
    var todo model.Todo
    err := c.JsonBodyValue(&todo)
    if err != nil {
        result(c.Context, nil, err)
    }
    data, err := model.Todo{}.Update(todo)
    result(c.Context, data, err)
}

// Delete 删除单条或多条数据
func (c *Todo) Delete() {
    id := c.QueryValue("id")
    if len(id) == 0 {
        result(c.Context, nil, errors.New("id is required"))
    }
    data, err := model.Todo{}.Delete(id)
    result(c.Context, data, err)
}
