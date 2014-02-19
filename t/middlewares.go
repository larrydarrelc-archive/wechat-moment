package t

import (
    "strconv"
    "net/http"
    "log"
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
    "github.com/astaxie/beego/orm"
)

// Check if the request carry logined token & id.
//
// Each request should carry `X-ID` and `X-TOKEN` in the header.
// If the token is valid, it will also map logined user (`*User`) into
// request handler.
func LoginRequired(resp http.ResponseWriter,
                   req *http.Request,
                   c martini.Context,
                   r render.Render) {
    banAccess := func() {
        r.JSON(http.StatusUnauthorized, Error("Login required."))
    }

    id, err := strconv.Atoi(req.Header.Get("X-ID"))
    if err != nil {
        log.Fatal("X-ID cannot be parsed into `int`.", req.Header.Get("X-ID"), err)
        banAccess()
        return
    }
    token := req.Header.Get("X-TOKEN")

    user := User{Id: id}
    ok, err := user.CheckLogin(token)
    if err != nil {
        log.Fatal("User check login failed.", id, token, err)
        banAccess()
        return
    }
    if !ok {
        log.Print("User login failed.", id, token)
        banAccess()
        return
    }

    // Map logined user into request context.
    err = orm.NewOrm().Read(&user)
    if err != nil {
        log.Fatal("Read user failed.", id, err)
        banAccess()
        return
    }
    c.Map(&user)
}
