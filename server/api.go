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
		url: "start",
		method: "POST",
		handler: testPost,
	},
	api {
		url: "debug",
		method: "POST",
		handler: debug,

	},
}

func test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	  })
  
	  cmd := TestCommand {
		Command: Command {
			typ : COMMAND_TEST,
		},
		desc : "test",
	  }

	  serv.PostCommand(cmd)
}

func testPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "test post",
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

func debug(c *gin.Context) {
	var cmd DebugCommand
	postCommand(c, &cmd)

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