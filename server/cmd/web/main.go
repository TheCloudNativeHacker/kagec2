package main

import (
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4/middleware"
	"github.com/thecloudnativehacker/kagec2/server/pkg/render.go"
)

const port = ":1323"

func main() {
	e := echo.New()
	renderer := &render.Template{}
	e.Renderer = renderer
	SetRoutes(e)
	e.Logger.Fatal(e.Start(port))
}
