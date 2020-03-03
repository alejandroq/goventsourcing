.PHONY: start
start:
	go run cmd/main.go

.PHONY: test
test:
	go test ./...

.PHONY: dependencytree
dependencytree:
	godepgraph -p github.com/google -s cmd/main.go | dot -Tpng -o .github/dependencytree.png
