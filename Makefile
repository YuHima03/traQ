SOURCES ?= $(shell find . -path "./vendor" -prune -o -type f -name "*.go" -print)

traQ: $(SOURCES)
	go build -ldflags "-X main.version=$$(git describe --tags --abbrev=0) -X main.revision=$$(git rev-parse --short HEAD)"

.PHONY: init
init:
	go mod download

.PHONY: genkey
genkey:
	mkdir -p ./dev/keys
	cd ./dev/keys && go run ../bin/gen_ec_pem.go

.PHONY: up-docker-test-db
up-docker-test-db:
	docker run --name traq-test-db -p 3100:3306 -e MYSQL_ROOT_PASSWORD=password -d mariadb:10.0.19 mysqld --character-set-server=utf8 --collation-server=utf8_general_ci
	sleep 5
	TEST_DB_PORT=3100 go run .circleci/init.go

.PHONY: down-docker-test-db
down-docker-test-db:
	docker rm -f -v traq-test-db

.PHONY: make-db-docs
make-db-docs:
	if [ -d "./docs/dbschema" ]; then \
		rm -r ./docs/dbschema; \
	fi
	TBLS_DSN="mysql://root:password@127.0.0.1:3002/traq" tbls doc

.PHONY: diff-db-docs
diff-db-docs:
	TBLS_DSN="mysql://root:password@127.0.0.1:3002/traq" tbls diff
