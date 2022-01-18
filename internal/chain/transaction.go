package chain

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	"math/big"
	"time"
)

// etcd register
func (ci *CessInfo) BuySpaceOnChain(Quantity, Expected int) (string, error) {
	var (
		err         error
		accountInfo types.AccountInfo
	)
	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		//err := recover()
		//if err != nil {
		//	fmt.Printf("[Error]Recover UserHoldSpaceDetails panic :%s\n", err)
		//}
	}()
	keyring, err := signature.KeyringPairFromSecret(ci.IdentifyAccountPhrase, 0)
	if err != nil {
		return "", errors.Wrap(err, "KeyringPairFromSecret err")
	}

	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return "", errors.Wrap(err, "GetMetadataLatest err")
	}

	if err != nil {
		return "", err
	}

	c, err := types.NewCall(meta, ci.TransactionName, types.NewU128(*big.NewInt(int64(Quantity))), types.NewU128(*big.NewInt(int64(Expected))))
	if err != nil {
		return "", errors.Wrap(err, "NewCall err")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return "", errors.Wrap(err, "NewExtrinsic err")
	}

	genesisHash, err := api.r.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return "", errors.Wrap(err, "GetBlockHash err")
	}

	rv, err := api.r.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return "", errors.Wrap(err, "GetRuntimeVersionLatest err")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey)
	if err != nil {
		return "", errors.Wrap(err, "CreateStorageKey err")
	}

	ok, err := api.r.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return "", errors.Wrap(err, "GetStorageLatest err")
	}
	if !ok {
		return "", errors.New("GetStorageLatest return value is empty")
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(keyring, o)
	if err != nil {
		return "", errors.Wrap(err, "Sign err")
	}

	// Do the transfer and track the actual status
	sub, err := api.r.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return "", errors.Wrap(err, "SubmitAndWatchExtrinsic err")
	}
	defer sub.Unsubscribe()

	timeout := time.After(10 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				return fmt.Sprintf("%#x", status.AsInBlock), nil
			}
		case <-timeout:
			return "", errors.Errorf("[%v] tx timeout", ci.TransactionName)
		default:
			time.Sleep(time.Second)
		}
	}
}
