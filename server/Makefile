FLAGS := -ldflags "-s -w" -trimpath
NOCGO := CGO_ENABLED=0

build::
	go vet && go fmt
	${NOCGO} go build ${FLAGS}

run::
	go run main.go

watch::
	git ls-files | entr -c make run