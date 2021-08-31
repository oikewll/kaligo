package database

// Expression is the struct for Set the expression string.
// 不会被转义的表达式
// e := &Expression{value: "COUNT(users.id)"}
type Expression struct {
    value string  // Raw expression string
}

// Value Get the expression value as a string.
func (e *Expression) Value() string {
    return e.value
}

