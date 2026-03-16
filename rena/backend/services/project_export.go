package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// sanitizeFilename removes or replaces invalid filename characters
func sanitizeFilename(filename string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitized := reg.ReplaceAllString(filename, "_")
	return sanitized
}

// ProjectExport represents the structure of exported project data
type ProjectExport struct {
	Metadata     ProjectMetadata `json:"metadata"`
	WorkspaceXML string          `json:"workspace_xml"`
	Code         string          `json:"code"`
}

// ProjectMetadata contains project metadata
type ProjectMetadata struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PackageName string `json:"package_name"`
	Platform    string `json:"platform"`
	Color       string `json:"color"`
	VersionCode int    `json:"version_code"`
	VersionName string `json:"version_name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// sanitizePackageName converts package name to valid directory path
// e.g., "com.example.app" -> "com/example/app"
func sanitizePackageName(packageName string) string {
	return strings.ReplaceAll(packageName, ".", "/")
}

// CompressProject creates a .rnp file (zip archive) containing project data
func CompressProject(exportData ProjectExport, outputPath string) error {
	// Create the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Create zip writer
	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// 1. Add project.properties
	propertiesContent := generateProjectProperties(exportData.Metadata)
	err = addToZip(zipWriter, "project.properties", []byte(propertiesContent))
	if err != nil {
		return err
	}

	// 2. Create src/{package_name}/ directory path
	srcPath := fmt.Sprintf("src/%s/", sanitizePackageName(exportData.Metadata.PackageName))

	// 3. Add workspace.bky file
	workspaceFilename := sanitizeFilename(exportData.Metadata.Name)
	err = addToZip(zipWriter, srcPath+workspaceFilename+".bky", []byte(exportData.WorkspaceXML))
	if err != nil {
		return err
	}

	// 4. Add workspace.scm file (generated code in Scheme format)
	scmContent := convertToScmFormat(exportData.Code, exportData.Metadata.Name)
	err = addToZip(zipWriter, srcPath+workspaceFilename+".scm", []byte(scmContent))
	if err != nil {
		return err
	}

	// 5. Create assets/ directory (add a placeholder file)
	err = addToZip(zipWriter, "assets/.gitkeep", []byte("# Assets directory for project resources\n"))
	if err != nil {
		return err
	}

	return nil
}

// DecompressProject reads a .rnp file and extracts project data
func DecompressProject(rnpPath string) (*ProjectExport, error) {
	// Open the zip file
	r, err := zip.OpenReader(rnpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open rnp file: %w", err)
	}
	defer r.Close()

	exportData := &ProjectExport{}

	// Iterate through files in the archive
	for _, f := range r.File {
		switch {
		case f.Name == "project.properties":
			content, err := readFileFromZip(f)
			if err != nil {
				return nil, fmt.Errorf("failed to read project.properties: %w", err)
			}
			exportData.Metadata = parseProjectProperties(string(content))

		case strings.HasSuffix(f.Name, ".bky"):
			content, err := readFileFromZip(f)
			if err != nil {
				return nil, fmt.Errorf("failed to read .bky file: %w", err)
			}
			exportData.WorkspaceXML = string(content)

		case strings.HasSuffix(f.Name, ".scm"):
			content, err := readFileFromZip(f)
			if err != nil {
				return nil, fmt.Errorf("failed to read .scm file: %w", err)
			}
			exportData.Code = convertFromScmFormat(string(content))

		case f.Name == "assets/.gitkeep" || strings.HasPrefix(f.Name, "assets/"):
			// Skip assets files for now, just acknowledge them
			continue
		}
	}

	return exportData, nil
}

// addToZip adds a file with content to the zip archive
func addToZip(zipWriter *zip.Writer, filename string, content []byte) error {
	writer, err := zipWriter.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file in zip: %w", err)
	}

	_, err = writer.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write content to zip: %w", err)
	}

	return nil
}

// readFileFromZip reads content from a file in the zip archive
func readFileFromZip(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, rc)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ExportProjectToBytes creates a .rnp file in memory and returns it as bytes
func ExportProjectToBytes(exportData ProjectExport) ([]byte, error) {
	var buf bytes.Buffer

	// Create zip writer
	zipWriter := zip.NewWriter(&buf)

	// 1. Add project.properties
	propertiesContent := generateProjectProperties(exportData.Metadata)
	err := addToZip(zipWriter, "project.properties", []byte(propertiesContent))
	if err != nil {
		return nil, err
	}

	// 2. Create src/{package_name}/ directory path
	srcPath := fmt.Sprintf("src/%s/", sanitizePackageName(exportData.Metadata.PackageName))

	// 3. Add workspace.bky file
	workspaceFilename := sanitizeFilename(exportData.Metadata.Name)
	err = addToZip(zipWriter, srcPath+workspaceFilename+".bky", []byte(exportData.WorkspaceXML))
	if err != nil {
		return nil, err
	}

	// 4. Add workspace.scm file
	scmContent := convertToScmFormat(exportData.Code, exportData.Metadata.Name)
	err = addToZip(zipWriter, srcPath+workspaceFilename+".scm", []byte(scmContent))
	if err != nil {
		return nil, err
	}

	// 5. Create assets/ directory
	err = addToZip(zipWriter, "assets/.gitkeep", []byte("# Assets directory for project resources\n"))
	if err != nil {
		return nil, err
	}

	// Close the zip writer
	err = zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %w", err)
	}

	return buf.Bytes(), nil
}

// ImportProjectFromBytes reads .rnp file from bytes and returns project data
func ImportProjectFromBytes(data []byte) (*ProjectExport, error) {
	reader := bytes.NewReader(data)

	// Open the zip file from bytes
	r, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to open rnp file from bytes: %w", err)
	}

	exportData := &ProjectExport{}

	// Iterate through files in the archive
	for _, f := range r.File {
		switch {
		case f.Name == "project.properties":
			content, err := readFileFromZip(f)
			if err != nil {
				return nil, fmt.Errorf("failed to read project.properties: %w", err)
			}
			exportData.Metadata = parseProjectProperties(string(content))

		case strings.HasSuffix(f.Name, ".bky"):
			content, err := readFileFromZip(f)
			if err != nil {
				return nil, fmt.Errorf("failed to read .bky file: %w", err)
			}
			exportData.WorkspaceXML = string(content)

		case strings.HasSuffix(f.Name, ".scm"):
			content, err := readFileFromZip(f)
			if err != nil {
				return nil, fmt.Errorf("failed to read .scm file: %w", err)
			}
			exportData.Code = convertFromScmFormat(string(content))

		case f.Name == "assets/.gitkeep" || strings.HasPrefix(f.Name, "assets/"):
			// Skip assets files
			continue
		}
	}

	return exportData, nil
}

// GetRNPFileInfo reads a .rnp file and returns basic info without full extraction
func GetRNPFileInfo(rnpPath string) (*ProjectMetadata, error) {
	// Open the zip file
	r, err := zip.OpenReader(rnpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open rnp file: %w", err)
	}
	defer r.Close()

	// Look for project.properties
	for _, f := range r.File {
		if f.Name == "project.properties" {
			content, err := readFileFromZip(f)
			if err != nil {
				return nil, fmt.Errorf("failed to read project.properties: %w", err)
			}
			metadata := parseProjectProperties(string(content))
			return &metadata, nil
		}
	}

	return nil, fmt.Errorf("project.properties not found in rnp file")
}

// Helper functions

// generateProjectProperties creates a Java properties file content
func generateProjectProperties(metadata ProjectMetadata) string {
	var sb strings.Builder
	sb.WriteString("# Rena Platform Project Configuration\n")
	sb.WriteString(fmt.Sprintf("project.id=%s\n", metadata.ID))
	sb.WriteString(fmt.Sprintf("project.name=%s\n", metadata.Name))
	sb.WriteString(fmt.Sprintf("project.package=%s\n", metadata.PackageName))
	sb.WriteString(fmt.Sprintf("project.platform=%s\n", metadata.Platform))
	sb.WriteString(fmt.Sprintf("project.color=%s\n", metadata.Color))
	sb.WriteString(fmt.Sprintf("project.version.code=%d\n", metadata.VersionCode))
	sb.WriteString(fmt.Sprintf("project.version.name=%s\n", metadata.VersionName))
	sb.WriteString(fmt.Sprintf("project.created.at=%s\n", metadata.CreatedAt))
	sb.WriteString(fmt.Sprintf("project.updated.at=%s\n", metadata.UpdatedAt))
	return sb.String()
}

// parseProjectProperties parses Java properties file content
func parseProjectProperties(content string) ProjectMetadata {
	metadata := ProjectMetadata{}
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "project.id":
			metadata.ID = value
		case "project.name":
			metadata.Name = value
		case "project.package":
			metadata.PackageName = value
		case "project.platform":
			metadata.Platform = value
		case "project.color":
			metadata.Color = value
		case "project.version.code":
			fmt.Sscanf(value, "%d", &metadata.VersionCode)
		case "project.version.name":
			metadata.VersionName = value
		case "project.created.at":
			metadata.CreatedAt = value
		case "project.updated.at":
			metadata.UpdatedAt = value
		}
	}

	return metadata
}

// convertToScmFormat converts React Native code to Scheme-like format
func convertToScmFormat(code string, projectName string) string {
	var sb strings.Builder
	sb.WriteString(";; Generated Scheme-like representation\n")
	sb.WriteString(fmt.Sprintf(";; Project: %s\n\n", projectName))
	sb.WriteString("(define app-component\n")
	sb.WriteString("  (begin\n")

	// Simple conversion - wrap code in scheme-like structure
	// This is a placeholder - you can implement proper conversion logic
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			sb.WriteString(fmt.Sprintf("    ;; %s\n", line))
		}
	}

	sb.WriteString("  ))\n")
	return sb.String()
}

// convertFromScmFormat extracts code from Scheme-like format
func convertFromScmFormat(scmContent string) string {
	var code strings.Builder
	lines := strings.Split(scmContent, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, ";;") && strings.Contains(trimmed, ";;") {
			// Extract comment content
			commentIdx := strings.Index(trimmed, ";;")
			if commentIdx >= 0 {
				content := strings.TrimSpace(trimmed[commentIdx+2:])
				if content != "" && !strings.HasPrefix(content, "Generated") &&
					!strings.HasPrefix(content, "Project:") &&
					content != "(define app-component" &&
					content != "(begin" &&
					content != "))" {
					code.WriteString(content + "\n")
				}
			}
		}
	}

	return code.String()
}
