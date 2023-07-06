default:
    @just --list

# build students-fall-2022 binary
build:
    go build -o students

# update go packages
update:
    go get -u

# run tests
test:
    go test -v ./... -covermode=atomic -coverprofile=coverage.out

# connect into the database file using sqlite
database:
    sqlite3 students.db

# run golangci-lint
lint:
    golangci-lint run -c .golangci.yml
