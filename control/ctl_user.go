package control

import (
	"net/http"
    "io"
    "log"
)
type User struct {
}

func (this *User) TestHome() {
    log.Print("TestUser...")
}

func (this *User) Index(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "User Index\n");
    log.Print("Index")
}

func (this *User) Register(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Home Register\n");
    log.Print("Register")
}

func (this *User) Login(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Login\n");
    log.Print("Login")
}


