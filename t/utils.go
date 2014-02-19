package t

import (
    "crypto/sha1"
    "fmt"
    "io"
)


type Delim string

func (delim Delim) Join(components []string) (joined string) {
    joined = ""
    delim_ := string(delim)
    last := len(components) - 1

    for i := range(components) {
        joined = joined + components[i]
        if i != last {
            joined = joined + delim_
        }
    }

    return
}

// A nop response function.
func Nop() (int, string) {
    return 200, "Hello, world!"
}

func Error(reason string) (map[string]string) {
    return map[string]string {
        "error": reason,
    }
}

func Success(reason string) (map[string]string) {
    return map[string]string {
        "message": reason,
    }
}

func Hash(raw string) (string) {
    h := sha1.New()
    io.WriteString(h, raw)
    return fmt.Sprintf("%x", h.Sum(nil))
}
