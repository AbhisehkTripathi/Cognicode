package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func AnalyzeProject(codeFilePath string) (string, error) {
	fmt.Printf("codeFilePath: %s\n", codeFilePath)
	file, err := os.Open(codeFilePath)
	if err != nil {
		return "", fmt.Errorf("cannot open code file: %w", err)
	}
	defer file.Close()

	formatTemplate := `
		# Project Technical Documentation

		## 1. Overview
		- Purpose of the project
		- High-level description of what it does

		## 2. Technology Stack
		- Languages used
		- Frameworks / Libraries
		- External Services (APIs, DBs, etc.)

		## 3. Architecture
		- High-level description (monolith, microservices, etc.)
		- Folder / module structure
		- Data flow or sequence diagram (if applicable)

		## 4. Setup & Installation
		- Prerequisites
		- Installation steps
		- How to run locally / deploy

		## 5. APIs
		- Endpoint details (method, path, description, parameters, response)

		## 6. Functions / Classes
		- Function name, inputs, outputs, purpose

		## 7. Error Handling
		- Common error codes
		- Known failure scenarios

		## 8. Usage Example
		- Sample request (curl / Python snippet)
		- Sample response

		## 9. Limitations
		- Known limitations
		- Model restrictions

		## 10. Future Improvements
		- Planned features
		- Possible optimizations
`

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, _ := w.CreateFormFile("code_file", codeFilePath)
	if _, err := io.Copy(fw, file); err != nil {
		return "", fmt.Errorf("failed to copy code file: %w", err)
	}
	_ = w.WriteField("format", formatTemplate)
	w.Close()

	url := "http://localhost:8000/analyze"
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not call analyze endpoint: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("agent error: %s", respBody)
	}

	doc := struct {
		Document string `json:"document"`
	}{}
	if err := json.Unmarshal(respBody, &doc); err != nil {
		return "", fmt.Errorf("invalid response from agent: %w", err)
	}
	fmt.Println("Document: ", doc.Document)
	fmt.Println(doc.Document)

	return doc.Document, nil
}
