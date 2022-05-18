package client

import (
	"cess-portal/conf"
	"cess-portal/internal/chain"
	"cess-portal/internal/logger"
	"cess-portal/internal/rpc"
	"cess-portal/module"
	"cess-portal/tools"
	"context"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/howeyc/gopass"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

/*
FileUpload means upload files to CESS system
path:The absolute path of the file to be uploaded
backups:Number of backups of files that need to be uploaded
PrivateKey:Encrypted password for uploaded files
*/
func FileUpload(path, backups, PrivateKey string) error {
	if len(PrivateKey) != 16 && len(PrivateKey) != 24 && len(PrivateKey) != 32 && len(PrivateKey) != 0 {
		fmt.Printf("[Error]The privatekey must be 16,24,32 bits long\n")
		return errors.New("[Error]The privatekey must be 16,24,32 bits long")
	}
	chain.Chain_Init()
	file, err := os.Stat(path)
	if err != nil {
		fmt.Printf("[Error]Please enter the correct file path!\n")
		return err
	}

	if file.IsDir() {
		fmt.Printf("[Error]Please do not upload the folder!\n")
		return err
	}

	spares, err := strconv.Atoi(backups)
	if err != nil {
		fmt.Printf("[Error]Please enter a correct integer!\n")
		return err
	}

	filehash, err := tools.CalcFileHash(path)
	if err != nil {
		fmt.Printf("[Error]There is a problem with the file, please replace it!\n")
		return err
	}

	fileid, err := tools.GetGuid(1)
	if err != nil {
		fmt.Printf("[Error]Create snowflake fail! error:%s\n", err)
		return err
	}
	var blockinfo module.FileUploadInfo
	blockinfo.Backups = backups
	blockinfo.FileId = fileid
	blockinfo.BlockSize = int32(file.Size())
	blockinfo.FileHash = filehash

	blocksize := 1024 * 1024
	blocktotal := 0

	f, err := os.Open(path)
	if err != nil {
		fmt.Println("[Error]This file was broken! ", err)
		return err
	}
	defer f.Close()
	filebyte, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("[Error]analyze this file error! ", err)
		return err
	}

	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindSchedulerInfoModule
	ci.ChainModuleMethod = chain.FindSchedulerInfoMethod
	schds, err := ci.GetSchedulerInfo()
	if err != nil {
		fmt.Println("[Error]Get scheduler randomly error! ", err)
		return err
	}
	//var filesize uint64
	fee := new(big.Int)

	ci.IdentifyAccountPhrase = conf.ClientConf.ChainData.IdAccountPhraseOrSeed
	ci.TransactionName = chain.UploadFileTransactionName

	//if file.Size()/1024 == 0 {
	//	filesize = 1
	//} else {
	//	filesize = uint64(file.Size() / 1024)
	//}
	fee.SetInt64(int64(0))

	AsInBlock, err := ci.UploadFileMetaInformation(fileid, file.Name(), filehash, PrivateKey == "", uint8(spares), uint64(file.Size()), fee)
	if err != nil {
		fmt.Printf("\n[Error]Upload file meta information error:%s\n", err)
		return err
	}
	fmt.Printf("\nFile meta info upload:%s ,fileid is:%s\n", AsInBlock, fileid)

	var client *rpc.Client
	for i, schd := range schds {
		wsURL := "ws://" + string(base58.Decode(string(schd.Ip)))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		client, err = rpc.DialWebsocket(ctx, wsURL, "")
		defer cancel()
		if err != nil {
			err = errors.New("Connect with scheduler timeout")
			fmt.Printf("%s[Tips]%sdialog with scheduler:%s fail! reason:%s\n", tools.Yellow, tools.Reset, string(base58.Decode(string(schd.Ip))), err)
			if i == len(schds)-1 {
				fmt.Printf("%s[Error]All scheduler is offline!!!%s\n", tools.Red, tools.Reset)
				logger.OutPutLogger.Sugar().Infof("\n%s[Error]All scheduler is offlien!!!%s\n", tools.Red, tools.Reset)
				return err
			}
			continue
		} else {
			break
		}
	}
	sp := sync.Pool{
		New: func() interface{} {
			return &rpc.ReqMsg{}
		},
	}
	commit := func(num int, data []byte) error {
		blockinfo.BlockNum = int32(num) + 1
		blockinfo.Data = data
		info, err := proto.Marshal(&blockinfo)
		if err != nil {
			fmt.Println("[Error]Serialization error, please upload again! ", err)
			logger.OutPutLogger.Sugar().Infof("[Error]Serialization error, please upload again! ", err)
			return err
		}
		reqmsg := sp.Get().(*rpc.ReqMsg)
		reqmsg.Body = info
		reqmsg.Method = module.UploadService
		reqmsg.Service = module.CtlServiceName

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		resp, err := client.Call(ctx, reqmsg)
		defer cancel()
		if err != nil {
			fmt.Printf("\n%s[Error]Failed to transfer file to scheduler,error:%s%s\n", tools.Red, err, tools.Reset)
			logger.OutPutLogger.Sugar().Infof("%s[Error]Failed to transfer file to scheduler,error:%s%s\n", tools.Red, err, tools.Reset)
			return err
		}

		var res rpc.RespBody
		err = proto.Unmarshal(resp.Body, &res)
		if err != nil {
			fmt.Printf("\n[Error]Error getting reply from schedule, transfer failed! ", err)
			logger.OutPutLogger.Sugar().Infof("[Error]Error getting reply from schedule, transfer failed! ", err)
			return err
		}
		if res.Code != 200 {
			fmt.Printf("\n[Error]Upload file fail!scheduler problem:%s\n", res.Msg)
			logger.OutPutLogger.Sugar().Infof("\n[Error]Upload file fail!scheduler problem:%s\n", res.Msg)
			os.Exit(conf.Exit_SystemErr)
		}
		sp.Put(reqmsg)
		return nil
	}

	if len(PrivateKey) != 0 {
		_, err = os.Stat(conf.ClientConf.PathInfo.KeyPath)
		if err != nil {
			err = os.Mkdir(conf.ClientConf.PathInfo.KeyPath, os.ModePerm)
			if err != nil {
				fmt.Printf("%s[Error]Create key path error :%s%s\n", tools.Red, err, tools.Reset)
				logger.OutPutLogger.Sugar().Infof("%s[Error]Create key path error :%s%s\n", tools.Red, err, tools.Reset)
				os.Exit(conf.Exit_SystemErr)
			}
		}

		os.Create(filepath.Join(conf.ClientConf.PathInfo.KeyPath, file.Name()) + ".pem")
		keyfile, err := os.OpenFile(filepath.Join(conf.ClientConf.PathInfo.KeyPath, file.Name())+".pem", os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("%s[Error]:Failed to save key%s error:%s\n", tools.Red, tools.Reset, err)
			logger.OutPutLogger.Sugar().Infof("%s[Error]:Failed to save key%s error:%s\n", tools.Red, tools.Reset, err)
			return err
		}
		_, err = keyfile.WriteString(PrivateKey)
		if err != nil {
			fmt.Printf("%s[Error]:Failed to write key to file:%s%s error:%s", tools.Red, filepath.Join(conf.ClientConf.PathInfo.KeyPath, (file.Name()+".pem")), tools.Reset, err)
			logger.OutPutLogger.Sugar().Infof("%s[Error]:Failed to write key to file:%s%s error:%s", tools.Red, filepath.Join(conf.ClientConf.PathInfo.KeyPath, (file.Name()+".pem")), tools.Reset, err)
			return err
		}

		encodefile, err := tools.AesEncrypt(filebyte, []byte(PrivateKey))
		if err != nil {
			fmt.Println("[Error]Encode the file fail ,error! ", err)
			return err
		}
		blocks := len(encodefile) / blocksize
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
			err = commit(i, block)
			if err != nil {
				bar.Finish()
				fmt.Printf("%s[Error]:Failed to upload the file%s error:%s\n", tools.Red, tools.Reset, err)
				return err
			}
		}
		bar.Finish()
	} else {
		fmt.Printf("%s[Tips]%s:upload file:%s without private key", tools.Yellow, tools.Reset, path)
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
			err = commit(i, block)
			if err != nil {
				bar.Finish()
				fmt.Printf("%s[Error]:Failed to upload the file%s error:%s\n", tools.Red, tools.Reset, err)
				return err
			}
		}
		bar.Finish()
	}
	fmt.Printf("%s[Success]%s:upload file:%s successful!", tools.Green, tools.Reset, path)
	return nil
}

