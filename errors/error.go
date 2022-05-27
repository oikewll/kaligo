package errors

type Error struct {
	Code int
	Msg  string
	Data any
}

// StatusNotAcceptable 客户端要的资源类型服务端无法提供, 比如要 XML, 但是接口只有 JSON
var (
	LOGIN_UNKNOWN = NewError(202, "User does not exist")
	LOGIN_ERROR   = NewError(203, "Wrong account or password")
	VALID_ERROR   = NewError(300, "Parameter error")
	ERROR         = NewError(400, "Operation failed")
	UNAUTHORIZED  = NewError(401, "You are not logged in.")
	NOT_FOUND     = NewError(404, "Resources do not exist")
	INNER_ERROR   = NewError(500, "System exception")
)

func (e *Error) Error() string { 
    return e.Msg 
}

func NewError(code int, msg string) *Error {
	return &Error{
		Msg:  msg,
		Code: code,
	}
}

func GetError(e *Error, data interface{}) *Error {
	return &Error{
		Msg:  e.Msg,
		Code: e.Code,
		Data: data,
	}
}
