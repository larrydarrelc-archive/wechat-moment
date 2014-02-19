package t

import (
    "path/filepath"
    "github.com/astaxie/beego/orm"
    _ "github.com/mattn/go-sqlite3"
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
)

func prepareDatabase(m *martini.ClassicMartini, config *Configs) {
    dbPath := config.Db
    if dbPath == "" {
        dbPath = ":memory:"
    } else {
        filepath.Abs(dbPath)
    }

    orm.RegisterDataBase("default", "sqlite3", dbPath)
}

func prepareViews(m *martini.ClassicMartini, config *Configs) {
    m.Get("/", Nop)

    m.Get("/poll", Nop)

    UserRoute(m)
    TRoute(m)
}

func prepareModels(m *martini.ClassicMartini, config *Configs) {
    orm.RegisterModel(new(User))
}

func Build(config *Configs) (*martini.ClassicMartini) {
    m := martini.Classic()

    m.Use(render.Renderer())
    m.Use(martini.Recovery())
    prepareViews(m, config)
    prepareDatabase(m, config)
    prepareModels(m, config)

    return m
}
