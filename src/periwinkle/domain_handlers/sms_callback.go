// Copyright 2015 Davis Webb
// Copyright 2015 Zhandos Suleimenov
// Copyright 2015 Luke Shumaker

package domain_handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"periwinkle"
	"sync"
)

type SmsStatus struct {
	MessageStatus string
	ErrorCode     string
	MessageSid    string
}

type SmsCallbackServer struct {
	Conf      *periwinkle.Cfg
	connsLock sync.Mutex
	conns     map[string]net.Conn
}

// server
func (server *SmsCallbackServer) Serve() (err error) {
	if server.conns == nil {
		server.conns = make(map[string]net.Conn)
	}
	listener, err := net.Listen("tcp", server.Conf.CallbackListen)
	if err != nil {
		return
	}
	log.Println("41")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept:", err)
			continue
		}
		go func() {
			reader := bufio.NewReader(conn)
			messageSID, _, err := reader.ReadLine()
			if err != nil {
				log.Println("Read:", err)
			}
			server.connsLock.Lock()
			server.conns[string(messageSID)] = conn
			server.connsLock.Unlock()
		}()
	}
}

// server
func (server *SmsCallbackServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	status := SmsStatus{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		log.Println(err)
	}

	status.MessageStatus = values.Get("MessageStatus")
	status.ErrorCode = values.Get("ErrorCode")
	status.MessageSid = values.Get("MessageSid")
	statusJSON, err := json.Marshal(status)
	log.Println("77", status.MessageStatus)
	server.connsLock.Lock()
	conn, ok := server.conns[status.MessageSid]
	if !ok {
		return
	}
	defer conn.Close()
	_, err = conn.Write(statusJSON)
	if err != nil {
		periwinkle.Logf("Couldnt write to the socket")
	}
	delete(server.conns, status.MessageSid)
	server.connsLock.Unlock()
}

// client
func SmsWaitForCallback(conf *periwinkle.Cfg, MessageSid string) (status SmsStatus, err error) {
	conn, err := net.Dial("tcp", conf.CallbackConnect)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = fmt.Fprintln(conn, MessageSid)
	if err != nil {
		return
	}
	log.Println("104")
	reader := bufio.NewReader(conn)
	statusJSON, _, err := reader.ReadLine()
	if err != nil {
		return
	}
	log.Println("110")
	err = json.Unmarshal([]byte(statusJSON), &status)
	return
}
