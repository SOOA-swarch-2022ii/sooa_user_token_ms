package main

import (
	"fmt"
	"github.com/SOOA-swarch-2022ii/sooa_user_token_ms/routes"
	"net/http"
)

func main() {
	fmt.Println("Inicializando microservicio de usuarios")
	enrutador := routes.Routes()
	http.ListenAndServe(":6665", enrutador)
}