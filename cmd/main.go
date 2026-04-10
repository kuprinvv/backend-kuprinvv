// @title          API бронирования переговорок
// @version        1.0
// @description    Сервис для управления переговорками, расписаниями и бронями.
//
// @host           localhost:8080
// @BasePath       /
//
// @securityDefinitions.apikey BearerAuth
// @in             header
// @name           Authorization
// @description    JWT-токен в формате: Bearer <token>
package main

import (
	"log"
	"test-backend-1-kuprinvv/internal/app"

	_ "test-backend-1-kuprinvv/docs"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
