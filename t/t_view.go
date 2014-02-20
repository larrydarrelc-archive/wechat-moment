package t

import (
    "log"
    "strconv"
    "net/http"
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
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
    m.Get("/t", getTimeline)
    m.Post("/t", Nop)

    m.Get("/t/:id", thisTweet, getTweet)
    m.Delete("/t/:id", LoginRequired, thisTweet, Nop)

    m.Put("/t/:id/like", Nop)

    m.Post("/t/:id/comment", Nop)
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
