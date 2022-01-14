package client

import "dapp_cess_client/internal/logger"

func FileUpload() {
	logger.Logch <- "there is File Upload command!"
	logger.OutPutLogger.Sugar().Infof("there is File Upload command!")
}
