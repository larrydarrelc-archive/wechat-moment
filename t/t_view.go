package t

import (
    "log"
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

    m.Get("/t/:id", Nop)
    m.Delete("/t/:id", Nop)

    m.Put("/t/:id/like", Nop)

    m.Get("/t/:id/comment", Nop)
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
