package t

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
    "github.com/astaxie/beego/orm"
)

// Routes:
//
//  Path            Method  Description
//  /user           POST    Create a user.
//  /user           PUT     Update user's profile.
//  /user/login     POST    Login a user.
//  /user/logout    GET     Logout a user.
//  /user/me        GET     Get current user's profile.
//  /user/:id       GET     Get user `/:id`'s profile.
func UserRoute(m *Application) {
    m.Post("/user", createUser)
    m.Put("/user", LoginRequired, updateUserProfile)
    m.Post("/user/login", loginUser)
    m.Get("/user/logout", LoginRequired, logoutUser)

    m.Get("/user/me", LoginRequired, getSelfProfile)
    m.Get("/user/:id", LoginRequired, getUserProfile)
}

// Create a user.
func createUser(req *http.Request, r render.Render) {
    o := orm.NewOrm()

    login := req.FormValue("login")
    if login == ""  {
        r.JSON(http.StatusForbidden,
               Error(fmt.Sprintf("Login name cannot be empty.")))
        return
    }
    name := req.FormValue("name")
    if name == ""  {
        r.JSON(http.StatusForbidden,
               Error(fmt.Sprintf("Name cannot be empty.")))
        return
    }
    password := req.FormValue("password")
    if len(password) < 5  {
        r.JSON(http.StatusForbidden,
               Error(fmt.Sprintf("Password aleast 5 character.")))
        return
    }

    user := User{Login: login}
    created, _, err := o.ReadOrCreate(&user, "Login");
    if err != nil {
        log.Print("Read user failed", user.Login, err)
        r.JSON(http.StatusForbidden, Error("Create user failed."))
        return
    }
    if !created {
        r.JSON(http.StatusConflict,
               Error(fmt.Sprintf("Login name %s already exists!", login)))
        return
    }

    user.Password = user.HashPassword(password)
    user.Name = name
    if _, err := o.Update(&user); err != nil {
        log.Print("Create user failed.", user.Name, user.Id, err)
        r.JSON(http.StatusForbidden, Error("Create user failed."))
        return
    }
    rv, err := user.Censor()
    if err != nil {
        log.Print("Censor user failed.", user.Id, err)
        r.JSON(http.StatusForbidden, Error("Create user failed."))
        return
    }
    r.JSON(http.StatusCreated, rv)
}

func getUserProfile(params martini.Params, r render.Render) {
    o := orm.NewOrm()

    id, err := strconv.Atoi(params["id"])
    if err != nil {
        log.Print("Cannot parse into `int`.", params["id"], err)
        r.JSON(http.StatusNotFound, Error("Read user profile failed."))
        return
    }
    user := User{Id: id}
    err = o.Read(&user)
    if err == orm.ErrNoRows {
        r.JSON(http.StatusNotFound,
               Error(fmt.Sprintf("User %d not exists", id)))
        return
    }
    rv, err := user.Censor()
    if err != nil {
        log.Print("User censor failed.", err, id)
        r.JSON(http.StatusForbidden, Error("Read user profile failed."))
        return
    }

    r.JSON(http.StatusOK, rv)
}

func getSelfProfile(u *User, r render.Render) {
    rv, err := u.Censor()
    if err != nil {
        log.Print("User censor failed.", err, u.Id)
        r.JSON(http.StatusForbidden, Error("Read user profile failed."))
        return
    }
    r.JSON(http.StatusOK, rv)
}

func updateUserProfile(user *User, req *http.Request, r render.Render) {
    name := req.FormValue("name")
    err := user.UpdateProfile(name)
    if err != nil {
        log.Print("Update user profile failed.", user.Id, err)
        r.JSON(http.StatusForbidden, Error("Update user profile failed."))
        return
    }
    r.JSON(http.StatusNoContent, "")
}

func loginUser(req *http.Request, r render.Render) {
    o := orm.NewOrm()

    var user User
    stat := o.QueryTable("user")
    stat = stat.Filter("login", req.FormValue("login"))
    stat = stat.Filter("password", user.HashPassword(req.FormValue("password")))
    if err := stat.One(&user); err == orm.ErrNoRows {
        r.JSON(http.StatusForbidden, Error("Name and password are mismatch!"))
        return
    } else if err != nil {
        log.Print("Read user failed.", err)
        r.JSON(http.StatusForbidden, Error("Login failed."))
        return
    }

    token, err := user.DoLogin()
    if err != nil {
        log.Print("User login failed.", err)
        r.JSON(http.StatusForbidden, Error("Login failed."))
        return
    }

    r.JSON(http.StatusOK, map[string]interface{} {
        "token": token,
    })
    return
}

func logoutUser(user *User, r render.Render) {
    if err := user.DoLogout(); err != nil {
        log.Print("User logout failed.", err)
        r.JSON(http.StatusForbidden, Error("Logout failed."))
        return
    }
}

// Check if the request carry logined token & id.
//
// Each request should carry `X-ID` and `X-TOKEN` in the header.
// If the token is valid, it will also map logined user (`*User`) into
// request handler.
func LoginRequired(req *http.Request, c martini.Context, r render.Render) {
    banAccess := func() {
        r.JSON(http.StatusUnauthorized, Error("Login required."))
    }

    login := req.Header.Get("X-LOGIN")
    token := req.Header.Get("X-TOKEN")

    var user *User
    user, err := GetUserByLogin(login)
    if err != nil {
        log.Print("Get user failed.", login, token, err)
        banAccess()
        return
    }
    ok, err := user.CheckLogin(token)
    if err != nil {
        log.Print("User check login failed.", login, token, err)
        banAccess()
        return
    }
    if !ok {
        log.Print("User login failed.", login, token)
        banAccess()
        return
    }

    // Map logined user into request context.
    c.Map(user)
}