/*
FileDownload means download file by file id
fileid:fileid of the file to download
*/
func FileDownload(fileid string) error {
	chain.Chain_Init()
	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindFileChainModule
	ci.ChainModuleMethod = chain.FindFileModuleMethod[0]
	fileinfo, err := ci.GetFileInfo(fileid)
	if err != nil {
		fmt.Printf("%s[Error]Get file:%s info fail:%s%s\n", tools.Red, fileid, err, tools.Reset)
		logger.OutPutLogger.Sugar().Infof("%s[Error]Get file:%s info fail:%s%s\n", tools.Red, fileid, err, tools.Reset)
		return err
	}
	if fileinfo.File_Name == nil {
		fmt.Printf("%s[Error]The fileid:%s used to find the file is incorrect, please try again%s\n", tools.Red, fileid, err, tools.Reset)
		logger.OutPutLogger.Sugar().Infof("%s[Error]The fileid:%s used to find the file is incorrect, please try again%s\n", tools.Red, fileid, err, tools.Reset)
		return err
	}
	if string(fileinfo.FileState) != "active" {
		fmt.Printf("%s[Tips]The file:%s has not been backed up, please try again later%s\n", tools.Yellow, fileid, tools.Reset)
		logger.OutPutLogger.Sugar().Infof("%s[Tips]The file:%s has not been backed up, please try again later%s\n", tools.Yellow, fileid, tools.Reset)
		return err
	}

	_, err = os.Stat(conf.ClientConf.PathInfo.InstallPath)
	if err != nil {
		err = os.Mkdir(conf.ClientConf.PathInfo.InstallPath, os.ModePerm)
		if err != nil {
			fmt.Printf("%s[Error]Create install path error :%s%s\n", tools.Red, err, tools.Reset)
			logger.OutPutLogger.Sugar().Infof("%s[Error]Create install path error :%s%s\n", tools.Red, err, tools.Reset)
			os.Exit(conf.Exit_SystemErr)
		}
	}
	_, err = os.Create(filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.File_Name[:])))
	if err != nil {
		fmt.Printf("%s[Error]Create installed file error :%s%s\n", tools.Red, err, tools.Reset)
		logger.OutPutLogger.Sugar().Infof("%s[Error]Create installed file error :%s%s\n", tools.Red, err, tools.Reset)
		os.Exit(conf.Exit_SystemErr)
	}
	installfile, err := os.OpenFile(filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.File_Name[:])), os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("%s[Error]:Failed to save key error:%s%s", tools.Red, err, tools.Reset)
		logger.OutPutLogger.Sugar().Infof("%s[Error]:Failed to save key error:%s%s", tools.Red, err, tools.Reset)
		return err
	}
	defer installfile.Close()

	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindSchedulerInfoModule
	ci.ChainModuleMethod = chain.FindSchedulerInfoMethod
	schds, err := ci.GetSchedulerInfo()
	if err != nil {
		fmt.Printf("%s[Error]Get scheduler list error:%s%s\n ", tools.Red, err, tools.Reset)
		return err
	}

	var client *rpc.Client
	for i, schd := range schds {
		wsURL := "ws://" + string(base58.Decode(string(schd.Ip)))
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		client, err = rpc.DialWebsocket(ctx, wsURL, "")
		defer cancel()
		if err != nil {
			err = errors.New("Connect with scheduler timeout")
			fmt.Printf("%s[Tips]%sdialog with scheduler:%s fail! reason:%s\n", tools.Yellow, tools.Reset, string(base58.Decode(string(schd.Ip))), err)
			if i == len(schds)-1 {
				fmt.Printf("%s[Error]All scheduler is offline!!!%s\n", tools.Red, tools.Reset)
				//logger.OutPutLogger.Sugar().Infof("\n%s[Error]All scheduler is offlien!!!%s\n", tools.Red, tools.Reset)
				return err
			}
			continue
		} else {
			break
		}
	}

	var wantfile module.FileDownloadReq
	var bar tools.Bar
	var getAllBar sync.Once
	sp := sync.Pool{
		New: func() interface{} {
			return &rpc.ReqMsg{}
		},
	}
	wantfile.FileId = fileid
	wantfile.WalletAddress = conf.ClientConf.ChainData.AccountPublicKey
	wantfile.Blocks = 1

	for {
		data, err := proto.Marshal(&wantfile)
		if err != nil {
			fmt.Printf("[Error]Marshal req file error:%s\n", err)
			logger.OutPutLogger.Sugar().Infof("[Error]Marshal req file error:%s\n", err)
			return err
		}
		req := sp.Get().(*rpc.ReqMsg)
		req.Method = module.DownloadService
		req.Service = module.CtlServiceName
		req.Body = data

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		resp, err := client.Call(ctx, req)
		cancel()
		if err != nil {
			fmt.Printf("[Error]Download file fail error:%s\n", err)
			logger.OutPutLogger.Sugar().Infof("[Error]Download file fail error:%s\n", err)
			return err
		}

		var respbody rpc.RespBody
		err = proto.Unmarshal(resp.Body, &respbody)
		if err != nil || respbody.Code != 200 {
			fmt.Printf("[Error]Download file from CESS error:%s. reply message:%s\n", err, respbody.Msg)
			logger.OutPutLogger.Sugar().Infof("[Error]Download file from CESS error:%s. reply message:%s\n", err, respbody.Msg)
			return err
		}
		var blockData module.FileDownloadInfo
		err = proto.Unmarshal(respbody.Data, &blockData)
		if err != nil {
			fmt.Printf("[Error]Download file from CESS error:%s\n", err)
			logger.OutPutLogger.Sugar().Infof("[Error]Download file from CESS error:%s\n", err)
			return err
		}

		_, err = installfile.Write(blockData.Data)
		if err != nil {
			fmt.Printf("%s[Error]:Failed to write file's block to file:%s%s error:%s\n", tools.Red, filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.File_Name[:])), tools.Reset, err)
			logger.OutPutLogger.Sugar().Infof("%s[Error]:Failed to write file's block to file:%s%s error:%s", tools.Red, filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.File_Name[:])), tools.Reset, err)
			return err
		}

		getAllBar.Do(func() {
			bar.NewOption(0, int64(blockData.BlockNum))
		})
		bar.Play(int64(blockData.Blocks))
		wantfile.Blocks++
		sp.Put(req)
		if blockData.Blocks == blockData.BlockNum {
			break
		}
	}

	bar.Finish()
	fmt.Printf("%s[OK]:File '%s' has been downloaded to the directory :%s%s\n", tools.Green, string(fileinfo.File_Name), filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.File_Name[:])), tools.Reset)
	//logger.OutPutLogger.Sugar().Infof("%s[OK]:File '%s' has been downloaded to the directory :%s%s", tools.Green,string(fileinfo.Filename),filepath.Join(conf.ClientConf.PathInfo.InstallPath,string(fileinfo.Filename[:])), tools.Reset)

	if !fileinfo.Public {
		fmt.Printf("%s[Warm]This is a private file, please enter the file password(If you don't want to decrypt, just press enter):%s\n", tools.Green, tools.Reset)
		fmt.Printf("Password:")
		filePWD, _ := gopass.GetPasswdMasked()
		if len(filePWD) == 0 {
			return nil
		}
		if len(filePWD) != 16 && len(filePWD) != 24 && len(filePWD) != 32 {
			return errors.New("[Error]The privatekey must be 16,24,32 bits long")
		}
		encodefile, err := ioutil.ReadFile(filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.File_Name[:])))
		if err != nil {
			fmt.Printf("%s[Error]:Decode file:%s fail%s error:%s\n", tools.Red, filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.File_Name[:])), tools.Reset, err)
			logger.OutPutLogger.Sugar().Infof("%s[Error]:Decode file:%s fail%s error:%s\n", tools.Red, filepath.Join(conf.ClientConf.PathInfo.InstallPath, string(fileinfo.File_Name[:])), tools.Reset, err)
			return err
		}
		decodefile, err := tools.AesDecrypt(encodefile, filePWD)
		if err != nil {
			fmt.Println("[Error]Dncode the file fail ,error! ", err)
			return err
		}
		err = installfile.Truncate(0)
		_, err = installfile.Seek(0, os.SEEK_SET)
		_, err = installfile.Write(decodefile[:])
	}

	return nil
}

