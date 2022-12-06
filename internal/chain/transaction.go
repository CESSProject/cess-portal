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
	"cess-portal/conf"
	"cess-portal/tools"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

func (c *chainClient) Register(ip, port string) (string, error) {
	var (
		txhash      string
		ipType      IpAddress
		accountInfo types.AccountInfo
	)

	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.IsChainClientOk() {
		c.SetChainState(false)
		return txhash, ERR_RPC_CONNECTION
	}
	c.SetChainState(true)

	if tools.IsIPv4(ip) {
		ipType.IPv4.Index = 0
		ips := strings.Split(ip, ".")
		for i := 0; i < len(ipType.IPv4.Value); i++ {
			temp, _ := strconv.Atoi(ips[i])
			ipType.IPv4.Value[i] = types.U8(temp)
		}
		temp, _ := strconv.Atoi(port)
		ipType.IPv4.Port = types.U16(temp)
	} else {
		return txhash, ERR_RPC_IP_FORMAT
	}

	call, err := types.NewCall(
		c.metadata,
		OssRegister,
		ipType.IPv4,
	)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewCall]")
	}

	ext := types.NewExtrinsic(call)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewExtrinsic]")
	}

	key, err := types.CreateStorageKey(
		c.metadata,
		pallet_System,
		account,
		c.keyring.PublicKey,
	)
	if err != nil {
		return txhash, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "[GetStorageLatest]")
	}

	if !ok {
		return txhash, ERR_RPC_EMPTY_VALUE
	}

	o := types.SignatureOptions{
		BlockHash:          c.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        c.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        c.runtimeVersion.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: c.runtimeVersion.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(c.keyring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "[Sign]")
	}

	// Do the transfer and track the actual status
	sub, err := c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		if !strings.Contains(err.Error(), "Priority is too low") {
			return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
		}
		var tryCount = 0
		for tryCount < 20 {
			o.Nonce = types.NewUCompactFromUInt(uint64(accountInfo.Nonce + types.NewU32(1)))
			// Sign the transaction
			err = ext.Sign(c.keyring, o)
			if err != nil {
				return txhash, errors.Wrap(err, "[Sign]")
			}
			sub, err = c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
			if err == nil {
				break
			}
			tryCount++
		}
	}
	if err != nil {
		return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
	}
	defer sub.Unsubscribe()
	timeout := time.After(c.timeForBlockOut)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := CessEventRecords{}
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				h, err := c.api.RPC.State.GetStorageRaw(c.keyEvents, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "[GetStorageRaw]")
				}

				types.EventRecordsRaw(*h).DecodeEventRecords(c.metadata, &events)

				if len(events.Oss_OssRegister) > 0 {
					return txhash, nil
				}
				return txhash, errors.New(ERR_Failed)
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "[sub]")
		case <-timeout:
			return txhash, ERR_RPC_TIMEOUT
		}
	}
}

