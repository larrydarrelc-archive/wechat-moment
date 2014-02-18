package t

import (
    "github.com/codegangsta/martini"
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
    m.Post("/user", Nop)
    m.Post("/user/login", Nop)
    m.Get("/user/logout", Nop)

    m.Get("/user/:id", Nop)
    m.Put("/user/:id", Nop)
}
