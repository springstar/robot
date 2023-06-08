package server

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"fmt"
)

type api struct {
	url string
	method string
	handler func(c *gin.Context)
}

var apis[] api = []api {
	api {
		url: "ping",
		method: "GET",
		handler: test,
	},
	api {
		url: "bench",
		method: "POST",
		handler: bench,
	},
	api {
		url: "stop",
		method: "GET",
		handler: stop,
	},
	api {
		url: "debug",
		method: "POST",
		handler: debug,

	},
	api {
		url: "report",
		method: "GET",
		handler: report,
	},
	api {
		url: "quit",
		method: "GET",
		handler: quit,
	},
}

func test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	  })
}

func bench(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "bench",
	  })

	var cmd BenchCommand
	cmd.Batch = 1

	fmt.Println("post bench command")
	
	if err := c.BindJSON(&cmd); err != nil {
		fmt.Println(err)
		return
	}

	serv.PostCommand(cmd)	  

}

func stop(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "stop",
	})

	
}

func debug(c *gin.Context) {
	var cmd DebugCommand
	postCommand(c, &cmd)

}

func quit(c *gin.Context) {
	
}

func report(c *gin.Context) {
	cmd := ReportCommand {

	}
	
	serv.PostCommand(cmd)

	res := <- serv.driver.rq

	c.JSON(http.StatusOK, gin.H{
		"message": res,
	  })


}

func postCommand(c *gin.Context, cmd iCommand) {
	c.JSON(http.StatusOK, gin.H{
		"message": "debug",
	})

	if err := c.BindJSON(cmd); err != nil {
		fmt.Println(err)
		return
	}
	
	serv.PostCommand(cmd)
}