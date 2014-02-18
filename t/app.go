package t

import (
    "github.com/codegangsta/martini"
)

func Build() (*martini.ClassicMartini) {
    m := martini.Classic()

    m.Get("/", Nop)
    m.Get("/poll", Nop)

    UserRoute(m)
    TRoute(m)

    return m
}
