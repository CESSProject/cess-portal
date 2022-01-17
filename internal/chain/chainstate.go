package chain

import (
	"dapp_cess_client/conf"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

// Get all miner information on the cess chain
func (ci *CessInfo) UserHoldSpaceDetails() (UserHoldSpaceDetails, error) {
	var (
		err  error
		data UserHoldSpaceDetails
	)
	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		err := recover()
		if err != nil {
			fmt.Printf("[Error]Recover UserHoldSpaceDetails panic fail :%s\n", err)
		}
	}()
	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}

	key, err := types.CreateStorageKey(meta, ci.ChainModule, ci.ChainModuleMethod, []byte(conf.ClientConf.ChainData.AccountPublicKey))
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", ci.ChainModule, ci.ChainModuleMethod)
	}

	_, err = api.r.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}
	return data, nil
}

func (userinfo UserHoldSpaceDetails) String() string {
	ret := "———————————————————You Purchased Space———————————————————\n"
	ret += "         PurchasedSpace:" + userinfo.PurchasedSpace.String() + "\n"
	ret += "         UsedSpace:" + userinfo.RemainingSpace.String() + "\n"
	ret += "         RemainingSpace:" + userinfo.RemainingSpace.String() + "\n"
	ret += "—————————————————————————————————————————————————————————"
	return ret
}
