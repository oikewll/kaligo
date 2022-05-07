package util

func CastSlice[T any, B any](in []T) []B {
    return MapSlice(in, func(t T) B { return any(t).(B) })
}

func MapSlice[T, U any](from []T, transformer func(T) U) []U {
    ret := make([]U, len(from))
    for i, v := range from {
        ret[i] = transformer(v)
    }
    return ret
}

func FlatSlice[T ~[]E, E any](from []T) T {
    ret := make([]E, 0)
    for _, v := range from {
        ret = append(ret, v...)
    }
    return ret
}
