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

	publickey, err := types.NewMultiAddressFromHexAccountID(conf.ClientConf.ChainData.AccountPublicKey)
	if err != nil {
		return data, err
	}
	key, err := types.CreateStorageKey(meta, ci.ChainModule, ci.ChainModuleMethod, publickey.AsID[:])
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
	ret += "                   PurchasedSpace:" + userinfo.PurchasedSpace.String() + "\n"
	ret += "                   UsedSpace:" + userinfo.RemainingSpace.String() + "\n"
	ret += "                   RemainingSpace:" + userinfo.RemainingSpace.String() + "\n"
	ret += "—————————————————————————————————————————————————————————"
	return ret
}

//Query the currently stored unit price.Calculation:(AvailableSpace - PurchasedSpace) / 1024 * 1000
func (ci *CessInfo) GetPurchasedSpace() (types.U128, error) {
	var (
		err  error
		data types.U128
	)
	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		err := recover()
		if err != nil {
			fmt.Printf("[Error]Recover UserHoldSpaceDetails panic :%s\n", err)
		}
	}()
	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}

	key, err := types.CreateStorageKey(meta, ci.ChainModule, ci.ChainModuleMethod)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", ci.ChainModule, ci.ChainModuleMethod)
	}

	_, err = api.r.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}
	return data, nil
}

//Query the currently stored unit price.Calculation:(AvailableSpace - PurchasedSpace) / 1024 * 1000
func (ci *CessInfo) GetAvailableSpace() (types.U128, error) {
	var (
		err  error
		data types.U128
	)
	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		err := recover()
		if err != nil {
			fmt.Printf("[Error]Recover UserHoldSpaceDetails panic :%s\n", err)
		}
	}()
	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}

	key, err := types.CreateStorageKey(meta, ci.ChainModule, ci.ChainModuleMethod)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", ci.ChainModule, ci.ChainModuleMethod)
	}

	_, err = api.r.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}
	return data, nil
}

func (ci *CessInfo) GetFileInfo(fileid string) (FileInfo, error) {
	var (
		err  error
		data FileInfo
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

	key, err := types.CreateStorageKey(meta, ci.ChainModule, ci.ChainModuleMethod, types.NewBytes([]byte(fileid))[:])
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", ci.ChainModule, ci.ChainModuleMethod)
	}

	_, err = api.r.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}
	return data, nil
}

func (ci *CessInfo) GetFileList() ([]FileList, error) {
	var (
		err  error
		data = make([]FileList, 0)
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
	publickey, err := types.NewMultiAddressFromHexAccountID(conf.ClientConf.ChainData.AccountPublicKey)
	if err != nil {
		return data, err
	}

	key, err := types.CreateStorageKey(meta, ci.ChainModule, ci.ChainModuleMethod, publickey.AsID[:])
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", ci.ChainModule, ci.ChainModuleMethod)
	}

	_, err = api.r.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}
	return data, nil
}

func (fileinfo FileInfo) String() string {
	ret := "———————————————————File Information———————————————————\n"
	ret += fmt.Sprintf("                  FileOwner:%x\n", string(fileinfo.Owner[:]))
	ret += fmt.Sprintf("                  Filename:%x\n", string(fileinfo.Filename[:]))
	ret += fmt.Sprintf("                  Filehash:%x\n", string(fileinfo.Filehash[:]))
	ret += fmt.Sprintf("                  Backups:%x\n", string(fileinfo.Backups))
	ret += fmt.Sprintf("                  Filesize:%f\n", float64(fileinfo.Filesize.Int64()))
	ret += fmt.Sprintf("                  Downloadfee:%f\n", float64(fileinfo.Downloadfee.Int64()))
	return ret
}
