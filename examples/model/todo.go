package model

import (
    "strings"
)

type Todo struct {
    Title string `json:"title"`
    Date  string `json:"date"`
    Done  bool   `json:"done"`
    Id    ID     `json:"id"`
}

func (m *Todo) table() string { return "todo" }

// List 分页获取数据列表
func (m Todo) List() ([]Todo, int, error) {
    var t []Todo
    _, err := DB.Select("id", "title", "date", "done").From(m.table()).Scan(&t).Execute()
    return t, 0, err
}

// Detail 获取单条数据详情
func (m Todo) Detail(id string) (t Todo, err error) {
    _, err = DB.Select("*").From(m.table()).Where("id", "=", id).Scan(&t).Execute()
    return
}

// Create 添加一条数据
func (m Todo) Create(t Todo) (ID, error) {
    q, err := DB.Insert(m.table(), []string{"title", "date"}).Values([]any{t.Title, t.Date}).Execute()
    return ID(q.LastInsertId), err
}

// Update 更新单条或多条数据
func (m Todo) Update(t Todo) (ID, error) {
    q, err := DB.Update(m.table()).Set(map[string]string{"title": t.Title, "date": t.Date}).Where("id", "=", t.Id).Execute()
    return ID(q.LastInsertId), err
}

// Delete 删除单条或多条数据
func (m Todo) Delete(ids string) (bool, error) {
    _, err := DB.Delete(m.table()).Where("id", "in", strings.Split(ids, ",")).Execute()
    return err != nil, err
}
