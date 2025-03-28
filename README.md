# certifigo

![action workflow](https://github.com/exageraldo/certifigo/actions/workflows/release.yml/badge.svg)

Uma ferramenta para geração automatizada de certificados personalizados para eventos, com suporte a templates configuráveis e envio de e-mails.

## Uso

### Gerando certificado de participação

#### Parâmetros Obrigatórios:
- `--name`: Nome da pessoa participante.
- `--email`: Endereço de e-mail da pessoa participante.
- `--event`: Nome do evento.
- `--loc`: Local onde o evento ocorreu.
- `--date`: Data do evento no formato `dd/mm/yyyy`.
- `--duration`: Duração do evento em horas.
- `--signature`: Nome da pessoa responsável pela assinatura do certificado.

#### Parâmetros Opcionais:
- `--signature-img`: Caminho para o arquivo da assinatura a ser utilizado no certificado.
- `--logo`: Caminho para o arquivo de logo a ser utilizado no certificado.
- `--notify`: Indica se o participante deve ser notificado por e-mail (flag opcional).
- `--config`: Caminho para o arquivo de configuração adicional no formato TOML.

```sh
certifigo generate attendee \
    --name="Nome da Pessoa Participante" \
    --email="nome@email.com" \
    --event="11º Nome do Evento" \
    --loc="Nome do Local" \
    --date="01/01/2024" \
    --duration=4 \
    --signature="Nome da Pessoa Assinante" \
    --logo="caminho/para/logo.png" \
    --notify \
    --config="configuracao.toml"
```



### Geração de certificado de palestrante

#### Parâmetros Obrigatórios:
- `--name`: Nome do palestrante.
- `--email`: E-mail do palestrante.
- `--talk-title`: Título da palestra.
- `--talk-duration`: Duração da palestra em minutos.
- `--event`: Nome do evento.
- `--loc`: Local do evento.
- `--date`: Data do evento no formato `dd/mm/yyyy`.
- `--duration`: Duração total do evento em horas.
- `--signature`: Nome da pessoa responsável pela assinatura do certificado.

#### Parâmetros Opcionais:
- `--signature-img`: Caminho para o arquivo da assinatura a ser utilizado no certificado.
- `--logo`: Caminho para o arquivo de logo a ser utilizado no certificado.
- `--attendee`: Indica se o palestrante também é participante do evento.
- `--notify`: Indica se o palestrante deve ser notificado por e-mail após a geração do certificado.
- `--config`: Caminho para o arquivo de configuração no formato TOML.

```sh
certifigo generate speaker \
    --name="Nome da Pessoa Palestrante" \
    --email="nome@email.com" \
    --talk-title="Algum titulo" \
    --talk-duration=30 \
    --event="11º Nome do Evento" \
    --loc="Nome do Local" \
    --date="01/01/2024" \
    --duration=4 \
    --signature="Nome da Pessoa Assinante" \
    --logo="caminho/para/logo.png" \
    --attendee \
    --notify \
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
talk_title="Algum Titulo"
talk_duration=30 # minutos
attendee=true # se a pessoa também for participante
notify=true # se deseja notificar a pessoa por email

[[speakers]]
# ...
```

Uma vez que o arquivo com as informações do evento foi criado, os certificados podem ser gerados com o comando:

```sh
certifigo generate from-file \
    --file="evento.toml" \
    --config="configuracao.toml"
```

#### Parâmetros Opcionais:
- `--config`: Caminho para o arquivo de configuração no formato TOML.

#### Parâmetros Obrigatórios:
- `--config`: Caminho para o arquivo de configuração no formato TOML.

### Definindo as credenciais para enviar email

A ferramenta utiliza o serviço de e-mail para enviar mensagens automatizadas. Para configurar o envio de e-mails, é necessário definir as seguintes variáveis de ambiente:

- `EMAIL_SENDER`: O endereço de e-mail que será utilizado como remetente.
- `EMAIL_PASSWORD`: A senha ou token de acesso do e-mail remetente.

**Nota:** Atualmente, apenas credenciais do Gmail são suportadas. Certifique-se de habilitar o acesso a aplicativos menos seguros ou configurar um token de acesso específico para o envio de emails via SMTP.

Certifique-se de que todas as variáveis de ambiente estejam corretamente configuradas antes de utilizar a funcionalidade de envio de email.

### Arquivo de configuração

O arquivo de configuração é um arquivo no formato TOML que define as configurações para a geração dos certificados. Ele é opcional em todos os comandos; caso não seja fornecido, as configurações padrões internas da ferramenta serão utilizadas. Para especificar um arquivo de configuração personalizado, utilize a flag `--config`.

```toml
# Arquivo de configuração padrão

certification_size="1600x800"

[background]
border_color = "#616161[100%]"
border_size=20
color="#000000[100%]"

[text]
fonts_dir="fonts/"
text_size=30
text_color = "#ffffff"
title_text_size=80
title_text_color = "#ffffff[100%]"
person_text_size=70

[validator]
min_length=8
max_length=11
text_size=20
text_color = "#ffffff[35%]"

[signature]
folder="signatures/"
line_length=22
img_size=100
text_size=60
text_color = "#ffffff[100%]"
title_size=15
title_color = "#ffffff[100%]"

[output]
folder="output/"
default_file_name="_output.json"

[attendee]
title = "CERTIFICADO DE PARTICIPAÇÃO"
body = """
participou do {{ .Event.Name }}, realizado no dia {{ .Event.Date }},
nas instalações da {{ .Event.Location }}, com carga horária total de {{ .Event.Duration }} horas.
"""
email_subject = "Seu certificado chegou!"
email_body = """
Olá, tudo bem?

Aqui está seu certificado de participação do evento {{.Event.Name}}

Att,
"""

[speaker]
title = "CERTIFICADO DE PALESTRANTE"
body = """
participou do {{ .Event.Name }}, realizado no dia {{ .Event.Date }},
nas instalações da {{ .Event.Location }}, com carga horária total de {{ .Event.Duration }} horas.
"""
email_subject = "Seu certificado chegou!"
email_body = """
Olá, tudo bem?

Aqui está seu certificado de participação do evento {{.Event.Name}}

Att,
"""
```
O atributo `certification_size` define as dimensões do canvas utilizado para a certificação. O valor é especificado no formato "largura x altura" (em pixels), onde "1600" representa a largura e "800" representa a altura. Certifique-se de ajustar este valor conforme necessário para atender aos requisitos de design ou resolução desejados.


Os valores das cores são definidas no formato hexadecimal, onde os primeiros seis caracteres representam as cores RGB (vermelho, verde, azul). Opcionalmente, pode-se adicionar um valor de opacidade (alpha) no formato `"[XX%]"`, onde XX é a porcentagem de opacidade desejada. Por exemplo, `"#000000[100%]"` representa preto com opacidade total, enquanto `"#616161[50%]"` seria um cinza com 50% de opacidade. Se o valor alpha não for especificado, assume-se opacidade total (100%). 

As variáveis dentro dos templates, como `{{ .Event.Name }}`, são placeholders que serão substituídos pelos valores correspondentes definidos no arquivo de configuração ou fornecidos durante a execução do comando. Além disso, os objetos disponíveis para uso nos templates são `certifigo.Event` e `certifigo.CertificateConfigFile`. Esses objetos fornecem acesso às informações do evento e às configurações do arquivo de configuração, respectivamente.

Por exemplo:
- `{{ .Event.Name }}` será substituído pelo nome do evento.
- `{{ .Event.Date }}` será substituído pela data do evento.
- `{{ .Event.Location }}` será substituído pelo local do evento.
- `{{ .Event.Duration }}` será substituído pela duração do evento.
- `{{ .Config.CanvaSize.Width }}` será substituído pela largura do canvas definida no valor padrão do atributo `certification_size`.
- `{{ .Config.CanvaSize.Height }}` será substituído pela altura do canvas definida no valor padrão do atributo `certification_size`.
- `{{ .Config.Attendee.EmailSubject }}` será substituído pelo assunto do e-mail definida no valor padrão do atributo `attendee.email_subject`.

Essas substituições são realizadas automaticamente pelo mecanismo de template da ferramenta, garantindo que os certificados e e-mails gerados contenham as informações corretas e personalizadas para cada participante ou palestrante.


Esses objetos permitem criar templates altamente personalizáveis, garantindo que os certificados e e-mails gerados sejam adaptados às necessidades específicas de cada evento.

## Desenvolvimento

```sh
# Clone o repositório
# $ git clone https://github.com/exageraldo/certifigo.git # (via HTTPS)
# $ gh repo clone exageraldo/certifigo # (via GitHub CLI)
$ git clone git@github.com:exageraldo/certifigo.git # (via SSH)

# Entre na pasta do projeto
$ cd certifigo

# Instale as dependências e construa a cli
$ make build

# Execute a cli
$ ./bin/certifigo --help
```


## Fontes

Estamos usando duas fontes do [`Google Fonts`](https://fonts.google.com) nesse projeto:
- [**Cedarville Cursive**](https://fonts.google.com/specimen/Cedarville+Cursive) - *Designed by [Kimberly Geswein](https://fonts.google.com/?query=Kimberly%20Geswein)*
- [**Open Sans**](https://fonts.google.com/specimen/Open+Sans) - *Designed by [Steve Matteson](https://fonts.google.com/?query=Steve%20Matteson)*