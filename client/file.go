package client

import (
	"cess-portal/conf"
	"cess-portal/internal/chain"
	"cess-portal/internal/logger"
	"cess-portal/internal/rpc"
	"cess-portal/module"
	"cess-portal/tools"
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"sync"
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
	wsURL := "ws://" + schedIp
	client, err := rpc.DialWebsocket(context.Background(), wsURL, "")
	if err != nil {
		fmt.Println("[Error]dialog with websocket fail:! ", err)
		return
	}
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
		reqmsg.Body = info
		reqmsg.Method = module.UploadService
		reqmsg.Service = module.CtlServiceName

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := client.Call(ctx, &reqmsg)
		defer cancel()
		if err != nil {
			fmt.Printf("\n%s[Error]Failed to transfer file to scheduler,error:%s%s", tools.Red, err, tools.Reset)
			logger.OutPutLogger.Sugar().Infof("%s[Error]Failed to transfer file to scheduler,error:%s%s", tools.Red, err, tools.Reset)
			os.Exit(conf.Exit_SystemErr)
		}

		var res rpc.Err
		err = proto.Unmarshal(resp.Body, &res)
		if err != nil {
			fmt.Printf("\n[Error]Error getting reply from schedule, transfer failed! ", err)
			logger.OutPutLogger.Sugar().Infof("[Error]Error getting reply from schedule, transfer failed! ", err)
			os.Exit(conf.Exit_SystemErr)
		}
		if res.Code != 0 {
			fmt.Printf("\n[Error]Upload file fail!scheduler problem! ")
			logger.OutPutLogger.Sugar().Infof("[Error]Upload file fail!scheduler problem! ")
			os.Exit(conf.Exit_SystemErr)
		}
	}

	if len(PrivateKey) != 0 {
		os.Create(f.Name() + ".pem")
		keyfile, err := os.OpenFile(f.Name()+".pem", os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("%s[Error]:Failed to save key%s error:%s", tools.Red, tools.Reset, err)
			logger.OutPutLogger.Sugar().Infof("%s[Error]:Failed to save key%s error:%s", tools.Red, tools.Reset, err)
			return
		}
		_, err = keyfile.WriteString(PrivateKey)
		if err != nil {
			fmt.Printf("%s[Error]:Failed to write key to file:%s%s error:%s", tools.Red, filepath.Join(conf.ClientConf.PathInfo.KeyPath, (f.Name()+".pem")), tools.Reset, err)
			logger.OutPutLogger.Sugar().Infof("%s[Error]:Failed to write key to file:%s%s error:%s", tools.Red, filepath.Join(conf.ClientConf.PathInfo.KeyPath, (f.Name()+".pem")), tools.Reset, err)
			return
		}

		encodefile, err := tools.AesEncrypt(filebyte, []byte(PrivateKey))
		if err != nil {
			fmt.Println("[Error]Encode the file fail ,error! ", err)
			return
		}
		blocks := len(encodefile) / blocksize
		fmt.Println(len(encodefile))
		if len(encodefile)%blocksize == 0 {
			blocktotal = blocks
		} else {
			blocktotal = blocks + 1
		}
		blockinfo.Blocks = int32(blocktotal)
		var bar tools.Bar
		bar.NewOption(0, int64(blocktotal))
		for i := 0; i < blocktotal; i++ {
			block := make([]byte, 0)
			if blocks != i {
				block = encodefile[i*blocksize : (i+1)*blocksize]
				bar.Play(int64(i + 1))
			} else {
				block = encodefile[i*blocksize:]
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
		for i := 0; i < blocktotal; i++ {
			block := make([]byte, 0)
			if blocks != i {
				block = filebyte[i*blocksize : (i+1)*blocksize]
				bar.Play(int64(i + 1))
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
		fmt.Printf("\n[Error]Upload file meta information error:%s", err)
		return
	}
	fmt.Printf("\nTransaction chain block number is:%s\n", AsInBlock)
}

func FileDownload(fileid string) {
	chain.Chain_Init()
	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindFileChainModule
	ci.ChainModuleMethod = chain.FindFileModuleMethod[0]
	fileinfo, err := ci.GetFileInfo(fileid)
	if err != nil {
		fmt.Printf("[Error]Get file:%s info fail:%s\n", fileid, err)
		logger.OutPutLogger.Sugar().Infof("[Error]Get file:%s info fail:%s\n", fileid, err)
		return
	}

	_, err = os.Stat(conf.ClientConf.PathInfo.InstallPath)
	if err != nil {
		err = os.Mkdir(conf.ClientConf.PathInfo.InstallPath, os.ModePerm)
		if err != nil {
			fmt.Printf("[Error]Create install path error :%s\n", err)
			logger.OutPutLogger.Sugar().Infof("[Error]Create install path error :%s\n", err)
			os.Exit(conf.Exit_SystemErr)
		}
	}
	_, err = os.Create(filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.Filename[:])))
	if err != nil {
		fmt.Printf("[Error]Create installed file error :%s\n", err)
		logger.OutPutLogger.Sugar().Infof("[Error]Create installed file error :%s\n", err)
		os.Exit(conf.Exit_SystemErr)
	}
	installfile, err := os.OpenFile(filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.Filename[:])), os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("%s[Error]:Failed to save key%s error:%s", tools.Red, tools.Reset, err)
		logger.OutPutLogger.Sugar().Infof("%s[Error]:Failed to save key%s error:%s", tools.Red, tools.Reset, err)
		return
	}
	defer installfile.Close()

	//todo:get scheduler
	schedIp := ""
	wsURL := "ws:" + schedIp
	client, err := rpc.DialWebsocket(context.Background(), wsURL, "")
	if err != nil {
		return
	}

	var wantfile module.FileDownloadReq
	var bar tools.Bar
	var getAllBar sync.Once
	wantfile.FileId = fileid
	wantfile.WalletAddress = conf.ClientConf.ChainData.WalletAddress
	wantfile.Blocks = 1
	for {
		data, err := proto.Marshal(&wantfile)
		if err != nil {
			fmt.Printf("[Error]Marshal req file error:%s\n", err)
			logger.OutPutLogger.Sugar().Infof("[Error]Marshal req file error:%s\n", err)
			return
		}
		req := &rpc.ReqMsg{
			Method:  module.DownloadService,
			Service: module.CtlServiceName,
			Body:    data,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resp, err := client.Call(ctx, req)
		cancel()
		if err != nil {
			fmt.Printf("[Error]Download file fail error:%s\n", err)
			logger.OutPutLogger.Sugar().Infof("[Error]Download file fail error:%s\n", err)
			return
		}
		if resp.Body == nil {
			break
		}
		var blockData module.FileDownloadInfo
		err = proto.Unmarshal(resp.Body, &blockData)
		if err != nil {
			fmt.Printf("[Error]Unmarshal resp file error:%s\n", err)
			logger.OutPutLogger.Sugar().Infof("[Error]Unmarshal resp file error:%s\n", err)
			return
		}
		_, err = installfile.Write(blockData.Data)
		if err != nil {
			fmt.Printf("%s[Error]:Failed to write file's block to file:%s%s error:%s", tools.Red, filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.Filename[:])), tools.Reset, err)
			logger.OutPutLogger.Sugar().Infof("%s[Error]:Failed to write file's block to file:%s%s error:%s", tools.Red, filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.Filename[:])), tools.Reset, err)
			return
		}

		getAllBar.Do(func() {
			bar.NewOption(0, int64(blockData.Blocks))
		})
		bar.Play(int64(blockData.BlockNum))
		wantfile.Blocks++
		if blockData.Blocks == blockData.BlockNum {
			break
		}
	}
	fmt.Printf("%s[OK]:File '%s' has been downloaded to the directory :%s%s", tools.Green, string(fileinfo.Filename), filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.Filename[:])), tools.Reset)
	logger.OutPutLogger.Sugar().Infof("%s[OK]:File '%s' has been downloaded to the directory :%s%s", tools.Green, string(fileinfo.Filename), filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.Filename[:])), tools.Reset)

	fmt.Printf("%s[Warm]This is a private file, please enter the file password:%s", tools.Green, tools.Reset)
	filePWD := ""
	fmt.Scanln(&filePWD)
	encodefile, err := ioutil.ReadAll(installfile)
	if err != nil {
		fmt.Printf("%s[Error]:Decode file:%s fail%s error:%s", tools.Red, filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.Filename[:])), tools.Reset, err)
		logger.OutPutLogger.Sugar().Infof("%s[Error]:Decode file:%s fail%s error:%s", tools.Red, filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.Filename[:])), tools.Reset, err)
		return
	}
	decodefile, err := tools.AesDecrypt(encodefile, []byte(filePWD))
	if err != nil {
		fmt.Println("[Error]Dncode the file fail ,error! ", err)
		return
	}
	err = installfile.Truncate(0)
	_, err = installfile.Seek(0, os.SEEK_SET)
	_, err = installfile.Write(decodefile[:])
	return
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
