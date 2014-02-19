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

func prepareMiddlewares(m *martini.ClassicMartini, config *Configs) {
    m.Use(render.Renderer())
    m.Use(martini.Recovery())
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

    prepareDatabase(m, config)
    prepareModels(m, config)
    prepareMiddlewares(m, config)
    prepareViews(m, config)

    return m
}
