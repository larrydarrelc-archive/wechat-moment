package t

import (
    "fmt"
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
//  /user/login     POST    Login a user.
//  /user/logout    GET     Logout a user.
//  /user/:id       GET     Get user `/:id`'s profile.
//  /user/:id       PUT     Update user's profile.
func UserRoute(m *martini.ClassicMartini) {
    m.Post("/user", createUser)
    m.Post("/user/login", loginUser)
    m.Get("/user/logout", LoginRequired, logoutUser)

    m.Get("/user/:id", LoginRequired, getUserProfile)
    m.Put("/user/:id", LoginRequired, updateUserProfile)
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
        panic(err)
    }
    if !created {
        r.JSON(http.StatusConflict,
               Error(fmt.Sprintf("Login name %s already exists!", login)))
        return
    }

    user.Password = user.HashPassword(password)
    user.Name = name
    if _, err := o.Update(&user); err != nil {
        panic(err)
    }
    r.JSON(http.StatusCreated, map[string]interface{} {
        "Id": user.Id,
        "Login": user.Login,
    })
}

func getUserProfile(params martini.Params, r render.Render) {
    o := orm.NewOrm()

    id, err := strconv.Atoi(params["id"])
    if err != nil {
        panic(err)
    }
    user := User{Id: id}
    err = o.Read(&user)
    if err == orm.ErrNoRows {
        r.JSON(http.StatusNotFound,
               Error(fmt.Sprintf("User %d not exists", id)))
        return
    }
    r.JSON(http.StatusOK, user.Censor())
}

func updateUserProfile(req *http.Request, params martini.Params, r render.Render) {
    o := orm.NewOrm()

    id, err := strconv.Atoi(params["id"])
    if err != nil {
        r.JSON(http.StatusNotFound,
               Error(fmt.Sprintf("User %s not exists", params["id"])))
        return
    }
    user := User{Id: id}
    err = o.Read(&user)
    if err == orm.ErrNoRows {
        r.JSON(http.StatusNotFound,
               Error(fmt.Sprintf("User %d not exists", id)))
        return
    }

    name := req.FormValue("name")
    user.Name = name
    _, err = o.Update(&user)
    if err != nil {
        panic(err)
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
        panic(err)
    }

    token, err := user.DoLogin()
    if err != nil {
        panic(err)
    }

    r.JSON(http.StatusOK, map[string]interface{} {
        "token": token,
    })
    return
}

func logoutUser(user *User) {
    if err := user.DoLogout(); err != nil {
        panic(err)
    }
}
