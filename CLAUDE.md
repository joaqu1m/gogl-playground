# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

`gogl-playground` is a Go 3D renderer / game-engine playground built on OpenGL 4.1 core via CGO (`go-gl/gl`, `go-gl/glfw`). It loads glTF 2.0 / `.glb` models (`qmuntal/gltf`) and renders them with a multi-light forward renderer and directional shadow mapping. Long-term goals (Vulkan/Metal backends, animations, CGO-cost research) live in `README.md`'s TODO.

> This file is the single source of truth for working in this repo (architecture, build/run/debug, conventions). `AGENTS.md` just points here.

**⚠️ Before you build, run, or debug the app, read `docs/claude-dev-workflow.md`.** It is the source of truth for how Claude builds/launches/kills the app, reads logs, and drives `dlv` — including hard-won gotchas: CGO commands must run via **PowerShell, not git-bash**; the three run+debug modes; and the precise process-kill filter.

## Commands

Builds require CGO and platform GL toolchains, so target via the `Makefile`:

- `make windows` — builds `game.exe`. Auto-detects the host via `$(OS)`: on Windows uses the native gcc (`-H=windowsgui -s -w`), on Linux/macOS cross-compiles with the MinGW toolchain + static link. On Windows there is no GNU `make` by default — use `mingw32-make windows` (bundled with the choco mingw) or `choco install make`.
- `make linux` — installs `libglfw3-dev libgl1-mesa-dev xorg-dev` if missing, then builds `game`.
- `make mac` — builds a macOS `game` binary.
- `go build ./cmd/game` — quick compile check (no platform-specific flags).
- `go test ./...` — runs tests. There is no test suite yet; add `_test.go` files beside the package under test.

There is no lint config beyond `gofmt -w`. Run the binary from the repo root — asset paths (`assets/*.glb`) and the `logs/` directory are resolved relative to the working directory.

### Native Windows build & debug

The `Makefile` targets cross-compile *for* Windows from Linux (devcontainer/CI). To build and debug *on* Windows natively — needed for delve and for a real GLFW window on the GPU, neither of which works in the headless container:

- Requires MinGW-w64 **gcc** on PATH (`choco install mingw`) and `CGO_ENABLED=1` (persisted via `go env -w`). go-gl needs the GCC toolchain, not MSVC. No system GL libs are needed on Windows (unlike Linux).
- Go is managed by **mise**; the project pins `go = "1.26.1"` in `mise.toml` (kept in sync with the `go.mod` directive and the devcontainer Go feature). A fresh clone must run `mise trust` once before the shim activates. The mise setting `go_set_gobin=false` keeps `go install` binaries in the stable `%USERPROFILE%\go\bin` (on PATH) instead of mise's versioned dir.
- Build from the repo root: `mingw32-make windows` (or `make windows` if GNU make is installed) produces the `-H=windowsgui` release binary. For iterating/debugging use `go build ./cmd/game` or `go run ./cmd/game` directly — those keep a console window for stdout/logs.
- Debug: driven through `dlv` (trace / scripted breakpoints) — see `docs/claude-dev-workflow.md` for the three run+debug modes.

## Architecture

**Entry point & main loop.** `cmd/game/main.go` constructs an `engine.App` via `NewApp`, registers models on `app.ModelManager`, and calls `app.Run()`. `Run` (`engine/app.go`) is the classic GLFW loop: compute delta time → `Camera.Update` → `Draw` → swap buffers → poll events. `cmd/game` is wiring/scene-setup only; keep engine logic out of it.

**Single-thread invariant.** `main.init()` calls `runtime.LockOSThread()`. OpenGL state is thread-affine, so all GL calls (including model loading, which allocates GL buffers/textures eagerly) must stay on the main goroutine. Do not spawn goroutines that touch GL.

**Render passes.** `App.Draw` (`engine/renderer.go`) runs an ordered slice of `RenderPass` (`engine/renderpass.go`), each binding a framebuffer then invoking its `DrawFunc`:
1. `shadowPass` — renders scene depth from the directional light's POV into `ShadowMap`'s depth-only FBO using the minimal `ShadowProgram`.
2. `geometryPass` — renders to the default framebuffer (framebuffer 0) with the main `ShaderProgram`, sampling the shadow map.

