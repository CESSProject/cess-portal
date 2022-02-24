package client

import (
	"c-portal/conf"
	"c-portal/internal/chain"
	"c-portal/internal/logger"
	"c-portal/tools"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

func FileUpload(filepath, downloadfee string) {
	chain.Chain_Init()
	file, err := os.Stat(filepath)
	if err != nil {
		fmt.Printf("[Error]Please enter the correct file path!\n")
		return
	}
	if file.IsDir() {
		fmt.Printf("[Error]Please do not upload the folder!\n")
		return
	}
	cosfee, err := strconv.Atoi(downloadfee)
	if err != nil {
		fmt.Printf("[Error]Please enter a correct integer!\n")
		return
	}
	filehash, err := tools.CalcFileHash(filepath)
	if err != nil {
		fmt.Printf("[Error]There is a problem with the file, please replace it!\n")
		return
	}
	if file.Size()/100 == 0 {
		fmt.Printf("[Error]The upload file is too small!\n")
		return
	}
	fileid, err := tools.GetGuid(1)
	if err != nil {
		fmt.Printf("[Error]Create snowflake fail! error:%s\n", err)
		return
	}
	//todo:send file to scheduler
	uploadinfo := map[string]string{
		"file":      file.Name(),
		"fileid":    fileid,
		"backupnum": "3",
	}

	status, err := tools.PostFile("", filepath, uploadinfo)
	if err != nil {
		fmt.Printf("[Error]Post file to scheduler fail,error:\n", err)
		return
	}
	if status != 200 {
		fmt.Printf("[Error]Scheduler can not receive the file ,status code:%d\n", status)
		return
	}

	var ci chain.CessInfo
	filesize := new(big.Int)
	fee := new(big.Int)

	ci.IdentifyAccountPhrase = conf.ClientConf.ChainData.IdAccountPhraseOrSeed
	ci.TransactionName = chain.UploadFileTransactionName

	filesize.SetInt64(file.Size() / 100)
	fee.SetInt64(int64(cosfee))

	AsInBlock, err := ci.UploadFileMetaInformation(fileid, file.Name(), filehash, filesize, fee)
	if err != nil {
		fmt.Printf("[Error]Upload file meta information error:%s", err)
		return
	}
	fmt.Printf("Transaction chain block number is:%s\n", AsInBlock)
}

type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func FileDownload(fileid string) {
	chain.Chain_Init()

	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindFileChainModule
	ci.ChainModuleMethod = chain.FindFileModuleMethod[0]
	data, err := ci.GetFileInfo(fileid)
	if err != nil {
		fmt.Printf("[Error]Get file:%s info fail:%s\n", fileid, err)
		logger.OutPutLogger.Sugar().Infof("[Error]Get file:%s info fail:%s\n", fileid, err)
		return
	}

	var fd = struct {
		Hash     string `json:"hash"`
		FileName string `json:"filename"`
	}{
		string(data.Filehash),
		string(data.Filename),
	}
	resp, err := tools.Post("", fd)
	if err != nil {
		fmt.Printf("[Error]System error:%s\n", err)
		logger.OutPutLogger.Sugar().Infof("[Error]System error:%s\n", err)
		return
	}
	var res RespMsg
	err = json.Unmarshal(resp, &res)
	if err != nil {
		fmt.Println("unmarshal error:", err)
	}
	fmt.Println("The file address is:", res.Data)
}
