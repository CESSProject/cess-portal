package client

import (
	"cess-portal/conf"
	"cess-portal/internal/chain"
	"cess-portal/internal/logger"
	"cess-portal/internal/rpc"
	"cess-portal/module"
	"cess-portal/tools"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func FileUpload(path, backups, PrivateKey string) {
	chain.Chain_Init()
	file, err := os.Stat(path)
	if err != nil {
		fmt.Printf("[Error]Please enter the correct file path!\n")
		return
	}

	if file.IsDir() {
		fmt.Printf("[Error]Please do not upload the folder!\n")
		return
	}

	spares, err := strconv.Atoi(backups)
	if err != nil {
		fmt.Printf("[Error]Please enter a correct integer!\n")
		return
	}

	filehash, err := tools.CalcFileHash(path)
	if err != nil {
		fmt.Printf("[Error]There is a problem with the file, please replace it!\n")
		return
	}

	if file.Size()/1024 == 0 {
		fmt.Printf("[Error]The upload file is too small!\n")
		return
	}

	fileid, err := tools.GetGuid(1)
	if err != nil {
		fmt.Printf("[Error]Create snowflake fail! error:%s\n", err)
		return
	}
	var blockinfo module.FileUploadInfo
	blockinfo.Backups = backups
	blockinfo.FileId = fileid
	blockinfo.BlockSize = int32(file.Size())
	blockinfo.FileHash = filehash

	block := make([]byte, 0)
	blocksize := 2048
	blocktotal := 0

	f, err := os.Open(path)
	if err != nil {
		fmt.Println("[Error]This file was broken! ", err)
		return
	}
	defer f.Close()
	filebyte, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("[Error]analyze this file error! ", err)
		return
	}

	//todo:get scheduler
	schedIp := ""

	commit := func(num int, data []byte) {
		blockinfo.BlockNum = int32(num) + 1
		blockinfo.Data = data
		info, err := proto.Marshal(&blockinfo)
		if err != nil {
			fmt.Println("[Error]Serialization error, please upload again! ", err)
			logger.OutPutLogger.Sugar().Infof("[Error]Serialization error, please upload again! ", err)
			return
		}
		var reqmsg rpc.ReqMsg
		reqmsg.Version = 0
		reqmsg.Body = info
		reqmsg.Method = module.UploadService
		reqmsg.Id = uint64(num)
		reqmsg.Service = module.CtlServiceName

		wsURL := "ws:" + schedIp
		client, err := rpc.DialWebsocket(context.Background(), wsURL, "")
		if err != nil {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := client.Call(ctx, &reqmsg)
		defer cancel()
		if err != nil {
			fmt.Printf("%s[Error]Failed to transfer file to scheduler,error:%s%s", tools.Red, err, tools.Reset)
			return
		}

		var res rpc.RespMsg
		err = proto.Unmarshal(resp.Body, &res)
		if err != nil {
			logger.OutPutLogger.Sugar().Infof("[Error]Error getting reply from schedule, transfer failed! ", err)
			return
		}
	}

	if len(PrivateKey) != 0 {
		os.Create(filepath.Join(conf.ClientConf.KeyInfo.KeyPath, (f.Name() + ".pem")))
		keyfile, err := os.OpenFile(filepath.Join(conf.ClientConf.KeyInfo.KeyPath, (f.Name()+".pem")), os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("%s[Error]:Failed to save key%s error:%s", tools.Red, tools.Reset, err)
			logger.OutPutLogger.Sugar().Infof("%s[Error]:Failed to save key%s error:%s", tools.Red, tools.Reset, err)
			return
		}
		_, err = keyfile.WriteString(PrivateKey)
		if err != nil {
			fmt.Printf("%s[Error]:Failed to write key to file:%s%s error:%s", tools.Red, filepath.Join(conf.ClientConf.KeyInfo.KeyPath, (f.Name()+".pem")), tools.Reset, err)
			logger.OutPutLogger.Sugar().Infof("%s[Error]:Failed to write key to file:%s%s error:%s", tools.Red, filepath.Join(conf.ClientConf.KeyInfo.KeyPath, (f.Name()+".pem")), tools.Reset, err)
			return
		}

		decodefile, err := tools.AesEncrypt(filebyte, []byte(PrivateKey))
		if err != nil {
			fmt.Println("[Error]Encode the file fail ,error! ", err)
			return
		}
		blocks := len(decodefile) / blocksize
		if len(decodefile)%blocksize == 0 {
			blocktotal = blocks
		} else {
			blocktotal = blocks + 1
		}
		blockinfo.Blocks = int32(blocktotal)
		var bar tools.Bar
		bar.NewOption(0, int64(blocktotal))
		for i := 0; i < blocks; i++ {
			if blocks != i {
				block = decodefile[i*blocksize : (i+1)*blocksize]
				bar.Play(int64(i))
			} else {
				block = decodefile[i*blocksize:]
				bar.Play(int64(i + 1))
			}
			commit(i, block)
		}
	} else {
		fmt.Printf("%s[tips]:upload file:%s without private key%s", tools.Yellow, path, tools.Reset)
		blocks := len(filebyte) / blocksize
		if len(filebyte)%blocksize == 0 {
			blocktotal = blocks
		} else {
			blocktotal = blocks + 1
		}
		blockinfo.Blocks = int32(blocktotal)
		var bar tools.Bar
		bar.NewOption(0, int64(blocktotal))
		for i := 0; i < blocks; i++ {
			if blocks != i {
				block = filebyte[i*blocksize : (i+1)*blocksize]
				bar.Play(int64(i))
			} else {
				block = filebyte[i*blocksize:]
				bar.Play(int64(i + 1))
			}
			commit(i, block)
		}
	}

	var ci chain.CessInfo
	filesize := new(big.Int)
	fee := new(big.Int)

	ci.IdentifyAccountPhrase = conf.ClientConf.ChainData.IdAccountPhraseOrSeed
	ci.TransactionName = chain.UploadFileTransactionName

	filesize.SetInt64(file.Size() / 1024)
	fee.SetInt64(int64(0))

	AsInBlock, err := ci.UploadFileMetaInformation(fileid, file.Name(), filehash, uint8(spares), filesize, fee)
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

func FileDelete(fileid string) {
	chain.Chain_Init()
	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.IdentifyAccountPhrase = conf.ClientConf.ChainData.IdAccountPhraseOrSeed
	ci.TransactionName = chain.DeleteFileTransactionName

	err := ci.DeleteFileOnChain(fileid)
	if err != nil {
		fmt.Printf("%s[Error]Delete file error:%s%s\n", tools.Red, tools.Reset, err)
		logger.OutPutLogger.Sugar().Infof("%s[Error]Delete file error:%s%s\n", tools.Red, tools.Reset, err)
		return
	} else {
		fmt.Printf("%s[Error]Delete fileid:%s successful!:%s\n", tools.Green, fileid, tools.Reset)
		logger.OutPutLogger.Sugar().Infof("%s[Error]Delete fileid:%s successful!:%s\n", tools.Green, fileid, tools.Reset)
		return
	}

}
