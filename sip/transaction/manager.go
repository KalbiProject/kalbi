package transaction

import (
	"errors"
	"sync"

	"github.com/KalbiProject/kalbi/log"
	"github.com/KalbiProject/kalbi/sip/message"
	"github.com/KalbiProject/kalbi/sip/method"
	"github.com/sirupsen/logrus"
)

//NewTransactionManager returns a new TransactionManager
func NewTransactionManager() *TransactionManager {
	txmng := new(TransactionManager)
	txmng.ClientTX = make(map[string]message.Transaction)
	txmng.ServerTX = make(map[string]message.Transaction)
	txmng.txLock = &sync.RWMutex{}

	return txmng
}

//TransactionManager handles SIP transactions
type TransactionManager struct {
	ServerTX        map[string]message.Transaction
	ClientTX        map[string]message.Transaction
	RequestChannel  chan message.Transaction
	ResponseChannel chan message.Transaction
	ListeningPoint  message.ListeningPoint
	txLock          *sync.RWMutex
}

// Handle runs TransManager
func (tm *TransactionManager) Handle(event message.SipEventObject) message.SipEventObject {

	message := event.GetSipMessage()

	if message.Req.StatusCode != nil {
		log.Log.Info("Client transaction")

		tx, exists, err := tm.FindClientTransaction(message)
		if exists && err == nil {
			tx.SetLastMessage(message)
			log.Log.Info("Client Transaction already exists")
		} else {
			tx = tm.NewClientTransaction(message)
		}

		tx.Receive(message)
		event.SetTransaction(tx)
		return event

	} else if message.Req.Method != nil {
		log.Log.Info("Server transaction")
		tx, exists, err := tm.FindServerTransaction(message)

		if exists && err == nil {
			tx.SetLastMessage(message)
			log.Log.Info("Server Transaction already exists")

		} else {
			tx = tm.NewServerTransaction(message)
		}

		tx.Receive(message)

		event.SetTransaction(tx)

		return event
	}
	return event
}

//FindServerTransaction finds transaction by SipMsg
func (tm *TransactionManager) FindServerTransaction(msg *message.SipMsg) (message.Transaction, bool, error) {
	//key := tm.MakeKey(*msg)
	if len(msg.Via) == 0 {
		log.Log.Error("Via Headers Missing")
		return nil, false, errors.New("missing via headers")
	}
	tm.txLock.RLock()
	tx, exists := tm.ServerTX[string(msg.Via[0].Branch)]
	tm.txLock.RUnlock()
	return tx, exists, nil
}

//FindClientTransaction finds transaction by SipMsg
func (tm *TransactionManager) FindClientTransaction(msg *message.SipMsg) (message.Transaction, bool, error) {
	//key := tm.MakeKey(*msg)
	if len(msg.Via) == 0 {
		log.Log.Error("Via Headers Missing")
		return nil, false, errors.New(("missung via headers"))
	}
	tm.txLock.RLock()
	tx, exists := tm.ClientTX[string(msg.Via[0].Branch)]
	tm.txLock.RUnlock()
	return tx, exists, nil
}

//FindServerTransactionByID finds transaction by id
func (tm *TransactionManager) FindServerTransactionByID(value string) (message.Transaction, bool) {
	//key := tm.MakeKey(*msg)
	tm.txLock.RLock()
	tx, exists := tm.ServerTX[value]
	tm.txLock.RUnlock()
	return tx, exists
}

//FindClientTransactionByID finds transaction by id
func (tm *TransactionManager) FindClientTransactionByID(value string) (message.Transaction, bool) {
	//key := tm.MakeKey(*msg)
	tm.txLock.RLock()
	tx, exists := tm.ClientTX[value]
	tm.txLock.RUnlock()
	return tx, exists
}

//PutServerTransaction stores a transaction
func (tm *TransactionManager) PutServerTransaction(tx message.Transaction) {
	tm.txLock.Lock()
	tm.ServerTX[tx.GetBranchID()] = tx
	tm.txLock.Unlock()
}

//PutClientTransaction stores a transaction
func (tm *TransactionManager) PutClientTransaction(tx message.Transaction) {
	tm.txLock.Lock()
	tm.ClientTX[tx.GetBranchID()] = tx
	tm.txLock.Unlock()
}

//DeleteServerTransaction removes a stored transaction
func (tm *TransactionManager) DeleteServerTransaction(branch string) {
	log.Log.Info("Deleting transaction with ID: " + branch)
	log.Log.WithFields(logrus.Fields{"transactions": len(tm.ServerTX)}).Debug("Current transaction count before DeleteTransaction() is called")
	tm.txLock.Lock()
	delete(tm.ServerTX, branch)
	tm.txLock.Unlock()
	log.Log.WithFields(logrus.Fields{"transactions": len(tm.ServerTX)}).Debug("Current transaction count after DeleteTransaction() is called")
}

//DeleteClientTransaction removes a stored transaction
func (tm *TransactionManager) DeleteClientTransaction(branch string) {
	log.Log.Info("Deleting transaction with ID: " + branch)
	log.Log.WithFields(logrus.Fields{"transactions": len(tm.ClientTX)}).Debug("Current transaction count before DeleteTransaction() is called")
	tm.txLock.Lock()
	delete(tm.ClientTX, branch)
	tm.txLock.Unlock()
	log.Log.WithFields(logrus.Fields{"transactions": len(tm.ClientTX)}).Debug("Current transaction count after DeleteTransaction() is called")
}

//MakeKey creates new transaction identifier
func (tm *TransactionManager) MakeKey(msg message.SipMsg) string {

	key := string(msg.Via[0].Branch)
	var _method string
	if msg.Req.Method != nil {
		if string(msg.Req.Method) == method.ACK {
			_method = method.INVITE
		} else {
			_method = string(msg.Req.Method)
		}
	} else {
		_method = string(msg.Cseq.Method)
	}

	key += _method
	return key
}

//NewClientTransaction builds new CLientTransaction
func (tm *TransactionManager) NewClientTransaction(msg *message.SipMsg) *ClientTransaction {

	tx := new(ClientTransaction)
	tx.SetListeningPoint(tm.ListeningPoint)
	tx.TransManager = tm

	tx.InitFSM(msg)

	tx.BranchID = string(msg.Via[0].Branch)
	tx.Origin = msg

	tm.PutClientTransaction(tx)
	return tx

}

//NewServerTransaction builds new ServerTransaction
func (tm *TransactionManager) NewServerTransaction(msg *message.SipMsg) *ServerTransaction {

	tx := new(ServerTransaction)
	tx.SetListeningPoint(tm.ListeningPoint)

	tx.TransManager = tm

	tx.InitFSM(msg)

	tx.BranchID = string(msg.Via[0].Branch)
	tx.Origin = msg

	tm.PutServerTransaction(tx)
	return tx

}
