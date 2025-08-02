package models

import "time"

type Project struct {
	Name         string                  `json:"name"`
	Type         string                  `json:"type"`
	Path         string                  `json:"path"`
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
