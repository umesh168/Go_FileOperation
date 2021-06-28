package filereader

import(
	"archive/zip"
	// "crypto/aes"
	"os"
	"TestFileApplication/dbcallhelper"
	"log"
	"io"
	"strings"
	"path/filepath"
	"fmt"
	"io/ioutil"
	"TestFileApplication/constants"
	"TestFileApplication/model"
	"TestFileApplication/utils"
)

type FileReader struct{
	FilePath string
	Destiation string
	DestiationFolderPath string
}

func (fr *FileReader) ProcessFile(graphDB string) (error){
	fileNode := model.Node{}
	err := fr.Unzip()
	if err != nil{
		log.Println(err)
		return err
	} else{
		//add file data in db
		// dbcallhelper.CrateNode()
		var files []string

		root := filepath.FromSlash(fr.DestiationFolderPath)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			panic(err)
		}
		fileNode.Name = filepath.Base(fr.DestiationFolderPath)
		fileNode.Id = dbcallhelper.GetNodeIdAsPerProductConvention(graphDB, fileNode.Name,"fileName")
		dbcallhelper.CrateNode(graphDB,fileNode)
		for index, file := range files {
			if index == 0{
				continue
			}
			fr.EncryptFileData(file)
			fr.RenameEncryptedFile(file)
		}
		fr.ZipWriter()
		
	}
	return nil
}

func (fr *FileReader) Unzip() (error){
	r, err := zip.OpenReader(fr.FilePath)
    if err != nil {
		log.Println(err)
        return err
    }
    defer func() {
        if err := r.Close(); err != nil {
            panic(err)
        }
    }()

	os.MkdirAll(fr.Destiation, 0755)
	rawFolderName := filepath.Base(fr.FilePath)
	rawFolderName = strings.Split(rawFolderName,filepath.Ext(fr.FilePath))[0]
	fr.DestiationFolderPath = fr.Destiation+"\\" + rawFolderName
    // Closure to address file descriptors issue with all the deferred .Close() methods
    extractAndWriteFile := func(f *zip.File) error {
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer func() {
            if err := rc.Close(); err != nil {
                panic(err)
            }
        }()

        path := filepath.Join(fr.Destiation, f.Name)
        // Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(fr.Destiation) + string(os.PathSeparator)) {
			//TODO : use log instead of fmt
            return fmt.Errorf("illegal file path: %s", path)
        }

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
        } else {
            os.MkdirAll(filepath.Dir(path), f.Mode())
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer func() {
                if err := f.Close(); err != nil {
                    panic(err)
                }
            }()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
        return nil
    }

    for _, f := range r.File {
        err := extractAndWriteFile(f)
        if err != nil {
            return err
        }
    }

    return nil
}

func (fr *FileReader)ZipWriter() {
	file, err := os.Create(fr.DestiationFolderPath+"_encrypted.zip")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    w := zip.NewWriter(file)
    defer w.Close()

    walker := func(path string, info os.FileInfo, err error) error {
        fmt.Printf("Crawling: %#v\n", path)
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        // Ensure that `path` is not absolute; it should not start with "/".
        // This snippet happens to work because I don't use 
        // absolute paths, but ensure your real-world code 
        // transforms path into a zip-root relative path.
        f, err := w.Create(path)
        if err != nil {
            return err
        }

        _, err = io.Copy(f, file)
        if err != nil {
            return err
        }

        return nil
    }
	root := filepath.FromSlash(fr.DestiationFolderPath)
    err = filepath.Walk(root, walker)
    if err != nil {
        panic(err)
    }
}
func (fr *FileReader) EncryptFileData(rawfilepath string) error{
	// Read Write Mode
	file, err := os.OpenFile(rawfilepath, os.O_RDWR, 0644)
	
	defer file.Close()
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	byteValue, err := ioutil.ReadAll(file)
	encryptedData := utils.EncryptToBase64(constants.SecretKey,string(byteValue))
	_, err = file.WriteAt([]byte(encryptedData), 0) // Write at 0 beginning
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	
    return nil
}

func (fr *FileReader) RenameEncryptedFile(rawfilepath string) error{
	// Read Write Mode
	rawFileName := filepath.Base(rawfilepath)
	rawFileName = strings.Split(rawFileName,filepath.Ext(rawfilepath))[0]
	newName := strings.ReplaceAll(rawfilepath,rawFileName,rawFileName+ "_encrypted")
    err := os.Rename(filepath.Join(rawfilepath), filepath.Join(newName))

	if err != nil {
		log.Println(err)
	}
    return nil
}