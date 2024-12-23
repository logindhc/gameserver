package cherryConnector

import (
	"errors"
	cfacade "gameserver/cherry/facade"
	clog "gameserver/cherry/logger"
	"net"
	"time"
)

type (
	TCPConnector struct {
		cfacade.Component
		Connector
		Options
	}
)

func (*TCPConnector) Name() string {
	return "tcp_connector"
}

func (t *TCPConnector) OnAfterInit() {
}

func (t *TCPConnector) OnStop() {
	t.Stop()
}

func NewTCP(address string, opts ...Option) *TCPConnector {
	if address == "" {
		clog.Warn("Create tcp connector fail. Address is null.")
		return nil
	}

	tcp := &TCPConnector{
		Options: Options{
			address:  address,
			certFile: "",
			keyFile:  "",
			chanSize: 256,
		},
	}

	for _, opt := range opts {
		opt(&tcp.Options)
	}

	tcp.Connector = NewConnector(tcp.chanSize)

	return tcp
}

func (t *TCPConnector) Start() {
	listener, err := t.GetListener(t.certFile, t.keyFile, t.address)
	if err != nil {
		clog.Fatalf("failed to listen: %s", err)
	}

	clog.Infof("Tcp connector listening at Address %s", t.address)
	if t.certFile != "" || t.keyFile != "" {
		clog.Infof("certFile = %s, keyFile = %s", t.certFile, t.keyFile)
	}

	t.Connector.Start()
	var tempDelay time.Duration
	for t.Running() {
		conn, _err := listener.Accept()
		if _err != nil {
			var ne net.Error
			if errors.As(_err, &ne) && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if _max := 1 * time.Second; tempDelay > _max {
					tempDelay = _max
				}
				clog.Infof("accept error: %v; retrying in %v", _err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0

		t.InChan(conn)
	}
}

func (t *TCPConnector) Stop() {
	t.Connector.Stop()
}
