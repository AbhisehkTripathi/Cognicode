package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/docx"

	"code-doc-tool/internal/models"
)

// saveDocxFile saves a godocx document to a specific path
func saveDocxFile(doc *docx.RootDoc, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return doc.Write(file)
}

type DocxGenerator struct{}

func ensureDir(dir string) error {
	if dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

func NewDocxGenerator() *DocxGenerator {
	return &DocxGenerator{}
}

func (dg *DocxGenerator) GenerateDocumentation(project *models.Project, outputPath string) error {
	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := ensureDir(outputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create new document
	doc, err := godocx.NewDocument()
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	// Add main title
	titlePara := doc.AddParagraph("")
	titlePara.AddText(fmt.Sprintf("%s - Code Documentation", project.Name))
	titlePara.Style("Title")

	// Add all custom sections in order
	dg.addCustomSections(doc, project)

	// Save document to outputPath
	if err := saveDocxFile(doc, outputPath); err != nil {
		return fmt.Errorf("failed to save document: %w", err)
	}

	return nil
}

func (dg *DocxGenerator) addProjectOverview(doc *docx.RootDoc, project *models.Project) {
	// Add heading
	headingPara := doc.AddParagraph("")
	headingPara.AddText("Project Overview")
	headingPara.Style("Heading1")

	// Add project details as paragraphs
	addKeyValue := func(key, value string) {
		para := doc.AddParagraph("")
		para.AddText(fmt.Sprintf("%s: %s", key, value))
	}

	addKeyValue("Project Name", project.Name)
	addKeyValue("Project Type", project.Type)
	addKeyValue("Total Files", fmt.Sprintf("%d", len(project.Files)))
	addKeyValue("Analysis Date", "Generated automatically")
}

func (dg *DocxGenerator) addDependencies(doc *docx.RootDoc, project *models.Project) {
	// Check if there are any dependencies
	hasDeps := false
	for _, deps := range project.Dependencies {
		if len(deps) > 0 {
			hasDeps = true
			break
		}
	}

	if !hasDeps {
		para := doc.AddParagraph("")
		para.AddText("No dependencies found in the project.")
		return
	}

	// Add heading
	headingPara := doc.AddParagraph("")
	headingPara.AddText("Dependencies")
	headingPara.Style("Heading1")

	// Add dependencies grouped by type
	for depType, depList := range project.Dependencies {
		if len(depList) == 0 {
			continue
		}

		// Add dependency type heading
		typePara := doc.AddParagraph("")
		typePara.AddText(depType)
		typePara.Style("Heading2")

		// Add each dependency
		for _, dep := range depList {
			para := doc.AddParagraph("")
			// Assuming Dependency has Name and Version fields
			para.AddText(fmt.Sprintf("• %s: %s", dep.Name, dep.Version))
		}
	}
}

func (dg *DocxGenerator) addDirectoryStructure(doc *docx.RootDoc, project *models.Project) {
	if len(project.Structure) == 0 {
		para := doc.AddParagraph("")
		para.AddText("No directory structure information available.")
		return
	}

	// Add heading
	headingPara := doc.AddParagraph("")
	headingPara.AddText("Directory Structure")
	headingPara.Style("Heading1")

	// Add a paragraph for each directory entry
	for _, dir := range project.Structure {
		dg.addDirectoryEntry(doc, dir, 0)
	}
}

func (dg *DocxGenerator) addFileAnalysis(doc *docx.RootDoc, project *models.Project) {
	if len(project.Files) == 0 {
		para := doc.AddParagraph("")
		para.AddText("No code files found for analysis.")
		return
	}

	// Add heading
	headingPara := doc.AddParagraph("")
	headingPara.AddText("File Analysis")
	headingPara.Style("Heading1")

	// Add files as a list
	for _, file := range project.Files {
		para := doc.AddParagraph("")
		para.AddText(fmt.Sprintf("%s (%d bytes)", file.Path, file.Size))
	}

	// Analyze by language
	langStats := make(map[string]int)
	for _, file := range project.Files {
		langStats[file.Language]++
	}

	// Add language statistics
	headingPara = doc.AddParagraph("")
	headingPara.AddText("Language Distribution")
	headingPara.Style("Heading2")

	// Add language statistics as a list
	total := len(project.Files)
	for lang, count := range langStats {
		percentage := float64(count) / float64(total) * 100

		para := doc.AddParagraph("")
		para.AddText(fmt.Sprintf("%s: %d files (%.1f%%)", lang, count, percentage))
	}
}

// addDirectoryEntry adds a directory entry with proper indentation
func (dg *DocxGenerator) addDirectoryEntry(doc *docx.RootDoc, dir models.DirectoryNode, level int) {
	// Create indentation
	indent := strings.Repeat("    ", level)

	// Add the directory name
	para := doc.AddParagraph("")
	para.AddText(indent + dir.Name)

	// Recursively add children
	for _, child := range dir.Children {
		dg.addDirectoryEntry(doc, child, level+1)
	}
}
func (dg *DocxGenerator) addCustomSections(doc *docx.RootDoc, p *models.Project) {
	// 1. Project Overview
	doc.AddHeading("1. Project Overview", 1)
	if p.Overview != "" {
		doc.AddParagraph(p.Overview)
	} else {
		doc.AddParagraph("No overview available.")
	}

	// 2. Architecture Diagram
	doc.AddHeading("2. Architecture Diagram", 1)
	if p.Architecture != "" {
		doc.AddParagraph(p.Architecture)
		// Optionally: embed image if supported
	} else {
		doc.AddParagraph("No architecture diagram provided.")
	}

	// 3. Tech Stack Summary
	doc.AddHeading("3. Tech Stack Summary", 1)
	if len(p.TechStack) > 0 {
		for _, tech := range p.TechStack {
			doc.AddParagraph("• " + tech)
		}
	} else {
		doc.AddParagraph("No tech stack information available.")
	}

	// 4. Folder Structure
	doc.AddHeading("4. Folder Structure", 1)
	if len(p.FolderStructure) > 0 {
		table := doc.AddTable()
		row := table.AddRow()
		row.AddCell().AddParagraph("Folder")
		row.AddCell().AddParagraph("Description")
		for folder, desc := range p.FolderStructure {
			row := table.AddRow()
			row.AddCell().AddParagraph(folder)
			row.AddCell().AddParagraph(desc)
		}
	} else {
		doc.AddParagraph("No folder structure information available.")
	}

	// 5. Setup Instructions
	doc.AddHeading("5. Setup Instructions", 1)
	if len(p.SetupInstructions) > 0 {
		for _, step := range p.SetupInstructions {
			doc.AddParagraph("• " + step)
		}
	} else {
		doc.AddParagraph("No setup instructions provided.")
	}

	// 6. API Reference
	doc.AddHeading("6. API Reference", 1)
	if len(p.APIEndpoints) > 0 {
		table := doc.AddTable()
		row := table.AddRow()
		row.AddCell().AddParagraph("Method")
		row.AddCell().AddParagraph("Path")
		row.AddCell().AddParagraph("Description")
		row.AddCell().AddParagraph("Example")
		for _, ep := range p.APIEndpoints {
			row := table.AddRow()
			row.AddCell().AddParagraph(ep.Method)
			row.AddCell().AddParagraph(ep.Path)
			row.AddCell().AddParagraph(strings.Join(ep.Middleware, " → "))
			row.AddCell().AddParagraph(ep.Handler)
			row.AddCell().AddParagraph(ep.CurlExample)
		}
	} else {
		doc.AddParagraph("No API reference available.")
	}

	// 7. Parsers Info
	doc.AddHeading("7. Parsers Info", 1)
	if len(p.ParsersInfo) > 0 {
		table := doc.AddTable()
		row := table.AddRow()
		row.AddCell().AddParagraph("Language")
		row.AddCell().AddParagraph("Parser Details")
		for lang, desc := range p.ParsersInfo {
			row := table.AddRow()
			row.AddCell().AddParagraph(lang)
			row.AddCell().AddParagraph(desc)
		}
	} else {
		doc.AddParagraph("No parser information available.")
	}

	// 8. Data Flow
	doc.AddHeading("8. Data Flow", 1)
	if p.DataFlow != "" {
		doc.AddParagraph(p.DataFlow)
	} else {
		doc.AddParagraph("No data flow information provided.")
	}

	// 9. External Services
	doc.AddHeading("9. External Services", 1)
	if len(p.ExternalServices) > 0 {
		for _, svc := range p.ExternalServices {
			doc.AddParagraph("• " + svc)
		}
	} else {
		doc.AddParagraph("No external services listed.")
	}

	// 10. Deployment Info
	doc.AddHeading("10. Deployment Info", 1)
	if len(p.DeploymentInfo) > 0 {
		for _, dep := range p.DeploymentInfo {
			doc.AddParagraph(dep)
		}
	} else {
		doc.AddParagraph("No deployment info provided.")
	}

	// 11. Future Roadmap
	doc.AddHeading("11. Future Roadmap", 1)
	if len(p.FutureRoadmap) > 0 {
		for _, item := range p.FutureRoadmap {
			doc.AddParagraph("• " + item)
		}
	} else {
		doc.AddParagraph("No roadmap provided.")
	}

	// 12. Common Issues
	doc.AddHeading("12. Common Issues", 1)
	if len(p.CommonIssues) > 0 {
		for _, issue := range p.CommonIssues {
			doc.AddParagraph("• " + issue)
		}
	} else {
		doc.AddParagraph("No common issues documented.")
	}

	// 13. Developer Notes
	doc.AddHeading("13. Developer Notes", 1)
	if len(p.DeveloperNotes) > 0 {
		for _, note := range p.DeveloperNotes {
			doc.AddParagraph(note)
		}
	} else {
		doc.AddParagraph("No developer notes provided.")
	}
}
