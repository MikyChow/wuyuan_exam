package main

import (
	ginserver "gin-server"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/server"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/store"
	"log"
	"os"
	"strconv"
)

func main() {
	sqlConn := os.Getenv("SqlConn")
	if len(sqlConn) <= 0 {
		sqlConn = "host=192.168.31.128 port=5432 user=postgres password=example dbname=wuyuan_exam sslmode=disable"
		log.Println("app start with default connection string !")
	}
	nodeID, err := strconv.Atoi(os.Getenv("nodeID"))
	if err != nil {
		log.Println("app start with default node id !")
		nodeID = 1
	}
	if !(nodeID >= 0 && nodeID <= 1024) {
		panic("node id must between 0 and 1024")
	}
	manager := manage.NewDefaultManager()

	// token store
	manager.MustTokenStorage(store.NewFileTokenStore("data.db"))

	// client store
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost",
	})
	manager.MapClientStorage(clientStore)
	ginserver.InitServer(manager)
	ginserver.SetAllowGetAccessRequest(true)
	ginserver.SetClientInfoHandler(server.ClientFormHandler)

	r := gin.Default()
	taskService, err := initTaskService(sqlConn, int64(nodeID), r)
	authGroup := r.Group("/oauth2")
	authGroup.GET("/token", ginserver.HandleTokenRequest)

	taskGroup := taskService.CreateTaskGroup()
	taskGroup.Use(ginserver.HandleTokenVerify())
	if err != nil {
		panic(err)
	}

	err = r.Run(":2021")
	panic(err)
}
