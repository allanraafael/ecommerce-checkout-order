FROM golang:1.17

WORKDIR /go/src/

# Copiando diretório sem os arquivos 'go.mod' e 'go.sum'
COPY ./product/ .

# Copiando manualmente os arquivos 'go.mod' e 'go.sum' para container
COPY go.mod ./go.mod
COPY go.sum ./go.sum

# Criando imagem no formato linux e passando argumentos opcionais
RUN GOOS=linux go build -ldflags="-s -w"
CMD ["./order"]
