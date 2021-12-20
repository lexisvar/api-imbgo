package main

import (
	"github.com/eefret/gomdb"
	"github.com/gin-gonic/gin"
)

var api = gomdb.Init("9b41e7cc")
var router = gin.Default()

func main() {
	router.Run("localhost:9090")
}
