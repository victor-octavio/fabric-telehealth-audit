package diagnosticlogs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type DiagnosticRecord struct {
	ID          string `json:"id,omitempty"`
	PatientID   string `json:"patient_id,omitempty"`
	PhysicianID string `json:"physician_id,omitempty"`
	Diagnosis   string `json:"diagnosis,omitempty"`
	Observation string `json:"observation,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

func (s *SmartContract) InsertDiagnostic(ctx contractapi.TransactionContextInterface, id, patientId, physicianId, diganosis, observation, createdAt, updatedAt string) error {
	exists, err := s.DiagnosticExists(ctx, id)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("diagnosis %s already exists", id)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	record := DiagnosticRecord{
		ID:          id,
		PatientID:   patientId,
		PhysicianID: physicianId,
		Diagnosis:   diganosis,
		Observation: observation,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	recordJson, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, recordJson)
}

func (s *SmartContract) DiagnosticExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	record, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}
	return record != nil, nil
}
