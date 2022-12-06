package client

import (
	"cess-portal/conf"
	"cess-portal/internal/chain"
	. "cess-portal/internal/logger"
	"cess-portal/tools"
	"encoding/json"
	"fmt"
	"log"
)

type FileInfo struct {
	State string   `json:"file_state"`
	Size  uint64   `json:"file_size"`
	Names []string `json:"file_names"`
}

type SpacePackage struct {
	chain.SpacePackage
	State string `json:"state"`
}

const LOG_TAG_FILEQUERY = "FileQuery"
const LOG_TAG_BUCKETQUERY = "BucketQuery"

func UserSpaceQuery() {
	spaceInfo, err := chain.ChainClient.GetUserSpaceMetadata(conf.PublicKey)
	if err != nil {
		if err.Error() == chain.ERR_Empty {
			Uld.Sugar().Errorf("[%v] No space info", LOG_TAG_FILEQUERY)
			log.Println("Please configure  the correct account seed")
			return
		}
		Uld.Sugar().Errorf("[%v] Get space info error:%v", LOG_TAG_FILEQUERY, err)
		log.Println("user space info query failed.")
		return
	}
	wrap := SpacePackage{spaceInfo, string(spaceInfo.State)}
	jbytes, err := json.Marshal(wrap)
	if err != nil {
		Uld.Sugar().Errorf("[%v] Marshal space info error:%v", LOG_TAG_FILEQUERY, err)
		log.Println("user space info query failed.")
		return
	}
	fmt.Println("space info of your account is as follow:")
	tools.ShowJsonData(jbytes, "  ")
	fmt.Printf("\nNote: the unit of space capacity is (B),the space validity period is calculated according to the block height.\n")
}

func FilelistQuery(bucketName string) {
	//verify bucket name
	if !tools.VerifyBucketName(bucketName) {
		Uld.Sugar().Errorf("[%v] Bucket name error", LOG_TAG_FILEQUERY)
		log.Println("Please configure  the correct bucket name")
		return
	}
	//query bucket info
	bucketInfo, err := chain.ChainClient.GetBucketInfo(conf.PublicKey, bucketName)
	if err != nil {
		if err.Error() == chain.ERR_Empty {
			Uld.Sugar().Errorf("[%v] No bucket info", LOG_TAG_FILEQUERY)
			log.Println("Please check your params or configure the correct account seed")
			return
		}
		Uld.Sugar().Errorf("[%v] Get bucket info error:%v", LOG_TAG_FILEQUERY, err)
		log.Println("File list query failed.")
		return
	}
	list := make([]string, len(bucketInfo.Objects_list))
	for i := 0; i < len(bucketInfo.Objects_list); i++ {
		list[i] = string(bucketInfo.Objects_list[i][:])
	}
	jbytes, err := json.Marshal(list)
	if err != nil {
		Uld.Sugar().Errorf("[%v] Marshal file list error:%v", LOG_TAG_FILEQUERY, err)
		log.Println("file list query failed.")
		return
	}
	fmt.Printf("file hash list of bucket \"%s\" is as follow :\n", bucketName)
	tools.ShowJsonData(jbytes, "  ")
}

func FilestateQuery(fid string) {
	filestate, err := chain.ChainClient.GetFileMetaInfo(fid) //GetFileMetaInfoOnChain(fid)
	if err != nil {
		if err.Error() == chain.ERR_Empty {
			Uld.Sugar().Errorf("[%v] No fid", LOG_TAG_FILEQUERY)
			log.Println("Please enter the correct fid")
			return
		}
		Uld.Sugar().Errorf("[%v] Get user file state error:%v", LOG_TAG_FILEQUERY, err)
		log.Println("File state query failed.")
		return
	}
	shortInfo := &FileInfo{}
	shortInfo.Size = uint64(filestate.Size)
	shortInfo.State = string(filestate.State)
	for _, v := range filestate.UserBriefs {
		shortInfo.Names = append(shortInfo.Names, string(v.File_name))
	}
	jbytes, err := json.Marshal(shortInfo)
	if err != nil {
		Uld.Sugar().Errorf("[%v] Marshal file info error:%v", LOG_TAG_FILEQUERY, err)
		log.Println("File state query failed.")
		return
	}
	fmt.Printf("The short info of the file is as follow:\n")
	tools.ShowJsonData(jbytes, "  ")
}

func BucketlistQuery() {
	bucketList, err := chain.ChainClient.GetBucketList(conf.PublicKey)
	if err != nil {
		if err.Error() == chain.ERR_Empty {
			Uld.Sugar().Errorf("[%v] No bucket info", LOG_TAG_BUCKETQUERY)
			log.Println("Please configure the correct account seed")
			return
		}
		Uld.Sugar().Errorf("[%v] Get bucket list error:%v", LOG_TAG_BUCKETQUERY, err)
		log.Println("bucket list query failed.")
		return
	}
	fmt.Printf("bucket list of your account:\n")
	buckets := make([]string, len(bucketList))
	for i, b := range bucketList {
		buckets[i] = string(b)
	}
	jbytes, err := json.Marshal(buckets)
	if err != nil {
		Uld.Sugar().Errorf("[%v] Marshal bucket list error:%v", LOG_TAG_FILEQUERY, err)
		log.Println("bucket list query failed.")
		return
	}
	tools.ShowJsonData(jbytes, "  ")
}
