FROM golang:1.22.1-alpine3.19

RUN mkdir -p ./loging_module
WORKDIR ./loging_module

COPY . .

RUN go mod init loging_module
RUN go get google.golang.org/grpc

CMD ["go", "main.go"]