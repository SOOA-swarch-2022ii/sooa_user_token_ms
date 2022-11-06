package controllers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/SOOA-swarch-2022ii/sooa_user_token_ms/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var TokenDBhandler *mongo.Collection = Cliente_mongo.Database("SOOA_users").Collection("tokens")

func Login(response http.ResponseWriter, request *http.Request) {

}

// create a token given a user username and a expiration date in seconds
func CreateToken(userun string, seconds int) models.Token {
	
	var nuevoToken models.Token
	nuevoToken.User = userun
	//token expires at the present date plus the seconds given
	nuevoToken.Expires = time.Now().Add(time.Second*time.Duration(seconds)).Format("2006-01-02 15:04:05")
	//generate a random string for the token
	nuevoToken.Token = randString(32)
	nuevoToken.Creation = time.Now().Format("2006-01-02 15:04:05")

	return nuevoToken
}

func randString(length int,) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().Unix())
	b := make([]byte, length)
	for i := range b {
	  b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
  }
  
