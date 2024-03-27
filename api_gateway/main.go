package main

import (
	pb "api_gateway/protos"
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", "localhost:50051", "the addres to connect to")

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("templates/create.html")

	router.POST("/", func(c *gin.Context) {
		message := c.PostForm("message")
		logLev := c.PostForm("logLevel")
		logLevel, _ := strconv.Atoi(logLev)
		levelStr := c.PostForm("levelStr")
		loging := CreateLoging(message, logLevel, levelStr)

		var getUUID, getLogStatus = logRequest(
			loging.UUID, loging.message,
			loging.logLavel, loging.levelStr,
			loging.project, loging.podName,
			loging.ip)

		c.JSON(200, gin.H{
			"UUID":      getUUID,
			"LogStatus": getLogStatus,
		})

	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.html", gin.H{})
	})

	router.Run()
}

func logRequest(UUID, message string, logLavel int, lavelStr, project, podName, ip string) (string, int32) {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewMyServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateLog(ctx, &pb.LogRequest{
		Uuid:     UUID,
		Message:  message,
		LogLavel: int32(logLavel),
		LavelStr: lavelStr,
		Project:  project,
		PodName:  podName,
		Ip:       ip,
	})
	if err != nil {
		log.Fatalf("cloud not greet: %v", err)
	}
	uuid := r.GetUuid()
	logStatus := r.GetLogStatus()
	return uuid, logStatus
}

type LogingStruct struct {
	UUID     string
	message  string
	logLavel int
	levelStr string
	project  string
	podName  string
	ip       string
}

func CreateLoging(message string, logLavel int, levelStr string) LogingStruct {
	uuid := uuid.NewString()
	hostname := os.Getenv("hostname")
	ProjectName := os.Getenv("PROJECT")
	IpAddres := GetLocalIP().String()
	return LogingStruct{
		UUID:     uuid,
		message:  message,
		logLavel: logLavel,
		levelStr: levelStr,
		project:  ProjectName,
		podName:  hostname,
		ip:       IpAddres,
	}
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	loclaAddress := conn.LocalAddr().(*net.UDPAddr)
	return loclaAddress.IP
}
