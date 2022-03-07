package chain

import (
	"c-portal/conf"
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"os"
	"sync"
	"time"
)

type mySubstrateApi struct {
	wlock sync.Mutex
	r     *gsrpc.SubstrateAPI
}

var api mySubstrateApi

func Chain_Init() {
	var err error

	api.r, err = gsrpc.NewSubstrateAPI(conf.ClientConf.ChainData.CessRpcAddr)
	if err != nil {
		fmt.Printf("[Error]Problem with chain rpc:%s\n", err)
		os.Exit(conf.Exit_ChainErr)
	}
	go substrateAPIKeepAlive()
}

func substrateAPIKeepAlive() {
	var (
		err     error
		count_r uint8  = 0
		peer    uint64 = 0
	)

	for range time.Tick(time.Second * 25) {
		if count_r <= 1 {
			peer, err = healthchek(api.r)
			if err != nil || peer == 0 {
				count_r++
			}
		}
		if count_r > 1 {
			count_r = 2
			api.r, err = gsrpc.NewSubstrateAPI(conf.ClientConf.ChainData.CessRpcAddr)
			if err != nil {

			} else {
				count_r = 0
			}
		}
	}
}

func healthchek(a *gsrpc.SubstrateAPI) (uint64, error) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("[Error]Recover healthchek panic fail :", err)
		}
	}()
	h, err := a.RPC.System.Health()
	return uint64(h.Peers), err
}

func (myapi *mySubstrateApi) getSubstrateApiSafe() {
	myapi.wlock.Lock()
	return
}

func (myapi *mySubstrateApi) releaseSubstrateApi() {
	myapi.wlock.Unlock()
	return
}
