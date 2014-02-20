package t

import (
    "time"
    "github.com/astaxie/beego/orm"
)

type Tweet struct {
    Id int `orm:"auto"`
    UserId int
    Text string `orm:"size(5000)"`
    Image string `orm:"size(256)"`
    CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
}

func (t *Tweet) Censor() (TypeModel, error) {
    user, err := t.GetUser()
    if err != nil {
        return nil, err
    }
    likes, err := t.GetLike()
    if err != nil {
        return nil, err
    }
    c, err := GetCommentsByTweetId(t.Id)
    if err != nil {
        return nil, err
    }
    var comments []TypeModel
    for i := range c {
        t, err := c[i].Censor()
        if err != nil {
            return nil, err
        }
        comments = append(comments, t)
    }

    return TypeModel {
        "Id": t.Id,
        "Text": t.Text,
        "Image": t.Image,
        "CreatedAt": t.CreatedAt,
        "User": user,
        "Likes": likes,
        "Comments": comments,
    }, nil
}

func (*Tweet) All() (rv []TypeModel, err error) {
    o := orm.NewOrm()

    var tweets []Tweet
    stat := o.Raw("SELECT * FROM `t` ORDER BY `t`.`created_at` DESC")
    _, err = stat.QueryRows(&tweets)
    if err != nil {
        return nil, err
    }
    for i := range tweets {
        r, err := tweets[i].Censor()
        if err != nil {
            return nil, err
        }
        rv = append(rv, r)
    }

    return rv, nil
}

func (t *Tweet) GetUser() (TypeModel, error) {
    user, err := GetUserById(t.UserId)
    if err != nil {
        return nil, err
    }

    rv, err := user.Censor()
    if err != nil {
        return nil, err
    }

    return rv, nil
}

func (t *Tweet) GetLike() (rv []TypeModel, err error) {
    o := orm.NewOrm()

    var users_id []int
    var like_times []time.Time

    stat := o.Raw("SELECT `created_at`, `user_id` FROM `t_like` WHERE `t_id` = ?", t.Id)
    _, err = stat.QueryRows(&like_times, &users_id)
    if err != nil {
        return nil, err
    }

    for i := range users_id {
        user, err := GetUserById(users_id[i])
        if err != nil {
            return nil, err
        }
        r, err := user.Censor()
        if err != nil {
            return nil, err
        }
        rv = append(rv, TypeModel {
            "User": r,
            "CreatedAt": like_times[i],
        })
    }

    return rv, nil
}

func (t *Tweet) Delete() (err error) {
    o := orm.NewOrm()

    stat := o.Raw("DELETE FROM `t_comment` WHERE `t_id` = ?", t.Id)
    _, err = stat.Exec()
    if err != nil {
        return err
    }

    stat = o.Raw("DELETE FROM `t_like` WHERE `t_id` = ?", t.Id)
    _, err = stat.Exec()
    if err != nil {
        return err
    }

    stat = o.Raw("DELETE FROM `t` WHERE `id` = ?", t.Id)
    _, err = stat.Exec()
    if err != nil {
        return err
    }

    return nil
}

func (t *Tweet) IsLike(u *User) (bool, error) {
    o := orm.NewOrm()

    stat := o.Raw(
        "SELECT COUNT(*) FROM `t_like` WHERE `t_id` = ? AND `user_id` = ?",
        t.Id, u.Id,
    )
    var count int
    err := stat.QueryRow(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func (t *Tweet) Like(u *User) (error) {
    o := orm.NewOrm()

    stat := o.Raw(
        "INSERT INTO `t_like` (`t_id`, `user_id`) VALUES (?, ?)",
        t.Id, u.Id,
    )
    _, err := stat.Exec()
    return err
}

func (t *Tweet) UnLike(u *User) (error) {
    o := orm.NewOrm()

    stat := o.Raw(
        "DELETE FROM `t_like` WHERE `t_id` = ? AND `user_id` = ?",
        t.Id, u.Id,
    )
    _, err := stat.Exec()
    return err
}

func (t *Tweet) CreateComment(content string, user *User) (error) {
    o := orm.NewOrm()

    stat := o.Raw(
        "INSERT INTO `t_comment` (`t_id`, `user_id`, `content`) VALUES (?, ?, ?)",
        t.Id, user.Id, content,
    )
    _, err := stat.Exec()
    return err
}

func GetTweetById(id int) (t *Tweet, err error) {
    o := orm.NewOrm()

    stat := o.Raw("SELECT * FROM `t` WHERE `id` = ?", id)
    err = stat.QueryRow(&t)
    if err != nil {
        return nil, err
    }

    return t, nil
}

type TweetComment struct {
    Id int `orm:"auto"`
    TweetId int
    UserId int
    Content string `orm:"size(5000)"`
    CreatedAt time.Time `orm:"auto_now_add";type(datetime)`
}

func (c *TweetComment) Censor() (TypeModel, error) {
    user, err := GetUserById(c.UserId)
    if err != nil {
        return nil, err
    }

    return TypeModel {
        "Id": c.Id,
        "TweetId": c.TweetId,
        "User": user,
        "Content": c.Content,
        "CreatedAt": c.CreatedAt,
    }, nil
}

func GetCommentsByTweetId(tId int) (rv []TweetComment, err error) {
    o := orm.NewOrm()

    stat := o.Raw(
        "SELECT * FROM `t_comment` WHERE `t_id` = ? ORDER BY `created_at` DESC",
        tId,
    )
    _, err = stat.QueryRows(&rv)
    if err != nil {
        return nil, err
    }
    return rv, nil
}
