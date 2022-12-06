package client

import (
	"cess-portal/conf"
	"cess-portal/internal/chain"
	"cess-portal/internal/erasure"
	"cess-portal/internal/hashtree"
	. "cess-portal/internal/logger"
	"cess-portal/internal/tcp"
	"cess-portal/tools"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	cesskeyring "github.com/CESSProject/go-keyring"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

const LOG_TAG_FILEUPLOAD = "UploadFile"
const LOG_TAG_FILEDELETE = "FileDelete"
const LOG_TAG_FILEDOWNLOAD = "FileDownload"
const ERR_404 = "Not found"

//File Upload

func FileUpload(fullpath, bucketName string) {
	fpath, fname := filepath.Split(fullpath)
	//set cache dir
	conf.FileCacheDir = fpath
	// Calc file state
	fstat, err := os.Stat(filepath.Join(fpath, fname))
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", LOG_TAG_FILEUPLOAD, err)
	}
	// Calc reedsolomon
	chunkPath, datachunkLen, rduchunkLen, err := erasure.ReedSolomon(filepath.Join(fpath, fname), fstat.Size())
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", LOG_TAG_FILEUPLOAD, err)
		log.Println("Client internal error, please try again or check the problems reported in the log")
		return
	}

	if len(chunkPath) != (datachunkLen + rduchunkLen) {
		Uld.Sugar().Infof("[%v] %v", LOG_TAG_FILEUPLOAD, "ReedSolomon failed")
		log.Println("Client internal error, please try again or check the problems reported in the log")
		return
	}
	// Calc merkle hash tree
	hTree, err := hashtree.NewHashTree(chunkPath)
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", LOG_TAG_FILEUPLOAD, err)
		log.Println("Client internal error, please try again or check the problems reported in the log")
		return
	}

	// Merkel root hash
	fileid := hex.EncodeToString(hTree.MerkleRoot())
	//save fileid
	newpath := filepath.Join(fpath, fileid)
	f, err := os.Create(newpath)
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", LOG_TAG_FILEUPLOAD, err)
		log.Println("Failed to save fileid, possibly due to insufficient permissions. you can check the log for details")
		return
	}
	f.Close()
	// Rename chunks with root hash
	var newChunksPath = make([]string, 0)
	if rduchunkLen == 0 {
		newChunksPath = append(newChunksPath, fileid)
	} else {
		for i := 0; i < len(chunkPath); i++ {
			var ext = filepath.Ext(chunkPath[i])
			var newchunkpath = filepath.Join(fpath, fileid+ext)
			os.Rename(chunkPath[i], newchunkpath)
			newChunksPath = append(newChunksPath, fileid+ext)
		}
	}
	//build a user brief
	pubkey, err := tools.DecodePublicKeyOfCessAccount(conf.C.AccountId)
	if err != nil {
		Uld.Sugar().Infof("[%v] %v", LOG_TAG_FILEUPLOAD, err)
		log.Println("Failed to decode public key from cess account,please check your config setting")
		return
	}
	userBrief := chain.UserBrief{
		User:        types.NewAccountID(pubkey),
		File_name:   types.Bytes(fname),
		Bucket_name: types.Bytes(bucketName),
	}
	// Declaration file
	txhash, err := chain.ChainClient.DeclarationFile(fileid, userBrief)
	if err != nil || txhash == "" {
		Uld.Sugar().Infof("[%v] %v", LOG_TAG_FILEUPLOAD, err)
		log.Println("Failed to upload file declaration. you can check the log for details")
		return
	}
	task_StoreFile(newChunksPath, LOG_TAG_FILEUPLOAD, fileid, fname, fstat.Size())
}

func task_StoreFile(fpath []string, logtag, fid, fname string, fsize int64) {
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", err)
		}
	}()
	var channel_1 = make(chan uint8, 1)
	Uld.Sugar().Infof("[%v] Start the file backup management process", fid)
	go uploadToStorage(channel_1, fpath, logtag, fid, fname, fsize)
	for {
		select {
		case result := <-channel_1:
			if result == 1 {
				go uploadToStorage(channel_1, fpath, logtag, fid, fname, fsize)
				time.Sleep(time.Second * 6)
			}
			if result == 2 {
				Uld.Sugar().Infof("[%v] File save successfully", fid)
				log.Println("Upload file success")
				return
			}
			if result == 3 {
				Uld.Sugar().Infof("[%v] File save failed", fid)
				log.Println("Upload file failed, please try again.")
				return
			}
		}
	}
}

