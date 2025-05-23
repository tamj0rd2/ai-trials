package main

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
)

//go:embed template.html
var templateFS embed.FS

// TemplateData holds the data structure needed for the HTML template
type TemplateData struct {
	Developers []string
	PairDays   map[string]map[string]int
}

// createTemplateData converts the raw pair days data into a template-friendly structure
func createTemplateData(pairDays map[Pair]int, devList []string) TemplateData {
	// Initialize the map of maps
	pairDaysMap := make(map[string]map[string]int)
	for _, dev := range devList {
		pairDaysMap[dev] = make(map[string]int)
	}

	// Fill in the data
	for pair, days := range pairDays {
		dev1, dev2 := pair[0], pair[1]
		pairDaysMap[dev1][dev2] = days
		pairDaysMap[dev2][dev1] = days
	}

	return TemplateData{
		Developers: devList,
		PairDays:   pairDaysMap,
	}
}

// generateHTMLFromTemplate generates HTML from the template file and writes to output file
func generateHTMLFromTemplate(data TemplateData, templatePath, outputPath string) error {
	// Parse the template from embedded file system
	tmpl, err := template.ParseFS(templateFS, templatePath)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing output file: %w", closeErr)
		}
	}()

	// Execute the template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	fmt.Println("Wrote", outputPath)

	// Open the created file
	if err := exec.Command("open", outputPath).Start(); err != nil {
		return fmt.Errorf("error opening HTML file: %w", err)
	}

	return nil
}
