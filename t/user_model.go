package t

import (
    "fmt"
    "time"
    "github.com/astaxie/beego/orm"
)

type User struct {
    Id int `orm:"auto"`
    Name string `orm:"size(32)"`
    Login string `orm:"size(32)"`
    Password string `orm:"size(64)"`
    Avatar string `orm:"size(256)"`
    CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
    UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func (u *User) GenerateToken() (string) {
    return Hash(fmt.Sprintf("%s%d%s", time.Now().String(), u.Id, u.Login))
}

func (*User) HashPassword(raw string) (string) {
    return Hash(raw)
}

func (u *User) DoLogin() (string, error) {
    o := orm.NewOrm()

    token := u.GenerateToken()
    stat, err := o.Raw("INSERT INTO `token` (`code`, `user_id`) VALUES (?, ?)").Prepare()
    if err != nil {
        return "", err
    }
    defer stat.Close()
    _, err = stat.Exec(token, u.Id)
    if err != nil {
        return "", err
    }
    return token, nil
}

func (u *User) DoLogout() (error) {
    o := orm.NewOrm()

    stat, err := o.Raw("DELETE FROM `token` WHERE `user_id` = ?").Prepare()
    if err != nil {
        return err
    }
    defer stat.Close()
    _, err = stat.Exec(u.Id)
    if err != nil {
        return err
    }
    return nil
}

func (u *User) CheckLogin(token string) (bool, error) {
    var code string

    o := orm.NewOrm()

    stat := o.Raw("SELECT `code` FROM `token` WHERE `code` = ? AND `user_id` = ?",
                  token, u.Id)
    err := stat.QueryRow(&code)
    if err != nil {
        return false, err
    }

    return code == token, nil
}
