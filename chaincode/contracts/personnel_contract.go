package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/chrisarmitage/hlf-chaincode-starfleet-personnel/api"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type PersonnelContract struct {
	contractapi.Contract
}

func (c *PersonnelContract) Name() string {
	return "PersonnelContract"
}

const (
	DocTypePersonnel = "personnel"
)

func personnelKey(personnelID string) string {
	return fmt.Sprintf("%s:%s", DocTypePersonnel, personnelID)
}

func (c *PersonnelContract) GetPersonnel(ctx contractapi.TransactionContextInterface, personnelID string) (*api.Personnel, error) {
	key := personnelKey(personnelID)

	personnelBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read personnel from world state: %w", err)
	}
	if personnelBytes == nil {
		return nil, fmt.Errorf("personnel with ID %s does not exist", personnelID)
	}

	var personnel *api.Personnel
	if err := json.Unmarshal(personnelBytes, &personnel); err != nil {
		return nil, fmt.Errorf("failed to unmarshal personnel data: %w", err)
	}

	return personnel, nil
}
