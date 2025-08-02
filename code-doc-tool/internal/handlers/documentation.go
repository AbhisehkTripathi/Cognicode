package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func DownloadDocumentation(c *fiber.Ctx) error {
	filename := c.Params("filename")

	// Validate filename
	if filename == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Filename is required",
		})
	}

	// Construct file path
	filePath := filepath.Join("./output", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": "Documentation not found",
		})
	}

	// Set headers for file download
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.SendFile(filePath)
}

func GetStatus(c *fiber.Ctx) error {
	jobID := c.Params("jobId")

	// Check if output file exists
	outputPath := fmt.Sprintf("./output/%s_documentation.docx", jobID)
	if _, err := os.Stat(outputPath); err == nil {
		return c.JSON(fiber.Map{
			"status":       "completed",
			"message":      "Documentation generated successfully",
			"download_url": fmt.Sprintf("/api/download/%s_documentation.docx", jobID),
		})
	}

	// Check if upload directory exists (processing)
	uploadPath := fmt.Sprintf("./uploads/%s", jobID)
	if _, err := os.Stat(uploadPath); err == nil {
		return c.JSON(fiber.Map{
			"status":  "processing",
			"message": "Documentation is being generated",
		})
	}

	return c.Status(404).JSON(fiber.Map{
		"error": "Job not found",
	})
}
