# Contribuindo para o Awesome Tech Lead

Obrigado pelo seu interesse em contribuir para o Awesome Tech Lead! Este
documento fornece as diretrizes para adicionar novos itens ao nosso catálogo de
conteúdo.

## Como Contribuir

1. Faça um fork deste repositório
2. Clone o fork para sua máquina local
3. Crie uma nova branch para suas alterações
4. Adicione seus itens ao arquivo `catalog.yml`
5. Faça commit das suas alterações
6. Envie um Pull Request

## Adicionando Itens ao Catálogo

Para adicionar um novo item ao catálogo, você deve editar o arquivo
`catalog.yml` e adicionar uma nova entrada seguindo o esquema abaixo:

```yaml
- url: https://exemplo.com/conteudo
  title: Título do Conteúdo
  author: Nome do Autor (OPCIONAL)
  type: article
  tags:
    - Arquitetura e Design
    - Excelência Técnica
  is_paid: false
  level: intermediate
  career_bands:
    - mid
    - senior
  duration: 30
  language: pt_br
```

> [!IMPORTANT]
> O arquivo README.md é gerado e você não deve modifica-lo manualmente.

## Esquema do Catálogo

| Campo          | Descrição                                               | Obrigatório |
| -------------- | ------------------------------------------------------- | ----------- |
| `url`          | URL onde o conteúdo pode ser acessado                   | Sim         |
| `title`        | O título do conteúdo                                    | Sim         |
| `authors`      | Autor(es) ou criador do conteúdo                        | Não         |
| `type`         | O tipo de conteúdo disponibilizado                      | Sim         |
| `tags`         | Lista de palavras-chave que categorizam o conteúdo      | Sim         |
| `is_paid`      | Indica se o conteúdo é pago ou gratuito (padrão: false) | Sim         |
| `level`        | O nível de dificuldade do conteúdo                      | Sim         |
| `career_bands` | Níveis de carreira para os quais o conteúdo é relevante | Sim         |
| `duration`     | Duração do conteúdo em minutos                          | Não         |
| `language`     | Idioma do conteúdo                                      | Sim         |

### Valores Permitidos

#### `type`

- `video`: Conteúdo em formato de vídeo
- `article`: Artigo ou publicação em blog
- `book`: Livro
- `podcast`: Episódio ou série de podcast
- `roadmap`: Roteiro de aprendizado
- `feed`: Feed de conteúdo
- `course`: Curso estruturado
- `workshop`: Workshop ou treinamento prático

#### `level`

- `beginner`: Para iniciantes
- `intermediate`: Para pessoas com conhecimento intermediário
- `advanced`: Para pessoas com conhecimento avançado

#### `career_bands`

- `junior`: Desenvolvedores júnior
- `mid`: Desenvolvedores plenos
- `senior`: Desenvolvedores seniores
- `tl`: Tech Leads
- `staff`: Staff Engineers
- `principal`: Principal Engineers

#### `language`

Código ISO do idioma (ex: `pt_br`, `en_us`, `es`, etc.)

## Processo de Revisão

Após enviar seu Pull Request, a equipe de mantenedores irá revisar suas
alterações. Podemos solicitar ajustes ou esclarecimentos antes de mesclar suas
contribuições.

## Gerando o README

O arquivo README.md é gerado automaticamente a partir dos dados do catálogo.
Após adicionar novos itens ao arquivo `catalog.yml`, você deve regenerar o
README.

Você pode gerar o README de duas formas:
- utilizando o docker (não necessita instalar o GO na sua máquina)
- gerando o arquivo utilizando o GO diretamente via terminal (necessita instalar o GO lang)

### Gerando o README utilizando o docker
1. Certifique-se de ter o docker (24.0 ou superior) instalado em sua máquina

primeiro faça o build da imagem docker usando o comando:

```bash
docker build -t iawesome-tech-lead .
```

Logo após isso, execute o comando abaixo, que vai executar um container temporário responsável por gerar o arquivo README.md

```bash
docker run --rm -v "$PWD":/app iawesome-tech-lead
```

### Gerando o README utilizando o GO

1. Certifique-se de ter o GO instalado em sua máquina (versão 1.16 ou superior)
2. Execute o comando:

```bash
make generate-readme
```

Estes comandos irão processar o arquivo `catalog.yml` e atualizar o README.md com
os novos itens do catálogo.

## Configurando o Ambiente de Desenvolvimento

### Pré-requisitos

#### Usando docker

- [Docker](https://docs.docker.com/get-started/get-docker/) (versão 24.0 ou superior)
- [Git](https://git-scm.com/downloads)

#### Usando GO

- [Go](https://golang.org/doc/install) (versão 1.16 ou superior)
- [Git](https://git-scm.com/downloads)

### Instalando Dependências

#### Configurando o Ambiente de Desenvolvimento usando docker

1. Clone o repositório:

```bash
git clone https://github.com/seu-usuario/awesome-tech-lead.git
cd awesome-tech-lead
```

2. Instale o Docker e o Docker compose:

Instale a versão correta pro seu sistema operacional, para isso verifique a documentação do [docker](https://docs.docker.com/get-started/get-docker/).

#### Configurando o Ambiente de Desenvolvimento usando GO

1. Clone o repositório:

```bash
git clone https://github.com/seu-usuario/awesome-tech-lead.git
cd awesome-tech-lead
```

2. Instale as dependências:

```bash
make setup
```

3. Verifique se tudo está funcionando corretamente gerando o README:

```bash
make generate-readme
```

Agradecemos sua contribuição para tornar este catálogo ainda mais valioso para a
comunidade TechLeads.club!
