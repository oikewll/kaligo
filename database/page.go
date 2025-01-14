// 分页工具类
package database

import (
    "github.com/owner888/kaligo/form"
)

type Page[T any] struct {
    Page    int64      `json:"page"`     // 当前页(1 开始)
    Size    int64      `json:"size"`     // 单页数
    Total   int64      `json:"total"`    // 总条数
    Data    []T        `json:"data"`     // 数据
    Table   form.Table `json:"table"`    // 表格
    HasPrev bool       `json:"has_prev"` // 是否有上一页
    HasNext bool       `json:"has_next"` // 是否有下一页
}

func (p *Page[T]) UpdateTotal(total int64) {
    p.Total = total
    p.HasPrev = p.Page > 1
    p.HasPrev = total > p.Page*p.Size
}

type PageResponse[T any] struct {
    Page  int64 `json:"page"`  // 当前页
    Size  int64 `json:"size"`  // 单页数
    Total int64 `json:"total"` // 总条数
    Data  []T   `json:"data"`  // 数据
}

// 将原始的Page结构体转换为前端需要的PageResponse结构体
func NewPageResponse[T any](page Page[T]) *PageResponse[T] {
    return &PageResponse[T]{
        Page:  page.Page,
        Size:  page.Size,
        Total: page.Total,
        Data:  page.Data,
    }
}

// Group 的时候 Count 可能不准, isCount 用于要不要算总数
func (page *Page[T]) SelectPage(db *DB, columns []any, table string, wrapper func(query *Query, isCount bool)) (err error) {
    where := func(isCount bool) func(q *Query) {
        return func(q *Query) { wrapper(q, isCount) }
    }
    // var model T
    db.Select(db.Expr("COUNT(*) AS `count`")).From(table).WhereWrapper(where(true)).Scan(&page.Total).Execute()
    page.UpdateTotal(page.Total)
    // DB.Model(&model).Where(wrapper).Count(&page.Total)
    if page.Total == 0 {
        // 没有符合条件的数据，直接返回一个T类型的空列表
        page.Data = []T{}
        return
    }
    size := page.Size
    offset := int((page.Page - 1) * size)

    // 查询结果可以直接存到Page的Data字段中，因为编译的时候page.Data是有确定类型的
    db.Select(columns...).From(table).WhereWrapper(where(false)).Offset(offset).Limit(int(size)).Scan(&page.Data).Execute()
    // e = DB.Model(&model).Where(wrapper).Scopes(Paginate(page)).Find(&page.Data).Error
    return
}

// Paginate加上T就行
// func Paginate[T any](page *Page[T]) func(db *Query) *Query {
//     return func(db *Query) *Query {
//         if page.CurrentPage <= 0 {
//             page.CurrentPage = 0
//         }
//         switch {
//         case page.PageSize > 100:
//             page.PageSize = 100
//         case page.PageSize <= 0:
//             page.PageSize = 10
//         }
//         page.Pages = page.Total / page.PageSize
//         if page.Total%page.PageSize != 0 {
//             page.Pages++
//         }
//         p := page.CurrentPage
//         if page.CurrentPage > page.Pages {
//             p = page.Pages
//         }
//         size := page.PageSize
//         offset := int((p - 1) * size)
//         return db.Offset(offset).Limit(int(size))
//     }
// }
