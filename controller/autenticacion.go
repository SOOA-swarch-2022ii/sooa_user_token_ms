package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/SOOA-swarch-2022ii/sooa_user_token_ms/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var TokenDBhandler *mongo.Collection = Cliente_mongo.Database("SOOA_users").Collection("tokens")

// login function that returns a token if the user is valid
func Login(response http.ResponseWriter, request *http.Request) {
	//print that service was called
	fmt.Println("login called")
	var login models.Login
	response.Header().Set("content-type", "application/json")
	err := json.NewDecoder(request.Body).Decode(&login)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"algo salio mal en login:` + err.Error() + ` "}`))
		return
	}
	//print login received
	fmt.Println(login)

	var usuario models.User
	filtro := bson.D{{"username", login.Username}}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	error_no_user := UserDBhandler.FindOne(ctx, filtro).Decode(&usuario)
	if error_no_user != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"usuario no encontrado para login"}`))
		return
	}
	//print usuario found
	fmt.Println(usuario)

	//hash login password and compare it to the one in the database
	if CheckPasswordHash(login.Password, usuario.Password) {
		//print that we get the same password
		fmt.Println("passwords match, creating token")
		//create token
		token := CreateToken(usuario.Username, 3600, usuario.Role)
		//print token
		fmt.Println(token)
		//insert token into the database
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := TokenDBhandler.InsertOne(ctx, token)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"algo salio mal en login token writer : ` + err.Error() + ` "}`))
			return
		}
		//return token
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(token)
	}else{
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(`{"error":"passwords do not match"}`))
	}

}

// create a token given a user username and a expiration date in seconds
func CreateToken(userun string, seconds int, role string) models.Token {

	var nuevoToken models.Token
	nuevoToken.User = userun
	//token expires at the present date plus the seconds given
	nuevoToken.Expires = time.Now().Add(time.Second * time.Duration(seconds)).Format("2006-01-02 15:04:05")
	//generate a random string for the token
	nuevoToken.Role = role
	nuevoToken.Token = randString(32)
	nuevoToken.Creation = time.Now().Format("2006-01-02 15:04:05")

	return nuevoToken
}

func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().Unix())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
