FROM golang:1.22.1-alpine3.19

RUN mkdir -p ./login_mod
WORKDIR ./login_mod

COPY . .

RUN go get google.golang.org/grpc

CMD ["go", "run", "server.go"]