package handlers

import (
	"fmt"
	"log"
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

func UploadCodebase(c *fiber.Ctx) error {
	// Get uploaded file
	file, err := c.FormFile("codebase")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Validate file type (zip, tar, tar.gz)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidArchive(ext) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid file type. Please upload .zip, .tar, or .tar.gz files",
		})
	}

	// Generate unique job ID
	jobID := uuid.New().String()

	// Create upload directory
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

	// Extract archive
	extractPath := fmt.Sprintf("./uploads/%s/extracted", jobID)
	if err := utils.ExtractArchive(filePath, extractPath); err != nil {
		log.Printf("Failed to extract archive for job %s: %v", jobID, err)
		return
	}
	log.Printf("Extraction complete for job %s", jobID)

	// Analyze codebase
	analyzer := services.NewFileAnalyzer()
	project, err := analyzer.AnalyzeProject(extractPath)
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
