// 分页工具类
package database

// "fmt"

type Page[T any] struct {
    CurrentPage int64 `json:"current_page"` // 当前页
    PageSize    int64 `json:"page_size"`    // 单页数
    Total       int64 `json:"total"`        // 总条数
    Pages       int64 `json:"pages"`        // 总页数
    Data        []T   `json:"data"`         // 数据
}

type PageResponse[T any] struct {
    CurrentPage int64 `json:"current_page"` // 当前页
    PageSize    int64 `json:"page_size"`    // 单页数
    Total       int64 `json:"total"`        // 总条数
    Pages       int64 `json:"pages"`        // 总页数
    Data        []T   `json:"data"`         // 数据
}

// 将原始的Page结构体转换为前端需要的PageResponse结构体
func NewPageResponse[T any](page Page[T]) *PageResponse[T] {
    return &PageResponse[T]{
        CurrentPage: page.CurrentPage,
        PageSize:    page.PageSize,
        Total:       page.Total,
        Pages:       page.Pages,
        Data:        page.Data,
    }
}

func (page *Page[T]) SelectPage(db *DB, columns []any, table string, wrapper func(*Query)) (e error) {
    // var model T
    db.Select(db.Expr("COUNT(*) AS `count`")).From(table).WhereWrapper(wrapper).Scan(&page.Total).Execute()
    // DB.Model(&model).Where(wrapper).Count(&page.Total)
    if page.Total == 0 {
        // 没有符合条件的数据，直接返回一个T类型的空列表
        page.Data = []T{}
        return
    }

    size := page.PageSize
    offset := int((page.CurrentPage - 1) * size)

    // 查询结果可以直接存到Page的Data字段中，因为编译的时候page.Data是有确定类型的
    db.Select(columns...).From(table).WhereWrapper(wrapper).Offset(offset).Limit(int(size)).Scan(&page.Data).Execute()
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
