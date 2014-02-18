package main

import (
    "github.com/codegangsta/martini"
)

func hello() (int, string) {
    return 404, "Hello, world!"
}

func main() {
    m := martini.Classic()

    m.Get("/", hello)

    m.Run()
}
