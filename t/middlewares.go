package t

import (
    "strconv"
    "net/http"
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
    "github.com/astaxie/beego/orm"
)

func LoginRequired(resp http.ResponseWriter, req *http.Request, c martini.Context, r render.Render) {
    banAccess := func() {
        r.JSON(http.StatusUnauthorized, Error("Login required."))
    }

    id, err := strconv.Atoi(req.Header.Get("X-ID"))
    if err != nil {
        banAccess()
        return
    }
    token := req.Header.Get("X-TOKEN")
    user := User{Id: id}
    ok, err := user.CheckLogin(token)
    if err != nil {
        banAccess()
        return
    }
    if !ok {
        banAccess()
        return
    }
    err = orm.NewOrm().Read(&user)
    if err != nil {
        banAccess()
        return
    }
    c.Map(&user)
}
