# gogl-playground

Game engine / 3d renderer em Go.

Nome provisório, a ideia por enquanto é criar um playground para experimentar e aprender sobre gráficos 3D em Go, com foco inicial em OpenGL, Vulkan e possivelmente Metal.

## Definições do escopo inicial

- **GLTF 2.0** como formato de modelo 3D principal, usando a biblioteca [qmuntal/gltf](https://github.com/qmuntal/gltf) para parsing dos modelos para o formato de dados interno do projeto.
- **OpenGL** como API gráfica inicial, com a possibilidade de adicionar suporte a Vulkan e Metal no futuro.
- **CGO** para integração com as APIs gráficas, mas com a intenção de minimizar o impacto de performance e explorar alternativas "Go puras" se o custo do CGO for relevante.
- **Multiplataforma**: o projeto deve ser capaz de compilar e rodar em Windows, Linux e macOS, com atenção especial às diferenças de implementação e suporte das APIs gráficas em cada plataforma.

## TODO

### Finalizar implementação OpenGL
- [ ] Implementar sistema de animações
- [ ] Implementar iluminação e sombreamento

### Build & Run multiplataforma
- [X] Testar compilação e execução em **Linux**
- [ ] Testar compilação e execução em **macOS**

### Pesquisa e Implementação Vulkan *(casado com a pesquisa de CGO)*
- [ ] Pesquisar a melhor forma de integrar Vulkan como substituto/alternativa ao OpenGL
- [ ] Se viável, criar uma **interface genérica** com métodos comuns para abstrair OpenGL e Vulkan
- [ ] Se a interface genérica for inviável, implementar os dois backends de forma independente momentaneamente, antes de decidirmos com qual está valendo mais a pena seguir
- [ ] Implementar e testar as features com o backend em Vulkan

### Pesquisa de performance CGO
> **Importante:** o impacto do CGO deve ser avaliado separadamente para **OpenGL** e **Vulkan**, pois o custo provavelmente é diferente em cada API.

- [ ] Medir o custo real de performance do CGO no contexto do projeto
- [ ] Caso o custo seja relevante, investigar a viabilidade de uma abordagem "Go puro" (sem CGO), começando pelo desenvolvimento em Windows

### Tela de Crash Report
- [ ] Implementar uma tela de crash report que:
  - Fecha o arquivo de log ao detectar o crash
  - Exibe os logs em formato de **tabela estruturada** (`timestamp` | `log level` | `log message`)
  - Disponibiliza também a visualização em **plaintext**

### Documentação pública da Engine conforme o projeto evolui
- [ ] Entender a melhor maneira de fazer isso - provavelmente com um go-to-markdown (como [esse aqui](https://github.com/princjef/gomarkdoc)) e depois para um site estático com hugo
