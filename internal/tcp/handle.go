package tcp

import (
	"cess-portal/conf"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Client interface {
	SendFile(fid string, fsize int64, pkey, signmsg, sign []byte) error
	RecvFile(fid string, fsize int64, pkey, signmsg, sign []byte) error
}

type NetConn interface {
	HandlerLoop()
	GetMsg() (*Message, bool)
	SendMsg(m *Message)
	Close() error
	IsClose() bool
}

type ConMgr struct {
	conn     NetConn
	dir      string
	fileName string

	sendFiles []string

	waitNotify chan bool
	stop       chan struct{}
}

func (c *ConMgr) handler() error {
	var (
		err      error
		recvFile *os.File
	)

	defer func() {
		recover()
		c.conn.Close()
		close(c.waitNotify)
		if recvFile != nil {
			_ = recvFile.Close()
		}
	}()

	for !c.conn.IsClose() {
		m, ok := c.conn.GetMsg()
		if !ok {
			return fmt.Errorf("Getmsg failed")
		}
		if m == nil {
			continue
		}

		switch m.MsgType {
		case MsgHead:
			switch cap(m.Bytes) {
			case conf.TCP_ReadBuffer:
				readBufPool.Put(m.Bytes)
			default:
			}
			c.conn.SendMsg(NewNotifyMsg(c.fileName, Status_Ok))
		case MsgFile:
			if recvFile == nil {
				recvFile, err = os.OpenFile(filepath.Join(c.dir, m.FileName), os.O_RDWR|os.O_TRUNC, os.ModePerm)
				if err != nil {
					c.conn.SendMsg(NewNotifyMsg("", Status_Err))
					time.Sleep(conf.TCP_Message_Interval)
					c.conn.SendMsg(NewCloseMsg("", Status_Err))
					time.Sleep(conf.TCP_Message_Interval)
					return err
				}
			}
			_, err = recvFile.Write(m.Bytes[:m.FileSize])
			if err != nil {
				c.conn.SendMsg(NewNotifyMsg("", Status_Err))
				time.Sleep(conf.TCP_Message_Interval)
				c.conn.SendMsg(NewCloseMsg("", Status_Err))
				time.Sleep(conf.TCP_Message_Interval)
				return err
			}
			switch cap(m.Bytes) {
			case conf.TCP_ReadBuffer:
				readBufPool.Put(m.Bytes)
			default:
			}
		case MsgEnd:
			switch cap(m.Bytes) {
			case conf.TCP_ReadBuffer:
				readBufPool.Put(m.Bytes)
			default:
			}
			info, err := recvFile.Stat()
			if err != nil {
				c.conn.SendMsg(NewNotifyMsg("", Status_Err))
				time.Sleep(conf.TCP_Message_Interval)
				c.conn.SendMsg(NewCloseMsg("", Status_Err))
				time.Sleep(conf.TCP_Message_Interval)
				return err
			}
			if info.Size() != int64(m.FileSize) {
				err = fmt.Errorf("file.size %v rece size %v \n", info.Size(), m.FileSize)
				c.conn.SendMsg(NewNotifyMsg("", Status_Err))
				time.Sleep(conf.TCP_Message_Interval)
				c.conn.SendMsg(NewCloseMsg("", Status_Err))
				time.Sleep(conf.TCP_Message_Interval)
				return err
			}
			recvFile.Close()
			recvFile = nil

		case MsgNotify:
			c.waitNotify <- m.Bytes[0] == byte(Status_Ok)
			switch cap(m.Bytes) {
			case conf.TCP_ReadBuffer:
				readBufPool.Put(m.Bytes)
			default:
			}

		case MsgClose:
			switch cap(m.Bytes) {
			case conf.TCP_ReadBuffer:
				readBufPool.Put(m.Bytes)
			default:
			}
			return errors.New("Close message")

		default:
			switch cap(m.Bytes) {
			case conf.TCP_ReadBuffer:
				readBufPool.Put(m.Bytes)
			default:
			}
			return errors.New("Invalid msgType")
		}
	}

	return err
}

func NewClient(conn NetConn, dir string, files []string) Client {
	return &ConMgr{
		conn:       conn,
		dir:        dir,
		sendFiles:  files,
		waitNotify: make(chan bool, 1),
		stop:       make(chan struct{}),
	}
}

func (c *ConMgr) SendFile(fid string, fsize int64, pkey, signmsg, sign []byte) error {
	c.conn.HandlerLoop()
	go func() {
		_ = c.handler()
	}()

	err := c.sendFile(fid, fsize, pkey, signmsg, sign)
	return err
}

func (c *ConMgr) RecvFile(fid string, fsize int64, pkey, signmsg, sign []byte) error {
	c.conn.HandlerLoop()
	go func() {
		_ = c.handler()
	}()
	err := c.recvFile(fid, fsize, pkey, signmsg, sign)
	return err
}

func (c *ConMgr) sendFile(fid string, fsize int64, pkey, signmsg, sign []byte) error {
	defer func() {
		c.conn.Close()
	}()

	var err error
	var lastmatrk bool

	for i := 0; i < len(c.sendFiles); i++ {
		if (i + 1) == len(c.sendFiles) {
			lastmatrk = true
		}
		err = c.sendSingleFile(filepath.Join(c.dir, c.sendFiles[i]), fid, fsize, lastmatrk, pkey, signmsg, sign)
		if err != nil {
			return err
		}
		if strings.Contains(c.sendFiles[i], ".") {
			os.Remove(filepath.Join(c.dir, c.sendFiles[i]))
		}
	}

	c.conn.SendMsg(NewCloseMsg(c.fileName, Status_Ok))
	time.Sleep(time.Second * 3)
	return err
}

func (c *ConMgr) recvFile(fid string, fsize int64, pkey, signmsg, sign []byte) error {
	defer func() {
		c.conn.Close()
	}()

	//log.Println("Ready to recvhead: ", fid)
	c.conn.SendMsg(NewRecvHeadMsg(fid, pkey, signmsg, sign))
	timerHead := time.NewTimer(time.Second * 5)
	defer timerHead.Stop()
	select {
	case ok := <-c.waitNotify:
		if !ok {
			return fmt.Errorf("send err")
		}
	case <-timerHead.C:
		return fmt.Errorf("wait server msg timeout")
	}

	_, err := os.Create(filepath.Join(c.dir, fid))
	if err != nil {
		c.conn.SendMsg(NewCloseMsg(fid, Status_Err))
		return err
	}
	//log.Println("Ready to recvfile: ", fid)
	c.conn.SendMsg(NewRecvFileMsg(fid))

	waitTime := fsize / 1024 / 10
	if waitTime < 5 {
		waitTime = 5
	}

	timerFile := time.NewTimer(time.Second * time.Duration(waitTime))
	defer timerFile.Stop()
	select {
	case ok := <-c.waitNotify:
		if !ok {
			return fmt.Errorf("send err")
		}
	case <-timerFile.C:
		return fmt.Errorf("wait server msg timeout")
	}
	c.conn.SendMsg(NewCloseMsg(fid, Status_Ok))
	time.NewTimer(time.Second * 3)
	return nil
}

func (c *ConMgr) sendSingleFile(filePath string, fid string, fsize int64, lastmark bool, pkey, signmsg, sign []byte) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("open file err %v \n", err)
		return err
	}

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	fileInfo, _ := file.Stat()

	//log.Println("Ready to write file: ", filePath)
	c.conn.SendMsg(NewHeadMsg(fileInfo.Name(), fid, lastmark, pkey, signmsg, sign))

	timerHead := time.NewTimer(10 * time.Second)
	defer timerHead.Stop()
	select {
	case ok := <-c.waitNotify:
		if !ok {
			return fmt.Errorf("send head msg err")
		}
	case <-timerHead.C:
		return fmt.Errorf("wait server msg timeout")
	}

	readBuf := sendBufPool.Get().([]byte)
	defer func() {
		sendBufPool.Put(readBuf)
	}()

	for !c.conn.IsClose() {
		n, err := file.Read(readBuf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		c.conn.SendMsg(NewFileMsg(c.fileName, n, readBuf[:n]))
	}

	c.conn.SendMsg(NewEndMsg(c.fileName, fid, uint64(fileInfo.Size()), uint64(fsize), lastmark))
	waitTime := fileInfo.Size() / 1024 / 10
	if waitTime < 10 {
		waitTime = 10
	}

	timerFile := time.NewTimer(time.Second * time.Duration(waitTime))
	defer timerFile.Stop()
	select {
	case ok := <-c.waitNotify:
		if !ok {
			return fmt.Errorf("send end msg err")
		}
	case <-timerFile.C:
		return fmt.Errorf("wait server msg timeout")
	}

	return nil
}
