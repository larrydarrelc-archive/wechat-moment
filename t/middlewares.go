package t

import (
    "strconv"
    "net/http"
    "github.com/codegangsta/martini"
    "github.com/astaxie/beego/orm"
)

func LoginRequired(resp http.ResponseWriter, req *http.Request, c martini.Context) {
    banAccess := func() {
        resp.WriteHeader(http.StatusUnauthorized)
    }

    id, err := strconv.Atoi(req.Header.Get("X-ID"))
    if err != nil {
        banAccess()
    }
    token := req.Header.Get("X-TOKEN")
    user := User{Id: id}
    ok, err := user.CheckLogin(token)
    if err != nil {
        banAccess()
    }
    if !ok {
        banAccess()
    }
    err = orm.NewOrm().Read(&user)
    if err != nil {
        banAccess()
    }
    c.Map(&user)
}
