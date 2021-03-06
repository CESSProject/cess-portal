package client

import (
	"cess-portal/conf"
	"cess-portal/internal/chain"
	"cess-portal/internal/logger"
	"cess-portal/tools"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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
walletaddress:wallet address
*/
func ObtainFromFaucet(walletaddress string) error {
	pubkey, err := tools.DecodeToPub(walletaddress, tools.ChainCessTestPrefix)
	if err != nil {
		fmt.Printf("[Error]The wallet address you entered is incorrect, please re-enter:%v\n", err.Error())
		return err
	}
	var ob = struct {
		Address string `json:"Address"`
	}{
		tools.PubBytesTo0XString(pubkey),
	}
	var res faucet
	resp, err := tools.Post(conf.ClientConf.ChainData.FaucetAddress, ob)
	if err != nil {
		fmt.Printf("[Error]Network problem, please check your network connection\n")
		logger.OutPutLogger.Sugar().Infof("[Error]Network problem, please check your network connection, error:%s\n", err)
		return err
	}
	err = json.Unmarshal(resp, &res)
	if err != nil {
		fmt.Printf("Incorrect response from faucet\n")
		logger.OutPutLogger.Sugar().Infof("Incorrect response from faucet,error:%s", err)
		return err
	}
	if res.Ans.Err != "" {
		fmt.Printf("[Error]get free token from faucet fail:%s\n", res.Ans.Err)
		logger.OutPutLogger.Sugar().Infof("[Error]get free token from faucet fail:%s\n", res.Ans.Err)
		return err
	}

	if res.Ans.AsInBlock {
		fmt.Printf("[Success]Obtain from faucet success\n")
		logger.OutPutLogger.Sugar().Infof("[Success]Obtain from faucet success\n")
	} else {
		fmt.Printf("[Fail]Obtain from faucet fail,Please wait 24 hours to get it again\n")
		logger.OutPutLogger.Sugar().Infof("[Fail]Obtain from faucet fail,Please wait 24 hours to get it again\n")
	}
	return nil
}

/*
Expansion means the purchase of storage capacity for the current customer
quantity:The amount of space to be purchased (1/G)
duration:Market for space that needs to be purchased (1/month)
expected:The expected number of prices when buying is required to prevent price fluctuations when buying. When it is 0, it means that any price can be accepted
*/
func Expansion(quantity, duration, expected int) error {
	chain.Chain_Init()
	if quantity == 0 || duration == 0 {
		fmt.Printf("[Error] Please enter the correct purchase number\n")
		return errors.New("Please enter the correct purchase number")
	}
	var ci chain.CessInfo
	ci.RpcAddr = conf.ClientConf.ChainData.CessRpcAddr
	ci.IdentifyAccountPhrase = conf.ClientConf.ChainData.IdAccountPhraseOrSeed
	ci.TransactionName = chain.BuySpaceTransactionName

	//Buying space on-chain, failure could mean running out of money
	err := ci.BuySpaceOnChain(quantity, duration, expected/1024)
	if err != nil {
		fmt.Printf("[Error] Failed to buy space, please check if you have enough money or check if there is enough space on the chain\n")
		logger.OutPutLogger.Sugar().Infof("[Error] Failed to buy space, please check if you have enough money or check if there is enough space on the chain,error:%v\n", err)
		return err
	}
	fmt.Printf("[Success]Buy space on chain success!\n")
	return nil
}
