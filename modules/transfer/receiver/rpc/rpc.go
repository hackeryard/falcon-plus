// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"
	"sync"

	"github.com/open-falcon/falcon-plus/modules/transfer/g"
)

// define global in package rpc
var clientIP string

// StartRpc Server: accept connections and do rpc codec
func StartRpc() {
	if !g.Config().Rpc.Enabled {
		return
	}

	addr := g.Config().Rpc.Listen
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("net.ResolveTCPAddr fail: %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		log.Println("rpc listening", addr)
	}

	server := rpc.NewServer()
	// RPC client using "Transfer.Update" to using the function in server
	// Note: register at begining, just once
	transferIns := new(Transfer)
	server.Register(transferIns)

	var mutex sync.Mutex
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept occur error:", err)
			continue
		}

		mutex.Lock()
		// 从conn这里获取jsonrpc客户端的ip地址
		clientIP = strings.Split(conn.RemoteAddr().String(), ":")[0]
		// 只能把IP作为全局变量共享 最后发现这样也不行 现在的做法是使用互斥锁
		// 一个conn对应一个地址 修改rpc的代码 使得IP能够作为参数传递
		log.Println("[info] new remote client addr: ", clientIP)

		transferIns.clientAddr = clientIP
		// @@remove go func
		// conn的信息一直都在
		// rpc的设计：直接加一个json的encodec 加了一层 就构成了jsonrpc
		// ServeCodec is like ServeConn but uses the specified codec to decode requests and encode responses.
		// 每一次调用ServeCodec 就带一个IP var
		go server.ServeCodec(jsonrpc.NewServerCodec(conn), clientIP)

		mutex.Unlock()

	}
}
