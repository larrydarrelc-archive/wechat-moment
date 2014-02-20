.PHONY: setup_test_db setup_test_server run_test

run_test:
	python2 tests.py

setup_test_server: setup_test_db
	go run tests.go

setup_test_db:
	rm -f data/data.test.sqlite3
	sqlite3 data/data.test.sqlite3 < migrations/1.schema.sql
	sqlite3 data/data.test.sqlite3 < migrations/2.login_name.sql
