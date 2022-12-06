package conf

import (
	"time"
)

// default set up about cess client
var (
	// base dir
	BaseDir = "./data"
	// keyfile dir
	PrivateKeyfile = BaseDir + "/.privateKey.pem"
	PublicKeyfile  = BaseDir + "/.publicKey.pem"
	// file cache dir
	FileCacheDir = BaseDir + "/cache"
	// log dir
	LogfileDir = BaseDir + "/logs"

	// random number valid time, the unit is minutes
	RandomValidTime = 5.0

	// the time to wait for the event, in seconds
	TimeToWaitEvents = time.Duration(time.Second * 15)

	// The validity period of the token, the default is 30 days
	ValidTimeOfToken = time.Duration(time.Hour * 24 * 30)

	// Valid Time Of Captcha
	ValidTimeOfCaptcha = time.Duration(time.Minute * 5)

	//
	SIZE_1KB int64 = 1024
	SIZE_1MB int64 = 1024 * SIZE_1KB
	SIZE_1GB int64 = 1024 * SIZE_1MB
)

const (
	// Tcp message interval
	TCP_Message_Interval = time.Duration(time.Millisecond * 10)
	// Number of tcp message caches
	TCP_Message_Send_Buffers = 10
	TCP_Message_Read_Buffers = 10
	//
	TCP_SendBuffer = 8192
	TCP_ReadBuffer = 12000
	//
	Tcp_Dial_Timeout = time.Duration(time.Second * 5)
)

/*
system set up
*/
const (
	Exit_Normal         = 0
	Exit_CmdLineParaErr = -1
	Exit_ConfErr        = -2
	Exit_ChainErr       = -3
	Exit_SystemErr      = -4
)

const MaxBackups = 6

var PublicKey []byte
