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

func CompactMapSliceE[T, U any, E error](from []T, transformer func(T) (U, E)) ([]U, E) {
    ret := make([]U, len(from))
    var err E
    for i, v := range from {
        r, e := transformer(v)
        if any(e) != nil {
            err = e
        } else {
            ret[i] = r
        }
    }
    return ret, err
}

func CompactMapSlice[T, U any, E bool](from []T, transformer func(T) (U, E)) ([]U, E) {
    ret := make([]U, len(from))
    var ok E
    for i, v := range from {
        r, o := transformer(v)
        if ok {
            ret[i] = r
        }
        ok = ok && o
    }
    return ret, ok
}

func FlatSlice[T ~[]E, E any](from []T) T {
    ret := make([]E, 0)
    for _, v := range from {
        ret = append(ret, v...)
    }
    return ret
}
