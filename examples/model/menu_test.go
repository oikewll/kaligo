package model

import (
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestMenuParse(t *testing.T) {
    menu, err := Menu{}.Load(strings.NewReader(testMenuStr))
    assert.NoError(t, err)
    assert.Equal(t, []Menu{}, menu)
}

func TestPermisionFilter(t *testing.T) {
    pStr := "GET-/content,POST-/content"
    menu, _ := Menu{}.Load(strings.NewReader(testMenuStr))
    p := Permission{}.Parse(pStr)
    root := &Menu{Children: menu}
    root = root.Permission(p)
    assert.Equal(t, Menu{}, root)
}

func TestPermision(t *testing.T) {
    pStr := "GET-/api/todo,GET-/api/todo/:id,POST-/api/todo,PUT-/api/todo/:id,DELETE-/api/todo/:id"
    assert.Equal(t, []Permission{}, Permission{}.Parse(pStr))
}

var testMenuStr = `
<xml>
    <!-- //这里的子菜单为隐性项目 -->
    <menu name='常用' icon='fa fa-th-list'>    
        <menu name='内容管理' icon='fa fa-file'>
            <menu name='内容列表' path='/content' method='GET' show='1' reload='1' top='1' />
            <menu name='内容添加' path='/content' method='POST' />
            <menu name='内容修改' path='/content/:id' method='PUT' />
            <menu name='内容删除' path='/content/:id' method='DELETE' />
            <menu name='内容详情' path='/content/:id' method='GET' />
            <!-- <menu name='{lang.menu_content_detail}' path='content' method='GET' /> -->
        </menu>
        <menu name='会员管理' icon='fa fa-users'>
            <menu name='会员列表' path='/member' method='GET' show='1' reload='1' />
            <menu name='会员添加' path='/member' method='POST' />
            <menu name='会员修改' path='/member/:id' method='PUT' />
            <menu name='会员删除' path='/member/:id' method='DELETE' />
            <menu name='会员详情' path='/member/:id' method='GET' />
        </menu>
    </menu>
</xml>
`

// var testMenu = Menu{{ID: 0, Path: "", Show: false, Method: "", Meta: Meta{Name: "常用", Icon: "fa fa-th-list"}, Children:[]Menu{
//          Menu{ID: 0, Path: "", Show: false, Method: "", Meta: Meta{Name: "内容管理", Icon: "fa fa-file"}, Children: []Menu{
//          Menu{ID:0, Path: "/content", Show: true, Method: "GET", Meta: Meta{Name: "内容列表", Icon: ""}, Children: []Menu(nil)
//      },
//      Menu{ID: 0, : "/content", Show: false, Method: "POST", Meta: Meta{Name: "内容添加", Icon: ""}, Children: []Menu(nil)}, Menu{ID: 0, Pathcontent/:id", Show: false, Method: "PUT", Meta: Meta{Name: "内容修改", Icon: ""}, Children: []Menu(nil)}, Menu{ID: 0, Path: "/ent/:id", Show: false, Method: "DELETE", Meta: Meta{Name: "内容删除", Icon: ""}, Children: []Menu(nil)}, Menu{ID: 0, Path: "/cont:id", Show: false, Method: "GET", Meta: Meta{Name: "内容详情", Icon: ""}, Children: []Menu(nil)}}}, Menu{ID:Path: "", Show: false, Method: "", Meta: Meta{Name: "会员管理", Icon: "fa fa-users"}, Children: []Menu{Menu{ID: 0, Path: "/member", Shorue, Method: "GET", Meta: Meta{Name: "会员列表", Icon: ""}, Children: []Menu(nil)}, Menu{ID: 0, Path: "/member", Show: f, Method: "POST", Meta: Meta{Name: "会员添加", Icon: ""}, Children: []Menu(nil)}, Menu{ID: 0, Path: "/member/:id", Show: f, Method: "PUT", Meta: Meta{Name: "会员修改", Icon: ""}, Children: []Menu(nil)}, Menu{ID: 0, Path: "/member/:id", Show: falsethod: "DELETE", Meta: Meta{Name: "会员删除", Icon: ""}, Children: []Menu(nil)}, Menu{ID: 0, Path: "/member/:id", Show: false, Me: "GET", Meta: Meta{Name: "会员详情", Icon: ""}, Children: []Menu(nil)}}}}}}
