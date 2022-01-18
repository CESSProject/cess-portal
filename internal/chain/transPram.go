package chain

import "github.com/centrifuge/go-substrate-rpc-client/v4/types"

var (
	//trade
	BuySpaceTransactionName = "FileBank.buy_space"

	//find
	PurchasedSpaceChainModule  = "FileBank"
	PurchasedSpaceModuleMethod = "UserHoldSpaceDetails"

	FindPriceChainModule  = "Sminer"
	FindPriceModuleMethod = []string{"PurchasedSpace", "AvailableSpace"}

	FindFileChainModule  = "FileBank"
	FindFileModuleMethod = []string{"File", "UserHoldFileList"}
)

type CessInfo struct {
	RpcAddr               string
	IdentifyAccountPhrase string
	TransactionName       string
	ChainModule           string
	ChainModuleMethod     string
}

type UserHoldSpaceDetails struct {
	PurchasedSpace types.U128 `json:"purchased_space"`
	UsedSpace      types.U128 `json:"used_space"`
	RemainingSpace types.U128 `json:"remaining_space"`
}

type FileInfo struct {
	Filename    types.Bytes8    `json:"filename"`
	Owner       types.AccountID `json:"owner"`
	Filehash    types.Bytes8    `json:"filehash"`
	Backups     types.U8        `json:"backups"`
	Filesize    types.U128      `json:"filesize"`
	Downloadfee types.U128      `json:"downloadfee"`
}

type FileList struct {
	Fileid types.U8 `json:"fileid"`
}