func (c *chainClient) Update(ip, port string) (string, error) {
	var (
		txhash      string
		ipType      IpAddress
		accountInfo types.AccountInfo
	)

	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.IsChainClientOk() {
		c.SetChainState(false)
		return txhash, ERR_RPC_CONNECTION
	}
	c.SetChainState(true)

	if tools.IsIPv4(ip) {
		ipType.IPv4.Index = 0
		ips := strings.Split(ip, ".")
		for i := 0; i < len(ipType.IPv4.Value); i++ {
			temp, _ := strconv.Atoi(ips[i])
			ipType.IPv4.Value[i] = types.U8(temp)
		}
		temp, _ := strconv.Atoi(port)
		ipType.IPv4.Port = types.U16(temp)
	} else {
		return txhash, ERR_RPC_IP_FORMAT
	}

	call, err := types.NewCall(
		c.metadata,
		OssUpdate,
		ipType.IPv4,
	)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewCall]")
	}

	ext := types.NewExtrinsic(call)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewExtrinsic]")
	}

	key, err := types.CreateStorageKey(
		c.metadata,
		pallet_System,
		account,
		c.keyring.PublicKey,
	)
	if err != nil {
		return txhash, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return txhash, ERR_RPC_EMPTY_VALUE
	}

	o := types.SignatureOptions{
		BlockHash:          c.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        c.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        c.runtimeVersion.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: c.runtimeVersion.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(c.keyring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "[Sign]")
	}

	// Do the transfer and track the actual status
	sub, err := c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		if !strings.Contains(err.Error(), "Priority is too low") {
			return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
		}
		var tryCount = 0
		for tryCount < 20 {
			o.Nonce = types.NewUCompactFromUInt(uint64(accountInfo.Nonce + types.NewU32(1)))
			// Sign the transaction
			err = ext.Sign(c.keyring, o)
			if err != nil {
				return txhash, errors.Wrap(err, "[Sign]")
			}
			sub, err = c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
			if err == nil {
				break
			}
			tryCount++
		}
	}
	if err != nil {
		return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
	}
	defer sub.Unsubscribe()
	timeout := time.After(c.timeForBlockOut)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := CessEventRecords{}
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				h, err := c.api.RPC.State.GetStorageRaw(c.keyEvents, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "[GetStorageRaw]")
				}

				types.EventRecordsRaw(*h).DecodeEventRecords(c.metadata, &events)

				if len(events.Oss_OssUpdate) > 0 {
					return txhash, nil
				}
				return txhash, errors.New(ERR_Failed)
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "[sub]")
		case <-timeout:
			return txhash, ERR_RPC_TIMEOUT
		}
	}
}

func (c *chainClient) CreateBucket(owner_pkey []byte, name string) (string, error) {
	var (
		txhash      string
		accountInfo types.AccountInfo
	)

	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.IsChainClientOk() {
		c.SetChainState(false)
		return txhash, ERR_RPC_CONNECTION
	}
	c.SetChainState(true)

	// b, err := types.Encode(name)
	// if err != nil {
	// 	return txhash, errors.Wrap(err, "[Encode]")
	// }
	call, err := types.NewCall(
		c.metadata,
		FileBank_CreateBucket,
		types.NewAccountID(owner_pkey),
		types.NewBytes([]byte(name)),
	)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewCall]")
	}

	ext := types.NewExtrinsic(call)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewExtrinsic]")
	}

	key, err := types.CreateStorageKey(
		c.metadata,
		pallet_System,
		account,
		c.keyring.PublicKey,
	)
	if err != nil {
		return txhash, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return txhash, ERR_RPC_EMPTY_VALUE
	}

	o := types.SignatureOptions{
		BlockHash:          c.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        c.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        c.runtimeVersion.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: c.runtimeVersion.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(c.keyring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "[Sign]")
	}

	// Do the transfer and track the actual status
	sub, err := c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		if !strings.Contains(err.Error(), "Priority is too low") {
			return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
		}
		var tryCount = 0
		for tryCount < 20 {
			o.Nonce = types.NewUCompactFromUInt(uint64(accountInfo.Nonce + types.NewU32(1)))
			// Sign the transaction
			err = ext.Sign(c.keyring, o)
			if err != nil {
				return txhash, errors.Wrap(err, "[Sign]")
			}
			sub, err = c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
			if err == nil {
				break
			}
			tryCount++
		}
	}
	if err != nil {
		return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
	}
	defer sub.Unsubscribe()
	timeout := time.After(c.timeForBlockOut)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := CessEventRecords{}
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				h, err := c.api.RPC.State.GetStorageRaw(c.keyEvents, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "[GetStorageRaw]")
				}

				types.EventRecordsRaw(*h).DecodeEventRecords(c.metadata, &events)

				if len(events.FileBank_CreateBucket) > 0 {
					return txhash, nil
				}
				return txhash, errors.New(ERR_Failed)
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "[sub]")
		case <-timeout:
			return txhash, ERR_RPC_TIMEOUT
		}
	}
}

