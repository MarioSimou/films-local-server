.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getSongs internal/functions/getSongs/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getSong internal/functions/getSong/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/postSong internal/functions/postSong/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/putSong internal/functions/putSong/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/deleteSong internal/functions/deleteSong/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/uploadSong internal/functions/uploadSong/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/putSongImage internal/functions/putSongImage/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

watcher:
	reflex -r '\.go$$' -s -- make build

local:
	sls offline start --dockerNetwork songs-local-server-sls-backend-network --useDocker