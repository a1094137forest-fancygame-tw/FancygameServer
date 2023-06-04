package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"

	"BackendServer/config"
	"BackendServer/constant"
	v1 "BackendServer/controllers/v1"
	"BackendServer/model"
)

func init() {
	constant.ReadConfig(".env")
	v1.Initialize()
}

func main() {
	r := gin.Default()
	m := melody.New()

	m.HandleMessageBinary(func(s *melody.Session, msg []byte) {
		messageCode := int(msg[1])
		log.Println("messageCode", messageCode)
		client, _ := s.Get("client")
		if client, ok := client.(v1.Client); ok {
			v1.GetMessage(&client, messageCode, msg[2:])
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		///go Ping(s)
		fmt.Println("New connection")
		client := v1.Client{
			IsWait:  false,
			Session: s,
		}

		s.Set("client", client)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		fmt.Println("Disconnected")
	})

	m.HandleError(func(s *melody.Session, err error) {
		log.Println("err", err)
	})

	r.GET("/", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	err := r.Run(config.Port)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func Ping(s *melody.Session) {
	for {
		req := model.Ping{
			MessageCode: 0,
			SystemTime:  int(time.Now().Unix()),
		}

		jsonD, _ := json.Marshal(req)
		s.Write(jsonD)
		time.Sleep(5 * time.Second)
	}
}
