package fileprocessor

type FileProcessor interface{
	Process() error
}

/*
	ProcessFile function take []FileProcessor as input
	call resoective process methos acording to file procesor  
*/
func ProcessFile(fp []FileProcessor){
	for _,reader := range fp{
		reader.Process()
	}
}

