package client

import (
	"cess-portal/conf"
	"cess-portal/internal/chain"
	"cess-portal/internal/logger"
	"cess-portal/tools"
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

/*
ObtainFromFaucet means to obtain tCESS for transaction spending through the faucet
pbk:wallet's public key
*/
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

/*
Expansion means the purchase of storage capacity for the current customer
quantity:The amount of space to be purchased (1/G)
duration:Market for space that needs to be purchased (1/month)
expected:The expected number of prices when buying is required to prevent price fluctuations when buying. When it is 0, it means that any price can be accepted
*/
func Expansion(quantity, duration, expected int) {
	chain.Chain_Init()

	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.IdentifyAccountPhrase = conf.ClientConf.ChainData.IdAccountPhraseOrSeed
	ci.TransactionName = chain.BuySpaceTransactionName

	err := ci.BuySpaceOnChain(quantity, duration, expected)
	if err != nil {
		fmt.Printf("[Error] Buy space on chain fail:%s\n", err)
		logger.OutPutLogger.Sugar().Infof("[Error] Buy space on chain fail:%s\n", err)
		return
	}
	fmt.Printf("[Success]Buy space on chain success!\n")
}
