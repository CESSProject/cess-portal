package test

import (
	"cess-portal/internal/rpc"
	"cess-portal/module"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net/http/httptest"
	"testing"
)

type clientService struct {
}

func (clientService) CtluploadAction(body []byte) (proto.Message, error) {
	var blockinfo module.FileUploadInfo
	var Err rpc.RespBody
	err := proto.Unmarshal(body, &blockinfo)
	if err != nil {
		Err.Code = -1
		Err.Msg = err.Error()
		return &Err, nil
	}
	fmt.Println(blockinfo.BlockNum)
	Err.Code = 0
	Err.Msg = ""
	return &Err, nil
}

func TestUploadService(t *testing.T) {
	srv := rpc.NewServer()
	srv.Register(module.CtlServiceName, clientService{})
	s := httptest.NewServer(srv.WebsocketHandler([]string{"*"}))
	defer s.Close()
	defer srv.Close()
	fmt.Println(s.URL)
}
