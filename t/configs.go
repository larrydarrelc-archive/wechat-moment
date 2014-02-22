package t

import (
    "os"
    "log"
    "path/filepath"
    "encoding/json"
)

type Configs struct {
    Db string
    Port int
    Host string
    Static struct {
        Directory string
        Prefix string
    }
}

func ReadConfigs(path string) (configs *Configs) {
    absPath, _ := filepath.Abs(path)
    file, err := os.Open(absPath)
    if err != nil {
        log.Fatal("Configuration file not exists.", absPath, err)
        panic(err)
    }

    configs = &Configs{}
    err = json.NewDecoder(file).Decode(&configs)
    if err != nil {
        log.Fatal("Configuration file parse failed.", absPath, err)
        panic(err)
    }

    return
}
