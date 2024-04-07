FLAGS := -ldflags "-s -w" -trimpath

build::
	go vet && go fmt
	go build ${FLAGS}

run:: build
	./detect -device="Logitech StreamCam: USB Audio (hw:1,0)"

watch::
	git ls-files | entr -rc make run