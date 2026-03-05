package main

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

func (s *SmartContract) ReadDiagnosis(ctx contractapi.TransactionContextInterface, id string) (*DiagnosticRecord, error) {
	recordJson, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, err
	}

	if recordJson == nil {
		return nil, fmt.Errorf("diagnosis not found, { id: %s }", id)
	}

	var record DiagnosticRecord

	if err := json.Unmarshal(recordJson, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (s *SmartContract) GetHistory(ctx contractapi.TransactionContextInterface, id string) ([]*DiagnosticRecord, error) {
	historyIterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, err
	}

	defer historyIterator.Close()
	var history []*DiagnosticRecord

	for historyIterator.HasNext() {
		entry, err := historyIterator.Next()
		if err != nil {
			return nil, err
		}

		var record DiagnosticRecord
		if err := json.Unmarshal(entry.Value, &record); err != nil {
			return nil, err
		}

		history = append(history, &record)
	}
	return history, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("error creating chaincode, err = %v\n", err)
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("an error occurred during the chaincode start up")
	}
}
