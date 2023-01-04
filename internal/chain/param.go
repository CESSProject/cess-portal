/*
   Copyright 2022 CESS scheduler authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package chain

// Pallets
const (
	pallet_FileBank    = "FileBank"
	pallet_FileMap     = "FileMap"
	pallet_Sminer      = "Sminer"
	pallet_SegmentBook = "SegmentBook"
	pallet_System      = "System"
	pallet_Oss         = "Oss"
)

// Pallet's method
const (
	// System
	account = "Account"
	events  = "Events"

	// Sminer
	allMinerItems = "AllMiner"
	minerItems    = "MinerItems"
	segInfo       = "SegInfo"

	// FileMap
	fileMetaInfo = "File"
	schedulerMap = "SchedulerMap"

	// FileBank
	fileBank_UserFilelist   = "UserHoldFileList"
	fileBank_Bucket         = "Bucket"
	fileBank_BucketList     = "UserBucketList"
	fileBank_BuySpace       = "BuySpace"
	fileBank_userOwnedSpace = "UserOwnedSpace"
	// Oss
	oss     = "Oss"
	Grantor = "AuthorityList"
)

// Extrinsics
const (
	// FileBank
	tx_FileBank_Update         = "FileBank.update"
	tx_FileBank_Upload         = "FileBank.upload"
	FileBank_CreateBucket      = "FileBank.create_bucket"
	FileBank_DeleteBucket      = "FileBank.delete_bucket"
	FileBank_DeleteFile        = "FileBank.delete_file"
	FileBank_UploadDeclaration = "FileBank.upload_declaration"
	FileBank_BuySpace          = "FileBank.buy_space"
	Oss_AuthSpace              = "Oss.authorize"
	Oss_CancelAuthorize        = "Oss.cancel_authorize"
	// Oss
	OssRegister = "Oss.register"
	OssUpdate   = "Oss.update"
)

const (
	FILE_STATE_ACTIVE  = "active"
	FILE_STATE_PENDING = "pending"
)

const (
	MINER_STATE_POSITIVE = "positive"
	MINER_STATE_FROZEN   = "frozen"
	MINER_STATE_EXIT     = "exit"
)
