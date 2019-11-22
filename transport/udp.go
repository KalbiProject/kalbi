package transport

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"Kalbi/sip/message"
    "Kalbi/log"
)

//ListenAndServe function is an endless loop for listening on the specified host and port
func ListenAndServe(Host string, Port int) {
    log.Log.Info("Starting Kalbi Server")
	log.Log.Info("Listening on "+ Host)



	buffer := make([]byte, 2048)

	addr := net.UDPAddr{
		IP:   net.ParseIP(Host),
		Port: Port,
	}

	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		log.Log.Error("Unable to bind to socket")
	}

	for {
		n, _, err := conn.ReadFromUDP(buffer)

	    fmt.Println(n)
		if err != nil {
			log.Log.Error("Unable to read from socket")
		}
        
		bytesbuffer := bytes.NewBuffer(buffer[:n])

		newreader := bufio.NewReader(bytesbuffer)

		request := message.ReadSIPRequest(newreader)
        
		log.Log.Info(request)

		if err != nil {
			panic(err)
		}

	}

}
