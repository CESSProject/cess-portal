package chain

import (
	"cess-portal/conf"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	"strconv"
)

//UserHoldSpaceDetails means to get specific information about user space
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
	PurchasedSpace, _ := strconv.Atoi(userinfo.PurchasedSpace.String())
	UsedSpace, _ := strconv.Atoi(userinfo.UsedSpace.String())
	RemainingSpace, _ := strconv.Atoi(userinfo.RemainingSpace.String())
	ret := "———————————————————You Purchased Space———————————————————\n"
	ret += "                   PurchasedSpace:" + strconv.Itoa(PurchasedSpace/1024) + "(MB)\n"
	ret += "                   UsedSpace:" + strconv.Itoa(UsedSpace/1024) + "(MB)\n"
	ret += "                   RemainingSpace:" + strconv.Itoa(RemainingSpace/1024) + "(MB)\n"
	ret += "—————————————————————————————————————————————————————————"
	return ret
}

//GetPurchasedSpace means the size of the space purchased by all customers of the whole CESS system
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

//GetAvailableSpace Means the purchaseable space of the whole CESS system
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

//GetFileInfo means to get the specific information of the file through the current fileid
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
	id, err := types.EncodeToBytes(fileid)

	key, err := types.CreateStorageKey(meta, ci.ChainModule, ci.ChainModuleMethod, types.NewBytes(id))
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", ci.ChainModule, ci.ChainModuleMethod)
	}

	_, err = api.r.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}
	return data, nil
}

//GetFileList means to get a list of all files of the current user
func (ci *CessInfo) GetFileList() ([][]byte, error) {
	var (
		err  error
		data = make([][]byte, 0)
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
	ret += fmt.Sprintf("                  Filename:%v\n", string(fileinfo.File_Name[:]))
	ret += fmt.Sprintf("                  Public:%v\n", fileinfo.Public)
	ret += fmt.Sprintf("                  Filehash:%v\n", string(fileinfo.FileHash[:]))
	ret += fmt.Sprintf("                  Backups:%v\n", fileinfo.Backups)
	ret += fmt.Sprintf("                  Filesize:%v\n", fileinfo.FileSize)
	ret += fmt.Sprintf("                  Downloadfee:%v\n", fileinfo.Downloadfee)
	return ret
}

//GetSchedulerInfo means to get all currently registered schedulers
func (ci *CessInfo) GetSchedulerInfo() ([]SchedulerInfo, error) {
	var (
		err  error
		data []SchedulerInfo
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

	//publickey, err := types.NewMultiAddressFromHexAccountID(conf.ClientConf.ChainData.AccountPublicKey)
	//if err != nil {
	//	return data, err
	//}
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
