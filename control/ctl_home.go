package control

import (
	"net/http"
    "io"
    "epooll/util"
)
type Home struct {
}

func (this *Home) Index(w http.ResponseWriter, r *http.Request) {
    c := util.RedisPool.Get()
    defer c.Close()
    io.WriteString(w, "Home Index\n");
}

