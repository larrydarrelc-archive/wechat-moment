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
    // Access path prefix.
    path string
    // Stored file name prefix.
    prefix string
}

// Map a uploader service into the Martini handler chain.
func UploadProvider(dest, path, prefix string) martini.Handler {
    uploader := NewUploader(dest, path, prefix)
    return func(c martini.Context) {
        c.Map(uploader)
    }
}

func NewUploader(dest, path, prefix string) (*Uploader) {
    return &Uploader{dest, path, prefix}
}

// Store a uploaded file.
func (u *Uploader) Store(upload multipart.File) (string, error) {
    name := u.newName()
    dest, err := os.Create(u.newStorePath(name))
    if err != nil {
        return "", err
    }
    defer dest.Close()

    if _, err = io.Copy(dest, upload); err != nil {
        return "", err
    }
    return u.newAccessPath(name), nil
}

// Generate a new file path.
func (u *Uploader) newStorePath(name string) (string) {
    // FIXME magic path delimiter
    return fmt.Sprintf("%s/%s", u.dest, name)
}

func (u *Uploader) newAccessPath(name string) (string) {
    return fmt.Sprintf("%s/%s", u.path, name)
}

// Generate a new file name.
func (u *Uploader) newName() (string) {
    return fmt.Sprintf("%s%s", u.prefix, generateName())
}

// Simply generate a random name.
// Don't take collison into consideration.
func generateName() (string) {
    rand.Seed(int64(time.Now().Second()))
    return Hash(fmt.Sprintf("%s%d", time.Now().String(), rand.Int()))
}
