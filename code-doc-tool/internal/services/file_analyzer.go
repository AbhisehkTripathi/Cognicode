package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"code-doc-tool/internal/models"
)

type FileAnalyzer struct{}

func NewFileAnalyzer() *FileAnalyzer {
	return &FileAnalyzer{}
}

func (fa *FileAnalyzer) AnalyzeProject(projectPath string) (*models.Project, error) {
	project := &models.Project{
		Name:         filepath.Base(projectPath),
		Path:         projectPath,
		Dependencies: make(map[string][]models.Dependency),
		Files:        []models.FileInfo{},
		Structure:    []models.DirectoryNode{},
	}

	// Detect project type
	project.Type = fa.detectProjectType(projectPath)

	// Parse dependencies
	if err := fa.parseDependencies(projectPath, project); err != nil {
		return nil, err
	}

	// Build directory structure
	if err := fa.buildDirectoryStructure(projectPath, project); err != nil {
		return nil, err
	}

	// Analyze files
	if err := fa.analyzeFiles(projectPath, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (fa *FileAnalyzer) detectProjectType(projectPath string) string {
	// Check for various project indicators
	indicators := map[string]string{
		"package.json":     "Node.js",
		"composer.json":    "PHP/Laravel",
		"requirements.txt": "Python",
		"go.mod":           "Go",
		"pom.xml":          "Java/Maven",
		"build.gradle":     "Java/Gradle",
		"Cargo.toml":       "Rust",
		"pubspec.yaml":     "Dart/Flutter",
	}

	for file, projectType := range indicators {
		if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
			return projectType
		}
	}

	return "Unknown"
}

func (fa *FileAnalyzer) parseDependencies(projectPath string, project *models.Project) error {
	switch project.Type {
	case "Node.js":
		return fa.parsePackageJSON(projectPath, project)
	case "PHP/Laravel":
		return fa.parseComposerJSON(projectPath, project)
	case "Python":
		return nil
		// return fa.parseRequirementsTxt(projectPath, project)
	case "Go":
		return nil
		// return fa.parseGoMod(projectPath, project)
	}
	return nil
}

func (fa *FileAnalyzer) parsePackageJSON(projectPath string, project *models.Project) error {
	packagePath := filepath.Join(projectPath, "package.json")
	data, err := os.ReadFile(packagePath)
	if err != nil {
		return err
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}

	// Parse dependencies
	for name, version := range pkg.Dependencies {
		project.Dependencies["production"] = append(project.Dependencies["production"],
			models.Dependency{Name: name, Version: version, Type: "production"})
	}

	for name, version := range pkg.DevDependencies {
		project.Dependencies["development"] = append(project.Dependencies["development"],
			models.Dependency{Name: name, Version: version, Type: "development"})
	}

	return nil
}

func (fa *FileAnalyzer) parseComposerJSON(projectPath string, project *models.Project) error {
	composerPath := filepath.Join(projectPath, "composer.json")
	data, err := os.ReadFile(composerPath)
	if err != nil {
		return err
	}

	var composer struct {
		Require    map[string]string `json:"require"`
		RequireDev map[string]string `json:"require-dev"`
	}

	if err := json.Unmarshal(data, &composer); err != nil {
		return err
	}

	// Parse dependencies
	for name, version := range composer.Require {
		project.Dependencies["production"] = append(project.Dependencies["production"],
			models.Dependency{Name: name, Version: version, Type: "production"})
	}

	for name, version := range composer.RequireDev {
		project.Dependencies["development"] = append(project.Dependencies["development"],
			models.Dependency{Name: name, Version: version, Type: "development"})
	}

	return nil
}

func (fa *FileAnalyzer) buildDirectoryStructure(projectPath string, project *models.Project) error {
	// Create a map to store nodes by their path
	nodes := make(map[string]*models.DirectoryNode)
	
	// First pass: create all nodes
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and common ignore patterns
		if fa.shouldSkip(path, info) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, _ := filepath.Rel(projectPath, path)
		if relPath == "." {
			nodes[relPath] = &models.DirectoryNode{
				Name:     filepath.Base(projectPath),
				Path:     relPath,
				IsDir:    true,
				Children: []models.DirectoryNode{},
			}
			return nil
		}

		nodes[relPath] = &models.DirectoryNode{
			Name:  info.Name(),
			Path:  relPath,
			IsDir: info.IsDir(),
			Size:  info.Size(),
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Second pass: build the tree structure
	for path, node := range nodes {
		if path == "." {
			// This is the root node
			project.Structure = append(project.Structure, *node)
			continue
		}

		// Find parent directory
		parentPath := filepath.Dir(path)
		if parentPath == "." {
			parentPath = "./"
		}

		if parent, exists := nodes[parentPath]; exists {
			parent.Children = append(parent.Children, *node)
		}
	}

	return err
}

func (fa *FileAnalyzer) analyzeFiles(projectPath string, project *models.Project) error {
	return filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || fa.shouldSkip(path, info) {
			return err
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if fa.isCodeFile(ext) {
			relPath, _ := filepath.Rel(projectPath, path)
			project.Files = append(project.Files, models.FileInfo{
				Name:      info.Name(),
				Path:      relPath,
				Extension: ext,
				Size:      info.Size(),
				Language:  fa.getLanguage(ext),
			})
		}

		return nil
	})
}

func (fa *FileAnalyzer) shouldSkip(path string, info os.FileInfo) bool {
	skipDirs := []string{"node_modules", "vendor", ".git", "dist", "build", "__pycache__"}
	skipFiles := []string{".DS_Store", "Thumbs.db"}

	name := info.Name()

	if info.IsDir() {
		for _, skipDir := range skipDirs {
			if name == skipDir {
				return true
			}
		}
	} else {
		for _, skipFile := range skipFiles {
			if name == skipFile {
				return true
			}
		}
	}

	return strings.HasPrefix(name, ".")
}

func (fa *FileAnalyzer) isCodeFile(ext string) bool {
	codeExts := []string{
		".js", ".ts", ".jsx", ".tsx", ".vue",
		".php", ".py", ".go", ".java", ".c", ".cpp",
		".html", ".css", ".scss", ".sass", ".less",
		".json", ".xml", ".yaml", ".yml", ".md",
	}

	for _, codeExt := range codeExts {
		if ext == codeExt {
			return true
		}
	}
	return false
}

func (fa *FileAnalyzer) getLanguage(ext string) string {
	langMap := map[string]string{
		".js":   "JavaScript",
		".ts":   "TypeScript",
		".jsx":  "React",
		".tsx":  "React TypeScript",
		".vue":  "Vue.js",
		".php":  "PHP",
		".py":   "Python",
		".go":   "Go",
		".java": "Java",
		".html": "HTML",
		".css":  "CSS",
		".scss": "SCSS",
		".json": "JSON",
		".md":   "Markdown",
	}

	if lang, exists := langMap[ext]; exists {
		return lang
	}
	return "Unknown"
}
