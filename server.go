package xfsgotoken

import (
	"xfstoken/logs"

	"github.com/gin-gonic/gin"
)

type server struct {
	ginEngine *gin.Engine
	log       logs.ILogger
}

func Start() {

}

func Stop() {

}

func setupRouter() *server {
	return &server{
		ginEngine: gin.Default(),
		log:       logs.NewLogger("server"),
	}
}
