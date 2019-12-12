package main

import (
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "./handler"
    "./libs"
    ext "./middleware"
)

func main() {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    libs.InitDB()

    e.GET("/hello", handler.Hello())
    e.POST("/login", handler.Login())
    r1 := e.Group("/restricted")
    r1.Use(middleware.JWT([]byte("secret1")))
    r1.POST("", handler.Restricted())

    r2 := e.Group("/reauth")
    config := middleware.JWTConfig{
        Claims:     &ext.MyClaim{},
        SigningKey: []byte("secret2"),
    }
    r2.Use(middleware.JWTWithConfig(config))
    r2.POST("", handler.ReAuth())

    e.Start(":3000")
}
