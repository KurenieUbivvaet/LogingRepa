FROM golang:1.22.1-alpine3.19

RUN mkdir -p ./api_gateway
WORKDIR ./api_gateway

COPY . .

RUN export PROJECT=name-0.0.1
RUN go get github.com/gin-gonic/gin github.com/google/uuid google.golang.org/grpc github.com/swaggo/swag/cmd/swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

CMD ["go","run", "main.go"]