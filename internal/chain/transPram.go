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

	FindPriceChainModule  = "Sminer"
	FindPriceModuleMethod = []string{"PurchasedSpace", "AvailableSpace"}

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
	DuplId    types.Bytes     `json:"dupl_id"`    //Backup id
	RandKey   types.Bytes     `json:"rand_key"`   //Random key
	SliceNum  types.U16       `json:"slice_num"`  //Number of slices
	FileSlice []FileSliceInfo `json:"file_slice"` //Slice information list
}

type FileSliceInfo struct {
	SliceId   types.Bytes   `json:"slice_id"`   //Slice id
	SliceSize types.U32     `json:"slice_size"` //Slice size
	SliceHash types.Bytes   `json:"slice_hash"` //Slice hash
	FileShard FileShardInfo `json:"file_shard"` //Shard information
}

type FileShardInfo struct {
	DataShardNum  types.U8      `json:"data_shard_num"`  //Number of data shard
	RedunShardNum types.U8      `json:"redun_shard_num"` //Number of redundant shard
	ShardHash     []types.Bytes `json:"shard_hash"`      //Shard hash list
	ShardAddr     []types.Bytes `json:"shard_addr"`      //Store miner service addr list
	Peerid        []types.U64   `json:"wallet_addr"`     //Store miner wallet addr list
}

type FileList struct {
	Fileid types.Bytes8 `json:"fileid"`
}
type SchedulerInfo struct {
	Ip             types.Bytes     `json:"ip"`
	Owner          types.AccountID `json:"stash_user"`
	ControllerUser types.AccountID `json:"controller_user"`
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

type Event_SegmentBook_ParamSet struct {
	Phase     types.Phase
	PeerId    types.U64
	SegmentId types.U64
	Random    types.U32
	Topics    []types.Hash
}

type Event_VPABCD_Submit_Verify struct {
	Phase     types.Phase
	PeerId    types.U64
	SegmentId types.U64
	Topics    []types.Hash
}

type Event_Sminer_TimedTask struct {
	Phase  types.Phase
	Topics []types.Hash
}

type Event_Sminer_Registered struct {
	Phase   types.Phase
	PeerAcc types.AccountID
	Staking types.U128
	Topics  []types.Hash
}

type Event_FileMap_RegistrationScheduler struct {
	Phase  types.Phase
	Acc    types.AccountID
	Ip     types.Bytes
	Topics []types.Hash
}

type Event_DeleteFile struct {
	Phase  types.Phase
	Acc    types.AccountID
	Fileid types.Bytes
	Topics []types.Hash
}

type Event_BuySpace struct {
	Phase  types.Phase
	Acc    types.AccountID
	Size   types.U128
	Fee    types.U128
	Topics []types.Hash
}

type Event_FileUpload struct {
	Phase  types.Phase
	Acc    types.AccountID
	Topics []types.Hash
}

type MyEventRecords struct {
	types.EventRecords

	SegmentBook_ParamSet          []Event_SegmentBook_ParamSet
	SegmentBook_VPASubmitted      []Event_VPABCD_Submit_Verify
	SegmentBook_VPBSubmitted      []Event_VPABCD_Submit_Verify
	SegmentBook_VPCSubmitted      []Event_VPABCD_Submit_Verify
	SegmentBook_VPDSubmitted      []Event_VPABCD_Submit_Verify
	SegmentBook_VPAVerified       []Event_VPABCD_Submit_Verify
	SegmentBook_VPBVerified       []Event_VPABCD_Submit_Verify
	SegmentBook_VPCVerified       []Event_VPABCD_Submit_Verify
	SegmentBook_VPDVerified       []Event_VPABCD_Submit_Verify
	Sminer_TimedTask              []Event_Sminer_TimedTask
	Sminer_Registered             []Event_Sminer_Registered
	FileMap_RegistrationScheduler []Event_FileMap_RegistrationScheduler

	FileBank_DeleteFile []Event_DeleteFile
	FileBank_BuySpace   []Event_BuySpace
	FileBank_FileUpload []Event_FileUpload

	ElectionProviderMultiPhase_UnsignedPhaseStarted []Event_UnsignedPhaseStarted
	ElectionProviderMultiPhase_SolutionStored       []Event_SolutionStored
}
