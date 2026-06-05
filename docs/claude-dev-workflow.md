# Guia de colaboração com o Claude Code — debug via DLV

Como o Claude auxilia no desenvolvimento deste projeto: setup, build/run, leitura de logs
e o fluxo de debug **conduzido pelo `dlv` (delve)**. Foi validado empiricamente neste
ambiente (Windows nativo + mise + MinGW). As capacidades e limites são reais, não suposições.

> O debug aqui é todo via `dlv` no terminal. Não usamos o debugger do VS Code (F5). Se um
> dev quiser montar o fluxo dele no editor depois, fica a critério dele — este guia não cobre.

Capacidades do Claude neste projeto:

| Capacidade | Consegue? |
|------------|-----------|
| Buildar (`make windows` / `go build`) | ✅ |
| Iniciar o app (em segundo plano) | ✅ |
| Fechar o app (matar processo) | ✅ |
| Ler logs do app (`logs/`) e o trace do `dlv` | ✅ |
| Conduzir `dlv trace` / breakpoints scriptados | ✅ (modo batch, não REPL contínuo) |
| Ver a janela renderizada (GPU) | ❌ (depende do usuário) |
| "Stream" de logs empurrado em tempo real | ❌ — há equivalente por evento (ver §4) |

---

## 0. Regra de ouro do ambiente

**Todo comando com CGO (build, `dlv`, `go run`) roda via PowerShell, não via Bash/git-bash.**
O git-bash exporta o `PATH` em formato Unix (`/c/...`), que o `cgo`/`gcc` (processos Windows
nativos) não entendem → erro `cgo.exe: exit status 2`. Via PowerShell o `PATH` está em
formato Windows e funciona. O Claude reconstrói o `PATH` do registro para enxergar
`gcc`/`dlv`/`make` mesmo em processos iniciados antes das instalações:

```powershell
$env:Path = [Environment]::GetEnvironmentVariable('Path','Machine') + ';' + [Environment]::GetEnvironmentVariable('Path','User')
```

---

## 1. Setup — detecção e remediação

No início de uma tarefa de build/debug, o Claude verifica e remedia o que faltar.

```powershell
go version                 # mise fornece o Go (pin em mise.toml = 1.26.1)
(Get-Command gcc).Source   # MinGW-w64 (CGO)
(Get-Command dlv).Source   # delve
(Get-Command make).Source  # GNU make (opcional; ha mingw32-make)
go env CGO_ENABLED         # precisa ser 1
mise trust --show          # mise.toml precisa estar "trusted"
```

Remediação por item:

- **mise não confia no `mise.toml`** (`not trusted`): `mise trust` (necessário uma vez por
  máquina em cada clone). Se um editor/language-server estava aberto antes, reinicie-o para
  largar o erro antigo.
- **`CGO_ENABLED` ≠ 1**: `go env -w CGO_ENABLED=1`.
- **`gcc` ausente**: **requer Administrador** — o Claude não eleva. Peça ao usuário:
  `choco install mingw -y` num PowerShell Admin. Sem GCC não há build CGO.
- **`make` ausente**: use `mingw32-make` (já vem com o mingw) ou `choco install make -y` (Admin).
- **`dlv` ausente**: `go install github.com/go-delve/delve/cmd/dlv@latest`
  (cai em `%USERPROFILE%\go\bin`; mantenha esse dir no PATH e `go_set_gobin=false` no mise).

> Quando um passo exigir Administrador, o Claude **pausa e pede ao usuário** para rodar num
> terminal elevado, depois valida.

---

## 2. Buildar, rodar e fechar

**Buildar:**
```powershell
make windows        # ou: mingw32-make windows  (auto-detecta Windows vs container)
go build ./cmd/game # iteracao rapida
```

**Rodar:** o app é GUI (janela GLFW) num loop infinito. O Claude **sempre roda em segundo
plano** (`run_in_background`) — primeiro plano travaria o Claude no loop. A janela abre **no
desktop do usuário**; o Claude **não vê** o que é renderizado.

**Fechar:** matar o processo com **filtro preciso** (lição aprendida: `game*` amplo chega a
matar serviços do Windows como `GameInputSvc`/`GameSDK` — nunca usar):

```powershell
Get-Process -ErrorAction SilentlyContinue |
  Where-Object { $_.ProcessName -in @('dlv','game','game-debug','go') -or $_.ProcessName -like '__debug_bin*' } |
  Stop-Process -Force
```

`dlv` gera binários `__debug_bin*` na raiz (já no `.gitignore`); o Claude os remove ao fim.

---