// Upload files to cess storage system
func uploadToStorage(ch chan uint8, fpath []string, logtag, fid, fname string, fsize int64) {
	defer func() {
		err := recover()
		if err != nil {
			ch <- 1
			Uld.Sugar().Infof("[panic]: [%v] [%v] %v", logtag, fpath, err)
		}
	}()

	var existFile = make([]string, 0)
	for i := 0; i < len(fpath); i++ {
		_, err := os.Stat(filepath.Join(conf.FileCacheDir, fpath[i]))
		if err != nil {
			continue
		}
		existFile = append(existFile, fpath[i])
	}
	msg := tools.GetRandomcode(16)

	kr, _ := cesskeyring.FromURI(conf.C.AccountSeed, cesskeyring.NetSubstrate{})
	// sign message
	sign, err := kr.Sign(kr.SigningContext([]byte(msg)))
	if err != nil {
		ch <- 1
		Uld.Sugar().Infof("[%v] %v", logtag, err)
		return
	}

	// Get all scheduler
	schds, err := chain.ChainClient.GetSchedulerList()
	if err != nil {
		ch <- 1
		Uld.Sugar().Infof("[%v] %v", logtag, err)
		return
	}

	tools.RandSlice(schds)

	for i := 0; i < len(schds); i++ {
		wsURL := fmt.Sprintf("%d.%d.%d.%d:%d",
			schds[i].Ip.Value[0],
			schds[i].Ip.Value[1],
			schds[i].Ip.Value[2],
			schds[i].Ip.Value[3],
			schds[i].Ip.Port,
		)
		log.Println("Will send to ", wsURL)
		conTcp, err := dialTcpServer(wsURL)
		if err != nil {
			Uld.Sugar().Error(fmt.Errorf("dial %v err: %v", wsURL, err))
			continue
		}
		srv := tcp.NewClient(tcp.NewTcp(conTcp), conf.FileCacheDir, existFile)
		err = srv.SendFile(fid, fsize, conf.PublicKey, []byte(msg), sign[:])
		if err != nil {
			Uld.Sugar().Infof("[%v] %v", logtag, err)
			continue
		}
		ch <- 2
		return
	}
	ch <- 1
}

// File Download

