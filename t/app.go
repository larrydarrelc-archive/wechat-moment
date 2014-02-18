package t

import (
    "github.com/codegangsta/martini"
)

func hello() (int, string) {
    return 200, "Hello, world!"
}

func Build() (*martini.ClassicMartini) {
    m := martini.Classic()

    m.Get("/", hello)

    return m
}
