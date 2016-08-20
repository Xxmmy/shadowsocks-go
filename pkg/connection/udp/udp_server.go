package udp

import (
	"reflect"
	"time"

	"shadowsocks-go/pkg/config"
	conn "shadowsocks-go/pkg/connection"
	encrypt "shadowsocks-go/pkg/connection"

	"github.com/golang/glog"
)

//UDPServer maintain a listener
type UDPServer struct {
	Config   *config.ConnectionInfo
	udpProxy *conn.Proxy
}

//NewUDPServer create a TCPServer
func NewUDPServer(cfg *config.ConnectionInfo) *UDPServer {
	return &UDPServer{
		Config: cfg,
	}
}

//Stop implement quit go routine
func (udpSrv *UDPServer) Stop() {
	glog.V(5).Infof("udp server close %v\r\n", udpSrv.Config)
	udpSrv.udpProxy.Stop()
	//udpSrv.wg.Wait()
}

//Run implement a new udp listener
func (udpSrv *UDPServer) Run() {

	password := udpSrv.Config.Password
	method := udpSrv.Config.EncryptMethod
	port := udpSrv.Config.Port
	auth := udpSrv.Config.EnableOTA
	timeout := time.Duration(udpSrv.Config.Timeout) * time.Second

	cipher, err := encrypt.NewCipher(method, password)
	if err != nil {
		glog.Fatalf("Error generating cipher for udp port: %d %v\n", port, err)
		return
	}

	proxy := conn.NewProxy(port, cipher.Copy(), auth, timeout)
	if proxy == nil {
		glog.Fatalf("listening upd port: %v error:%v\r\n", port, err)
		return
	}
	udpSrv.udpProxy = proxy

	go proxy.RunProxy()
}

func (udpSrv *UDPServer) Compare(client *config.ConnectionInfo) bool {
	return reflect.DeepEqual(udpSrv.Config, client)
}