func (c *chainClient) DeleteBucket(owner_pkey []byte, name string) (string, error) {
	var (
		txhash      string
		accountInfo types.AccountInfo
	)
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.IsChainClientOk() {
		c.SetChainState(false)
		return txhash, ERR_RPC_CONNECTION
	}
	c.SetChainState(true)

	call, err := types.NewCall(
		c.metadata,
		FileBank_DeleteBucket,
		types.NewAccountID(owner_pkey),
		types.NewBytes([]byte(name)),
	)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewCall]")
	}

	ext := types.NewExtrinsic(call)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewExtrinsic]")
	}

	key, err := types.CreateStorageKey(
		c.metadata,
		pallet_System,
		account,
		c.keyring.PublicKey,
	)
	if err != nil {
		return txhash, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return txhash, ERR_RPC_EMPTY_VALUE
	}

	o := types.SignatureOptions{
		BlockHash:          c.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        c.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        c.runtimeVersion.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: c.runtimeVersion.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(c.keyring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "[Sign]")
	}

	// Do the transfer and track the actual status
	sub, err := c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		if !strings.Contains(err.Error(), "Priority is too low") {
			return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
		}
		var tryCount = 0
		for tryCount < 20 {
			o.Nonce = types.NewUCompactFromUInt(uint64(accountInfo.Nonce + types.NewU32(1)))
			// Sign the transaction
			err = ext.Sign(c.keyring, o)
			if err != nil {
				return txhash, errors.Wrap(err, "[Sign]")
			}
			sub, err = c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
			if err == nil {
				break
			}
			tryCount++
		}
	}
	if err != nil {
		return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
	}
	defer sub.Unsubscribe()
	timeout := time.After(c.timeForBlockOut)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := CessEventRecords{}
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				h, err := c.api.RPC.State.GetStorageRaw(c.keyEvents, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "[GetStorageRaw]")
				}

				types.EventRecordsRaw(*h).DecodeEventRecords(c.metadata, &events)

				if len(events.FileBank_DeleteBucket) > 0 {
					return txhash, nil
				}
				return txhash, errors.New(ERR_Failed)
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "[sub]")
		case <-timeout:
			return txhash, ERR_RPC_TIMEOUT
		}
	}
}

func (c *chainClient) DeclarationFile(filehash string, user UserBrief) (string, error) {
	var (
		txhash      string
		accountInfo types.AccountInfo
	)

	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.IsChainClientOk() {
		c.SetChainState(false)
		return txhash, ERR_RPC_CONNECTION
	}
	c.SetChainState(true)

	var hash FileHash
	if len(filehash) != len(hash) {
		return txhash, errors.New("invalid filehash")
	}
	for i := 0; i < len(hash); i++ {
		hash[i] = types.U8(filehash[i])
	}

	call, err := types.NewCall(
		c.metadata,
		FileBank_UploadDeclaration,
		hash,
		user,
	)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewCall]")
	}

	ext := types.NewExtrinsic(call)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewExtrinsic]")
	}

	key, err := types.CreateStorageKey(
		c.metadata,
		pallet_System,
		account,
		c.keyring.PublicKey,
	)
	if err != nil {
		return txhash, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return txhash, ERR_RPC_EMPTY_VALUE
	}

	o := types.SignatureOptions{
		BlockHash:          c.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        c.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        c.runtimeVersion.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: c.runtimeVersion.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(c.keyring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "[Sign]")
	}

	// Do the transfer and track the actual status
	sub, err := c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		if !strings.Contains(err.Error(), "Priority is too low") {
			return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
		}
		var tryCount = 0
		for tryCount < 20 {
			o.Nonce = types.NewUCompactFromUInt(uint64(accountInfo.Nonce + types.NewU32(1)))
			// Sign the transaction
			err = ext.Sign(c.keyring, o)
			if err != nil {
				return txhash, errors.Wrap(err, "[Sign]")
			}
			sub, err = c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
			if err == nil {
				break
			}
			tryCount++
		}
	}
	if err != nil {
		return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
	}
	defer sub.Unsubscribe()
	timeout := time.After(c.timeForBlockOut)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := CessEventRecords{}
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				h, err := c.api.RPC.State.GetStorageRaw(c.keyEvents, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "[GetStorageRaw]")
				}

				types.EventRecordsRaw(*h).DecodeEventRecords(c.metadata, &events)

				if len(events.FileBank_UploadDeclaration) > 0 {
					return txhash, nil
				}
				return txhash, errors.New(ERR_Failed)
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "[sub]")
		case <-timeout:
			return txhash, ERR_RPC_TIMEOUT
		}
	}
}

