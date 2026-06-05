# O alvo `windows` funciona em dois ambientes:
#   - Windows nativo (OS=Windows_NT): usa o gcc local (MinGW-w64), GOOS/GOARCH ja
#     sao windows/amd64, e o shell e o cmd.exe (sem prefixo de env var no estilo sh).
#   - Linux/macOS (container/CI): cross-compila com o toolchain MinGW e link estatico.
ifeq ($(OS),Windows_NT)
  WIN_ENV     :=
  WIN_LDFLAGS := -H=windowsgui -s -w
else
  WIN_ENV     := CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64
  WIN_LDFLAGS := -H=windowsgui -s -w -extldflags '-static'
endif

linux: # not tested in a real Linux environment
	@dpkg -s libglfw3-dev libgl1-mesa-dev xorg-dev >/dev/null 2>&1 || \
		(sudo apt-get update && sudo apt-get install -y libglfw3-dev libgl1-mesa-dev xorg-dev)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o game ./cmd/game

mac: # not tested in a real macOS environment
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o game ./cmd/game

# export (recurso do GNU make) injeta CGO_ENABLED no ambiente do comando
# independente do shell, funcionando tanto no cmd.exe quanto no sh.
windows: export CGO_ENABLED=1
windows:
	$(WIN_ENV) go build -ldflags "$(WIN_LDFLAGS)" -o game.exe ./cmd/game
