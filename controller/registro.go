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
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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

	//hash the password given
	hashedPassword, err := HashPassword(usuario.Password)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"error hashing":"` + err.Error() + ` "}`))
		return
	}
	usuario.Password = hashedPassword
	//print the user
	fmt.Println(usuario)

	//inserta o actualiza el usuario en la base de datos, basado en el id, si no existe lo crea
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //contexto de tiempo de espera de 10 segundos para la conexion
	//result, err := UserDBhandler.InsertOne(ctx, usuario)


	result_register_user, err := UserDBhandler.ReplaceOne(ctx, bson.M{"username": usuario.Username}, usuario, options.Replace().SetUpsert(true))
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"algo salio mal en la base de datos mongo:` + err.Error() + ` "}`))
		return
	}
	result_reg_user_marshal, _ := json.Marshal(result_register_user.UpsertedID)

	fmt.Println(`usuario creado en:` + strings.Replace(string(result_reg_user_marshal), `"`, ``, 2))

	//create a token for the user and save it in the tokens database for the first time for 2 hours
	var token models.Token = CreateToken(usuario.Username, 7200, usuario.Role)
	result_token, err_token := TokenDBhandler.InsertOne(context.TODO(), token)
	if err_token != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"algo salio mal en la base de datos de tokens:` + err.Error() + ` "}`))
		return
	}
	result_token_marshal, _ := json.Marshal(result_token.InsertedID)

	fmt.Println(`usuario creado en:` + strings.Replace(string(result_token_marshal), `"`, ``, 2) + ` con token:` + strings.Replace(string(result_token_marshal), `"`, ``, 2))

	response.WriteHeader(http.StatusOK)
	JSON_TOKEN := `{"token":"` + token.Token + `"}`
	response.Write([]byte(JSON_TOKEN))

}

func GetUserByUN(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	//id, _ := primitive.ObjectIDFromHex(params["id"])
	var nombre string = params["username"]
	fmt.Println("buscando nombre " + nombre)
	var usuario models.User
	filtro := bson.D{{"username", nombre}}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := UserDBhandler.FindOne(ctx, filtro).Decode(&usuario)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"algo salio mal en GetUserByUN : ` + err.Error() + ` "}`))
		return
	}
	response.WriteHeader(http.StatusOK)
	//de json a struct golang
	//eliminate password field
	usuario.Password = ""
	json.NewEncoder(response).Encode(usuario)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
