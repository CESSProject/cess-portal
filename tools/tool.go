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

var (
	Black        = string([]byte{27, 91, 57, 48, 109})
	Red          = string([]byte{27, 91, 57, 49, 109})
	Green        = string([]byte{27, 91, 57, 50, 109})
	Yellow       = string([]byte{27, 91, 57, 51, 109})
	Blue         = string([]byte{27, 91, 57, 52, 109})
	Magenta      = string([]byte{27, 91, 57, 53, 109})
	Cyan         = string([]byte{27, 91, 57, 54, 109})
	White        = string([]byte{27, 91, 57, 55, 59, 52, 48, 109})
	Reset        = string([]byte{27, 91, 48, 109})
	DisableColor = false
)

type Bar struct {
	percent int64
	cur     int64
	total   int64
	rate    string
	graph   string
}

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

func PostFileChunks(url, filepath string, params map[string]string) (status int, err error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	//for key, val := range params {
	//	m.WriteField(key, val)
	//}
	go func() {
		defer w.Close()
		defer m.Close()
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("Fail to open the file,error:%s", err)
			return
		}
		part, err := m.CreateFormFile("file", params["file"])
		if err != nil {
			fmt.Printf("Failed to create form file,error:%s", err)
			return
		}
		for key, val := range params {
			m.WriteField(key, val)
		}
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			fmt.Printf("Failed to send file chunks,error:%s\n", err)
			return
		}
	}()
	resp, err := http.Post(url, m.FormDataContentType(), r)
	if err != nil {
		return resp.StatusCode, err
	}
	return
}

func PostFile(url string, filepath string, params map[string]string) (status int, err error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", params["file"])
	if err != nil {
		return 0, err
	}
	src, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	_, err = io.Copy(part, src)
	if err != nil {
		return 0, err
	}
	for key, val := range params {
		writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return 0, err
	}
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return 0, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	// content, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	logger.ErrLogger.Sugar().Errorf("%v", err)
	// 	return  err
	// }
	// fmt.Println(string(content))
	return resp.StatusCode, err
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

func (bar *Bar) NewOption(start, total int64) {
	bar.cur = start
	bar.total = total
	if bar.graph == "" {
		bar.graph = "█"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph
	}
}

func (bar *Bar) getPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}

func (bar *Bar) Play(cur int64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\r[%-50s]%3d%%  %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}
func (bar *Bar) Finish() {
	fmt.Println()
}
