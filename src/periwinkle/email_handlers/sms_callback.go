// Copyright 2015 Davis Webb
// Copyright 2015 Zhandos Suleimenov
// Copyright 2015 Luke Shumaker

package handlers

import (
	"net/http"
	"sync"
	"net"
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"bufio"
	"net/url"
)

type SmsStatus struct{
	MessageStatus string
	ErrorCode string
	MessageSid string
}

type SmsCallbackServer struct {
	connsLock sync.Mutex
	conns map[string]net.Conn
}

// server
func (server *SmsCallbackServer) Serve() (err error) {
	if server.conns == nil {
		server.conns = make(map[string]net.Conn)
	}
	listener, err := net.Listen("tcp", ":42586")
	if err != nil {
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Accept: %v\n", err)
			continue
		}
		go func() {
			reader := bufio.NewReader(conn)
			message_sid, _, err := reader.ReadLine()
			if err != nil {
				defer conn.Close()
				fmt.Fprintf(os.Stderr, "read: %v\n", err)
			}
			server.connsLock.Lock()
			server.conns[string(message_sid)] = conn
			server.connsLock.Unlock()
		}()
	}
}

// server
func (server *SmsCallbackServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	status := SmsStatus{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("%v", err)
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Printf("%v", err)
	}

	status.MessageStatus = values.Get("MessageStatus")
	status.ErrorCode = values.Get("ErrorCode")
	status.MessageSid = values.Get("MessageSid")
	status_json, err := json.Marshal(status)

	server.connsLock.Lock()
	conn, ok := server.conns[status.MessageSid]
	if !ok {
		return
	}
	defer conn.Close()
	_, err = conn.Write(status_json)
	// TODO: check err
	delete(server.conns, status.MessageSid)
	server.connsLock.Unlock()

	// TODO: respond to the HTTP request (empty body or whatever)
}

// client
func SmsWaitForCallback(MessageSid string) (status SmsStatus, err error) {
	conn, err := net.Dial("tcp", "localhost:42586")
	defer conn.Close()
	if err != nil {
		return
	}
	_, err = fmt.Fprintln(conn, MessageSid)
	if err != nil {
		return
	}
	reader := bufio.NewReader(conn)
	status_json, _, err := reader.ReadLine()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(status_json), &status)
	return
}

