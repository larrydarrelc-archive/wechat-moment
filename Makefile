.PHONY: setup_test_db setup_test_server run_test deploy install

run_test:
	python2 tests.py

setup_test_server: setup_test_db
	go run tests.go

setup_test_db:
	rm -f data/data.test.sqlite3
	sqlite3 data/data.test.sqlite3 < migrations/1.schema.sql
	sqlite3 data/data.test.sqlite3 < migrations/2.login_name.sql

deploy:
	mv server server-old
	mv server-new server

install:
	export GOPATH=`pwd`
	rm -rf pkg
	rm -rf src/github.com/{astaxie,mattn,codegangsta}
	go get github.com/astaxie/beego/orm
	go get github.com/mattn/go-sqlite3
	go get github.com/codegangsta/martini
	go get github.com/codegangsta/martini-contrib/render
	go build server.go
