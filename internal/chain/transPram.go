package chain

var (
	BuySpaceChainModule     = "Sminer"
	BuySpaceTransactionName = "Sminer.buy_space"
	BuySpaceModuleMethod    = ""
)

type CessInfo struct {
	RpcAddr               string
	IdentifyAccountPhrase string
	TransactionName       string
	ChainModule           string
	ChainModuleMethod     string
}
