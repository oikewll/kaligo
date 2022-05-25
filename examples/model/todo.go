package model

import (
    "strings"
)

type Todo struct {
    Title string `db:"title" json:"title"`
    Date  string `db:"date" json:"date"`
    Done  bool   `db:"done" json:"done"`
    Id    ID     `db:"id" json:"id"`
}

func (m *Todo) Table() string { return "todo" }

// List 分页获取数据列表
func (m Todo) List() ([]Todo, int, error) {
    var t []Todo
    _, err := DB.Select("id", "title", "date", "done").From(m.Table()).Scan(&t).Execute()
    return t, 0, err
}

// Detail 获取单条数据详情
func (m Todo) Detail(id string) (t Todo, err error) {
    _, err = DB.Select("*").From(m.Table()).Where("id", "=", id).Scan(&t).Execute()
    return
}

// Create 添加一条数据
func (m Todo) Create(t Todo) (ID, error) {
    q, err := DB.Insert(m.Table(), []string{"title", "date"}).Values([]any{t.Title, t.Date}).Execute()
    return ID(q.LastInsertId), err
}

// Update 更新单条或多条数据
func (m Todo) Update(t Todo) (ID, error) {
    q, err := DB.Update(m.Table()).Set(map[string]string{"title": t.Title, "date": t.Date}).Where("id", "=", t.Id).Execute()
    return ID(q.LastInsertId), err
}

// Delete 删除单条或多条数据
func (m Todo) Delete(ids string) (bool, error) {
    _, err := DB.Delete(m.Table()).Where("id", "in", strings.Split(ids, ",")).Execute()
    return err != nil, err
}
