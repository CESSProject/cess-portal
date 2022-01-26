package chain

import "github.com/centrifuge/go-substrate-rpc-client/v4/types"

var (
	//trade
	BuySpaceTransactionName   = "FileBank.buy_space"
	UploadFileTransactionName = "FileBank.upload"

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
	Filename    types.Bytes     `json:"filename"` //
	Owner       types.AccountID `json:"owner"`    //
	Filehash    types.Bytes     `json:"filehash"` //
	Backups     types.U8        `json:"backups"`
	Filesize    types.U128      `json:"filesize"`
	Downloadfee types.U128      `json:"downloadfee"`
}

type FileList struct {
	Fileid types.Bytes8 `json:"fileid"`
}

//On-chain event analysis param
type Event_UnsignedPhaseStarted struct {
	Phase  types.Phase
	Round  types.U32
	Topics []types.Hash
}
type Event_SolutionStored struct {
	Phase            types.Phase
	Election_compute types.ElectionCompute
	Prev_ejected     types.Bool
	Topics           []types.Hash
}

type MyEventRecords struct {
	types.EventRecords

	//
	ElectionProviderMultiPhase_UnsignedPhaseStarted []Event_UnsignedPhaseStarted
	ElectionProviderMultiPhase_SolutionStored       []Event_SolutionStored
}
