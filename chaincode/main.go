package main

import (
	"log"
	"os"

	"github.com/chrisarmitage/hlf-chaincode-starfleet-personnel/chaincode/contracts"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type serverConfig struct {
	CCID    string
	Address string
}

func main() {
	config := loadConfig()

	chaincode, err := contractapi.NewChaincode(
		&contracts.PersonnelContract{},
	)

	if err != nil {
		log.Panicf("error creating chaincode: %s", err)
	}

	server := &shim.ChaincodeServer{
		CCID:     config.CCID,
		Address:  config.Address,
		CC:       chaincode,
		TLSProps: getTLSProperties(),
	}

	log.Println("Starting chaincode server...")
	if err := server.Start(); err != nil {
		log.Panicf("error starting chaincode: %s", err)
	}
}

func loadConfig() *serverConfig {
	ccid := os.Getenv("CHAINCODE_CCID")
	if ccid == "" {
		log.Panic("CHAINCODE_CCID environment variable is required")
	}

	address := os.Getenv("CHAINCODE_ADDRESS")
	if address == "" {
		log.Panic("CHAINCODE_ADDRESS environment variable is required")
	}

	log.Println("=== Config ===")
	log.Printf("  CHAINCODE_CCID: %s", ccid)
	log.Printf("  CHAINCODE_ADDRESS: %s", address)

	return &serverConfig{
		CCID:    ccid,
		Address: address,
	}
}

func getTLSProperties() shim.TLSProperties {
	// not tested yet
	return shim.TLSProperties{
		Disabled: true,
	}
}
