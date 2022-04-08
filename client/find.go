package client

import (
	"cess-portal/conf"
	"cess-portal/internal/chain"
	"cess-portal/internal/logger"
	"cess-portal/tools"
	"fmt"
)

/*
FindPurchasedSpace means to query the space that the current user has purchased and the space that has been used
*/
func FindPurchasedSpace() error {
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
FindPrice means to get real-time price of storage space
*/
func FindPrice() error {
	chain.Chain_Init()

	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindPriceChainModule

	ci.ChainModuleMethod = chain.FindPriceModuleMethod[0]
	AllPurchased, err := ci.GetPurchasedSpace()
	if err != nil {
		fmt.Printf("[Error]Get all purchased fail:%s\n", err)
		logger.OutPutLogger.Sugar().Infof("[Error]Get all purchased fail::%s\n", err)
		return err
	}

	ci.ChainModuleMethod = chain.FindPriceModuleMethod[1]
	AllAvailable, err := ci.GetAvailableSpace()
	if err != nil {
		fmt.Printf("[Error]Get all available fail:%s\n", err)
		logger.OutPutLogger.Sugar().Infof("[Error]Get all available fail:%s\n", err)
		return err
	}

	var purc int64
	var ava int64
	if AllPurchased.Int != nil {
		purc = AllPurchased.Int64()
	}
	if AllAvailable.Int != nil {
		ava = AllAvailable.Int64()
	}
	if purc == ava {
		fmt.Printf("[Success]All space has been bought,The current storage price is:+∞ per (MB)\n")
		logger.OutPutLogger.Sugar().Infof("[Success]All space has been bought,The current storage price is:+∞ per (MB)\n")
		return err
	}

	result := (1024 / float64((ava - purc))) * 1000

	fmt.Printf("[Success]The current storage price is:%f per (MB)\n", result)
	logger.OutPutLogger.Sugar().Infof("[Success]The current storage price is:%f per (MB)\n", result)
	return nil
}

/*
FindFile means to query the files uploaded by the current user
fileid:fileid of the file to look for
*/
func FindFile(fileid string) error {
	chain.Chain_Init()

	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindFileChainModule

	if fileid != "" {
		ci.ChainModuleMethod = chain.FindFileModuleMethod[0]
		data, err := ci.GetFileInfo(fileid)
		if err != nil {
			fmt.Printf("[Error]Get file:%s info fail:%s\n", fileid, err)
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
			fmt.Printf("[Error]Get file list fail:%s\n", err)
			logger.OutPutLogger.Sugar().Infof("[Error]Get file list fail:%s\n", err)
			return err
		}
		for _, fileinfo := range data {
			fmt.Printf("%s\n", string(fileinfo))
		}
	}
	return nil
}
