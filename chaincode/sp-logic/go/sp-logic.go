/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a picture
type SmartContract struct {
	contractapi.Contract
}

// Picture describes basic details of what makes up a picture
type Picture struct {
	Make        string `json:"make"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	Owner       string `json:"owner"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Picture
}

// InitLedger adds a base set of pictures to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	pictures := []Picture{
		Picture{Make: "A", Category: "Nature", Subcategory: "Rock", Owner: "A"},
		Picture{Make: "A", Category: "Nature", Subcategory: "Sky", Owner: "A"},
		Picture{Make: "A", Category: "Animal", Subcategory: "Lion", Owner: "A"},
		Picture{Make: "B", Category: "Nature", Subcategory: "Rock", Owner: "B"},
		Picture{Make: "B", Category: "Animal", Subcategory: "Snake", Owner: "B"},
		Picture{Make: "B", Category: "Insect", Subcategory: "Spider", Owner: "B"},
		Picture{Make: "C", Category: "Nature", Subcategory: "Rock", Owner: "C"},
		Picture{Make: "C", Category: "Nature", Subcategory: "Sky", Owner: "C"},
		Picture{Make: "C", Category: "Nature", Subcategory: "Planet", Owner: "C"},
		Picture{Make: "extD", Category: "Nature", Subcategory: "Rock", Owner: "C"},
	}

	for i, picture := range pictures {
		pictureAsBytes, _ := json.Marshal(picture)
		err := ctx.GetStub().PutState("PICTURE"+strconv.Itoa(i), pictureAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreatePicture(ctx contractapi.TransactionContextInterface, pictureNumber string, make string, category string, subcategory string, owner string) error {
	picture := Picture{
		Make:        make,
		Category:    category,
		Subcategory: subcategory,
		Owner:       owner,
	}

	pictureAsBytes, _ := json.Marshal(picture)

	return ctx.GetStub().PutState(pictureNumber, pictureAsBytes)
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryPicture(ctx contractapi.TransactionContextInterface, pictureNumber string) (*Picture, error) {
	pictureAsBytes, err := ctx.GetStub().GetState(pictureNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if pictureAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", pictureNumber)
	}

	picture := new(Picture)
	_ = json.Unmarshal(pictureAsBytes, picture)

	return picture, nil
}

// QueryAllCars returns all cars found in world state
func (s *SmartContract) QueryAllPictures(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		picture := new(Picture)
		_ = json.Unmarshal(queryResponse.Value, picture)

		queryResult := QueryResult{Key: queryResponse.Key, Record: picture}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeCarOwner updates the owner field of car with given id in world state
func (s *SmartContract) ChangePictureOwner(ctx contractapi.TransactionContextInterface, pictureNumber string, newOwner string) error {
	picture, err := s.QueryPicture(ctx, pictureNumber)

	if err != nil {
		return err
	}

	picture.Owner = newOwner

	pictureAsBytes, _ := json.Marshal(picture)

	return ctx.GetStub().PutState(pictureNumber, pictureAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create sp-logic chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting sp-logic chaincode: %s", err.Error())
	}
}
