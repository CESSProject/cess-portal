package tools

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func Post(url string, para interface{}) ([]byte, error) {
	body, err := json.Marshal(para)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	var resp = new(http.Response)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return respBody, err
	}
	return nil, err
}

func PostFile(url, filepath string) (status int, err error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("Fail to open the file,error:%s", err)
			return
		}
		part, err := m.CreateFormFile("File", file.Name())
		if err != nil {
			fmt.Printf("Failed to create form file,error:%s", err)
			return
		}
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			fmt.Printf("Failed to send file chunks,error:%s", err)
			return
		}
	}()
	resp, err := http.Post(url, m.FormDataContentType(), r)
	if err != nil {
		return resp.StatusCode, err
	}
	return
}

func CalcFileHash(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func GetGuid(num int64) (string, error) {
	node, err := snowflake.NewNode(num)
	if err != nil {
		return "", err
	}

	id := node.Generate()
	return id.String(), nil
}
