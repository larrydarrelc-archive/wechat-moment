package main

import (
    "net/http"
    "fmt"
    "github.com/larrydarrelc/t"
)

func main() {
    config := t.ReadConfigs("configs/default.json")
    m := t.Build(config)

    dest := fmt.Sprintf("%s:%d", config.Host, config.Port)
    err := http.ListenAndServe(dest, m)
    if err != nil {
        panic(err)
    }
}
