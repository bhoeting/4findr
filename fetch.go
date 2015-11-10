package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// ClassData is the JSON data structure for a class
type ClassData struct {
	ShortName string  `json:"ShortName"`
	Professor string  `json:"name"`
	Number    string  `json:"number"`
	Title     string  `json:"Title"`
	GPA       float64 `json:"avggpa"`
}

func main() {
	for index := range subjects {
		if err := subjects[index].fetchAndSetDeptID(); err != nil {
			fmt.Print(err.Error())
			continue
		}
	}

	if err := fetchClassData(); err != nil {
		log.Fatal(err.Error())
	}
}

func fetchClassData() error {
	for _, subject := range subjects {
		classesDataHTTPResponse, err := http.Get(fmt.Sprintf(classesURI, subject.Short, subject.DeptID))
		if err != nil {
			return err
		}

		classesJSON, err := ioutil.ReadAll(classesDataHTTPResponse.Body)

		if err != nil {
			return err
		}

		var classesData []ClassData
		json.Unmarshal(classesJSON, &classesData)
	}

	return nil
}

func (subject *Subject) fetchAndSetDeptID() error {
	type DeptID struct {
		DID int `json:"did"`
	}

	deptIDHTTPResponse, err := http.Get(fmt.Sprintf(deptIdURI, subject.Short))
	if err != nil {
		log.Fatalf("Failed to get departement ID: %v", err.Error())
	}

	body, _ := ioutil.ReadAll(deptIDHTTPResponse.Body)

	var deptIDs []DeptID
	json.Unmarshal(body, &deptIDs)

	if len(deptIDs) > 0 && deptIDs[0].DID > -1 {
		subject.DeptID = uint(deptIDs[0].DID)
		return nil
	}

	return fmt.Errorf("Subject %v does not have a deptartment ID.", subject.Short)
}
