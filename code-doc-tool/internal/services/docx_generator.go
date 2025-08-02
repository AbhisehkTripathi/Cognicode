package services

import (
	"fmt"
	"strings"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/docx"

	"code-doc-tool/internal/models"
)

type DocxGenerator struct{}

func NewDocxGenerator() *DocxGenerator {
	return &DocxGenerator{}
}

func (dg *DocxGenerator) GenerateDocumentation(project *models.Project, outputPath string) error {
	// Create new document
	doc, err := godocx.NewDocument()
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	// Add title
	titlePara := doc.AddParagraph("")
	titlePara.AddText(fmt.Sprintf("%s - Code Documentation", project.Name))

	// Add project overview
	dg.addProjectOverview(doc, project)

	// Add dependencies section
	dg.addDependencies(doc, project)

	// Add directory structure
	dg.addDirectoryStructure(doc, project)

	// Add file analysis
	dg.addFileAnalysis(doc, project)

	// Save document
	if err := doc.Save(); err != nil {
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
