package tools

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/bwmarrin/snowflake"
	"golang.org/x/crypto/blake2b"
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
var (
	SSPrefix            = []byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45}
	PolkadotPrefix      = []byte{0x00}
	KsmPrefix           = []byte{0x02}
	KatalPrefix         = []byte{0x04}
	PlasmPrefix         = []byte{0x05}
	BifrostPrefix       = []byte{0x06}
	EdgewarePrefix      = []byte{0x07}
	KaruraPrefix        = []byte{0x08}
	ReynoldsPrefix      = []byte{0x09}
	AcalaPrefix         = []byte{0x0a}
	LaminarPrefix       = []byte{0x0b}
	PolymathPrefix      = []byte{0x0c}
	SubstraTEEPrefix    = []byte{0x0d}
	KulupuPrefix        = []byte{0x10}
	DarkPrefix          = []byte{0x11}
	DarwiniaPrefix      = []byte{0x12}
	StafiPrefix         = []byte{0x14}
	DockTestNetPrefix   = []byte{0x15}
	DockMainNetPrefix   = []byte{0x16}
	ShiftNrgPrefix      = []byte{0x17}
	SubsocialPrefix     = []byte{0x1c}
	PhalaPrefix         = []byte{0x1e}
	RobonomicsPrefix    = []byte{0x20}
	DataHighwayPrefix   = []byte{0x21}
	CentrifugePrefix    = []byte{0x24}
	MathMainPrefix      = []byte{0x27}
	MathTestPrefix      = []byte{0x28}
	SubstratePrefix     = []byte{0x2a}
	ChainXPrefix        = []byte{0x2c}
	ChainCessTestPrefix = []byte{0x50, 0xac}
)

type Bar struct {
	percent   int64
	cur       int64
	total     int64
	rate      string
	graph     string
	calibrate float64
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

//Get file unique identifier
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
	if bar.percent != last {
		add := (float64(int(bar.percent-last)) / 100) * 50
		bar.calibrate += add - float64(int(add))
		for i := 0; i < int(add); i++ {
			bar.rate += bar.graph
		}
		if int(bar.calibrate) > 0 {
			for i := 0; i < int(bar.calibrate); i++ {
				bar.rate += bar.graph
			}
			bar.calibrate = 0
		}
	}
	fmt.Printf("\r[%-50s]%3d%%  %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}
func (bar *Bar) Finish() {
	fmt.Println()
}

//prefix: chain.SubstratePrefix
func EncodeByPubHex(publicHex string, prefix []byte) (string, error) {
	publicKeyHash, err := hex.DecodeString(publicHex)
	if err != nil {
		return "", err
	}
	return Encode(publicKeyHash, prefix)
}

func DecodeToPub(address string, prefix []byte) ([]byte, error) {
	err := VerityAddress(address, prefix)
	if err != nil {
		return nil, errors.New("Invalid addrss")
	}
	data := base58.Decode(address)
	if len(data) != (34 + len(prefix)) {
		return nil, errors.New("base58 decode error")
	}
	return data[len(prefix) : len(data)-2], nil
}

func PubBytesTo0XString(pub []byte) string {
	return fmt.Sprintf("%#x", pub)
}

func PubBytesToString(b []byte) string {
	s := ""
	for i := 0; i < len(b); i++ {
		tmp := fmt.Sprintf("%#02x", b[i])
		s += tmp[2:]
	}
	return s
}

func Encode(publicKeyHash []byte, prefix []byte) (string, error) {
	if len(publicKeyHash) != 32 {
		return "", errors.New("public hash length is not equal 32")
	}
	payload := appendBytes(prefix, publicKeyHash)
	input := appendBytes(SSPrefix, payload)
	ck := blake2b.Sum512(input)
	checkum := ck[:2]
	address := base58.Encode(appendBytes(payload, checkum))
	if address == "" {
		return address, errors.New("base58 encode error")
	}
	return address, nil
}

func appendBytes(data1, data2 []byte) []byte {
	if data2 == nil {
		return data1
	}
	return append(data1, data2...)
}

func VerityAddress(address string, prefix []byte) error {
	decodeBytes := base58.Decode(address)
	if len(decodeBytes) != (34 + len(prefix)) {
		return errors.New("base58 decode error")
	}
	if decodeBytes[0] != prefix[0] {
		return errors.New("prefix valid error")
	}
	pub := decodeBytes[len(prefix) : len(decodeBytes)-2]

	data := append(prefix, pub...)
	input := append(SSPrefix, data...)
	ck := blake2b.Sum512(input)
	checkSum := ck[:2]
	for i := 0; i < 2; i++ {
		if checkSum[i] != decodeBytes[32+len(prefix)+i] {
			return errors.New("checksum valid error")
		}
	}
	if len(pub) != 32 {
		return errors.New("decode public key length is not equal 32")
	}
	return nil
}
