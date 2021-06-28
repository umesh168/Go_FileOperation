package fileprocessor

import(
	"io/ioutil"
	"os"
	"log"
)

//TextFileReader is struct which hold FilePath string object
type TextFileReader struct{
	FilePath string
}

/*
	Process: process method will read file from reciever TextFileReader
				get byte data from ReadAll method 
				if erroor == nil then then return nil error and string data read from file 
				else return error and set data as empty string  
*/
func (jfd *TextFileReader) Process()(error, string){
	file, err := os.Open(jfd.FilePath)	// open file from given path
	
	if err != nil {
        log.Println(os.Stderr, "Failed to open file:", err)	// thow/return error if invalid path 
        return err,""
    }

	byteValue, err := ioutil.ReadAll(file)	// ReadAll method reads error or EOF and returns the data it read
	if err == nil{
		log.Println("got error while reading file", err)
		return nil,string(byteValue)
	} else {
		log.Println("got error while reading file", err)
		return err,""
	}
}