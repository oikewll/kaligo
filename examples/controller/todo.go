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
// @Summary 获取所有 Todo
// @tags    todo
// @Param   page       query integer false "当前页数, 1开始"
// @Param   size       query integer false "当前页数, 默认20"
// @Success 200 {object} []model.Todo
// @Router  /todo [get]
func (c *Todo) List() {
    data, _, err := model.Todo{}.List()
    result(c.Context, data, err)
}

// Detail 获取单条数据详情
// @Summary Detail 获取单条数据详情
// @tags    todo
// @Param   id       path integer true "Todo ID"
// @Success 200 {object} model.Todo
// @Router  /todo/:id [GET]
func (c *Todo) Detail() {
    id := c.QueryValue("id")
    if len(id) == 0 {
        result(c.Context, nil, errors.New("id is required"))
    }
    data, err := model.Todo{}.Detail(id)
    result(c.Context, data, err)
}

// Create 添加一条数据
// @Summary Create 添加一条数据
// @tags    todo
// @Param   todo formData model.Todo true "Todo"
// @Success 200 {object} model.Todo
// @Router  /todo [POST]
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
// @Summary Update 更新单条或多条数据
// @tags    todo
// @Param   todo formData model.Todo true "Todo"
// @Success 200 {object} []model.Todo
// @Router  /todo [PUT]
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
// @Summary Delete 删除单条或多条数据
// @tags    todo
// @Param   id       query integer false "Todo ID"
// @Success 200 {integer} integer
// @Router  /todo [DELETE]
func (c *Todo) Delete() {
    id := c.QueryValue("id")
    if len(id) == 0 {
        result(c.Context, nil, errors.New("id is required"))
    }
    data, err := model.Todo{}.Delete(id)
    result(c.Context, data, err)
}
