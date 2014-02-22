package t

import (
    "os"
    "io"
    "fmt"
    "time"
    "math/rand"
    "mime/multipart"
    "github.com/codegangsta/martini"
)

// A upload files helper.
type Uploader struct {
    // Store destination.
    dest string
    // Stored file name prefix.
    prefix string
}

// Map a uploader service into the Martini handler chain.
func UploadProvider(dest, prefix string) martini.Handler {
    uploader := NewUploader(dest, prefix)
    return func(c martini.Context) {
        c.Map(uploader)
    }
}

func NewUploader(dest, prefix string) (*Uploader) {
    return &Uploader{dest, prefix}
}

// Store a uploaded file.
func (u *Uploader) Store(upload multipart.File) (error) {
    dest, err := os.Create(u.NewPath())
    if err != nil {
        return err
    }
    defer dest.Close()

    if _, err = io.Copy(dest, upload); err != nil {
        return err
    }
    return nil
}

// Generate a new file path.
func (u *Uploader) NewPath() (string) {
    // FIXME magic path delimiter
    return fmt.Sprintf("%s/%s%s", u.dest, u.prefix, generateName())
}

// Simply generate a random name.
// Don't take collison into consideration.
func generateName() (string) {
    rand.Seed(int64(time.Now().Second()))
    return Hash(fmt.Sprintf("%s%d", time.Now().String(), rand.Int()))
}
