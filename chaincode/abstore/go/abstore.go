/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ABstore Chaincode 구현
type ABstore struct {
	contractapi.Contract
}

// 기존 ABstore용 상수
var Admin = "Admin"

// 대출 요청 구조체
type LoanRequest struct {
	ID           string `json:"id"`
	Requester    string `json:"requester"`
	Amount       int    `json:"amount"`
	DurationDays int    `json:"durationDays"`
	Status       string `json:"status"`   // "Pending", "Active"
	Provider     string `json:"provider"` // 수락한 대출 제공자
	StartTime    int64  `json:"startTime"` // 대출 시작 시간 (Unix timestamp)
}

// =============== 기존 ABstore 기능 ===============

// 초기화
func (t *ABstore) Init(ctx contractapi.TransactionContextInterface, A string, Aval int, B string, Bval int) error {
	fmt.Println("ABstore Init")
	err := ctx.GetStub().PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(Admin, []byte("0"))
	if err != nil {
		return err
	}

	return nil
}

// 송금
func (t *ABstore) Invoke(ctx contractapi.TransactionContextInterface, A, B string, X int) error {
	var Aval, Bval, Adminval int

	Avalbytes, err := ctx.GetStub().GetState(A)
	if err != nil {
		return fmt.Errorf("failed to get state for %s", A)
	}
	if Avalbytes == nil {
		return fmt.Errorf("entity %s not found", A)
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := ctx.GetStub().GetState(B)
	if err != nil {
		return fmt.Errorf("failed to get state for %s", B)
	}
	if Bvalbytes == nil {
		return fmt.Errorf("entity %s not found", B)
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	Adminvalbytes, err := ctx.GetStub().GetState(Admin)
	if err != nil {
		return fmt.Errorf("failed to get state for Admin")
	}
	if Adminvalbytes == nil {
		return fmt.Errorf("Admin not found")
	}
	Adminval, _ = strconv.Atoi(string(Adminvalbytes))

	Aval = Aval - X
	Bval = Bval + (X - X/10)
	Adminval = Adminval + (X / 10)

	err = ctx.GetStub().PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(Admin, []byte(strconv.Itoa(Adminval)))
	if err != nil {
		return err
	}

	return nil
}

// 삭제
func (t *ABstore) Delete(ctx contractapi.TransactionContextInterface, A string) error {
	err := ctx.GetStub().DelState(A)
	if err != nil {
		return fmt.Errorf("failed to delete state for %s", A)
	}
	return nil
}

// 단일 조회
func (t *ABstore) Query(ctx contractapi.TransactionContextInterface, A string) (string, error) {
	Avalbytes, err := ctx.GetStub().GetState(A)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to get state for %s", A))
	}
	if Avalbytes == nil {
		return "", errors.New(fmt.Sprintf("nil amount for %s", A))
	}
	return string(Avalbytes), nil
}

// 전체 조회
func (t *ABstore) GetAllQuery(ctx contractapi.TransactionContextInterface) ([]string, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var wallet []string
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		jsonResp := "{\"Name\":\"" + queryResponse.Key + "\",\"Amount\":\"" + string(queryResponse.Value) + "\"}"
		wallet = append(wallet, jsonResp)
	}
	return wallet, nil
}

// =============== 깐부 대출 기능 ===============

// 대출 요청 생성
func (t *ABstore) CreateLoanRequest(ctx contractapi.TransactionContextInterface, id, requester string, amount, durationDays int) error {
	exists, err := t.LoanRequestExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("loan request %s already exists", id)
	}

	loan := LoanRequest{
		ID:           id,
		Requester:    requester,
		Amount:       amount,
		DurationDays: durationDays,
		Status:       "Pending",
		Provider:     "",
		StartTime:    0,
	}

	loanJSON, err := json.Marshal(loan)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, loanJSON)
}

// 대출 요청 승인
func (t *ABstore) ApproveLoanRequest(ctx contractapi.TransactionContextInterface, id, provider string) error {
	loanJSON, err := ctx.GetStub().GetState(id)
	if err != nil || loanJSON == nil {
		return fmt.Errorf("loan request %s does not exist", id)
	}

	var loan LoanRequest
	err = json.Unmarshal(loanJSON, &loan)
	if err != nil {
		return err
	}

	if loan.Status != "Pending" {
		return fmt.Errorf("loan request %s is not pending", id)
	}

	loan.Status = "Active"
	loan.Provider = provider
	loan.StartTime = time.Now().Unix()

	updatedLoanJSON, err := json.Marshal(loan)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, updatedLoanJSON)
}

// 대출 요청 삭제
func (t *ABstore) DeleteLoanRequest(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := t.LoanRequestExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("loan request %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// 대출 요청 수정
func (t *ABstore) UpdateLoanRequest(ctx contractapi.TransactionContextInterface, id string, newAmount, newDurationDays int) error {
	loanJSON, err := ctx.GetStub().GetState(id)
	if err != nil || loanJSON == nil {
		return fmt.Errorf("loan request %s does not exist", id)
	}

	var loan LoanRequest
	err = json.Unmarshal(loanJSON, &loan)
	if err != nil {
		return err
	}

	if loan.Status != "Pending" {
		return fmt.Errorf("cannot update loan request %s because it is not pending", id)
	}

	loan.Amount = newAmount
	loan.DurationDays = newDurationDays

	updatedLoanJSON, err := json.Marshal(loan)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, updatedLoanJSON)
}

// 단일 대출 요청 조회
func (t *ABstore) QueryLoanRequest(ctx contractapi.TransactionContextInterface, id string) (*LoanRequest, error) {
	loanJSON, err := ctx.GetStub().GetState(id)
	if err != nil || loanJSON == nil {
		return nil, fmt.Errorf("loan request %s does not exist", id)
	}

	var loan LoanRequest
	err = json.Unmarshal(loanJSON, &loan)
	if err != nil {
		return nil, err
	}

	return &loan, nil
}

// 전체 대출 요청 조회
func (t *ABstore) QueryAllLoanRequests(ctx contractapi.TransactionContextInterface) ([]*LoanRequest, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var loans []*LoanRequest
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var loan LoanRequest
		err = json.Unmarshal(queryResponse.Value, &loan)
		// LoanRequest 구조로 Unmarshal이 실패하면 무시 (기존 ABstore 데이터 때문)
		if err == nil {
			loans = append(loans, &loan)
		}
	}
	return loans, nil
}

// 대출 요청 존재 여부 확인
func (t *ABstore) LoanRequestExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	loanJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}
	return loanJSON != nil, nil
}

// 메인 함수
func main() {
	cc, err := contractapi.NewChaincode(new(ABstore))
	if err != nil {
		panic(err.Error())
	}

	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting ABstore chaincode: %s", err)
	}
}
