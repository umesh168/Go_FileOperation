// Package restapi implements mainly Parsing inputs to all APIs and building their json Responses using new/old methods in existing PricingModelDataProvider.go.
// ASSUMPTIONS
// graphName will not be maintain a static map  AUTH TOKEN to graphDB
// 	For each new client we will issue a jwt auth token and make a entry in this map
// 	Make sure tokens have no expiry for now….Plan a token
// In API wherever need to pass graphDB then use from that static  map

// Proposing local api level structs instead of reusing our existing structs because these input/output jsons will change DRASTICALLY

package restapi

import (
	"TestFileApplication/constants"
	"encoding/json"
	"io/ioutil"
	"os"
	"io"
	"net/http"
	"TestFileApplication/filereader"
	"strings"
	"time"
	"TestFileApplication/utils"
	"github.com/golang/glog"
	"log"
	"github.com/pkg/errors"
)

type FlieAPIHandler struct{}

type ResponseObject struct {
	Message string
	Result  string
	Status  string
}

var AuthTokenToGraphDB = map[string]string{
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJjbGllbnQiOiJ0ZXN0Y2xpZW50IiwiZXhwIjoxNjU2NDI2NDEwfQ.ZVYXQnGrTxwY7HIc_6_sd9vhpzmukrY2JwgrtMOoO1Q": "fileoperationgraph",
}

/*
	This method handle mainly routing of APIs
		 Parsing inputs to all APIs and building their json Responses
		using new/old methods in existing PricingModelDataProvider
*/
func (pah *FlieAPIHandler) FlieAPIHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		w.WriteHeader(200) // Success
		return
	}
	var responseObject = map[string]string{}
	glog.V(3).Info("RequestURI :  ", r.RequestURI)
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		return
	}

	isAuthorized, err := isAuthorized(w, r)
	if !isAuthorized {
		responseObject["Result"] = "[]"
		responseObject["Message"] = err
		responseObject["Status"] = "0"
	} else {

		switch endpoint(r.RequestURI) {
		case "resolvePriceModels":
			new(FlieAPIHandler).updateAndSaveFileData(w, r)

		default:
			new(FlieAPIHandler).invalidRequestResponse(w, r)
		}

	}

}

func(path *FlieAPIHandler) updateAndSaveFileData(w http.ResponseWriter, r *http.Request){
	glog.V(3).Info("updateAndSaveFileData [Started]", time.Now().Format("2006-01-02 15:04:05.000000"))
	authToken := r.Header["Authorized-Token"][0]
	var response = ResponseObject{}
	type resultObject struct {
	}
	type payload struct {
		FileName     string `json:"filename,"`
	}

	result := resultObject{}
	//TODO : jwt token auth integration
	//graphDB > need to save and update file info 
	_, isAuthTokenToGraphDBMappingPresent := getDBFromToken(authToken)
	if isAuthTokenToGraphDBMappingPresent {
		// get the body of our POST request
		// unmarshal this into a new payload struct
		// pass price model name to getCostFactors
		body, err := ioutil.ReadAll(r.Body)
		r.ParseMultipartForm(10 << 20)

		// Get handler for filename, size and headers
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			log.Println("Error Retrieving the File")
			log.Println(err)
			return
		}

		defer file.Close()
		
		// Create file
		dst, err := os.Create(constants.UPLOAD_FILE_PATH + handler.Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		if err != nil {
			errors.Wrap(err, "read body failed")
			glog.Error("read body failed: ", err)
		}

		var requestObject payload
		err = json.Unmarshal(body, &requestObject)
		if err != nil {
			errors.Wrap(err, "Invalid Body")
			glog.Error("getfile info : Parse body failed due to ", err)
		} else {
			
			// call file reader 
			if err != nil {
				glog.Error(err)
				return
			}
			// defer file.Close()
			filepathToStoreReceiedFile := dst.Name()
			fileRadear := filereader.FileReader{}
			fileRadear.FilePath = filepathToStoreReceiedFile
			fileRadear.Unzip()

			dataJSONBytes, _ := json.Marshal(result)
			response.Result = string(dataJSONBytes)
			// response.Status = statusCode
		}

	} else {
		response.Message = "Invalid token"
		response.Status = "0"
	}

	glog.V(4).Info(response)
	glog.V(3).Info("status : ", response.Status, " message=", response.Message)
	glog.V(3).Info("updateAndSaveFileData Response Method [Completed] : " + time.Now().Format("2006-01-02 15:04:05.000000"))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

/*
	invalidRequestResponse are this request whose end point is invalid OR not added yet
		for such request we are not processign and thowing error Requested URL is not valid
*/
func (pah *FlieAPIHandler) invalidRequestResponse(w http.ResponseWriter, r *http.Request) {
	var response = ResponseObject{}
	response.Message = "Requested URL is not valid"
	response.Status = "0"
	glog.V(4).Info(response)
	glog.V(3).Info("pricingHandler Response Method [Completed] : " + time.Now().Format("2006-01-02 15:04:05.000000"))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func endpoint(requestURI string) string {
	requestURIlist := strings.Split(requestURI, "/")
	return requestURIlist[len(requestURIlist)-1]
}

func getDBFromToken(authToken string) (string, bool) {
	glog.V(3).Info("getDBFromToken [Started]", time.Now().Format("2006-01-02 15:04:05.000000"))
	glog.V(3).Info("getDBFromToken inputs authToken=", authToken)

	// we are maintaining static AuthTokenToGraphDB map
	// every unique token is mapping to database name
	// 	example token ABC map to pggraph
	//			token XYZ map to icicigraph
	//				if valid token mapping present then it will retuen
	//				else it will return invalid token in response
	if grapgDB, isAuthTokenToGraphDBMappingPresent := AuthTokenToGraphDB[authToken]; isAuthTokenToGraphDBMappingPresent {
		glog.V(3).Info("result  grapgDB=", grapgDB, " isAuthTokenToGraphDBMappingPresent:", isAuthTokenToGraphDBMappingPresent)
		glog.V(3).Info("getDBFromToken [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))

		return grapgDB, isAuthTokenToGraphDBMappingPresent
	}

	glog.V(3).Info("result  grapgDB=", "", " isAuthTokenToGraphDBMappingPresent:", false)
	glog.V(3).Info("getDBFromToken [Completed]", time.Now().Format("2006-01-02 15:04:05.000000"))
	return "", false
}


func getAesDecryptedStringList(inputEncryptedList []string) []string {
	descruptedStringList := []string{}
	for _, encryptedString := range inputEncryptedList {
		descruptedStringList = append(descruptedStringList, utils.DecryptBase64(constants.AES_ENCRYPTION_KEY, encryptedString))
	}
	return descruptedStringList
}

func getShortLiveAesDecryptedStringList(inputEncryptedList []string) ([]string, bool) {
	descruptedStringList := []string{}

	for _, encryptedString := range inputEncryptedList {
		isValidToken, descruptedString := utils.DecryptBase64ShorLiveToken(constants.AES_ENCRYPTION_KEY, encryptedString)
		if isValidToken {
			descruptedStringList = append(descruptedStringList, descruptedString)
		} else {
			return []string{}, false
		}
	}
	return descruptedStringList, true
}

func getStringObject(object interface{}) string {
	objectBytes, _ := json.Marshal(object)
	stringObject := string(objectBytes)
	return stringObject
}
