package transport

import (
	"github.com/KalbiProject/kalbi/log"
	"github.com/KalbiProject/kalbi/sip/message"
	"strconv"
)

// NewTransportListenPoint creates listen
func NewTransportListenPoint(protocol string, host string, port int) message.ListeningPoint {
	switch protocol {
	case "udp":
		log.Log.Info("Creating UDP listening point")
		listner := new(UDPTransport)
		log.Log.Info("Binding to " + host + ":" + strconv.Itoa(port))
		listner.Build(host, port)
		return listner
	case "tcp":
		log.Log.Info("Creating TCP listening point")
		listner := new(UDPTransport)
		log.Log.Info("Binding to " + host + ":" + strconv.Itoa(port))
		listner.Build(host, port)
		return listner
	case "tls":
		log.Log.Info("Creating TLS listening point")
		listner := new(UDPTransport)
		log.Log.Info("Binding to " + host + ":" + strconv.Itoa(port))
		listner.Build(host, port)
		return listner
	default:
		log.Log.Info("Unknown protocol specified")
		panic("Unknown protocol specified")

	}

}
