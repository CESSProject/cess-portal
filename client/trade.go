package client

import (
	"dapp_cess_client/conf"
	"dapp_cess_client/internal/chain"
	"dapp_cess_client/internal/logger"
	"dapp_cess_client/tools"
	"encoding/json"
	"fmt"
)

type faucet struct {
	Ans    answer `json:"Result"`
	Status string `json:"Status"`
}
type answer struct {
	Err       string `json:"Err"`
	AsInBlock string `json:"AsInBlock"`
}

func ObtainFromFaucet(pbk string) {
	var ob = struct {
		Address string `json:"Address"`
	}{
		pbk,
	}
	var res faucet
	resp, err := tools.Post(conf.ClientConf.ChainData.FaucetAddress, ob)
	if err != nil {
		fmt.Printf("[Error]System error:%s\n", err)
		logger.OutPutLogger.Sugar().Infof("[Error]System error:%s\n", err)
		return
	}
	err = json.Unmarshal(resp, &res)
	if err != nil {
		fmt.Println("unmarshal error:", err)
	}
	if res.Ans.Err != "" {
		fmt.Printf("[Error]Obtain from faucet fail:%s\n", res.Ans.Err)
		logger.OutPutLogger.Sugar().Infof("[Error]Obtain from faucet fail:%s\n", res.Ans.Err)
		return
	}
	fmt.Printf("[Success]Obtain from faucet success,AsInBlock is:%s\n", res.Ans.AsInBlock)
	logger.OutPutLogger.Sugar().Infof("[Success]Obtain from faucet success,AsInBlock is:%s\n", res.Ans.AsInBlock)
}

func Expansion(quantity string) {
	chain.Chain_Init()

	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.IdentifyAccountPhrase = conf.ClientConf.ChainData.IdAccountPhraseOrSeed
	ci.ChainModule = chain.BuySpaceChainModule
	ci.TransactionName = chain.BuySpaceTransactionName
	ci.ChainModuleMethod = chain.BuySpaceModuleMethod
	ci.BuySpaceOnChain(quantity)
}
