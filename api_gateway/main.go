package main

import (
	_ "api_gateway/docs"
	pb "api_gateway/protos"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	GRPCHost string `json:"gRPC_host"`
	GRPCPort int    `json:"gRPC_port"`
	ApiHost  string `json:"api_host"`
	ApiPort  int    `json:"api_port"`
}

var addr = flag.String("addr", "localhost:50051", "the addres to connect to")

//	@title			Журналирование
//	@version		0.0.1
//	@description	Здесь по идее должен быть списочек.

//	@contact.name	KurenieUbivvaet

// @host		localhost:8000
// @BasePath	/
func main() {
	configFile, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	jsonError := json.Unmarshal(data, &config)
	if jsonError != nil {
		log.Fatal(jsonError)
	}

	fmt.Println(config)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.POST("/", MainPost)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.html", gin.H{})
	})

	router.POST("/outlog", PostOutlog)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ApiAdr := ":" + strconv.Itoa(config.ApiPort)
	Error := router.Run(ApiAdr)
	if Error != nil {
		panic("[Error] failed to start Gin server due to: " + Error.Error())
	}
}

// PostOutLog godoc
// @Summary Summarize your endpoint
// @Schemes
// @Description логи из удалённого сервера
// @Router /outlog [post]
func PostOutlog(c *gin.Context) {
	UUID := c.PostForm("UUID")
	massage := c.PostForm("massage")
	logLev := c.PostForm("LogLevel")
	logLevel, _ := strconv.Atoi(logLev)
	levelStr := c.PostForm("levelStr")
	Project := c.PostForm("project")
	PodName := c.PostForm("podName")
	Ip := c.PostForm("ip")
	var getUUID, getLogStatus = logRequest(UUID, massage, logLevel, levelStr, Project, PodName, Ip)
	c.JSON(200, gin.H{
		"UUID":      getUUID,
		"LogStatus": getLogStatus,
	})
}

// MainPost godoc
// @Summary Summarize your endpoint
// @Schemes
// @Description ручные удалённого сервера
// @Router / [post]
func MainPost(c *gin.Context) {
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
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Fatal("UUID v1 не сгенерирован!")
	}
	return LogingStruct{
		UUID:     uuid.String(),
		message:  message,
		logLavel: logLavel,
		levelStr: levelStr,
		project:  os.Getenv("PROJECT"),
		podName:  os.Getenv("hostname"),
		ip:       GetLocalIP().String(),
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
