package fileprocessor


import (
	"log"
	"testing"
	"time"
	"fmt"
)

/*
	ProcessFile Testcase
	where you can pass test of diff type of file data
	file processor call process method according to inerface it recieved
*/
func TestProcessFile(t *testing.T) {
	startTime := time.Now().Format("2006-01-02 15:04:05.000000")

	fileData1 := JSONFileReader{
		FilePath:"test.json",
	}

	fileData2 := TextFileReader{
		FilePath:"TextFile.txt",
	}
	fileProcessorList := []FileProcessor{&fileData1,&fileData2}

	ProcessFile(fileProcessorList)
	log.Println("Response Method [Start] : " + startTime)
	log.Println("Response Method [Completed] : " + time.Now().Format("2006-01-02 15:04:05.000000"))
}


/*
	Process Testcase
	here we using diff file with table driven test case
	where we pass diff tast file and run then in one shot 
	for loop will execute all test cases data and print response according response of test case
*/
func TestProcess(t *testing.T) {
	startTime := time.Now().Format("2006-01-02 15:04:05.000000")

	fileData1 := JSONFileReader{
		FilePath:"test.json",
	}

	fileData2 := TextFileReader{
		FilePath:"TextFile.txt",
	}

	fileProcessorList := []FileProcessor{&fileData1,&fileData2}

	for index, tt := range fileProcessorList {
        testname := fmt.Sprintf("%d", index)
        t.Run(testname, func(t *testing.T) {
            response := tt.Process()
            if response != nil {
                t.Error("Processing error :", response)
            }
        })
    }

	log.Println("TestProcess Method [Start] : " + startTime)
	log.Println("TestProcess Method [Completed] : " + time.Now().Format("2006-01-02 15:04:05.000000"))
}