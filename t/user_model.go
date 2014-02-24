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

func (u *User) UpdateProfile(name string) (error) {
    o := orm.NewOrm()

    stat := o.Raw("UPDATE `user` SET `name` = ? WHERE `id` = ?", name, u.Id)
    _, err := stat.Exec()
    return err
}

func (u *User) GetTweets() (rv []TypeModel, err error) {
    o := orm.NewOrm()

    var ids []int
    stat := o.Raw("SELECT `id` FROM `t` WHERE `user_id` = ?", u.Id)
    _, err = stat.QueryRows(&ids)
    if err != nil {
        if err == orm.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    for i := range ids {
        tweet, err := GetTweetById(ids[i])
        if err != nil {
            return nil, err
        }
        // XXX Don't call `GetTweets` in tweet.Censor.
        censored, err := tweet.Censor()
        if err != nil {
            return nil, err
        }
        rv = append(rv, censored)
    }
    if rv == nil {
        rv = []TypeModel{}
    }

    return rv, nil
}

func (u *User) GetTimeline() (rv TypeModel, err error) {
    o := orm.NewOrm()

    authorIds, err := u.GetFriendIds()
    if err != nil {
        return nil, err
    }

    // Include user himself.
    authorIds = append(authorIds, u.Id)

    var (
        tweets []Tweet
        timeline []TypeModel
    )
    stat := o.Raw("SELECT * FROM `t` WHERE `user_id` IN ?", authorIds)
    _, err = stat.QueryRows(&tweets)

    for i := range tweets {
        censored, err := tweets[i].Censor()
        if err != nil {
            return nil, err
        }
        timeline = append(timeline, censored)
    }
    if timeline == nil {
        timeline = []TypeModel{}
    }

    return TypeModel{"t": timeline}, nil
}

func (u *User) AddFriend(friend *User) (error) {
    hasFriend, err := u.HasFriend(friend)
    if err != nil || hasFriend {
        return err
    }

    o := orm.NewOrm()
    stat := o.Raw("INSERT INTO `user_friend` (`user_a_id`, `user_b_id`) VALUES (?, ?)", u.Id, friend.Id)
    _, err = stat.Exec()
    if err != nil {
        return err;
    }

    return friend.AddFriend(u)
}

func (u *User) RemoveFriend(friend *User) (error) {
    o := orm.NewOrm()

    stat := o.Raw("DELETE FROM `user_friend` WHERE `user_a_id` = ? AND `user_b_id` = ?", u.Id, friend.Id)
    _, err := stat.Exec()
    if err != nil {
        return err
    }

    stat = o.Raw("DELETE FROM `user_friend` WHERE `user_a_id` = ? AND `user_b_id` = ?", friend.Id, u.Id)
    _, err = stat.Exec()

    return err
}

func (u *User) HasFriend(friend *User) (bool, error) {
    o := orm.NewOrm()

    var count int
    stat := o.Raw("SELECT COUNT(*) FROM `user_friend` WHERE `user_a_id` = ? AND `user_b_id` = ?", u.Id, friend.Id)
    err := stat.QueryRow(&count)
    if err != nil {
        return false, err
    }

    return count > 0, nil
}

func (u *User) GetFriendIds() (rv []int, err error) {
    o := orm.NewOrm()

    stat := o.Raw("SELECT `user_b_id` FROM `user_friend` WHERE `user_a_id` = ?", u.Id)
    _, err = stat.QueryRows(&rv)
    if err != nil && err != orm.ErrNoRows {
        if err == orm.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    return
}

func (u *User) GetFriends() (rv []TypeModel, err error) {
    friendIds, err := u.GetFriendIds()
    if err != nil {
        return nil, err
    }

    for i := range friendIds {
        friend, err := GetUserById(friendIds[i])
        if err != nil {
            return nil, err
        }
        censored, err := friend.Censor()
        if err != nil {
            return nil, err
        }
        rv = append(rv, censored)
    }

    return
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
