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

	phoneBookRepository := repository.NewGormPhoneBookRepository()
	phoneBookService := services.NewPhoneBookService(phoneBookRepository)
	phoneBookController := controllers.PhoneBookController{PhoneBookService: phoneBookService, UserService: userService}

	smsRepository := repository.NewGormMessageRepository()
	smsService := services.NewMessageService(smsRepository, userRepository)
	smsController := controllers.NewMessageController(userService, smsService)

	adminRepository := repository.NewGormAdminRepository()
	adminService := services.NewAdminService(adminRepository)
	adminController := controllers.NewAdminController(adminService)

	userGroup := e.Group("/sms/user")
	userGroup.POST("/signup", userController.Signup)
	userGroup.GET("/login", userController.Login)
	userGroup.GET("/me", userController.GetUserByToken)
	userGroup.GET("/buy", userController.GetPhoneNumbersToBuy)
	userGroup.POST("/logout", userController.LogOut)
	userGroup.PUT("/buy/:id", userController.BuyNumber)
	userGroup.PUT("/rent/:id", userController.PlaceRent)
	userGroup.DELETE("/dropRent/:id", userController.DropRent)
	userGroup.POST("/contacts/add", userController.AddContact)
	userGroup.DELETE("/:id/contacts/:contactID", userController.DeleteContact)
	userGroup.PUT("/:id/contacts/:contactID", userController.UpdateContact)
	userGroup.PUT("/main/:id", userController.SetMainNumber)

	paymentGroup := e.Group("/sms/payment")
	paymentGroup.GET("/pay/:amount", paymentController.AddBalance)
	paymentGroup.GET("/verify/", paymentController.VerifyPayment)

	phoneBookGroup := e.Group("/sms/phonebook")
	phoneBookGroup.POST("", phoneBookController.CreatePhoneBook)
	phoneBookGroup.GET("/:id", phoneBookController.GetPhoneBookByID)
	phoneBookGroup.PUT("/:id/:nid", phoneBookController.AddNumberToPhoneBook)
	phoneBookGroup.DELETE("/:id", phoneBookController.DeletePhoneBook)

	smsGroup := e.Group("/sms/send")
	smsGroup.POST("/single", smsController.SingleMessage)
	smsGroup.POST("/periodic", smsController.PeriodicMessage)

	adminGroup := e.Group("/sms/admin", middlewares.IsAdmin)
	adminGroup.POST("/addNumber", adminController.AddNumber)
	adminGroup.PUT("/suspend/:id", adminController.SuspendUser)
	adminGroup.PUT("/unsuspend/:id", adminController.UnSuspendUser)

	return &App{
		E: e,
	}
}

func (a *App) Start(addr string) error {
	a.E.Logger.Fatal(a.E.Start(addr))
	return nil
}
