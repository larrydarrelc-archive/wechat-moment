package t

import (
    "log"
    "strconv"
    "net/http"
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
    "github.com/astaxie/beego/orm"
)

// Routes:
//
//  Path            Method  Description
//  /t              GET     User's main timeline.
//  /t              POST    Create a new status.
//  /t/:id          GET     Get a status.
//  /t/:id          DELETE  Delete a status.
//  /t/:id/like     PUT     Like a status
//  /t/:id/comment  GET     Get a status' comments
//  /t/:id/comment  POST    Create a comment for status `:id`
func TRoute(m *martini.ClassicMartini) {
    m.Get("/t", LoginRequired, getTimeline)
    m.Post("/t", LoginRequired, createTweet)

    m.Get("/t/:id", thisTweet, getTweet)
    m.Delete("/t/:id", LoginRequired, thisTweet, deleteTweet)

    m.Put("/t/:id/like", LoginRequired, thisTweet, likeTweet)

    m.Post("/t/:id/comment", LoginRequired, thisTweet, createTweetComment)
}

func getTimeline(r render.Render) {
    var tweet Tweet

    tweets, err := tweet.All()
    if err != nil {
        log.Print("Read all tweets failed.", err)
        r.JSON(http.StatusForbidden, Error("Read timeline failed."))
        return
    }
    r.JSON(http.StatusOK, tweets)
}

func getTweet(t *Tweet, r render.Render) {
    tweet, err := t.Censor()
    if err != nil {
        log.Print("Censor tweet failed.", t.Id, err)
        r.JSON(http.StatusNotFound, Error("Tweet not found."))
        return
    }

    r.JSON(http.StatusOK, tweet)
}

func createTweet(req *http.Request, u *User, r render.Render) {
    text := req.FormValue("text")
    if text == "" {
        log.Print("Tweet text validate failed.", u.Id)
        r.JSON(http.StatusForbidden, Error("Text cannot be empty."))
        return
    }

    // TODO Remove magic number.
    if len(text) > 5000 {
        log.Print("Tweet text validate failed.", u.Id)
        r.JSON(http.StatusForbidden, Error("Text too long."))
        return
    }

    o := orm.NewOrm()
    tweet := Tweet {UserId: u.Id, Text: text}
    _, err := o.Insert(tweet)
    if err != nil {
        log.Print("Tweet create failed.", u.Id, text, err)
        r.JSON(http.StatusForbidden, Error("Tweet create failed."))
        return
    }
    r.JSON(http.StatusCreated, tweet)
}

func deleteTweet(tweet *Tweet, u *User, r render.Render) {
    if tweet.UserId != u.Id {
        log.Print("Tweet user mismatch", tweet.Id, u.Id)
        r.JSON(http.StatusForbidden, Error("Delete tweet failed."))
        return
    }

    err := tweet.Delete()
    if err != nil {
        log.Print("Delete tweet failed.", tweet.Id, u.Id, err)
        r.JSON(http.StatusForbidden, Error("Delete tweet failed."))
        return
    }
    log.Print("Tweet was deleted.", tweet.Id, u.Id)
    r.JSON(http.StatusNoContent, "")
}

func likeTweet(tweet *Tweet, u *User, r render.Render) {
    like, err := tweet.IsLike(u)
    if err != nil {
        log.Print("Check is like tweet failed.", tweet.Id, u.Id, err)
        r.JSON(http.StatusForbidden, Error("Like tweet failed."))
        return
    }
    if like {
        if err = tweet.Like(u); err != nil {
            log.Print("Like tweet failed.", tweet.Id, u.Id, err)
            r.JSON(http.StatusForbidden, Error("Like tweet failed."))
            return
        }
        log.Print("Tweet was liked.", tweet.Id, u.Id)
    } else {
        if err = tweet.UnLike(u); err != nil {
            log.Print("Unlike tweet failed.", tweet.Id, u.Id, err)
            r.JSON(http.StatusForbidden, Error("Unlike tweet failed."))
            return
        }
        log.Print("Tweet was unliked.", tweet.Id, u.Id)
    }

    r.JSON(http.StatusNoContent, "")
}

func createTweetComment(req *http.Request,
                        tweet *Tweet,
                        u *User,
                        r render.Render) {
    content := req.FormValue("content")
    if content == "" {
        log.Print("Comment validation failed.", tweet.Id, u.Id)
        r.JSON(http.StatusForbidden, Error("Create comment failed."))
        return
    }
    // TODO Remove maginc number.
    if len(content) > 5000 {
        log.Print("Comment validation failed.", tweet.Id, u.Id)
        r.JSON(http.StatusForbidden, Error("Comment too long."))
        return
    }

    if err := tweet.CreateComment(content, u); err != nil {
        log.Print("Create comment failed.", tweet.Id, u.Id, err)
        r.JSON(http.StatusForbidden, Error("Create comment failed."))
        return
    }
    rv, err := tweet.Censor()
    if err != nil {
        log.Print("Censor tweet failed.", tweet.Id, u.Id, err)
        r.JSON(http.StatusForbidden, Error("Create comment failed."))
        return
    }

    r.JSON(http.StatusCreated, rv)
}

// Map requested tweet to context.
func thisTweet(params martini.Params, r render.Render, c martini.Context) {
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        log.Print("Cannot parse int `int`.", params["id"], err)
        r.JSON(http.StatusNotFound, Error("Tweet not found."))
        return
    }
    tweet, err := GetTweetById(id)
    if err != nil {
        log.Print("Read tweet failed.", id, err)
        r.JSON(http.StatusNotFound, Error("Tweet not found."))
        return
    }
    c.Map(tweet)
}
