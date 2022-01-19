package client

import (
	"dapp_cess_client/conf"
	"dapp_cess_client/internal/chain"
	"dapp_cess_client/tools"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

func FileUpload(filepath, downloadfee string) {
	chain.Chain_Init()
	file, err := os.Stat(filepath)
	if err != nil {
		fmt.Printf("[Error]Please enter the correct file path!\n")
		return
	}
	if file.IsDir() {
		fmt.Printf("[Error]Please do not upload the folder!\n")
		return
	}
	cosfee, err := strconv.Atoi(downloadfee)
	if err != nil {
		fmt.Printf("[Error]Please enter a correct integer!\n")
		return
	}
	filehash, err := tools.CalcFileHash(filepath)
	if err != nil {
		fmt.Printf("[Error]There is a problem with the file, please replace it!\n")
		return
	}
	if file.Size()/100 == 0 {
		fmt.Printf("[Error]The upload file is too small!\n")
		return
	}
	fileid, err := tools.GetGuid(1)
	if err != nil {
		fmt.Printf("[Error]Create snowflake fail! error:%s\n", err)
		return
	}

	var ci chain.CessInfo
	filesize := new(big.Int)
	fee := new(big.Int)

	ci.IdentifyAccountPhrase = conf.ClientConf.ChainData.IdAccountPhraseOrSeed
	ci.TransactionName = chain.UploadFileTransactionName

	filesize.SetInt64(file.Size() / 100)
	fee.SetInt64(int64(cosfee))
	fmt.Println(file.Size() / 100)

	AsInBlock, err := ci.UploadFileMetaInformation(fileid, file.Name(), filehash, filesize, fee)
	if err != nil {
		fmt.Printf("[Error]Upload file meta information error:%s", err)
		return
	}
	fmt.Printf("Transaction chain block number is:%s\n", AsInBlock)
}
