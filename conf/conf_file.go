package conf

type Configfile struct {
	RpcAddr     string `toml:"RpcAddr"`
	AccountSeed string `toml:"AccountSeed"`
	AccountId   string `toml:"AccountId"`
}

var C = new(Configfile)
var ConfigFilePath string

const ConfigFile_Templete = `
#The rpc address of the chain node
RpcAddr           = ""
#Phrase or seed for wallet account
AccountSeed       = ""
#wallet account of cess 
AccountId = ""
`
