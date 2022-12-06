package client

import (
	"cess-portal/conf"
	"cess-portal/internal/chain"
	. "cess-portal/internal/logger"
	"cess-portal/tools"
	"fmt"
	"log"
)

const LOG_TAG_BUCKETCREATE = "BucketCreate"

func BucketCreate(bucketName string) {
	if !tools.VerifyBucketName(bucketName) {
		Uld.Sugar().Errorf("[%v] Bucket name error", LOG_TAG_BUCKETCREATE)
		log.Println("Please configure  the correct bucket name")
		return
	}
	txHash, err := chain.ChainClient.CreateBucket(conf.PublicKey, bucketName)
	if err != nil {
		Uld.Sugar().Errorf("[%v] Create bucket error:%v", LOG_TAG_BUCKETCREATE, err)
		log.Println("Create bucket failed.")
		return
	}
	fmt.Println("Create bucket success. Tx hash:", txHash)
}

func BucketDelete(bucketName string) {
	if !tools.VerifyBucketName(bucketName) {
		Uld.Sugar().Errorf("[%v] Bucket name error", LOG_TAG_BUCKETCREATE)
		log.Println("Please configure  the correct bucket name")
		return
	}
	txHash, err := chain.ChainClient.DeleteBucket(conf.PublicKey, bucketName)
	if err != nil {
		Uld.Sugar().Errorf("[%v] Delete bucket error:%v", LOG_TAG_BUCKETCREATE, err)
		log.Println("Delete bucket failed.")
		return
	}
	fmt.Println("Delete bucket success. Tx hash:", txHash)
}
