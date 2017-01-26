package main

import (
  "errors"
  "fmt"
  "encoding/json"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  "github.com/hyperledger/fabric/core/crypto/primitives"
  "github.com/op/go-logging"
)

const (
  tableColumn       = "ChatLog"
)

var myLogger = logging.MustGetLogger("supplier-dealer")

type SupplierDealerChaincode struct {
}

type ChatLog struct {
  sender string
  message  string
}

func (t *SupplierDealerChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  myLogger.Debug("Init Chaincode...........")

  // if len(args) != 0 {
  //   return nil, errors.New("Incorrect number of arguments. Expecting 0")
  // }

  // Create ChatLog table
  chatLogErr := stub.CreateTable("ChatLog", []*shim.ColumnDefinition{
    &shim.ColumnDefinition{Name: "Sender", Type: shim.ColumnDefinition_STRING, Key: false},
    &shim.ColumnDefinition{Name: "Message", Type: shim.ColumnDefinition_STRING, Key: false},
  })
  if chatLogErr != nil {
    return nil, errors.New("Failed creating ChatLog table.")
  }

  // Create Order table
  orderErr := stub.CreateTable("Order", []*shim.ColumnDefinition{
    &shim.ColumnDefinition{Name: "ProductName", Type: shim.ColumnDefinition_STRING, Key: false},
    &shim.ColumnDefinition{Name: "DeliveryAddress", Type: shim.ColumnDefinition_STRING, Key: false},
  })
  if orderErr != nil {
    return nil, errors.New("Failed creating Order table.")
  }

  myLogger.Debug("Init Chaincode...done")

  return nil, nil
}

func (t *SupplierDealerChaincode) sendMessage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  sender := args[0]
  message := args[1]
  ok, err := stub.InsertRow(tableColumn, shim.Row{
    Columns: []*shim.Column{
      &shim.Column{Value: &shim.Column_String_{String_: sender}},
      &shim.Column{Value: &shim.Column_String_{String_: message}},
    },
  })

  if !ok && err == nil {
    myLogger.Errorf("system error %v", err)
    return nil, errors.New("Cannot send the message")
  }

  return nil, nil
}

func (t *SupplierDealerChaincode) readMessages(stub shim.ChaincodeStubInterface, args []string) ( []byte, error) {
  var columns []shim.Column
  var chatLogs []ChatLog

  rowsChan, err := stub.GetRows(tableColumn, columns)
  if err != nil {
    return nil, err
  }

  for row := range rowsChan {
    chatLog := ChatLog{sender: row.Columns[0].GetString_(), message: row.Columns[1].GetString_()}
    chatLogs = append(chatLogs, chatLog)
  }

  return json.Marshal(chatLogs)
}

func (t *SupplierDealerChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

  // Handle different functions
  if function == "sendMessage" {
    return t.sendMessage(stub, args)
  }
  return nil, errors.New("Received unknown function invocation")
}

func (t *SupplierDealerChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

  if function == "readMessages" {
    return t.readMessages(stub, args)
  }
  return nil, errors.New("Received unknown function query invocation with function " + function)
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(SupplierDealerChaincode))
	if err != nil {
		fmt.Printf("Error starting SupplierDealerChaincode: %s", err)
	}
}
