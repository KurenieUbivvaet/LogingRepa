package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	pb "loging_module/protos"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedMyServiceServer
}

func (s *server) CreateLog(ctx context.Context, in *pb.LogRequest) (*pb.LogResponse, error) {
	status := CreateLoging(in.GetProject(), in.GetPodName(),
		in.GetIp(), strconv.Itoa(int(in.GetLogLavel())), in.GetLavelStr(), in.GetMessage())
	return &pb.LogResponse{
		Uuid:      in.GetUuid(),
		LogStatus: int32(status),
	}, nil
}

func CreateLoging(project, podName, ip, logLavel, lavelStr, returned string) int {
	t := time.Now()
	serverEventDatetime := fmt.Sprintf("%d-%02d-%02d		%02d:%02d:%02d,%03d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.UnixMilli()%1000)
	f, err := os.OpenFile("files/loging.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errors.Is(err, os.ErrNotExist) {
		os.Create("files/loging.log")
	}
	f.WriteString("serverEventDatetime='" + serverEventDatetime + "' project='" + project + "' podName='" + podName + "' ip=" +
		ip + "' logLavel='" + logLavel + "' lavelStr='" + lavelStr + "' returned='" + returned + "'\n")
	defer f.Close()
	status := 201
	if err != nil {
		status = 501
	}
	return status
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMyServiceServer(s, &server{})
	log.Printf("Сервер запущен: %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
