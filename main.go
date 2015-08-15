package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gansidui/cofile/filestore"
)

var (
	ErrInvalidParam = errors.New("invalid param")

	ErrStatusCodes = map[error]int{
		ErrInvalidParam: http.StatusBadRequest,
	}
)

// retcode: 0：继续上传，1：上传成功，2：上传失败
const (
	RetContinue = 0
	RetSucceed  = 1
	RetFailed   = 2
)

var cofilePath string = "./DATA"

func main() {
	os.MkdirAll(cofilePath, os.ModePerm)

	http.HandleFunc("/upload", handleUpload)
	http.ListenAndServe(":8989", nil)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	// 获取各个字段
	ver := r.Header.Get("ver")

	filesize, err := strconv.ParseInt(r.Header.Get("filesize"), 10, 64)
	if handleInvalidParam(w, err) {
		log.Println("invalid filesize")
		return
	}
	chunksize, err := strconv.ParseInt(r.Header.Get("chunksize"), 10, 64)
	if handleInvalidParam(w, err) {
		log.Println("invalid chunksize")
		return
	}
	offset, err := strconv.ParseInt(r.Header.Get("offset"), 10, 64)
	if handleInvalidParam(w, err) {
		log.Println("invalid offset")
		return
	}
	filetype := r.Header.Get("filetype")
	filemd5 := r.Header.Get("filemd5")

	// 检查chunksize是否等于ContentLength
	if int64(chunksize) != r.ContentLength {
		log.Printf("chunksize[%v] ContentLength[%v]\r\n", chunksize, r.ContentLength)
		handleInvalidParam(w, errors.New("chunksize not equal ContentLength"))
		return
	}

	log.Printf("{ver[%v] filetype[%v] filesize[%v] filemd5[%v] chunksize[%v] offset[%v]}\r\n",
		ver, filetype, filesize, filemd5, chunksize, offset)

	info := &filestore.FileInfo{}
	info.ID = filemd5
	info.Type = filetype
	info.Size = filesize
	info.Offset = offset

	// TODO 其他校验

	fs := filestore.NewFileStore(cofilePath)
	epOffset, isCompleted, err := fs.NewUpload(info)
	if err != nil {
		respond(w, ver, epOffset, RetFailed, http.StatusText(http.StatusInternalServerError))
		return
	}
	if isCompleted {
		respond(w, ver, epOffset, RetSucceed, "ok")
		return
	}
	if epOffset != offset {
		respond(w, ver, epOffset, RetContinue, "ok")
		return
	}

	n, isCompleted, err := fs.WriteChunk(filemd5, epOffset, r.Body)
	if err != nil {
		respond(w, ver, epOffset, RetFailed, http.StatusText(http.StatusInternalServerError))
		return
	}
	if isCompleted {
		respond(w, ver, epOffset, RetSucceed, "ok")
		return
	}

	respond(w, ver, epOffset+n, RetContinue, "ok")
}

func handleInvalidParam(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err)
		w.Header().Set("errmsg", ErrInvalidParam.Error())
		w.WriteHeader(ErrStatusCodes[ErrInvalidParam])
		return true
	}
	return false
}

func respond(w http.ResponseWriter, ver string, offset, retcode int64, errmsg string) {
	w.Header().Set("ver", ver)
	w.Header().Set("offset", strconv.FormatInt(offset, 10))
	w.Header().Set("retcode", strconv.FormatInt(retcode, 10))
	w.Header().Set("errmsg", errmsg)
}
