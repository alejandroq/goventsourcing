.PHONY: start
start:
	go run cmd/main.go

.PHONY: dependencytree
dependencytree:
	godepgraph -s cmd/main.go | dot -Tpng -o dependencytree.png
