linux: # not tested in a real Linux environment
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o game ./cmd/game

mac: # not tested in a real macOS environment
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o game ./cmd/game

windows:
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags "-H=windowsgui -s -w -extldflags '-static'" -o game.exe ./cmd/game
