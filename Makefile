test:
	go test -v *.go

test/watch:
	@ls *.go | entr -c -s 'go test -failfast -v ./*.go && notify "💚" || notify "🛑"'

coverage/html:
	go test -v -cover -coverprofile=c.out
	go tool cover -html=c.out

clean:
	rm -f c.out
