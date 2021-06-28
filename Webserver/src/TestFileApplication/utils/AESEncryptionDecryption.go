package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	b64 "encoding/base64"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

// var secretKey = "dFoeiAXDQwfjuUHNBnLcpjR71S3mxq8O3mDd2AySYYq7WWyMg3E1WwOsyvwQ"

func EncryptToBase64(secretKey string, plainText string) string {
	salt := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err.Error())
	}

	key, iv := __DeriveKeyAndIv(secretKey, string(salt))

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	pad := __PKCS7Padding([]byte(plainText), block.BlockSize())
	ecb := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(pad))
	ecb.CryptBlocks(encrypted, pad)

	return b64.StdEncoding.EncodeToString([]byte("Salted__" + string(salt) + string(encrypted)))
}

func EncryptToBase64ShortLiveToken(secretKey string, plainText string) string {
	formatedDate := Format("M/D/YYYY hh:mm:ss pm", time.Now().UTC().Add(time.Minute*time.Duration(30)))
	plainText = plainText + "~#$" + formatedDate
	// fmt.Println("inputText : " + plainText)
	return EncryptToBase64(secretKey, plainText)
}

func EncryptToBase64ShortLiveTokenWithDuration(secretKey string, plainText string, durationInMinutes int) string {
	formatedDate := Format("M/D/YYYY hh:mm:ss pm", time.Now().UTC().Add(time.Minute*time.Duration(durationInMinutes)))
	plainText = plainText + "~#$" + formatedDate
	return EncryptToBase64(secretKey, plainText)
}

func DecryptBase64(secretKey string, base64EncryptedText string) string {
	ct, _ := b64.StdEncoding.DecodeString(base64EncryptedText)
	if len(ct) < 16 || string(ct[:8]) != "Salted__" {
		return ""
	}

	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			log.Println(err)
		}
	}()

	salt := ct[8:16]
	ct = ct[16:]
	key, iv := __DeriveKeyAndIv(secretKey, string(salt))

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	cbc := cipher.NewCBCDecrypter(block, []byte(iv))
	dst := make([]byte, len(ct))
	cbc.CryptBlocks(dst, ct)

	return string(__PKCS7Trimming(dst))
}

func DecryptBase64ShorLiveToken(secretKey string, base64EncryptedText string) (bool, string) {
	decryptString := DecryptBase64(secretKey, base64EncryptedText)
	dateFormat := "M/D/YYYY hh:mm:ss pm"
	if len(decryptString) > 0 {
		decryptStringArray := strings.Split(decryptString, "~#$")
		fmt.Println("descryptedString : " + decryptString)
		if len(decryptStringArray) > 1 {
			//check if duration is expire or not : duration for key = 20 min
			tokenTime, err := Parse(dateFormat, decryptStringArray[1])
			loc, _ := time.LoadLocation("UTC")
			tokenTime = tokenTime.In(loc)
			if err == nil {
				formatedDate := Format(dateFormat, time.Now().UTC().Add(time.Minute*time.Duration(0)))
				// fmt.Println("currentTime : " + formatedDate)
				currentTime, err := Parse(dateFormat, formatedDate)
				if err == nil {
					if tokenTime.After(currentTime) {
						return true, decryptStringArray[0]
					} else {
						return false, ""
					}
				}
			}
		}
	}
	return false, ""
}

func __PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func __PKCS7Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func __DeriveKeyAndIv(secretKey string, salt string) (string, string) {
	salted := ""
	dI := ""

	for len(salted) < 48 {
		md := md5.New()
		md.Write([]byte(dI + secretKey + salt))
		dM := md.Sum(nil)
		dI = string(dM[:16])
		salted = salted + dI
	}

	key := salted[0:32]
	iv := salted[32:48]

	return key, iv
}
