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
    o := orm.NewOrm()

    stat := o.Raw("SELECT `code` FROM `token` WHERE `code` = ? AND `user_id` = ?",
                  token, u.Id)
    var code string
    err := stat.QueryRow(&code)
    if err != nil {
        return false, err
    }

    return code == token, nil
}

func (u *User) SetAvatar(path string) (error) {
    o := orm.NewOrm()

    stat := o.Raw("UPDATE `user` SET `avatar` = ? WHERE `id` = ?", path, u.Id)
    _, err := stat.Exec()
    return err
}

// Hide some secret field.
func (u User) Censor() (TypeModel, error) {
    return TypeModel {
        "Id": u.Id,
        "Name": u.Name,
        "Avatar": u.Avatar,
        "CreatedAt": u.CreatedAt,
        "UpdatedAt": u.UpdatedAt,
    }, nil
}

func (u *User) UpdateProfile(name string) (error) {
    o := orm.NewOrm()

    stat := o.Raw("UPDATE `user` SET `name` = ? WHERE `id` = ?", name, u.Id)
    _, err := stat.Exec()
    return err
}

func GetUserById(id int) (user *User, err error) {
    o := orm.NewOrm()

    stat := o.Raw("SELECT * FROM `user` WHERE `id` = ?", id)
    err = stat.QueryRow(&user)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func GetUserByLogin(login string) (user *User, err error) {
    o := orm.NewOrm()

    stat := o.Raw("SELECT * FROM `user` WHERE `login` = ?", login)
    err = stat.QueryRow(&user)
    if err != nil {
        return nil, err
    }
    return user, nil
}
