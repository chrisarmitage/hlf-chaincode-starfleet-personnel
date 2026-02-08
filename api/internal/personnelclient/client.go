package personnelclient

import (
	"encoding/json"
	"fmt"

	"github.com/chrisarmitage/hlf-chaincode-starfleet-personnel/domain"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type PersonnelClient struct {
	contract *client.Contract
}

var ErrInvalidPersonnelID = fmt.Errorf("invalid personnel ID")

func NewPersonnelClient(contract *client.Contract) *PersonnelClient {
	return &PersonnelClient{
		contract: contract,
	}
}

func (c *PersonnelClient) GetPersonnel(personnelID string) (*domain.Personnel, error) {
	if personnelID == "" {
		return nil, ErrInvalidPersonnelID
	}

	result, err := c.contract.EvaluateTransaction("PersonnelContract:GetPersonnel", personnelID)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	var personnel *domain.Personnel
	if err := json.Unmarshal(result, &personnel); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return personnel, nil
}

func (c *PersonnelClient) EnrollCadet(personnelID, name, campus string) (*domain.Personnel, error) {
	if personnelID == "" {
		return nil, ErrInvalidPersonnelID
	}

	result, err := c.contract.SubmitTransaction("PersonnelContract:EnrollCadet", personnelID, name, campus)
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction: %w", err)
	}

	var personnel *domain.Personnel
	if err := json.Unmarshal(result, &personnel); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return personnel, nil
}
