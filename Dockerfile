FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
RUN go mod download

ENV CONFIG_PATH=/app/.config
COPY . /app

COPY ./config ${CONFIG_PATH}

RUN apk add --no-cache make
RUN make init

RUN go build -o app ./cmd/main.go

EXPOSE ${PORT}

CMD ["./app"]
