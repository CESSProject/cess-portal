package chain

import "github.com/centrifuge/go-substrate-rpc-client/v4/types"

var (
	//trade
	BuySpaceTransactionName   = "FileBank.buy_space"
	UploadFileTransactionName = "FileBank.upload"
	DeleteFileTransactionName = "FileBank.delete_file"

	//find
	PurchasedSpaceChainModule  = "FileBank"
	PurchasedSpaceModuleMethod = "UserHoldSpaceDetails"

	FindPriceChainModule  = "FileBank"
	FindPriceModuleMethod = "UnitPrice"

	FindFileChainModule  = "FileBank"
	FindFileModuleMethod = []string{"File", "UserHoldFileList"}

	FindSchedulerInfoModule = "FileMap"
	FindSchedulerInfoMethod = "SchedulerMap"
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
	//FileId      types.Bytes         `json:"acc"`         //File id
	File_Name   types.Bytes         `json:"file_name"`   //File name
	FileSize    types.U64           `json:"file_size"`   //File size
	FileHash    types.Bytes         `json:"file_hash"`   //File hash
	Public      types.Bool          `json:"public"`      //Public or not
	UserAddr    types.AccountID     `json:"user_addr"`   //Upload user's address
	FileState   types.Bytes         `json:"file_state"`  //File state
	Backups     types.U8            `json:"backups"`     //Number of backups
	Downloadfee types.U128          `json:"downloadfee"` //Download fee
	FileDupl    []FileDuplicateInfo `json:"file_dupl"`   //File backup information list
}
type FileDuplicateInfo struct {
	MinerId   types.U64
	BlockNum  types.U32
	ScanSize  types.U32
	Acc       types.AccountID
	MinerIp   types.Bytes
	DuplId    types.Bytes
	RandKey   types.Bytes
	BlockInfo []BlockInfo
}
type BlockInfo struct {
	BlockIndex types.Bytes
	BlockSize  types.U32
}

type FileList struct {
	Fileid types.Bytes8 `json:"fileid"`
}
type SchedulerInfo struct {
	Ip             types.Bytes     `json:"ip"`
	Owner          types.AccountID `json:"stash_user"`
	ControllerUser types.AccountID `json:"controller_user"`
}
