package client

import (
	"cess-portal/conf"
	"cess-portal/internal/chain"
	. "cess-portal/internal/logger"
	"log"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

const LOG_TAG_PURCHASE = "Purchase"

func StoragePurchase(size uint32) {
	txhash, err := chain.ChainClient.BuySpace(types.NewU32(size))
	if err != nil {
		if err.Error() == chain.ERR_Empty {
			Uld.Sugar().Infof("[%v] Empty account", LOG_TAG_PURCHASE)
			log.Println("Account not found")
		} else {
			Uld.Sugar().Infof("[%v] Buy space error: %v", LOG_TAG_PURCHASE, err)
			log.Println("Buy space failed,please check whether your account balance is sufficient")
		}
		return
	}
	log.Println("Buy space success. Tx hash:", txhash)
}

func SpaceAuthorize() {
	txhash, err := chain.ChainClient.AuthorizeSpace(conf.PublicKey)
	if err != nil {
		if err.Error() == chain.ERR_Empty {
			Uld.Sugar().Infof("[%v] Empty account", LOG_TAG_PURCHASE)
			log.Println("Account not found")
		} else {
			Uld.Sugar().Infof("[%v] Authorize space error: %v", LOG_TAG_PURCHASE, err)
			log.Println("Authorize space failed,please configure the correct account seed")
		}
		return
	}
	log.Println("Authorize space success. Tx hash:", txhash)
}

func AuthCancel() {
	txhash, err := chain.ChainClient.CancelAuth()
	if err != nil {
		if err.Error() == chain.ERR_Empty {
			Uld.Sugar().Infof("[%v] Empty account", LOG_TAG_PURCHASE)
			log.Println("Account not found")
		} else {
			Uld.Sugar().Infof("[%v] Buy space error: %v", LOG_TAG_PURCHASE, err)
			log.Println("cancel space Authorizition failed,please configure the correct account seed")
		}
		return
	}
	log.Println("Cancel space Authorizition success. Tx hash:", txhash)
}