## 3. Logs

O `libs/logger` escreve em **arquivo**, não no stdout: `logs/debug_<timestamp>.log`
(nível DEBUG, mantém os 5 mais recentes). É a principal fonte que o Claude lê.

```powershell
Get-ChildItem logs\debug_*.log | Sort-Object LastWriteTime | Select-Object -Last 1
```

---

## 4. Os três modos de RUN + DEBUG

### Modo 1 — Claude autônomo (não-interativo)

Claude builda, inicia em background, deixa rodar, mata e lê os logs/trace — **sem
interação humana com a janela**.

- **Pega:** bugs de startup, carregamento de modelo, render do primeiro frame, crash imediato.
- **Não pega:** nada que dependa de o usuário mexer (câmera, input, estados ao longo do tempo).
- **Mais autônomo**, porém o mais limitado em cobertura.

### Modo 2 — Claude conduz o `dlv`, usuário interage

Claude inicia o app **sob `dlv trace`** (em background do lado do Claude, mas a **janela fica
visível e interativa** pro usuário). O usuário mexe livremente; o `dlv` grava o trace num
arquivo. Quando o usuário **fecha o app ou diz "peguei o bug"**, o Claude lê o trace + os
`logs/` e analisa.

- **Ponto-chave:** o `dlv trace` exige que o Claude defina **antes** quais funções rastrear
  (regex). Por isso costuma ser **iterativo**:
  1. usuário descreve a suspeita ("a sombra some quando ando pra trás");
  2. Claude escolhe o que rastrear (ex: `directionalLightSpaceMatrix|shadowPass`);
  3. usuário reproduz;
  4. Claude lê o trace e refina ou conclui.
- É o modo mais poderoso para bugs **interativos** sem o usuário precisar saber de `dlv`.

### Modo 3 — Usuário roda e reporta

O usuário inicia o app na mão (com ou sem `dlv`), observa, e relata: *"rodei aqui, vi algo
estranho no X, vê se acha."* O Claude trabalha sobre o **código + os logs** que o usuário
compartilhar. Útil quando o usuário já sabe onde está o problema ou quando é puramente visual.

### Ferramentas de `dlv` que o Claude usa

**`dlv trace`** — loga entrada/saída de funções que casam um regex. Não-interativo; o Claude
lê direto. Validado:
```powershell
dlv trace ./cmd/game 'LoadGLB'
# > LoadGLB("assets/dead_by_daylight_-_eleven.glb")
# >> LoadGLB => (*GLTFModel(0x...), error nil)
```

**Breakpoints scriptados** — `dlv debug`/`dlv exec` com `--init` contendo `break`,
`continue`, `print`, `stack`. O Claude define o que capturar, roda, lê. Não é um REPL ao vivo
mantido entre mensagens — é batch: define → roda → lê → ajusta.

⚠️ **Cache frio:** a 1ª recompilação do `dlv` (com `-N -l`) pode falhar com
`cgo.exe: exit status 2`. Pré-aquecer resolve:
```powershell
go build -gcflags="all=-N -l" -o $null ./cmd/game
```
Ou evitar o recompile do dlv buildando nós e usando `dlv exec`:
```powershell
go build -gcflags="github.com/joaqu1m/gogl-playground/...=-N -l" -o game-debug.exe ./cmd/game
dlv exec game-debug.exe
```

### Sobre "tempo real"

O Claude **não** recebe stream contínuo empurrado. Mas tem o equivalente **por evento**:

- **`run_in_background`**: processo roda solto escrevendo num arquivo; o Claude lê o arquivo
  a qualquer momento (snapshot atual) e é notificado quando o processo termina.
- **Monitor (espera por condição)**: o Claude "acorda" quando um padrão aparece no arquivo —
  ex.: linha `ERROR` ou um tracepoint específico. É o mais perto de tempo real.

**Na maioria dos bugs nem é necessário tempo real:** o usuário reproduz e o Claude lê o
arquivo depois (Modo 1 e 2). Tempo real só importa em bugs de timing/interação contínua, e
aí o Monitor cobre o gatilho.

---

## 5. Divisão de trabalho (resumo)

| Tipo de bug / tarefa | Modo | Quem conduz |
|----------------------|------|-------------|
| Startup / load / crash imediato | 1 | **Claude** sozinho |
| Lógica/estado dependente de interação | 2 | **Claude** rastreia, usuário reproduz |
| Visual / o usuário já localizou | 3 | **Usuário** reporta, Claude analisa código+logs |
| Build / setup / limpeza | — | **Claude** (pausando em passos Admin) |