func (c *chainClient) DeleteFile(owner_pkey []byte, filehash string) (string, error) {
	var (
		txhash      string
		accountInfo types.AccountInfo
	)

	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.IsChainClientOk() {
		c.SetChainState(false)
		return txhash, ERR_RPC_CONNECTION
	}
	c.SetChainState(true)

	var hash FileHash
	if len(filehash) != len(hash) {
		return txhash, errors.New("invalid filehash")
	}
	for i := 0; i < len(hash); i++ {
		hash[i] = types.U8(filehash[i])
	}

	call, err := types.NewCall(
		c.metadata,
		FileBank_DeleteFile,
		types.NewAccountID(owner_pkey),
		hash,
	)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewCall]")
	}

	ext := types.NewExtrinsic(call)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewExtrinsic]")
	}

	key, err := types.CreateStorageKey(
		c.metadata,
		pallet_System,
		account,
		c.keyring.PublicKey,
	)
	if err != nil {
		return txhash, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return txhash, ERR_RPC_EMPTY_VALUE
	}

	o := types.SignatureOptions{
		BlockHash:          c.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        c.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        c.runtimeVersion.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: c.runtimeVersion.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(c.keyring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "[Sign]")
	}

	// Do the transfer and track the actual status
	sub, err := c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		if !strings.Contains(err.Error(), "Priority is too low") {
			return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
		}
		var tryCount = 0
		for tryCount < 20 {
			o.Nonce = types.NewUCompactFromUInt(uint64(accountInfo.Nonce + types.NewU32(1)))
			// Sign the transaction
			err = ext.Sign(c.keyring, o)
			if err != nil {
				return txhash, errors.Wrap(err, "[Sign]")
			}
			sub, err = c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
			if err == nil {
				break
			}
			tryCount++
		}
	}
	if err != nil {
		return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
	}
	defer sub.Unsubscribe()
	timeout := time.After(c.timeForBlockOut)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := CessEventRecords{}
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				h, err := c.api.RPC.State.GetStorageRaw(c.keyEvents, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "[GetStorageRaw]")
				}

				types.EventRecordsRaw(*h).DecodeEventRecords(c.metadata, &events)

				if len(events.FileBank_DeleteFile) > 0 {
					return txhash, nil
				}
				return txhash, errors.New(ERR_Failed)
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "[sub]")
		case <-timeout:
			return txhash, ERR_RPC_TIMEOUT
		}
	}
}

func (c *chainClient) BuySpace(count types.U32) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%v", tools.RecoverError(err))
		}
	}()

	var txhash string
	var accountInfo types.AccountInfo

	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.IsChainClientOk() {
		c.SetChainState(false)
		return txhash, ERR_RPC_CONNECTION
	}
	c.SetChainState(true)
	newCall, err := types.NewCall(c.metadata, FileBank_BuySpace, count)
	if err != nil {
		return txhash, errors.Wrap(err, "NewCall")
	}

	ext := types.NewExtrinsic(newCall)
	if err != nil {
		return txhash, errors.Wrap(err, "NewExtrinsic")
	}
	key, err := types.CreateStorageKey(c.metadata, "System", "Account", conf.PublicKey, nil)
	if err != nil {
		return txhash, errors.Wrap(err, "CreateStorageKey")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "GetStorageLatest")
	}

	if !ok {
		return txhash, errors.New(ERR_Empty)
	}

	o := types.SignatureOptions{
		BlockHash:          c.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        c.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        c.runtimeVersion.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: c.runtimeVersion.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(c.keyring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "Sign")
	}

	// Do the transfer and track the actual status
	sub, err := c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return txhash, errors.Wrap(err, "SubmitAndWatchExtrinsic")
	}

	defer sub.Unsubscribe()
	timeout := time.After(conf.TimeToWaitEvents)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := CessEventRecords{}
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				h, err := c.api.RPC.State.GetStorageRaw(c.keyEvents, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "GetStorageRaw")
				}
				err = types.EventRecordsRaw(*h).DecodeEventRecords(c.metadata, &events)
				if err != nil {
					log.Printf("[%v]Decode event err:%v", txhash, err)
				}

				if len(events.FileBank_BuySpace) > 0 {
					return txhash, nil
				}

				return txhash, errors.New(ERR_Failed)
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "<-sub")
		case <-timeout:
			return txhash, errors.New(ERR_Timeout)
		}
	}
}

