FROM golang:1.16

WORKDIR /app

ENV CONFIG=docker

COPY ./ /app

RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 5001

ENTRYPOINT CompileDaemon --build="go build cmd/api_gateway_service/main.go" --command=./main