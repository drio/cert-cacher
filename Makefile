PRJ=cert-cacher
BINS=$(PRJ).linux.amd64 $(PRJ).darwin.arm64 $(PRJ).windows.amd64

test:
	go test -v *.go

test/watch:
	@ls *.go | entr -c -s 'go test -failfast -v ./*.go && notify "💚" || notify "🛑"'

coverage/html:
	go test -v -cover -coverprofile=c.out
	go tool cover -html=c.out

build: $(BINS)

cert-cacher.linux.amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $@ .

cert-cacher.darwin.arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $@ .

cert-cacher.windows.amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $@ .

clean:
	rm -f c.out $(BINS) cert-cacher $(BINS)

.PHONY: module
module:
	go mod init github.com/drio/cert-cacher
	go mod tidy	
