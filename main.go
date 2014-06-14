package main

import (
	"encoding/json"
	"flag"
	"net"
	"os"
	"os/exec"
	"os/signal"

	"github.com/kr/pty"
	"github.com/nitrous-io/nitrous-common/logger"
)

var (
	respWtf       = []byte("{\"error\":\"wtf\"}\n")
	respConnected = []byte("{\"connected\":true}\n")
	respAttached  = []byte("{\"attached\":true}\n")
)

type Handshake struct {
	ContainerID string `json:"attach"`
	User        string `json:"user"`
}

type ErrorResp struct {
	Error string `json:"error"`
}

type Attacher struct {
	logger *logger.Logger
}

func (a *Attacher) handle(conn net.Conn) {
	defer conn.Close()

	_, err := conn.Write(respConnected)
	if err != nil {
		return
	}

	var h Handshake
	in := json.NewDecoder(conn)
	out := json.NewEncoder(conn)

	err = in.Decode(&h)
	if err != nil {
		conn.Write(respWtf)
		return
	}

	if h.ContainerID == "" || h.User == "" {
		conn.Write(respWtf)
		return
	}

	cmd := exec.Command("lxc-attach", "-n", h.ContainerID, "su", "-", h.User)
	pty, err := pty.Start(cmd)
	if err != nil {
		a.logger.Error(err)
		out.Encode(&ErrorResp{Error: err.Error()})
		return
	}
	defer pty.Close()

	_, err = conn.Write(respAttached)
	if err != nil {
		a.logger.Error(err)
		return
	}

	c := make(chan bool, 1)
	stop := false

	// conn
	go func() {
		cmd := make([]string, 2)
		for !stop {
			err := in.Decode(&cmd)
			if err != nil {
				out.Encode(&ErrorResp{Error: err.Error()})
				c <- true
				return
			}

			if cmd[0] == "disconnect" {
				c <- true
				return
			}

			if cmd[0] == "d" && cmd[1] != "" {
				_, err := pty.Write([]byte(cmd[1]))
				if err != nil {
					out.Encode(&ErrorResp{Error: err.Error()})
					c <- true
					return
				}
			}
		}
	}()

	// pty
	go func() {
		for !stop {
			buf := make([]byte, 1024)
			n, err := pty.Read(buf)
			if err != nil {
				out.Encode(&ErrorResp{Error: err.Error()})
				c <- true
				return
			}
			err = out.Encode([]string{"d", string(buf[:n])})
			if err != nil {
				out.Encode(&ErrorResp{Error: err.Error()})
				c <- true
				return
			}
		}
	}()

	stop = <-c
}

func main() {
	flSockPath := flag.String("s", "/var/run/attacher.sock", "Specify the path for attcher's unix socket")
	flag.Parse()

	logger := &logger.Logger{Level: logger.LogLevelDebug}

	logger.Infof("attacher (pid: %v) listening on %v", os.Getpid(), *flSockPath)
	listener, err := net.Listen("unix", *flSockPath)
	if err != nil {
		logger.FatalExit(err)
	}

	defer func() {
		logger.Info("Cleaning up...")
		listener.Close()
		os.Remove(*flSockPath)
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	attacher := &Attacher{logger: logger}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Error("Accept failed!", err)
				return
			}

			go attacher.handle(conn)
		}
	}()

	<-sig
}
