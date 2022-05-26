package util

// Tree 树形结构
type Tree interface {
    GetChildren() []Tree
    SetChildren(Tree, []Tree)
}

// filterTree 树形结构过滤，深度优先，先过滤叶子结点
func FilterTree[T any](root *T, getter func(node *T) []T, setter func(node *T, children []T), filter func(node *T) bool) *T {
    var children []T
    for _, v := range getter(root) {
        n := FilterTree(&v, getter, setter, filter)
        if n != nil {
            children = append(children, *n)
        }
    }
    setter(root, children)
    ok := filter(root)
    if ok {
        return root
    }
    return nil
}

// forTree 树形结构中序遍历
func ForTree[T any](root *T, getter func(node *T) []T, setter func(node *T, children []T), each func(node *T)) {
    each(root)
    var children []T
    for _, v := range getter(root) {
        ForTree(&v, getter, setter, each)
        children = append(children, v)
    }
    setter(root, children)
}

// forTree 树形结构遍历（先遍历孩子节点）
func ForTreeChild[T any](root *T, getter func(node *T) []T, setter func(node *T, children []T), each func(node *T)) {
    var children []T
    for _, v := range getter(root) {
        ForTreeChild(&v, getter, setter, each)
        children = append(children, v)
    }
    setter(root, children)
    each(root)
}
