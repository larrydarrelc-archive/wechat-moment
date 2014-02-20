package main

import (
    "fmt"
    "time"
    "net/http"
    "net/url"
    "encoding/json"
    "github.com/larrydarrelc/t"
)

type Config struct {
    Url string
    token string
}

func StartServer(config *t.Configs) {
    m := t.Build(config)

    dest := fmt.Sprintf("%s:%d", config.Host, config.Port)
    err := http.ListenAndServe(dest, m)
    if err != nil {
        panic(err)
    }
}

func GetToken(config *Config) {
    type ResponseMessage struct {
        Token string
    }

    dest := fmt.Sprintf("%s/user/login", config.Url)
    var message ResponseMessage

    payload := url.Values{
        "login": {"testlogin"},
        "password": {"testpassword"},
    }
    resp, err := http.PostForm(dest, payload)
    if err != nil {
        panic(err)
    }
    dec := json.NewDecoder(resp.Body)
    err = dec.Decode(&message)
    if err != nil {
        panic(err)
    }
    fmt.Println(message)
    config.token = message.Token
}

func TestCreateUser(config Config) {
    type ResponseMessage struct {
        Id int
        Login string
    }

    dest := fmt.Sprintf("%s/user", config.Url)
    var message ResponseMessage

    payload := url.Values{
        "name": {"testname"},
        "login": {"testlogin"},
        "password": {"testpassword"},
    }
    resp, err := http.PostForm(dest, payload)
    if err != nil {
        panic(err)
    }
    dec := json.NewDecoder(resp.Body)
    err = dec.Decode(&message)
    if err != nil {
        panic(err)
    }
    fmt.Println(message)
    assert(resp.StatusCode == 201, "Create user")
    assert(message.Id == 1, "Create user")
    assert(message.Login == "testlogin", "Create user")
    resp.Body.Close()

    resp, err = http.PostForm(dest, payload)
    if err != nil {
        panic(err)
    }
    dec = json.NewDecoder(resp.Body)
    err = dec.Decode(&message)
    if err != nil {
        panic(err)
    }
    assert(resp.StatusCode == 409, "Duplicate user")

    payload.Set("login", "testlogin2")
    payload.Set("name", "")
    resp, err = http.PostForm(dest, payload)
    if err != nil {
        panic(err)
    }
    dec = json.NewDecoder(resp.Body)
    err = dec.Decode(&message)
    if err != nil {
        panic(err)
    }
    assert(resp.StatusCode == 403, "Empty name")

    payload.Set("login", "testlogin2")
    payload.Set("name", "testname2")
    payload.Set("password", "1234")
    resp, err = http.PostForm(dest, payload)
    if err != nil {
        panic(err)
    }
    dec = json.NewDecoder(resp.Body)
    err = dec.Decode(&message)
    if err != nil {
        panic(err)
    }
    assert(resp.StatusCode == 403, "Short password")
}

func TestGetUserProfile(config Config) {
    type ResponseMessage struct {
        Id int
        Login string
        Name string
        Password string
    }

    var message ResponseMessage
    dest := fmt.Sprintf("%s/user/%d", config.Url, 1)

    resp, err := http.Get(dest)
    if err != nil {
        panic(err)
    }
    assert(resp.StatusCode == 200, "Get user profile")
    dec := json.NewDecoder(resp.Body)
    err = dec.Decode(&message)
    if err != nil {
        panic(err)
    }
    assert(message.Id == 1, "Get user profile")
    assert(message.Name == "testname", "Get user profile")
    assert(message.Login == "testlogin", "Get user profile")
    assert(message.Password == "", "Get user profile")

    dest = fmt.Sprintf("%s/user/%d", config.Url, 404)
    resp, err = http.Get(dest)
    if err != nil {
        panic(err)
    }
    assert(resp.StatusCode == 404, "Get non exist user profile")
}

func TestUpdateUserProfile(config Config) {
    type ResponseMessage struct {
        Id int
        Login string
        Name string
        Password string
    }

    var message ResponseMessage
    dest := fmt.Sprintf("%s/user/%d", config.Url, 1)

    payload := url.Values{
        "name": {"newname"},
    }
    resp, err := http.PostForm(dest, payload)
    if err != nil {
        panic(err)
    }
    dec := json.NewDecoder(resp.Body)
    err = dec.Decode(&message)
    if err != nil {
        panic(err)
    }
    fmt.Println(message)
    assert(resp.StatusCode == 204, "Update user profile")
    assert(message.Name == "newname", "Update user profile")
    resp.Body.Close()
}

func TestUserLogin(config Config) {
    dest := fmt.Sprintf("%s/user/login", config.Url)

    payload := url.Values{
        "login": {"testlogin"},
        "password": {"testpassword"},
    }
    resp, err := http.PostForm(dest, payload)
    if err != nil {
        panic(err)
    }
    assert(resp.StatusCode == 200, "Test user login")

    payload.Set("password", "not_my_password")
    resp, err = http.PostForm(dest, payload)
    if err != nil {
        panic(err)
    }
    assert(resp.StatusCode == 403, "Test user login fail")
}

func main() {
    serverConfig := t.ReadConfigs("configs/testing.json")
    StartServer(serverConfig)
    config := Config {
        Url: fmt.Sprintf("http://%s:%d", serverConfig.Host, serverConfig.Port),
    }
    time.Sleep(200 * time.Millisecond)
    TestCreateUser(config)
    GetToken(&config)
    TestGetUserProfile(config)
    //TestUpdateUserProfile(config)
    TestUserLogin(config)
}

func assert(ass bool, hint string) {
    if ass {
        fmt.Println("Passed")
    } else {
        fmt.Println("Not passed", hint)
    }
}
