package middlewares

import "github.com/labstack/echo/v4/middleware"

var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte("secret"),
})
