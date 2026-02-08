package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/chrisarmitage/hlf-chaincode-starfleet-personnel/domain"
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

func (c *PersonnelContract) GetPersonnel(ctx contractapi.TransactionContextInterface, personnelID string) (*domain.Personnel, error) {
	key := personnelKey(personnelID)

	personnelBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read personnel from world state: %w", err)
	}
	if personnelBytes == nil {
		return nil, fmt.Errorf("personnel with ID %s does not exist", personnelID)
	}

	var personnel *domain.Personnel
	if err := json.Unmarshal(personnelBytes, &personnel); err != nil {
		return nil, fmt.Errorf("failed to unmarshal personnel data: %w", err)
	}

	return personnel, nil
}

func (c *PersonnelContract) EnrollCadet(ctx contractapi.TransactionContextInterface, personnelID, name, campus string) (*domain.Personnel, error) {
	if personnelID == "" {
		return nil, fmt.Errorf("personnelID is required")
	}
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if campus == "" {
		return nil, fmt.Errorf("campus is required")
	}

	existingPersonnel, err := c.GetPersonnel(ctx, personnelID)
	if err == nil && existingPersonnel != nil {
		return nil, fmt.Errorf("personnel with ID %s already exists", personnelID)
	}

	personnel := &domain.Personnel{
		PersonnelID: personnelID,
		Name:        name,
		Rank:        domain.RankCadet,
		Campus:      campus,
		Status:      domain.StatusActive,
	}

	personnelBytes, err := json.Marshal(personnel)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal personnel: %v", err)
	}

	if err := ctx.GetStub().PutState(personnelKey(personnelID), personnelBytes); err != nil {
		return nil, fmt.Errorf("failed to put personnel state: %v", err)
	}

	return personnel, nil
}