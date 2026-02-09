package personnelclient

import (
	"encoding/json"
	"fmt"
	"time"

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

func (c *PersonnelClient) CompleteTraining(recordID, personnelID, campus, trainingCode, completedAt, issuedBy string) (*domain.Training, error) {
	// Parameter validation
	if recordID == "" {
		return nil, fmt.Errorf("recordID is required")
	}
	if personnelID == "" {
		return nil, ErrInvalidPersonnelID
	}
	if campus == "" {
		return nil, fmt.Errorf("campus is required")
	}
	if trainingCode == "" {
		return nil, fmt.Errorf("trainingCode is required")
	}
	if completedAt == "" {
		return nil, fmt.Errorf("completedAt is required")
	}
	if issuedBy == "" {
		return nil, fmt.Errorf("issuedBy is required")
	}

	// Check for valid completedAt format (ISO 8601)
	if _, err := time.Parse(time.RFC3339, completedAt); err != nil {
		return nil, fmt.Errorf("completedAt must be in ISO 8601 / RFC3339 format: %w", err)
	}

	result, err := c.contract.SubmitTransaction(
		"PersonnelContract:CompleteTraining",
		recordID,
		personnelID,
		campus,
		trainingCode,
		completedAt,
		issuedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction: %w", err)
	}

	var training *domain.Training
	if err := json.Unmarshal(result, &training); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return training, nil
}
