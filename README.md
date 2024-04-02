# Suaçuna CLI

![action workflow](https://github.com/exageraldo/suacuna-cli/actions/workflows/release.yml/badge.svg)

## Uso

### Arquivo de configuração

O arquivo de configuração é um arquivo TOML que contém as configurações padrão para a geração de certificados. O arquivo de configuração é opcional em todos os comandos, se não for fornecido, as configurações padrão serão usadas. O arquivo de configuração pode ser fornecido com a flag `--config`.

```toml
# Arquivo de configuração padrão

# Text
TextSize=30
TextColor={R=255, G=255, B=255, A=255}

# Title
AttendanceTitle="CERTIFICADO DE PARTICIPAÇÃO"
SpeakerTitle="CERTIFICADO DE PALESTRANTE"
TitleTextSize=80
TitleTextColor={R=255, G=255, B=255, A=255}

# Person
PersonTextSize=70

# Validator
ValidatorMinLength=8
ValidatorMaxLength=11
ValidatorTextSize=20
ValidatorTextColor={R=255, G=255, B=255, A=85} # same as text color with 1/3 of the alpha

# Signature
SignatureDir="signatures/"
SignatureLineLength=22
SignatureImgSize=100
SignatureTextSize=60
SignatureTextColor={R=255, G=255, B=255, A=255}
SignatureTitleSize=15
SignatureTitleColor={R=255, G=255, B=255, A=255}

# Output
OutputDir="output/"
DefaultFileName="_output.json"
```

### Gerando certificado de participação

```sh
suacuna-cli generate attendee \
    --name="Nome da Pessoa Participante" \
    --email="nome@email.com" \
    --event="11º Nome do Evento" \
    --loc="Nome do Local" \
    --date="01/01/2024" \ # dd/mm/yyyy
    --duration=4 \ # horas
    --signature="Nome da Pessoa Assinante" \
    --notify \ # se deseja notificar a pessoa por email
    --config="configuracao.toml"
```

### Geração de certificado de palestrante

```sh
suacuna-cli generate speaker \
    --name="Nome da Pessoa Palestrante" \
    --email="nome@email.com" \
    --talk-title="Algum titulo" \
    --talk-duration=30 \ # minutos
    --event="11º Nome do Evento" \
    --loc="Nome do Local" \
    --date="01/01/2024" \ # dd/mm/yyyy
    --duration=4 \ # horas
    --signature="Nome da Pessoa Assinante" \
    --attendee \ # se a pessoa também for participante
    --notify \ # se deseja notificar a pessoa por email
    --config="configuracao.toml"
```

### Geração de certificado usando arquivo com informações do evento

Arquivo com informações do evento

```toml
# evento.toml

[event]
name="1º Nome do Evento"
location="Nome do Local"
date="01/01/2024" # dd/mm/yyyy
duration=4 # horas
signature="Nome da Pessoa Assinante"

[[attendees]]
name="Nome da Pessoa Participante"
email="nome@email.com"
notify=true # se deseja notificar a pessoa por email

[[attendees]]
# ...

[[speakers]]
name="Nome da Pessoa Palestrante"
email="nome@email.com"
talkTitle="Algum Titulo"
talkDuration=30 # minutos
attendee=true # se a pessoa também for participante
notify=true # se deseja notificar a pessoa por email

[[speakers]]
# ...
```

Uma vez que o arquivo com as informações do evento foi criado, os certificados podem ser gerados com o comando:

```sh
suacuna-cli generate from-file \
    --file="evento.toml" \
    --config="configuracao.toml"
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