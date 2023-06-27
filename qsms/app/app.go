package app

import (
	"github.com/labstack/echo/v4"
	"qsms/controllers"
	"qsms/repository"
	"qsms/services"
)

type App struct {
	E *echo.Echo
}

func NewApp() *App {
	e := echo.New()

	userRepository := repository.NewGormUserRepository()
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	paymentRepository := repository.NewGormPaymentRepository()
	paymentService := services.NewPaymentService(paymentRepository)
	paymentController := controllers.NewPaymentController(userService, paymentService)

	userGroup := e.Group("/sms/user")
	userGroup.POST("/signup", userController.Signup)
	userGroup.GET("/login", userController.Login)
	userGroup.GET("/me", userController.GetUserByToken)
	userGroup.POST("/logout", userController.LogOut)

	paymentGroup := e.Group("/sms/payment")
	paymentGroup.GET("/pay/:amount", paymentController.AddBalance)
	paymentGroup.GET("/verify/", paymentController.VerifyPayment)

	return &App{
		E: e,
	}
}

func (a *App) Start(addr string) error {
	a.E.Logger.Fatal(a.E.Start(addr))
	return nil
}
