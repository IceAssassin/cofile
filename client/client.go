package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gansidui/cofile/utils"
)

func main() {
	testfile := "testdata.txt"

	file, err := os.Open(testfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	filemd5, filesize := utils.GetMd5FromFile(testfile)
	data := make([]byte, 256)
	offset := 0

	for {
		n, err := file.Read(data)
		if err != nil {
			break
		}
		if !postChunk(filemd5, filesize, offset, offset+n-1, data) {
			break
		}
		offset += n
	}
}

func postChunk(filemd5 string, filesize, rangestart, rangeend int, data []byte) bool {
	client := &http.Client{}
	body := bytes.NewReader(data)

	req, err := http.NewRequest("post", "http://127.0.0.1:8989/upload", body)
	if err != nil {
		log.Println(err)
		return false
	}

	req.Header.Add("ver", strconv.Itoa(0))
	req.Header.Add("filetype", "txt")
	req.Header.Add("filemd5", filemd5)
	req.Header.Add("filesize", strconv.Itoa(filesize))
	req.Header.Add("rangestart", strconv.Itoa(rangestart))
	req.Header.Add("rangeend", strconv.Itoa(rangeend))

	rsp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	defer rsp.Body.Close()

	fmt.Println(rsp.Header.Get("flag"))
	fmt.Println(rsp.Header.Get("errmsg"))
	fmt.Println(rsp.StatusCode)

	return true
}
