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

import (
	"sync"
	"sync/atomic"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var ChainClient Chainer

type Chainer interface {
	// Getpublickey returns its own public key
	GetPublicKey() []byte
	// Getpublickey returns its own private key
	GetMnemonicSeed() string
	// NewAccountId returns the account id
	NewAccountId(pubkey []byte) types.AccountID
	// GetSyncStatus returns whether the block is being synchronized
	GetSyncStatus() (bool, error)
	// GetChainStatus returns chain status
	GetChainStatus() bool
	// Getstorageminerinfo is used to get the details of the miner
	GetStorageMinerInfo(pkey []byte) (MinerInfo, error)
	// Getallstorageminer is used to obtain the AccountID of all miners
	GetAllStorageMiner() ([]types.AccountID, error)
	// GetFileMetaInfo is used to get the meta information of the file
	GetFileMetaInfo(fid string) (FileMetaInfo, error)
	// GetCessAccount is used to get the account in cess chain format
	GetCessAccount() (string, error)
	// GetAccountInfo is used to get account information
	GetAccountInfo(pkey []byte) (types.AccountInfo, error)

	// GetSchedulerList is used to get information about all schedules
	GetSchedulerList() ([]SchedulerInfo, error)
	// GetBucketList is used to obtain all buckets of the user
	GetBucketList(owner_pkey []byte) ([]types.Bytes, error)
	// GetBucketInfo is used to query bucket details
	GetBucketInfo(owner_pkey []byte, name string) (BucketInfo, error)
	// GetGrantor is used to query the user's space grantor
	GetGrantor(pkey []byte) (types.AccountID, error)
	// GetState is used to obtain OSS status information
	GetState(pubkey []byte) (string, error)
	//GetUserSpaceMetadata is used to query the user's space info
	GetUserSpaceMetadata(owner_pkey []byte) (SpacePackage, error)
	// Register is used to register oss services
	Register(ip, port string) (string, error)
	// Update is used to update the communication address of the scheduling service
	Update(ip, port string) (string, error)
	// CreateBucket is used to create a bucket for users
	CreateBucket(owner_pkey []byte, name string) (string, error)
	// DeleteBucket is used to delete buckets created by users
	DeleteBucket(owner_pkey []byte, name string) (string, error)
	//
	DeleteFile(owner_pkey []byte, filehash string) (string, error)
	//
	DeclarationFile(filehash string, user UserBrief) (string, error)
	//
	BuySpace(count types.U32) (string, error)
	//
	CancelAuth() (string, error)
	//
	AuthorizeSpace(owner_pkey []byte) (string, error)
}

type chainClient struct {
	lock            *sync.Mutex
	api             *gsrpc.SubstrateAPI
	chainState      *atomic.Bool
	metadata        *types.Metadata
	runtimeVersion  *types.RuntimeVersion
	keyEvents       types.StorageKey
	genesisHash     types.Hash
	keyring         signature.KeyringPair
	rpcAddr         string
	timeForBlockOut time.Duration
}

func NewChainClient(rpcAddr, secret string, t time.Duration) (Chainer, error) {
	var (
		err error
		cli = &chainClient{}
	)
	cli.api, err = gsrpc.NewSubstrateAPI(rpcAddr)
	if err != nil {
		return nil, err
	}
	cli.metadata, err = cli.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}
	cli.genesisHash, err = cli.api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}
	cli.runtimeVersion, err = cli.api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, err
	}
	cli.keyEvents, err = types.CreateStorageKey(
		cli.metadata,
		pallet_System,
		events,
		nil,
	)
	if err != nil {
		return nil, err
	}
	if secret != "" {
		cli.keyring, err = signature.KeyringPairFromSecret(secret, 0)
		if err != nil {
			return nil, err
		}
	}
	cli.lock = new(sync.Mutex)
	cli.chainState = &atomic.Bool{}
	cli.chainState.Store(true)
	cli.timeForBlockOut = t
	cli.rpcAddr = rpcAddr
	return cli, nil
}

func (c *chainClient) IsChainClientOk() bool {
	err := healthchek(c.api)
	if err != nil {
		c.api = nil
		cli, err := reconnectChainClient(c.rpcAddr)
		if err != nil {
			return false
		}
		c.api = cli
		return true
	}
	return true
}

func (c *chainClient) SetChainState(state bool) {
	c.chainState.Store(state)
}

func (c *chainClient) GetChainState() bool {
	return c.chainState.Load()
}

func (c *chainClient) NewAccountId(pubkey []byte) types.AccountID {
	return types.NewAccountID(pubkey)
}

func reconnectChainClient(rpcAddr string) (*gsrpc.SubstrateAPI, error) {
	return gsrpc.NewSubstrateAPI(rpcAddr)
}

func healthchek(a *gsrpc.SubstrateAPI) error {
	defer func() {
		recover()
	}()
	_, err := a.RPC.System.Health()
	return err
}