/*
FileDelete means to delete the file from the CESS system by the file id
fileid:fileid of the file that needs to be deleted
*/
func FileDelete(fileid string) error {
	chain.Chain_Init()
	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.IdentifyAccountPhrase = conf.ClientConf.ChainData.IdAccountPhraseOrSeed
	ci.TransactionName = chain.DeleteFileTransactionName

	err := ci.DeleteFileOnChain(fileid)
	if err != nil {
		fmt.Printf("%s[Error]Delete file error:%s%s\n", tools.Red, tools.Reset, err)
		logger.OutPutLogger.Sugar().Infof("%s[Error]Delete file error:%s%s\n", tools.Red, tools.Reset, err)
		return err
	} else {
		fmt.Printf("%s[OK]Delete fileid:%s success!%s\n", tools.Green, fileid, tools.Reset)
		logger.OutPutLogger.Sugar().Infof("%s[OK]Delete fileid:%s success!%s\n", tools.Green, fileid, tools.Reset)
		return nil
	}

}

/*
FileDecrypt means that if the file is not decrypted when downloading the file, it can be decrypted by this method
When you download the file if it is not decode, you can decode it this way
*/
func FileDecrypt(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		fmt.Printf("%s[Error]There is no such file, please confirm the correct location of the file, please enter the absolute path of the file%s\n", tools.Red, tools.Reset)
		logger.OutPutLogger.Sugar().Infof("%s[Error]There is no such file, please confirm the correct location of the file, please enter the absolute path of the file%s\n", tools.Red, tools.Reset)
		return err
	}

	fmt.Println("Please enter the file's password:")
	fmt.Print(">")
	psw, _ := gopass.GetPasswdMasked()
	if len(psw) != 16 && len(psw) != 24 && len(psw) != 32 {
		return errors.New("[Error]The password must be 16,24,32 bits long")
	}
	encodefile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%s[Error]Failed to read file, please check file integrity%s\n", tools.Red, tools.Reset)
		logger.OutPutLogger.Sugar().Infof("%s[Error]Failed to read file, please check file integrity%s\n", tools.Red, tools.Reset)
		return err
	}

	decodefile, err := tools.AesDecrypt(encodefile, psw)
	if err != nil {
		fmt.Printf("%s[Error]File decode failed, please check your password! error:%s%s ", tools.Red, err, tools.Reset)
		return err
	}
	filename := filepath.Base(path)
	//The decoded file is saved to the download folder, if the name is the same, the original file will be deleted
	if path == filepath.Join(conf.ClientConf.PathInfo.InstallPath, filename) {
		err = os.Remove(path)
		if err != nil {
			fmt.Printf("%s[Error]An error occurred while saving the decoded file! error:%s%s ", tools.Red, err, tools.Reset)
			return err
		}
	}
	fileinfo, err := os.Create(filepath.Join(conf.ClientConf.PathInfo.InstallPath, filename))
	if err != nil {
		fmt.Printf("%s[Error]An error occurred while saving the decoded file! error:%s%s ", tools.Red, err, tools.Reset)
		return err
	}
	defer fileinfo.Close()
	_, err = fileinfo.Write(decodefile)
	if err != nil {
		fmt.Printf("%s[Error]Failed to save decrypted content to file! error:%s%s ", tools.Red, err, tools.Reset)
		return err
	}

	fmt.Printf("%s[Success]The file was decrypted successfully and the file has been saved to:%s%s ", tools.Green, filepath.Join(conf.ClientConf.PathInfo.InstallPath, filename), tools.Reset)

	return nil
}
