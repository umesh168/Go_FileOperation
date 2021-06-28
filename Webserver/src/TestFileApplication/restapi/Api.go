package restapi

import (
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
)

var signingKey = []byte("kronfs*014#44$$$kjklsjs")

type Response struct {
	Result []map[string]string
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	//https://flaviocopes.com/golang-enable-cors/
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Cache-Control, X-CSRF-Token, X-Auth-Token, X-Requested-With, Authorization,Authorized-Token,Access-Control-Allow-Origin,Access-Control-Allow-Headers")
}

func isAuthorized(w http.ResponseWriter, r *http.Request) (bool, string) {
	glog.V(3).Info("isAuthorized [Started]", time.Now().Format("2006-01-02 15:04:05.000000"))
	if r.Header["Authorized-Token"] != nil {

		token, err := jwt.Parse(r.Header["Authorized-Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				glog.V(3).Info("parseToken failed")
				glog.V(3).Info("isAuthorized [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
				return nil, fmt.Errorf("There was an error")
			}
			return signingKey, nil
		})

		if err != nil {
			log.Printf("Unauthorized access token: %v", r.Header["Authorized-Token"][0])
			glog.V(3).Info("Unauthorized access token: ", r.Header["Authorized-Token"][0])
			glog.V(3).Info("isAuthorized [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
			w.WriteHeader(401) // Wrong password or username, Return 401.
			return false, err.Error()
		}

		if token.Valid {
			glog.V(3).Info("validation:", true)
			glog.V(3).Info("isAuthorized [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
			return true, ""
		}
	}
	glog.V(3).Info("validation: ", false, " : Something went wrong while parsing token")
	glog.V(3).Info("isAuthorized [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
	return false, "Something went wrong while parsing token"
}

func generateToken() string {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
	}
	return validToken
}

func GenerateJWT() (string, error) {
	glog.V(3).Info("GenerateJWT [Started]", time.Now().Format("2006-01-02 15:04:05.000000"))
	//Token will expire in 1 year
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "testclient"
	claims["exp"] = time.Now().Add(time.Hour * 8760).Unix()

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		glog.Error("Something Went Wrong: ", err.Error())
		glog.V(3).Info("GenerateJWT [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	glog.V(3).Info("tokenString=", tokenString)
	glog.V(3).Info("GenerateJWT [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
	return tokenString, nil
}