`BeginFrame` clears the screen before passes run. Adding a pass (e.g. an HDR/post pass) means appending to the `passes` slice and supplying a `DrawFunc`.

**Shaders are embedded Go strings** in `engine/shader.go` (note the `+ "\x00"` null terminators required by the GL API) — there are no external `.glsl` files. The main fragment shader implements Blinn-Phong for a `lights[MAX_LIGHTS]` (=8) array where one unified `Light` struct covers directional/point/spot via a `type` discriminant; `engine/light.go` mirrors this struct on the Go side. Only the first directional light casts shadows (`shadowCasterLight`).

**Uniforms** go through `UniformCache` (`engine/uniform.go`), which lazily resolves and caches `glGetUniformLocation` per program. Each program (`ShaderProgram`, `ShadowProgram`) has its own cache, both held in `App.DrawingContext`. Use the `SetUniform*` free functions rather than calling `gl.Uniform*` directly.

**Model loading & the scene graph.** `domain/model` wraps the low-level loader: `Model` pairs an `entities.Transform` (the placement you set in `main.go`) with a `gltfloader.GLTFModel`. `gltfloader.LoadGLB` (`libs/gltfloader/loader.go`) walks the glTF scene graph recursively (`processNode`), accumulating each node's local TRS into a baked world transform stored per `GLTFMesh.Transform`, and uploads interleaved `pos(3)+normal(3)+uv(2)` vertex buffers plus textures to GL. At draw time the final matrix is `Model.MeshWorldMatrix` = entity transform × mesh node transform. Missing normals are generated flat; external-URI textures are unsupported (only embedded/buffer-view images).

**Matrix conventions.** Everything is **column-major**. `gmath.Mat4` is `[16]float32`; the loader uses a parallel set of raw `[16]float32` helpers (`mat4fMul`, `composeTRS`). `gmath.MatMul(a, b)` computes `a*b`, so transforms read right-to-left (e.g. `projection * view * model`, `T * R * S`). Transform composition lives in `entities.Transform.ToMat4`.

**Resource lifecycle.** Every GL-owning type has a `Destroy()` that deletes its buffers/textures/framebuffers. `App.Shutdown` (called at loop exit) cascades through `ModelManager.Destroy` → each mesh, plus the shadow map and shader programs. When adding GL resources, wire them into this teardown chain.

**Camera** is an interface (`engine/camera.go`): `ViewMatrix`, `ProjectionMatrix(aspect)`, `Position`, `Update(dt)`. `NewApp` installs an `FPSCamera` (`engine/fps_camera.go`) that captures the cursor and handles WASD + mouse-look; swap implementations by assigning `app.Camera`.

**Logging.** `libs/logger` is a singleton initialized in `init()` that writes timestamped files to `logs/` (keeping the newest 5) at DEBUG level. Use `logger.Infof/Debugf/Errorf`; `logger.Fatalf` logs then `os.Exit(1)` and is the engine's convention for unrecoverable GL/setup failures.

## Conventions

**Style.** Standard Go formatting — run `gofmt -w` on changed files; tabs for indentation; exported identifiers in `CamelCase`, package-private in `camelCase`. Prefer small, focused packages with descriptive filenames (`renderpass.go`, `model_manager.go`). Keep `cmd/game/main.go` to wiring / scene setup / startup only. No linter beyond `gofmt`.

**Tests.** No suite yet; ship new logic with table-driven `_test.go` files beside the package they exercise (e.g. `engine/renderer_test.go`). Run `go test ./...` before a PR, and add regression tests when fixing math/loader/renderer bugs.

**Commits & PRs.** Conventional prefixes scoped by package: `feat(engine): …`, `fix(window): …`, `refactor(renderer): …`. Imperative messages. PRs: short summary, the reason for the change, linked issues, and screenshots/video for any rendering or window-behavior change. Call out platform-specific build or CGO assumptions.
