compile:
	git stash -u
	gox -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}"

fmt:
	find ! -path "./vendor/*" -name "*.go" -exec go fmt {} \;

gometalinter:
	gometalinter -D gotype --vendor --deadline=240s --dupl-threshold=200 -e '_string' -j 5 ./...

run-tests:
	./test.sh

run-quick-tests:
	go test -v $(shell glide nv)

test-all: gometalinter run-tests

test-package:
	go test -race -cover -coverprofile=/tmp/chyle github.com/antham/chyle/$(pkg)
	go tool cover -html=/tmp/chyle
