package t

import (
    "os"
    "path/filepath"
    "encoding/json"
)

type Configs struct {
    Db string
    Port int
    Host string
}

func ReadConfigs(path string) (configs *Configs) {
    absPath, _ := filepath.Abs(path)
    file, err := os.Open(absPath)
    if err != nil {
        //panic("Load configuration failed.")
        panic(err)
    }

    configs = &Configs{}
    err = json.NewDecoder(file).Decode(&configs)
    if err != nil {
        //panic("Load configuration failed.")
        panic(err)
    }

    return
}
