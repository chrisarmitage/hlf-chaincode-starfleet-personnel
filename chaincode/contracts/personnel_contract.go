package contracts

import (
	"encoding/json"
	"fmt"
	"time"

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
	DocTypeTraining  = "training"
)

func personnelKey(personnelID string) string {
	return fmt.Sprintf("%s:%s", DocTypePersonnel, personnelID)
}

func trainingKey(recordID string) string {
	return fmt.Sprintf("%s:%s", DocTypeTraining, recordID)
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
		Rank:        domain.PersonnelRankCadet,
		Campus:      campus,
		Status:      domain.PersonnelStatusActive,
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

func (c *PersonnelContract) CompleteTraining(ctx contractapi.TransactionContextInterface, recordID, personnelID, campus, trainingCode, completedAt, issuedBy string) (*domain.Training, error) {
	// Parameter validation
	if recordID == "" {
		return nil, fmt.Errorf("recordID is required")
	}
	if personnelID == "" {
		return nil, fmt.Errorf("personnelID is required")
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

	// Existing record check
	existingTraining, err := ctx.GetStub().GetState(trainingKey(recordID))
	if err != nil {
		return nil, fmt.Errorf("failed to check existing training state: %w", err)
	}
	if existingTraining != nil {
		return nil, fmt.Errorf("training record with ID [%s] already exists", recordID)
	}

	personnel, err := c.GetPersonnel(ctx, personnelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get personnel: %w", err)
	}

	if personnel.Status != domain.PersonnelStatusActive {
		return nil, fmt.Errorf("cannot complete training for personnel with status [%s]", personnel.Status)
	}

	if personnel.Campus != campus {
		return nil, fmt.Errorf("personnel is not enrolled in campus [%s] (current campus [%s])", campus, personnel.Campus)
	}

	hasTraining, err := c.personnelHasTraining(ctx, personnelID, trainingCode)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing training: %w", err)
	}
	if hasTraining {
		return nil, fmt.Errorf("personnel has already completed training with code [%s]", trainingCode)
	}

	training := &domain.Training{
		RecordID:     recordID,
		PersonnelID:  personnelID,
		Campus:       campus,
		TrainingCode: trainingCode,
		CompletedAt:  completedAt,
		IssuedBy:     issuedBy,
		Status:       domain.TrainingStatusCompleted,
	}

	// Store primary record
	trainingBytes, err := json.Marshal(training)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal training: %w", err)
	}

	if err := ctx.GetStub().PutState(trainingKey(recordID), trainingBytes); err != nil {
		return nil, fmt.Errorf("failed to put training state: %w", err)
	}

	// Composite index: by personnel (timeline)
	// Pattern `training_byPersonnel~SF-12345~2024-01-01T12:00:00Z~TR-987`
	// Allows "All training for this personnel", "Ordered history"
	byPersonnelKey, err := ctx.GetStub().CreateCompositeKey(
		"training_byPersonnel",
		[]string{personnelID, completedAt, recordID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key byPersonnel: %w", err)
	}

	if err := ctx.GetStub().PutState(byPersonnelKey, []byte{0x00}); err != nil {
		return nil, fmt.Errorf("failed to put state for composite key byPersonnel: %w", err)
	}

	// Composite index: byTrainingCode (qualification checks)
	// Pattern `training_byCode~ENG-WARP-201~SF-12345~TR-987`
	// Allows "Who has completed ENG-WARP-201?", "Promotion validation"
	byCodeKey, err := ctx.GetStub().CreateCompositeKey(
		"training_byCode",
		[]string{trainingCode, personnelID, recordID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key byCode: %w", err)
	}

	if err := ctx.GetStub().PutState(byCodeKey, []byte{0x00}); err != nil {
		return nil, fmt.Errorf("failed to put state for composite key byCode: %w", err)
	}

	return training, nil
}

func (c *PersonnelContract) personnelHasTraining(ctx contractapi.TransactionContextInterface, personnelID, trainingCode string) (bool, error) {
	training, err := c.getTrainingByCodeForPersonnel(ctx, trainingCode, personnelID)
	if err != nil {
		return false, err
	}
	return training != nil, nil
}

func (c *PersonnelContract) getTrainingByCodeForPersonnel(ctx contractapi.TransactionContextInterface, trainingCode, personnelID string) (*domain.Training, error) {
	// Query composite index for this personnel and training code
	// Pattern `training_byCode~ENG-WARP-201~SF-12345~TR-987`
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(
		"training_byCode",
		[]string{trainingCode, personnelID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query composite key byCode: %w", err)
	}
	defer iterator.Close()

	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate composite key byCode: %w", err)
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(response.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to split composite key byCode: %w", err)
		}

		if len(compositeKeyParts) != 3 {
			continue // skip invalid keys
		}

		recordID := compositeKeyParts[2]

		trainingBytes, err := ctx.GetStub().GetState(trainingKey(recordID))
		if err != nil {
			return nil, fmt.Errorf("failed to get training state: %w", err)
		}
		if trainingBytes == nil {
			continue // skip if training record is missing
		}

		var training domain.Training
		if err := json.Unmarshal(trainingBytes, &training); err != nil {
			return nil, fmt.Errorf("failed to unmarshal training data: %w", err)
		}

		if training.Status == domain.TrainingStatusCompleted {
			return &training, nil
		}
	}

	return nil, nil
}
