package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var (
	ErrInvalidParam = errors.New("invalid param")

	ErrStatusCodes = map[error]int{
		ErrInvalidParam: http.StatusBadRequest,
	}
)

func main() {
	http.HandleFunc("/upload", handleUpload)
	http.ListenAndServe(":8989", nil)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	ver, err := strconv.Atoi(r.Header.Get("ver"))
	if err != nil {
		w.Header().Set("errmsg", ErrInvalidParam.Error())
		w.WriteHeader(ErrStatusCodes[ErrInvalidParam])
		return
	}

	fmt.Println(ver)
	fmt.Println(r.Header.Get("filetype"))
	fmt.Println(r.Header.Get("filemd5"))
	fmt.Println(r.Header.Get("filesize"))
	fmt.Println(r.Header.Get("rangestart"))
	fmt.Println(r.Header.Get("rangeend"))

	w.Header().Set("ver", r.Header.Get("ver"))
	w.Header().Set("rangestart", r.Header.Get("rangestart"))
	w.Header().Set("rangeend", r.Header.Get("rangeend"))

	rangeend, err := strconv.Atoi(r.Header.Get("rangeend"))
	if err != nil {
		w.Header().Set("errmsg", ErrInvalidParam.Error())
		w.WriteHeader(ErrStatusCodes[ErrInvalidParam])
		return
	}

	filesize, err := strconv.Atoi(r.Header.Get("filesize"))
	if err != nil {
		w.Header().Set("errmsg", ErrInvalidParam.Error())
		w.WriteHeader(ErrStatusCodes[ErrInvalidParam])
		return
	}

	if rangeend < filesize-1 {
		w.Header().Set("flag", "0")
	} else {
		w.Header().Set("flag", "1")
	}

}
