package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	db "github.com/SOOA-swarch-2022ii/sooa_user_token_ms/db"
	"github.com/SOOA-swarch-2022ii/sooa_user_token_ms/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var Cliente_mongo = db.ConnectDB()
var UserDBhandler *mongo.Collection = Cliente_mongo.Database("SOOA_users").Collection("users")

func NewUser(response http.ResponseWriter, request *http.Request) {
	
	//create a new user from the request body and insert it into the database 
	var usuario models.User
	response.Header().Set("content-type", "application/json")
	err := json.NewDecoder(request.Body).Decode(&usuario)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"algo salio mal en el formato json de la asignatura:` + err.Error() + ` "}`))
		return
	}

	result, err := UserDBhandler.InsertOne(context.TODO(), usuario)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"algo salio mal en la base de datos mongo:` + err.Error() + ` "}`))
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	//user creation and database upload for user is done

	//create a token for the user and save it in the tokens database for the first time
	var token models.Token = CreateToken(usuario.Username, 3600)
	result_token, err_token := TokenDBhandler.InsertOne(context.TODO(), token)
	if err_token != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"algo salio mal en la base de datos mongo:` + err.Error() + ` "}`))
		return
	}
	result_token_marshal, _ := json.Marshal(result_token.InsertedID)

	fmt.Println(`usuario creado en:` + strings.Replace(string(res), `"`, ``, 2)+` con token:`+strings.Replace(string(result_token_marshal), `"`, ``, 2))

	response.WriteHeader(http.StatusOK)
	JSON_TOKEN:=`{"token":"`+token.Token+`"}`
	response.Write([]byte(JSON_TOKEN))

}

func GetUserByUN(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	//id, _ := primitive.ObjectIDFromHex(params["id"])
	var nombre string = params["username"]
	fmt.Println("buscando nombre " + nombre)
	var asignatura models.User
	filtro := bson.D{{"name", nombre}}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := UserDBhandler.FindOne(ctx, filtro).Decode(&asignatura)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"algo salio mal en Asigsbyid: ` + err.Error() + ` "}`))
		return
	}
	response.WriteHeader(http.StatusOK)
	//de json a struct golang
	json.NewEncoder(response).Encode(asignatura)
}
