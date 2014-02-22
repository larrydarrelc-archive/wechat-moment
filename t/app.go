package t

import (
    "path/filepath"
    "github.com/astaxie/beego/orm"
    _ "github.com/mattn/go-sqlite3"
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
)

type Application struct {
    *martini.Martini
    martini.Router
    config *Configs
}

func (app *Application) prepareDatabase() {
    dbPath := app.config.Db
    if dbPath == "" {
        dbPath = ":memory:"
    } else {
        filepath.Abs(dbPath)
    }

    orm.RegisterDataBase("default", "sqlite3", dbPath)
}

func (app *Application) prepareMiddlewares() {
    if app.config.Static.Directory != "" {
        app.Use(martini.Static(
            app.config.Static.Directory,
            martini.StaticOptions{
                Prefix: app.config.Static.Prefix,
                SkipLogging: true,
            }),
        )
    }
    app.Use(martini.Logger())
    app.Use(render.Renderer())
    app.Use(martini.Recovery())
}

func (app *Application) prepareViews() {
    app.Get("/", Nop)

    app.Get("/poll", Nop)

    UserRoute(app)
    TRoute(app)
}

func (app *Application) prepareModels() {
    orm.RegisterModel(new(User))
    orm.RegisterModel(new(Tweet))
    orm.RegisterModel(new(TweetComment))
}

func Build(config *Configs) (*Application) {
    m := martini.New()
    r := martini.NewRouter()
    m.Action(r.Handle)
    app := &Application{m, r, config}

    app.prepareDatabase()
    app.prepareModels()
    app.prepareMiddlewares()
    app.prepareViews()

    return app
}
