package client

import (
	"cess-portal/conf"
	"cess-portal/internal/chain"
	"cess-portal/internal/logger"
	"cess-portal/tools"
	"fmt"
	"strconv"
)

/*
QueryPurchasedSpace means to query the space that the current user has purchased and the space that has been used
*/
func QueryPurchasedSpace() error {
	chain.Chain_Init()

	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.ChainModule = chain.PurchasedSpaceChainModule
	ci.ChainModuleMethod = chain.PurchasedSpaceModuleMethod

	userinfo, err := ci.UserHoldSpaceDetails()
	if err != nil {
		fmt.Printf("[Error]Get user data fail:%s\n", err)
		logger.OutPutLogger.Sugar().Infof("[Error]Get user data fail:%s\n", err)
		return err
	}
	fmt.Println(userinfo)
	return nil
}

/*
QueryPrice means to get real-time price of storage space
*/
func QueryPrice() error {
	chain.Chain_Init()

	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindPriceChainModule

	ci.ChainModuleMethod = chain.FindPriceModuleMethod
	Price, err := ci.GetPrice()
	if err != nil {
		fmt.Printf("%s[Error]%sGet price fail\n", tools.Red, tools.Reset)
		logger.OutPutLogger.Sugar().Infof("%s[Error]%sGet price fail::%s\n", tools.Red, tools.Reset, err)
		return err
	}
	PerGB, _ := strconv.ParseFloat(fmt.Sprintf("%.12f", float64(Price.Int64()*int64(1024))/float64(1000000000000)), 64)
	fmt.Printf("[Success]The current storage price is:%.12f TCESS/GB)\n", PerGB)
	logger.OutPutLogger.Sugar().Infof("[Success]The current storage price is:%.12f TCESS/GB)\n", PerGB)
	return nil
}

/*
QueryFile means to query the files uploaded by the current user
fileid:fileid of the file to look for
*/
func QueryFile(fileid string) error {
	chain.Chain_Init()

	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindFileChainModule

	if fileid != "" {
		ci.ChainModuleMethod = chain.FindFileModuleMethod[0]
		data, err := ci.GetFileInfo(fileid)
		if err != nil {
			fmt.Printf("[Error]Get file:%s info fail\n", fileid)
			logger.OutPutLogger.Sugar().Infof("[Error]Get file:%s info fail:%s\n", fileid, err)
			return err
		}
		fmt.Println(data)
		if len(data.File_Name) == 0 {
			fmt.Printf("%s[Tips]This file:%s may have been deleted by someone%s\n", tools.Yellow, fileid, tools.Reset)
		}
	} else {
		ci.ChainModuleMethod = chain.FindFileModuleMethod[1]
		data, err := ci.GetFileList()
		if err != nil {
			fmt.Printf("[Error]Get file list fail\n")
			logger.OutPutLogger.Sugar().Infof("[Error]Get file list fail:%s\n", err)
			return err
		}
		for _, fileinfo := range data {
			fmt.Printf("%s\n", string(fileinfo))
		}
	}
	return nil
}
