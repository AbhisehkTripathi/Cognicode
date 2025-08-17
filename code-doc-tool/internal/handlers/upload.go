package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"code-doc-tool/internal/services"
	"code-doc-tool/internal/utils"
)

type UploadResponse struct {
	JobID   string `json:"job_id"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func CollectSourceFiles(root string, exts []string) ([]string, error) {
	var files []string
	extMap := map[string]bool{}
	for _, e := range exts {
		extMap[strings.ToLower(e)] = true
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if extMap[strings.ToLower(filepath.Ext(path))] {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}

func UploadCodebase(c *fiber.Ctx) error {
	// Get uploaded file
	file, err := c.FormFile("codebase")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidArchive(ext) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid file type. Please upload .zip, .tar, or .tar.gz files",
		})
	}

	jobID := uuid.New().String()

	uploadPath := fmt.Sprintf("./uploads/%s", jobID)
	if err := utils.CreateDir(uploadPath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}

	// Save uploaded file
	filePath := filepath.Join(uploadPath, file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to save uploaded file",
		})
	}

	// Process asynchronously
	go processCodebase(jobID, filePath, file.Filename)

	return c.JSON(UploadResponse{
		JobID:   jobID,
		Message: "File uploaded successfully. Processing started.",
		Status:  "processing",
	})
}
func processCodebase(jobID, filePath, filename string) {
	log.Printf("Starting processing for job %s", jobID)

	extractPath := fmt.Sprintf("./uploads/%s/extracted", jobID)
	if err := utils.ExtractArchive(filePath, extractPath); err != nil {
		log.Printf("Failed to extract archive for job %s: %v", jobID, err)
		return
	}
	log.Printf("Extraction complete for job %s", extractPath)

	// Collect code files (.py, .js, .ts, .php, .go, ... add others as needed)
	exts := []string{".py", ".js", ".ts", ".php", ".go"}
	codeFiles, err := CollectSourceFiles(extractPath, exts)
	if err != nil || len(codeFiles) == 0 {
		log.Printf("No source files found for job %s: %v", jobID, err)
		return
	}

	// Analyze files (could aggregate, or select main if preferred)
	var docs []string
	for _, codeFile := range codeFiles {
		log.Printf("Analyzing file: %s", codeFile)
		doc, err := services.AnalyzeProject(codeFile)
		if err != nil {
			log.Printf("File analysis failed for %s: %v", codeFile, err)
			continue
		}
		docs = append(docs, doc)
	}

	// Combine all docs into one (simple join, or make a section per file)
	combinedDoc := strings.Join(docs, "\n\n---\n\n")

	// Generate documentation file (save as .docx, or markdown, as you wish)
	generator := services.NewDocxGenerator()
	outputPath := fmt.Sprintf("./output/%s_documentation.docx", jobID)
	if err := generator.GenerateDocumentation(combinedDoc, outputPath); err != nil {
		log.Printf("Failed to generate documentation for job %s: %v", jobID, err)
		return
	}
	log.Printf("Documentation generated successfully for job %s", jobID)

	utils.CleanupDir(fmt.Sprintf("./uploads/%s", jobID))
}

func processCodebaseOld(jobID, filePath, filename string) {
	log.Printf("Starting processing for job %s", jobID)

	extractPath := fmt.Sprintf("./uploads/%s/extracted", jobID)
	if err := utils.ExtractArchive(filePath, extractPath); err != nil {
		log.Printf("Failed to extract archive for job %s: %v", jobID, err)
		return
	}
	log.Printf("Extraction complete for job %s", extractPath)

	// Analyze codebase
	project, err := services.AnalyzeProject(extractPath)
	if err != nil {
		log.Printf("Failed to analyze project for job %s: %v", jobID, err)
		return
	}
	log.Printf("Analysis complete for job %s: %+v", jobID, project)

	// Generate documentation
	generator := services.NewDocxGenerator()
	outputPath := fmt.Sprintf("./output/%s_documentation.docx", jobID)
	if err := generator.GenerateDocumentation(project, outputPath); err != nil {
		log.Printf("Failed to generate documentation for job %s: %v", jobID, err)
		return
	}
	log.Printf("Documentation generated successfully for job %s", jobID)

	utils.CleanupDir(fmt.Sprintf("./uploads/%s", jobID))
}

func isValidArchive(ext string) bool {
	validExts := []string{".zip", ".tar", ".gz"}
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}
