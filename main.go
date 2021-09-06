package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/api/v1/query/:query", query)
	router.Run("localhost:8080")
}

func query(c *gin.Context) {
	query := c.Param("query")
	connection, _ := ExasolDriver{}.Open("exa:localhost:8563;user=sys;password=<pass>;encryption=0;usetls=0")
	json, _ := connection.simpleExec(query)
	connection.Close()
	c.IndentedJSON(http.StatusOK, json)
}
