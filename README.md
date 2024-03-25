# Suaçuna CLI

![action workflow](https://github.com/exageraldo/suacuna-cli/actions/workflows/release.yml/badge.svg)

## Uso

Arquivo de configuração

```toml
# configuracao.toml

[CANVA]
SIZE={W=1600, H=800}
OVERLAY_MARGIN_SIZE=20.0

# Colors
BACKGROUND_COLOR={R=92, G=255, B=230, A=255}
OVERLAY_COLOR={R=0, G=0, B=0, A=180}
TEXT_COLOR={R=255, G=255, B=255, A=255}

# Text
ATTENDANCE_TITLE="CERTIFICADO DE PARTICIPAÇÃO"
SPEAKER_TITLE="CERTIFICADO DE PALESTRANTE"
SIGNATURE_LINE_LENGTH=22

# Paths
OUTPUT_DIR="output/"
SIGNATURES_DIR="signatures/"
FONTS_DIR="fonts/"
```

Gerando certificado de participação

```sh
./suacuna-cli generate attendee \
    --name="Nome da Pessoa Participante" \
    --email="nome@email.com" \
    --event="11º Nome do Evento" \
    --loc="Nome do Local" \
    --date="01/01/2024" \
    --duration=4 \
    --signature="Nome da Pessoa Assinante" \
    --config="configuracao.toml"
```

Geração de certificado de palestrante

```sh
./suacuna-cli generate speaker \
    --name="Nome da Pessoa Palestrante" \
    --email="nome@email.com" \
    --talk-title="Algum titulo" \
    --talk-duration=30 \
    --event="11º Nome do Evento" \
    --loc="Nome do Local" \
    --date="01/01/2024" \
    --duration=4 \
    --signature="Nome da Pessoa Assinante" \
    --attendee \
    --config="configuracao.toml"
```

Geração de certificado usando arquivo de configuração

```sh
./suacuna-cli generate from-file \
    --file="evento.toml" \
    --config="configuracao.toml"
```

Arquivo com informações do evento

```toml
# evento.toml

[event]
name="1º Nome do Evento"
location="Nome do Local"
date="01/01/2024"
duration=4
signature="Nome da Pessoa Assinante"

[[attendees]]
name="Nome da Pessoa Participante"
email="nome@email.com"
notify=true

[[attendees]]
# ...

[[speakers]]
name="Nome da Pessoa Palestrante"
email="nome@email.com"
talkTitle="Algum Titulo"
talkDuration=30
notify=true

[[speakers]]
# ...
```



## Desenvolvimento

```sh
# Clone o repositório
# $ git clone https://github.com/exageraldo/suacuna-cli.git # (via HTTPS)
# $ gh repo clone exageraldo/suacuna-cli # (via GitHub CLI)
$ git clone git@github.com:exageraldo/suacuna-cli.git # (via SSH)

# Entre na pasta do projeto
$ cd suacuna-cli

# Instale as dependências e construa a cli
$ make build

# Execute a cli
$ ./bin/suacuna-cli --help
```


## Leia mais

- [Fontes](fonts/README.md)
- [Releases](https://github.com/exageraldo/suacuna-cli/releases)