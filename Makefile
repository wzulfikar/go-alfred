build-linux:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w' -o ./bin/luxtag-knowledge .

push:
	scp ./bin/luxtag-knowledge luxtag-internal:~/bin/

