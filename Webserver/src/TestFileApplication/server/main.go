package main

import(
	"net/http"
	"github.com/golang/glog"	
	"TestFileApplication/constants"
	"TestFileApplication/restapi"
	"log"
	"time"
)

func main() {
	// log.Println(restapi.GenerateJWT())
	startTime := time.Now().Format("2006-01-02 15:04:05.000000")
	http.HandleFunc("/pricingService/", new(restapi.FlieAPIHandler).FlieAPIHandler)
	glog.V(3).Info("Listening on port" + constants.WEBSERVER_PORT)
	log.Println("Listening on port" + constants.WEBSERVER_PORT)
	log.Println("server started init :",startTime)
	log.Println("server end init :",time.Now().Format("2006-01-02 15:04:05.000000"))
	http.ListenAndServe(":"+constants.WEBSERVER_PORT, nil)
	
	// a. Create a private key-
	// openssl genrsa -des3 -out server.key 2048
	// b. Remove its passphrase-
	// openssl rsa -in server.key -out server.key
	// c. Create a CSR (Certificate Signing Request):
	// openssl req -new -key server.key -out server.csr
	// d. Use this CSR to obtain a valid certificate from a certificate authority or generate a self-signed certificate with the following command.
	// openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
}
