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
	AsInBlock bool   `json:"AsInBlock"`
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

	if res.Ans.AsInBlock {
		fmt.Printf("[Success]Obtain from faucet success\n")
		logger.OutPutLogger.Sugar().Infof("[Success]Obtain from faucet success\n")
	} else {
		fmt.Printf("[Fail]Obtain from faucet fail,Please wait 24 hours to get it again\n")
		logger.OutPutLogger.Sugar().Infof("[Fail]Obtain from faucet fail,Please wait 24 hours to get it again\n")
	}

}

func Expansion(quantity, expected int) {
	chain.Chain_Init()

	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.IdentifyAccountPhrase = conf.ClientConf.ChainData.IdAccountPhraseOrSeed
	ci.TransactionName = chain.BuySpaceTransactionName

	AsInBlock, err := ci.BuySpaceOnChain(quantity, expected)
	if err != nil {
		fmt.Printf("[Error] Buy space on chain fail:%s\n", err)
		logger.OutPutLogger.Sugar().Infof("[Error] Buy space on chain fail:%s\n", err)
		return
	}
	fmt.Printf("Transaction chain block number is:%s\n", AsInBlock)
}
