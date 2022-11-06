package routes

import (
	control "github.com/SOOA-swarch-2022ii/sooa_user_token_ms/controller"
	"github.com/gorilla/mux"
)

// Routes -> define endpoints
func Routes() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/register-user/", control.NewUser).Methods("POST")
	//router.HandleFunc("/user/un={username}", control.GetUserByUN).Methods("GET")
	//router.HandleFunc("/user/login/{username}", control.Login).Methods("GET")


	/*router.HandleFunc("/subjects/name={name}", control.GetsbName).Methods("GET")

	router.HandleFunc("/subject/{id}", control.UpdateSB).Methods("PUT")
	router.HandleFunc("/course/{id}", control.UpdateCO).Methods("PUT")

	router.HandleFunc("/subject/{id}", control.DeleteSB).Methods("DELETE")
	router.HandleFunc("/course/{id}", control.DeleteCO).Methods("DELETE")*/

	return router
}
