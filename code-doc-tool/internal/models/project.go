package models

import "time"

type APIEndpoint struct {
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	Middleware  []string `json:"middleware"`
	Handler     string   `json:"handler"`
	Description string   `json:"description"`
	CurlExample string   `json:"curl_example"`
}

type Project struct {
	Name              string            `json:"name"`
	Type              string            `json:"type"`
	Path              string            `json:"path"`
	Overview          string            `json:"overview"`
	TechStack         []string          `json:"tech_stack"`
	Architecture      string            `json:"architecture"`
	FolderStructure   map[string]string `json:"folder_structure"`
	SetupInstructions []string          `json:"setup_instructions"`
	APIEndpoints      []APIEndpoint     `json:"api_endpoints"`
	ParsersInfo       map[string]string `json:"parsers_info"`
	DataFlow          string            `json:"data_flow"`
	ExternalServices  []string          `json:"external_services"`
	DeploymentInfo    []string          `json:"deployment_info"`
	FutureRoadmap     []string          `json:"future_roadmap"`
	CommonIssues      []string          `json:"common_issues"`
	DeveloperNotes    []string          `json:"developer_notes"`

	Dependencies map[string][]Dependency `json:"dependencies"`
	Files        []FileInfo              `json:"files"`
	Structure    []DirectoryNode         `json:"structure"`
	CreatedAt    time.Time               `json:"created_at"`
}

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"`
}

type FileInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Language  string `json:"language"`
}

type DirectoryNode struct {
	Name     string          `json:"name"`
	Path     string          `json:"path"`
	IsDir    bool            `json:"is_dir"`
	Size     int64           `json:"size"`
	Children []DirectoryNode `json:"children,omitempty"`
}

type Job struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Progress  int       `json:"progress"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