func FileDownload(fid, cacheDir string) {
	conf.FileCacheDir = cacheDir
	_, err := os.Stat(conf.FileCacheDir)
	if err != nil {
		err = os.MkdirAll(conf.FileCacheDir, os.ModeDir)
		if err != nil {
			Uld.Sugar().Infof("[%v] %v", LOG_TAG_FILEDOWNLOAD, err)
			return
		}
	}
	// //clear cache
	fpath := filepath.Join(conf.FileCacheDir, fid)
	_, err = os.Stat(fpath)
	if err == nil {
		os.Remove(fpath)
	}
	// file meta info
	fmeta, err := chain.ChainClient.GetFileMetaInfo(fid) //GetFileMetaInfoOnChain(fid)
	if err != nil {
		if err.Error() == chain.ERR_Empty {
			Uld.Sugar().Errorf("[%v] Get file metadata err: %v", LOG_TAG_FILEDOWNLOAD, err)
			log.Println("Get file metadata failed,please ensure that you have configured the correct account or passed in the fileid of.")
			return
		}
		Uld.Sugar().Errorf("[%v] %v", LOG_TAG_FILEDOWNLOAD, err)
		log.Println("Get file metadata failed.")
		return
	}
	r := len(fmeta.BlockInfo) / 3
	d := len(fmeta.BlockInfo) - r
	down_count := 0
	for i := 0; i < len(fmeta.BlockInfo); i++ {
		// Download the file from the scheduler service
		fname := filepath.Join(conf.FileCacheDir, string(fmeta.BlockInfo[i].BlockId[:]))
		if len(fmeta.BlockInfo) == 1 {
			fname = fname[:(len(fname) - 4)]
		}
		mip := fmt.Sprintf("%d.%d.%d.%d:%d",
			fmeta.BlockInfo[i].MinerIp.Value[0],
			fmeta.BlockInfo[i].MinerIp.Value[1],
			fmeta.BlockInfo[i].MinerIp.Value[2],
			fmeta.BlockInfo[i].MinerIp.Value[3],
			fmeta.BlockInfo[i].MinerIp.Port,
		)
		err = downloadFromStorage(fname, int64(fmeta.BlockInfo[i].BlockSize), mip)
		if err != nil {
			Uld.Sugar().Errorf("[%v] Downloading %drd shard err: %v", LOG_TAG_FILEDOWNLOAD, i, err)
			log.Printf("Download shard %d failed,please try again.\n", i)
		} else {
			down_count++
		}
		if down_count >= d {
			break
		}
	}
	log.Println("info", conf.FileCacheDir, fid, d, r)
	err = erasure.ReedSolomon_Restore(conf.FileCacheDir, fid, d, r, uint64(fmeta.Size))
	if err != nil {
		Uld.Sugar().Errorf("[%v] ReedSolomon_Restore: %v", LOG_TAG_FILEDOWNLOAD, err)
		log.Println("Restore reedSolomon failed,please try again.")
		return
	}

	if r > 0 {
		fstat, err := os.Stat(fpath)
		if err != nil {
			Uld.Sugar().Errorf("[%v] %v", LOG_TAG_FILEDOWNLOAD, err)
			log.Println("download file failed.")
			return
		}
		if uint64(fstat.Size()) > uint64(fmeta.Size) {
			tempfile := fpath + ".temp"
			copyFile(fpath, tempfile, int64(fmeta.Size))
			os.Remove(fpath)
			os.Rename(tempfile, fpath)
		}
	}
	//delete file slice and rename file
	for i := 0; i < d; i++ {
		os.Remove(fmt.Sprintf("%s.00%d", fpath, i))
	}
	newPath := filepath.Join(conf.FileCacheDir, string(fmeta.UserBriefs[0].File_name))
	os.Rename(fpath, newPath)
	log.Println("Download file success.")
}

// Download files from cess storage service
func downloadFromStorage(fpath string, fsize int64, mip string) error {
	fsta, err := os.Stat(fpath)
	if err == nil {
		if fsta.Size() == fsize {
			return nil
		} else {
			os.Remove(fpath)
		}
	}

	msg := tools.GetRandomcode(16)

	kr, _ := cesskeyring.FromURI(conf.C.AccountSeed, cesskeyring.NetSubstrate{})
	// sign message
	sign, err := kr.Sign(kr.SigningContext([]byte(msg)))
	if err != nil {
		return err
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", mip)
	if err != nil {
		return err
	}

	conTcp, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	srv := tcp.NewClient(tcp.NewTcp(conTcp), conf.FileCacheDir, nil)
	return srv.RecvFile(filepath.Base(fpath), fsize, conf.PublicKey, []byte(msg), sign[:])
}

func copyFile(src, dst string, length int64) error {
	srcfile, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer srcfile.Close()
	dstfile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer dstfile.Close()

	var buf = make([]byte, 64*1024)
	var count int64
	for {
		n, err := srcfile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		count += int64(n)
		if count < length {
			dstfile.Write(buf[:n])
		} else {
			tail := count - length
			if n >= int(tail) {
				dstfile.Write(buf[:(n - int(tail))])
			}
		}
	}

	return nil
}

//File Delete

func FileDelete(fid string) {
	if fid == "" {
		Uld.Sugar().Errorf("[%v] No fid", LOG_TAG_FILEDELETE)
		log.Println("Please enter the correct fid")
		return
	}
	//Delete files in cesss storage service
	txhash, err := chain.ChainClient.DeleteFile(chain.ChainClient.GetPublicKey(), fid)
	if txhash == "" {
		Err.Sugar().Errorf("[%sv] %v", LOG_TAG_FILEDELETE, err)
		log.Println("delete file in cess storage service failed.")
		return
	}
	log.Println("Delete file success,the Tx hash is", txhash)
}

func dialTcpServer(address string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	dialer := net.Dialer{Timeout: conf.Tcp_Dial_Timeout}
	netCon, err := dialer.Dial("tcp", tcpAddr.String())
	if err != nil {
		return nil, err
	}
	conTcp, ok := netCon.(*net.TCPConn)
	if !ok {
		return nil, errors.New("network conversion failed")
	}
	return conTcp, nil
}