func (c *chainClient) AuthorizeSpace(owner_pkey []byte) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%v", tools.RecoverError(err))
		}
	}()

	var txhash string
	var accountInfo types.AccountInfo

	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.IsChainClientOk() {
		c.SetChainState(false)
		return txhash, ERR_RPC_CONNECTION
	}
	c.SetChainState(true)
	newCall, err := types.NewCall(c.metadata, Oss_AuthSpace, types.NewAccountID(owner_pkey))
	if err != nil {
		return txhash, errors.Wrap(err, "NewCall")
	}

	ext := types.NewExtrinsic(newCall)
	if err != nil {
		return txhash, errors.Wrap(err, "NewExtrinsic")
	}
	key, err := types.CreateStorageKey(c.metadata, "System", "Account", conf.PublicKey, nil)
	if err != nil {
		return txhash, errors.Wrap(err, "CreateStorageKey")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "GetStorageLatest")
	}

	if !ok {
		return txhash, errors.New(ERR_Empty)
	}

	o := types.SignatureOptions{
		BlockHash:          c.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        c.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        c.runtimeVersion.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: c.runtimeVersion.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(c.keyring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "Sign")
	}

	// Do the transfer and track the actual status
	sub, err := c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return txhash, errors.Wrap(err, "SubmitAndWatchExtrinsic")
	}

	defer sub.Unsubscribe()
	timeout := time.After(conf.TimeToWaitEvents)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := CessEventRecords{}
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				h, err := c.api.RPC.State.GetStorageRaw(c.keyEvents, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "GetStorageRaw")
				}
				err = types.EventRecordsRaw(*h).DecodeEventRecords(c.metadata, &events)
				if err != nil {
					log.Printf("[%v]Decode event err:%v", txhash, err)
				}

				if len(events.Oss_Authorize) > 0 {
					return txhash, nil
				}

				return txhash, errors.New(ERR_Failed)
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "<-sub")
		case <-timeout:
			return txhash, errors.New(ERR_Timeout)
		}
	}
}

func (c *chainClient) CancelAuth() (string, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%v", tools.RecoverError(err))
		}
	}()

	var txhash string
	var accountInfo types.AccountInfo

	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.IsChainClientOk() {
		c.SetChainState(false)
		return txhash, ERR_RPC_CONNECTION
	}
	c.SetChainState(true)
	newCall, err := types.NewCall(c.metadata, Oss_CancelAuthorize)
	if err != nil {
		return txhash, errors.Wrap(err, "NewCall")
	}

	ext := types.NewExtrinsic(newCall)
	if err != nil {
		return txhash, errors.Wrap(err, "NewExtrinsic")
	}
	key, err := types.CreateStorageKey(c.metadata, "System", "Account", conf.PublicKey, nil)
	if err != nil {
		return txhash, errors.Wrap(err, "CreateStorageKey")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "GetStorageLatest")
	}

	if !ok {
		return txhash, errors.New(ERR_Empty)
	}

	o := types.SignatureOptions{
		BlockHash:          c.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        c.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        c.runtimeVersion.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: c.runtimeVersion.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(c.keyring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "Sign")
	}

	// Do the transfer and track the actual status
	sub, err := c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return txhash, errors.Wrap(err, "SubmitAndWatchExtrinsic")
	}

	defer sub.Unsubscribe()
	timeout := time.After(conf.TimeToWaitEvents)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := CessEventRecords{}
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				h, err := c.api.RPC.State.GetStorageRaw(c.keyEvents, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "GetStorageRaw")
				}
				err = types.EventRecordsRaw(*h).DecodeEventRecords(c.metadata, &events)
				if err != nil {
					log.Printf("[%v]Decode event err:%v", txhash, err)
				}

				if len(events.Oss_CancelAuthorize) > 0 {
					return txhash, nil
				}

				return txhash, errors.New(ERR_Failed)
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "<-sub")
		case <-timeout:
			return txhash, errors.New(ERR_Timeout)
		}
	}
}
