package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/larrydarrelc/t"
)

func main() {
    config := t.ReadConfigs("configs/default.json")
    m := t.Build(config)

    dest := fmt.Sprintf("%s:%d", config.Host, config.Port)
    log.Print(fmt.Sprintf("Start listening on %s", dest))
    err := http.ListenAndServe(dest, m)
    if err != nil {
        panic(err)
    }
}
