package main

import (
	"bytes"
	"errors"
	"io"
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
	var offset int64 = 0

	for {
		n, err := file.Read(data)
		if err != nil {
			break
		}

		epOffset, isCompleted, err := postChunk("1.0", "txt", filemd5, filesize,
			int64(n), offset, bytes.NewReader(data[:n]))

		if err != nil {
			// for test
			log.Fatal(err)
			break
		}
		if isCompleted {
			log.Println("upload success")
			break
		}

		offset += int64(n)

		if offset != epOffset {
			file.Seek(epOffset, os.SEEK_SET)
			offset = epOffset
		}
	}
}

func postChunk(ver, filetype, filemd5 string,
	filesize, chunksize, offsest int64, body io.Reader) (int64, bool, error) {

	client := &http.Client{}
	req, err := http.NewRequest("post", "http://127.0.0.1:8989/upload", body)
	if err != nil {
		return 0, false, err
	}

	req.Header.Add("ver", ver)
	req.Header.Add("filetype", filetype)
	req.Header.Add("filemd5", filemd5)
	req.Header.Add("filesize", strconv.FormatInt(filesize, 10))
	req.Header.Add("chunksize", strconv.FormatInt(chunksize, 10))
	req.Header.Add("offset", strconv.FormatInt(offsest, 10))

	rsp, err := client.Do(req)
	if err != nil {
		return 0, false, err
	}
	defer rsp.Body.Close()

	log.Printf("{%v, %v, %v, %v, %v}\r\n", rsp.StatusCode, rsp.Header.Get("ver"),
		rsp.Header.Get("offset"), rsp.Header.Get("retcode"), rsp.Header.Get("errmsg"))

	epOffsest, err := strconv.ParseInt(rsp.Header.Get("offset"), 10, 64)
	if err != nil {
		return 0, false, err
	}

	retcode := rsp.Header.Get("retcode")
	if retcode == "0" {
		return epOffsest, false, nil

	} else if retcode == "1" {
		return epOffsest, true, nil

	} else {
		return epOffsest, false, errors.New(rsp.Header.Get("errmsg"))
	}

}
