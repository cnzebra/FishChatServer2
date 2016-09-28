package server

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/protocol"
	"github.com/oikomi/FishChatServer2/server/gateway/client"
)

type Server struct {
	Server *libnet.Server
}

func New() (s *Server) {
	s = &Server{}
	return
}

func (s *Server) sessionLoop(client *client.Client) {
	for {
		reqData, err := client.Session.Receive()
		if err != nil {
			glog.Error(err)
		}
		if reqData != nil {
			baseCMD := &protocol.Base{}
			if err = proto.Unmarshal(reqData, baseCMD); err != nil {
				if err = client.Session.Send(&protocol.Error{
					ErrCode: ecode.ServerErr.Uint32(),
					ErrStr:  ecode.ServerErr.String(),
				}); err != nil {
					glog.Error(err)
				}
				continue
			}
			if err = client.Parse(baseCMD.Cmd, reqData); err != nil {
				glog.Error(err)
				continue
			}
		}
	}
}

func (s *Server) Loop() {
	for {
		session, err := s.Server.Accept()
		if err != nil {
			glog.Error(err)
		}
		go s.sessionLoop(client.New(session))
	}
}
