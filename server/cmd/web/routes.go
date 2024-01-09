package main

import (
	_ "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/thecloudnativehacker/kagec2/server/pkg/handlers"
)

func SetRoutes(e *echo.Echo) {
	//front-end pages
	unAuthGroup := e.Group("/")
	unAuthGroup.GET("", handlers.RootHandler)
	unAuthGroup.GET("login", handlers.GetLoginPage)
	unAuthGroup.POST("login", handlers.Login)

	apiGroup := e.Group("/api")
	//tasks assigned to agent
	apiGroup.GET("/tasks", handlers.GetTasks)
	apiGroup.GET("/tasks/", handlers.GetTasks)
	apiGroup.GET("/tasks/:id", handlers.GetTask)
	apiGroup.POST("/tasks", handlers.AddTask)
	apiGroup.POST("/tasks/", handlers.AddTask)
	// apiGroup.DELETE("/tasks", handlers.DeleteTask)
	// apiGroup.DELETE("/tasks/", handlers.DeleteTask)

	//results from tasks performed by implant
	apiGroup.GET("/results", handlers.GetResults)
	apiGroup.GET("/results/", handlers.GetResults)
	apiGroup.GET("/results/:id", handlers.GetResult)
	apiGroup.POST("/results", handlers.AddResult)
	apiGroup.POST("/results/", handlers.AddResult)

	//task history for tasks performed by implant
	apiGroup.GET("/taskhistory", handlers.GetTaskHistory)
	apiGroup.GET("/taskhistory/", handlers.GetTaskHistory)
	apiGroup.GET("/taskhistory/:id", handlers.GetTaskHistoryById)
	apiGroup.POST("/taskhistory", handlers.AddTaskHistory)
	apiGroup.POST("/taskhistory/", handlers.AddTaskHistory)
	//agent information
	// apiGroup.GET("/implants", handlers.GetImplants)
	// apiGroup.GET("/implants/", handlers.GetImplants)
	// apiGroup.GET("/implants/:id", handlers.GetImplant)
	// apiGroup.POST("/taskhistory", handlers.AddImplant)
	// apiGroup.POST("/taskhistory/", handlers.AddImplant)
}
