BINS=cert-cacher.arm64.osx cert-cacher.amd64.linux cert-cacher.amd64.windows

test:
	go test -v *.go

test/watch:
	@ls *.go | entr -c -s 'go test -failfast -v ./*.go && notify "💚" || notify "🛑"'

coverage/html:
	go test -v -cover -coverprofile=c.out
	go tool cover -html=c.out

build: $(BINS)
$(BINS):
	go build -o $@

clean:
	rm -f c.out $(BINS) cert-cacher

.PHONY: module
module:
	go mod init github.com/drio/cert-cacher
	go mod tidy	
