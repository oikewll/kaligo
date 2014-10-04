package control

import (
	"net/http"
    "fmt"
    "io"
    "time"
    "encoding/json"
    "github.com/garyburd/redigo/redis"
    "epooll/util"
)

type Index struct {
}

func (this *Index) Index(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Index --- Index\n");
}

func (this *Index) Login(w http.ResponseWriter, r *http.Request) {
    //io.WriteString(w, "Index --- Login\n");

    data := map[string]interface{} {
        "type":r.FormValue("type"),
        "uuid":r.FormValue("uuid"),
        "channel":r.FormValue("channel"),
        "user_id":r.FormValue("user_id"),
        "game_id":r.FormValue("game_id"),
        "game_name":r.FormValue("game_name"),
        "ip":r.FormValue("ip"),
        "utma":r.FormValue("utma"),
        "resolution":r.FormValue("resolution"),
        "device":r.FormValue("device"),
        "system":r.FormValue("system"),
        "net":r.FormValue("net"),
        "operators":r.FormValue("operators"),
        "location":r.FormValue("location"),
        "isnew":r.FormValue("isnew"),
        "version":r.FormValue("version"),
        "addtime":time.Now().Unix(),
    }
    dbValue, _ := json.Marshal(data)

    c := util.RedisPool.Get()
    defer c.Close()

    logstr := fmt.Sprintf("%s\n", dbValue)
    if ok, err := util.PutFile("./login.log", logstr, 1); !ok {
        panic(err)
    }
    dbKey := "netgame:login"
    if ok, err := redis.Bool(c.Do("LPUSH", dbKey, dbValue)); !ok {
        panic(err)
    }
    
    retValue := map[string]interface{} {
        "state":true,
        "msg":"",
    }
    retJson, _ := json.Marshal(retValue)
    retStr := fmt.Sprintf("%s", retJson);
    io.WriteString(w, retStr);
}


func (this *Index) Role_login(w http.ResponseWriter, r *http.Request) {

    data := map[string]interface{} {
        "type":r.FormValue("role_login"),
        "uuid":r.FormValue("uuid"),
        "channel":r.FormValue("channel"),
        "role_id":r.FormValue("role_id"),
        "role_name":r.FormValue("role_name"),
        "role_level":r.FormValue("role_level"),
        "server_id":r.FormValue("server_id"),
        "server_name":r.FormValue("server_name"),
        "addtime":time.Now().Unix(),
    }
    dbValue, _ := json.Marshal(data)

    c := util.RedisPool.Get()
    defer c.Close()

    logstr := fmt.Sprintf("%s\n", dbValue)
    if ok, err := util.PutFile("./role_login.log", logstr, 1); !ok {
        panic(err)
    }
    dbKey := "netgame:role_login"
    if ok, err := redis.Bool(c.Do("LPUSH", dbKey, dbValue)); !ok {
        panic(err)
    }
    
    retValue := map[string]interface{} {
        "state":true,
        "msg":"",
    }
    retJson, _ := json.Marshal(retValue)
    retStr := fmt.Sprintf("%s", retJson);
    io.WriteString(w, retStr);

}
