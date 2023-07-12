package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"qsms/controllers"
	"qsms/middlewares"
	"qsms/repository"
	"qsms/services"
)

type App struct {
	E *echo.Echo
}

func NewApp() *App {
	e := echo.New()

	purchaseRepository := repository.NewGormPurchaseRepository()
	purchaseService := services.NewPurchaseService(purchaseRepository)

	err := purchaseService.RegisterRentingTasks()
	if err != nil {
		fmt.Println(err)
	}

	userRepository := repository.NewGormUserRepository()
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService, purchaseService)

	paymentRepository := repository.NewGormPaymentRepository()
	paymentService := services.NewPaymentService(paymentRepository)
	paymentController := controllers.NewPaymentController(userService, paymentService)

	smsRepository := repository.NewGormMessageRepository()
	smsService := services.NewMessageService(smsRepository, userRepository)
	smsController := controllers.NewMessageController(userService, smsService)

	err = smsService.RegisterMessagingSchedules()
	if err != nil {
		fmt.Println(err)
	}

	adminRepository := repository.NewGormAdminRepository()
	adminService := services.NewAdminService(adminRepository)
	adminController := controllers.NewAdminController(adminService)

	userGroup := e.Group("/sms/user")
	userGroup.POST("/signup", userController.Signup)
	userGroup.GET("/login", userController.Login)
	userGroup.GET("/me", userController.GetUserByToken)
	userGroup.GET("/buy", userController.GetPhoneNumbersToBuy)
	userGroup.GET("/logout", userController.LogOut)
	userGroup.PUT("/buy/:id", userController.BuyNumber, middlewares.NotSuspended)
	userGroup.PUT("/rent/:id", userController.PlaceRent, middlewares.NotSuspended)
	userGroup.DELETE("/dropRent/:id", userController.DropRent, middlewares.NotSuspended)
	userGroup.PUT("/main/:id", userController.SetMainNumber, middlewares.NotSuspended)

	contactGroup := userGroup.Group("/contact", middlewares.NotSuspended)
	contactGroup.POST("/add", userController.AddContact)
	contactGroup.DELETE("/delete/:contactID", userController.DeleteContact)

	phoneBookGroup := userGroup.Group("/phonebook", middlewares.NotSuspended)
	phoneBookGroup.POST("/create", userController.CreatePhoneBook)
	phoneBookGroup.PUT("/addNumber/:id/:num", userController.AddNumberToPhoneBook)
	phoneBookGroup.DELETE("/delete/:id", userController.DeletePhoneBook)

	templateGroup := userGroup.Group("/template", middlewares.NotSuspended)
	templateGroup.POST("/create", userController.AddTemplate)
	templateGroup.DELETE("/delete/:id", userController.DeleteTemplate)

	paymentGroup := e.Group("/sms/payment")
	paymentGroup.GET("/pay/:amount", paymentController.AddBalance, middlewares.NotSuspended)
	paymentGroup.GET("/verify", paymentController.VerifyPayment)

	smsGroup := e.Group("/sms/send", middlewares.NotSuspended)
	smsGroup.POST("/single", smsController.SingleMessage)
	smsGroup.POST("/periodic", smsController.PeriodicMessage)

	adminGroup := e.Group("/sms/admin", middlewares.IsAdmin)
	adminGroup.POST("/addNumber", adminController.AddNumber)
	adminGroup.POST("/addBadWord", adminController.AddBadWord)
	adminGroup.PUT("/setFee", adminController.SetFee)
	adminGroup.PUT("/suspend/:id", adminController.SuspendUser)
	adminGroup.PUT("/unsuspend/:id", adminController.UnSuspendUser)
	adminGroup.GET("/count/:id", adminController.CountUserMessages)
	adminGroup.GET("/search", adminController.SearchMessages)

	return &App{
		E: e,
	}
}

func (a *App) Start(addr string) error {
	a.E.Logger.Fatal(a.E.Start(addr))
	return nil
}
