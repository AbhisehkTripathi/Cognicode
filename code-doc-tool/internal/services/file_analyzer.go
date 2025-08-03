package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code-doc-tool/internal/models"
)

type FileAnalyzer struct {
	parsers map[string]Parser
}


type BaseParser struct{}

func NewFileAnalyzer() *FileAnalyzer {
	return &FileAnalyzer{
		parsers: map[string]Parser{
			"Node.js":     &NodeJSParser{},
			"Python":      &PythonParser{},
			"PHP/Laravel": &PHPParser{},
			"Go":          &GoParser{},
			"Java":        &JavaParser{},
			"Ruby/Rails":  &RubyParser{},
			".NET":        &DotNetParser{},
			"Generic":     &GenericParser{},
		},
	}
}

func (fa *FileAnalyzer) AnalyzeProject(projectPath string) (*models.Project, error) {
	project := &models.Project{
		Name:         filepath.Base(projectPath),
		Path:         projectPath,
		Dependencies: make(map[string][]models.Dependency),
	}

	// 1. Detect project type and set parser
	project.Type = fa.detectProjectType(projectPath)
	parser, exists := fa.parsers[project.Type]
	if !exists {
		parser = fa.parsers["Generic"]
	}

	// 2. Generate Project Overview
	project.Overview = parser.GenerateOverview(project.Name, projectPath)

	// 3. Extract Key Concepts
	project.KeyConcepts = parser.GetKeyConcepts()

	// 4. Generate Architecture Diagram (textual)
	project.Architecture = parser.GetArchitecture()

	// 5. Analyze Tech Stack
	project.TechStack = parser.AnalyzeTechStack(projectPath)

	// 6. Document Folder Structure
	project.FolderStructure = parser.AnalyzeFolderStructure(projectPath)

	// 7. Generate Setup Instructions
	project.SetupInstructions = parser.GetSetupInstructions()

	// 8. Extract API Endpoints
	project.APIEndpoints = parser.ExtractAPIEndpoints(projectPath)

	// 9. Document Parsers
	project.ParsersInfo = parser.GetParserDetails()

	// 10. Describe Data Flow
	project.DataFlow = parser.GetDataFlow()

	// 11. Detect External Services
	project.ExternalServices = parser.DetectExternalServices(projectPath)

	// 12. Generate Deployment Info
	project.DeploymentInfo = parser.GetDeploymentInfo(projectPath)

	// 13. Create Roadmap
	project.FutureRoadmap = parser.GetRoadmap()

	// 14. Document Common Issues
	project.CommonIssues = parser.GetCommonIssues()

	// 15. Generate Developer Notes
	project.DeveloperNotes = parser.GetDeveloperNotes(projectPath)

	// 16. Parse Dependencies
	if err := parser.ParseDependencies(projectPath, project); err != nil {
		return nil, fmt.Errorf("dependency parsing failed: %v", err)
	}

	// 17. Build Sample Output
	project.SampleOutput = parser.GenerateSampleOutput(projectPath)

	return project, nil
}

// Parser interface defines all required analysis methods
type Parser interface {
	GenerateOverview(name, projectPath string) string
	GetKeyConcepts() []string
	GetArchitecture() string
	AnalyzeTechStack(projectPath string) []string
	AnalyzeFolderStructure(projectPath string) map[string]string
	GetSetupInstructions() []string
	ExtractAPIEndpoints(projectPath string) []models.APIEndpoint
	GetParserDetails() map[string]string
	GetDataFlow() string
	DetectExternalServices(projectPath string) []string
	GetDeploymentInfo(projectPath string) []string
	GetRoadmap() []string
	GetCommonIssues() []string
	GetDeveloperNotes(projectPath string) []string
	ParseDependencies(projectPath string, project *models.Project) error
	GenerateSampleOutput(projectPath string) string
}

// Base parser with common functionality
type BaseParser struct{}

func (p *BaseParser) readREADME(projectPath string) string {
	readmeFiles := []string{"README.md", "readme.md", "README.txt", "readme.txt", "README"}

	for _, filename := range readmeFiles {
		filePath := filepath.Join(projectPath, filename)
		if content, err := ioutil.ReadFile(filePath); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") && len(line) > 10 {
					if len(line) > 150 {
						return line[:150] + "..."
					}
					return line
				}
			}
		}
	}
	return ""
}

func (p *BaseParser) shouldSkip(path string, info os.FileInfo) bool {
	skipDirs := []string{"node_modules", "vendor", ".git", "dist", "build", "__pycache__", ".idea", ".vscode"}
	skipFiles := []string{".DS_Store", "Thumbs.db"}

	name := info.Name()

	if info.IsDir() {
		for _, skipDir := range skipDirs {
			if strings.EqualFold(name, skipDir) {
				return true
			}
		}
	} else {
		for _, skipFile := range skipFiles {
			if strings.EqualFold(name, skipFile) {
				return true
			}
		}
	}

	return strings.HasPrefix(name, ".")
}

// Example language-specific parser implementation (NodeJS)
type NodeJSParser struct {
	BaseParser
}

func (p *NodeJSParser) GenerateOverview(name, projectPath string) string {
	readme := p.readREADME(projectPath)
	if readme != "" {
		return fmt.Sprintf("%s: %s", name, readme)
	}
	return fmt.Sprintf("%s: Node.js application with modern JavaScript/TypeScript stack", name)
}

func (p *NodeJSParser) GetKeyConcepts() []string {
	return []string{
		"AST-based parsing using Esprima for JavaScript",
		"Package.json analysis for dependencies",
		"Route extraction from Express/Fastify/Koa applications",
		"ES Module and CommonJS support",
	}
}

func (p *NodeJSParser) GetArchitecture() string {
	return `Client → Router → Middleware → Controller → Service → Model → Database
│
├── API Layer (Express/Fastify/Koa)
├── Business Logic (Services)
└── Data Access (Models/Repositories)`
}

// Implement all other required methods for NodeJSParser...

// Similar implementations for other language parsers...
// PythonParser, PHPParser, GoParser, JavaParser, RubyParser, DotNetParser, GenericParser

func (fa *FileAnalyzer) detectProjectType(projectPath string) string {
	// Enhanced project type detection
	indicators := map[string][]string{
		"Node.js":     {"package.json", "node_modules"},
		"Python":      {"requirements.txt", "setup.py", "Pipfile"},
		"PHP/Laravel": {"composer.json", "artisan", "vendor"},
		"Go":          {"go.mod", "go.sum"},
		"Java":        {"pom.xml", "build.gradle", "src/main/java"},
		"Ruby/Rails":  {"Gemfile", "config.ru", "app/models"},
		".NET":        {"*.csproj", "*.sln", "Program.cs"},
	}

	for projectType, files := range indicators {
		for _, file := range files {
			if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
				return projectType
			}
		}
	}

	// Fallback to file extension analysis
	extensions := make(map[string]int)
	filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".js", ".ts", ".jsx", ".tsx":
			extensions["Node.js"]++
		case ".py":
			extensions["Python"]++
		case ".php":
			extensions["PHP/Laravel"]++
		case ".go":
			extensions["Go"]++
		case ".java":
			extensions["Java"]++
		case ".rb":
			extensions["Ruby/Rails"]++
		case ".cs":
			extensions[".NET"]++
		}
		return nil
	})

	if len(extensions) > 0 {
		maxCount := 0
		detectedType := "Generic"
		for lang, count := range extensions {
			if count > maxCount {
				maxCount = count
				detectedType = lang
			}
		}
	return detectedType
}

func (p *BaseParser) extractPathFromContent(content string) string {
	quoteChars := []rune{'\'', '"', '`'}
	
	for _, quote := range quoteChars {
		quoteStr := string(quote)
		startIdx := strings.Index(content, quoteStr)
		
		for startIdx != -1 {
			// Find closing quote
			endIdx := -1
			remaining := content[startIdx+1:]
			escaped := false
			
			for i, ch := range remaining {
				if escaped {
					escaped = false
					continue
				}
				if ch == '\\' {
					escaped = true
					continue
				}
				if ch == quote {
					endIdx = i
					break
				}
			}
			
			if endIdx == -1 {
				break
			}
			
			path := remaining[:endIdx]
			path = strings.SplitN(path, "?", 2)[0]
			path = strings.SplitN(path, "(", 2)[0]
			path = strings.TrimSpace(path)
			
			if path != "" {
				return path
			}
			
			// Find next potential quoted section
			nextStart := strings.Index(remaining[endIdx+1:], quoteStr)
			if nextStart == -1 {
				break
			}
			startIdx = startIdx + 1 + endIdx + 1 + nextStart
		}
	}
	return ""
}



// detectExternalServices checks for common service indicators in config files
func (p *BaseParser) detectExternalServices(projectPath string, configFiles map[string][]string) []string {
	var services []string
	serviceMap := map[string]string{
		"postgres":    "PostgreSQL",
		"mysql":       "MySQL",
		"mongodb":     "MongoDB",
		"redis":       "Redis",
		"aws":         "AWS",
		"stripe":      "Stripe",
		"sendgrid":    "SendGrid",
		"twilio":      "Twilio",
		"firebase":    "Firebase",
		"elastic":     "Elasticsearch",
		"kafka":       "Kafka",
		"rabbitmq":    "RabbitMQ",
		"googlecloud": "Google Cloud",
		"azure":       "Microsoft Azure",
	}

	for configFile, keywords := range configFiles {
		filePath := filepath.Join(projectPath, configFile)
		if content, err := ioutil.ReadFile(filePath); err == nil {
			contentStr := strings.ToLower(string(content))
			for keyword, service := range serviceMap {
				if strings.Contains(contentStr, keyword) {
					services = append(services, service)
				}
			}
			for _, keyword := range keywords {
				if strings.Contains(contentStr, strings.ToLower(keyword)) {
					services = append(services, keyword)
				}
			}
		}
	}

	return p.removeDuplicates(services)
}

func (p *BaseParser) removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}

// analyzeFolderStructure generates a map of directory purposes
func (p *BaseParser) analyzeFolderStructure(projectPath string) map[string]string {
	structure := make(map[string]string)
	dirDescriptions := map[string]string{
		"src":          "Main source code files",
		"lib":          "Library and utility code",
		"test":         "Test files",
		"config":       "Configuration files",
		"public":       "Publicly accessible assets",
		"assets":       "Static assets (images, styles)",
		"migrations":   "Database migration files",
		"models":       "Data models and schemas",
		"controllers":  "Application controllers",
		"middleware":   "Middleware functions",
		"routes":       "Route definitions",
		"services":     "Business logic services",
		"utils":        "Utility functions",
		"views":        "Templates and views",
		"dist":         "Compiled/transpiled output",
		"build":        "Build artifacts",
		"docs":         "Documentation files",
		"scripts":      "Utility scripts",
		"seeds":        "Database seed data",
		"fixtures":     "Test fixtures",
		"locales":      "Localization files",
		"logs":         "Application logs",
		"tmp":          "Temporary files",
	}

	filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(projectPath, path)
		if relPath == "." {
			return nil
		}

		dirName := info.Name()
		if desc, exists := dirDescriptions[strings.ToLower(dirName)]; exists {
			structure[relPath+"/"] = desc
		} else {
			structure[relPath+"/"] = "Project-specific directory"
		}

		return nil
	})

	return structure
}

// NodeJSParser implements Parser for Node.js projects
type NodeJSParser struct {
	BaseParser
}

func (p *NodeJSParser) GenerateOverview(name, projectPath string) string {
	readme := p.readREADME(projectPath)
	if readme != "" {
		return fmt.Sprintf("%s: %s", name, readme)
	}
	return fmt.Sprintf("%s: Node.js application with modern JavaScript/TypeScript stack", name)
}

func (p *NodeJSParser) GetKeyConcepts() []string {
	return []string{
		"AST-based parsing using Esprima for JavaScript",
		"Package.json analysis for dependencies",
		"Route extraction from Express/Fastify/Koa applications",
		"ES Module and CommonJS support",
		"Middleware architecture",
		"Environment configuration",
	}
}

func (p *NodeJSParser) GetArchitecture() string {
	return `Client → Router → Middleware → Controller → Service → Model → Database
│
├── API Layer (Express/Fastify/Koa)
├── Business Logic (Services)
└── Data Access (Models/Repositories)`
}

func (p *NodeJSParser) AnalyzeTechStack(projectPath string) []string {
	techStack := []string{"Node.js"}

	// Check for framework files
	frameworkFiles := map[string]string{
		"package.json":       "Check dependencies",
		"next.config.js":     "Next.js",
		"nuxt.config.js":     "Nuxt.js",
		"vue.config.js":      "Vue.js",
		"angular.json":       "Angular",
		"svelte.config.js":   "Svelte",
		"remix.config.js":    "Remix",
		"nest-cli.json":      "NestJS",
		"serverless.yml":     "Serverless Framework",
		"webpack.config.js":  "Webpack",
		"rollup.config.js":   "Rollup",
		"vite.config.js":     "Vite",
		"jest.config.js":     "Jest",
		"mocha.opts":         "Mocha",
	}

	for file, tech := range frameworkFiles {
		if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
			techStack = append(techStack, tech)
		}
	}

	// Check package.json for frameworks
	packagePath := filepath.Join(projectPath, "package.json")
	if content, err := ioutil.ReadFile(packagePath); err == nil {
		var pkg struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if json.Unmarshal(content, &pkg) == nil {
			for dep := range pkg.Dependencies {
				switch dep {
				case "express":
					techStack = append(techStack, "Express")
				case "koa":
					techStack = append(techStack, "Koa")
				case "fastify":
					techStack = append(techStack, "Fastify")
				case "react":
					techStack = append(techStack, "React")
				case "typescript":
					techStack = append(techStack, "TypeScript")
				}
			}
		}
	}

	return p.removeDuplicates(techStack)
}

func (p *NodeJSParser) AnalyzeFolderStructure(projectPath string) map[string]string {
	baseStructure := p.BaseParser.analyzeFolderStructure(projectPath)
	
	// Add Node.js specific folder descriptions
	nodeSpecific := map[string]string{
		"node_modules/": "Node.js dependencies",
		"types/":        "TypeScript type definitions",
		"__tests__/":    "Jest test files",
		".next/":        "Next.js build output",
		".nuxt/":        "Nuxt.js build output",
	}
	
	for path, desc := range nodeSpecific {
		if _, exists := baseStructure[path]; !exists {
			baseStructure[path] = desc
		}
	}
	
	return baseStructure
}

func (p *NodeJSParser) GetSetupInstructions() []string {
	return []string{
		"Install Node.js (version specified in .nvmrc or package.json if available)",
		"Run `npm install` to install dependencies",
		"Create .env file based on .env.example",
		"Run `npm run dev` to start development server",
		"Run `npm test` to execute tests",
	}
}

func (p *NodeJSParser) ExtractAPIEndpoints(projectPath string) []models.APIEndpoint {
	var endpoints []models.Endpoint

	// Check common route files
	routeFiles := []string{
		"src/routes/*.js",
		"src/routes/*.ts",
		"routes/*.js",
		"routes/*.ts",
		"app.js",
		"app.ts",
		"server.js",
		"server.ts",
	}

	for _, pattern := range routeFiles {
		matches, _ := filepath.Glob(filepath.Join(projectPath, pattern))
		for _, file := range matches {
			content, _ := ioutil.ReadFile(file)
			endpoints = append(endpoints, p.extractExpressRoutes(string(content))...)
			endpoints = append(endpoints, p.extractFastifyRoutes(string(content))...)
		}
	}

	if len(endpoints) == 0 {
		return []models.APIEndpoint{
			{
				Method:      "GET",
				Path:        "/",
				Description: "Default home endpoint",
				CurlExample: "curl http://localhost:3000",
			},
		}
	}

	return endpoints
}

func (p *NodeJSParser) extractExpressRoutes(content string) []models.APIEndpoint {
	var endpoints []models.APIEndpoint
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ".get(") || strings.Contains(line, ".post(") ||
			strings.Contains(line, ".put(") || strings.Contains(line, ".delete(") {
			
			var method string
			switch {
			case strings.Contains(line, ".get("):
				method = "GET"
			case strings.Contains(line, ".post("):
				method = "POST"
			case strings.Contains(line, ".put("):
				method = "PUT"
			case strings.Contains(line, ".delete("):
				method = "DELETE"
			default:
				continue
			}

			path := p.extractPathFromContent(line)
			if path == "" {
				continue
			}

			endpoints = append(endpoints, models.APIEndpoint{
				Method:      method,
				Path:        path,
				Description: fmt.Sprintf("Express %s endpoint", method),
				CurlExample: fmt.Sprintf("curl -X %s http://localhost:3000%s", method, path),
			})
		}
	}

	return endpoints
}

func (p *NodeJSParser) GetParserDetails() map[string]string {
	return map[string]string{
		"JavaScript/TypeScript": "AST parsing using Esprima, extracts routes, functions, classes",
		"package.json":          "JSON parsing for dependencies, scripts, and metadata",
		"Configuration Files":    "Analyzes next.config.js, nuxt.config.js, etc.",
	}
}

func (p *NodeJSParser) GetDataFlow() string {
	return `1. HTTP Request → 2. Server (Express/Fastify) → 3. Middleware → 
4. Route Handler → 5. Service Layer → 6. Data Access → 7. Database → 
8. Response Formatter → 9. HTTP Response`
}

func (p *NodeJSParser) DetectExternalServices(projectPath string) []string {
	configFiles := map[string][]string{
		"package.json":    {"mongoose", "sequelize", "typeorm", "redis", "pg", "mysql2"},
		".env":           {"DATABASE_URL", "REDIS_URL", "AWS_", "STRIPE_"},
		"docker-compose.yml": {"postgres", "redis", "mongo"},
	}

	return p.detectExternalServices(projectPath, configFiles)
}

func (p *NodeJSParser) GetDeploymentInfo(projectPath string) []string {
	deployment := []string{
		"Node.js applications can be deployed using:",
		"- PM2 process manager for production",
		"- Docker containers",
		"- Serverless platforms (AWS Lambda, Vercel, Netlify)",
		"- Traditional VPS with Nginx reverse proxy",
	}

	// Check for specific deployment configs
	if _, err := os.Stat(filepath.Join(projectPath, "Dockerfile")); err == nil {
		deployment = append(deployment, "\nDetected Dockerfile - can be built with `docker build .`")
	}

	if _, err := os.Stat(filepath.Join(projectPath, "serverless.yml")); err == nil {
		deployment = append(deployment, "\nDetected serverless.yml - deploy with `serverless deploy`")
	}

	return deployment
}

func (p *NodeJSParser) GetRoadmap() []string {
	return []string{
		"Add TypeScript type analysis",
		"Improve Express route detection",
		"Add GraphQL API documentation",
		"Support NestJS architecture analysis",
		"Add dependency vulnerability scanning",
	}
}

func (p *NodeJSParser) GetCommonIssues() []string {
	return []string{
		"Missing .env variables - check .env.example",
		"Node version mismatch - use nvm or check package.json engines",
		"Port already in use - change PORT in .env",
		"Module not found - run npm install",
		"ESLint/prettier conflicts - check .eslintrc and .prettierrc",
	}
}

func (p *NodeJSParser) GetDeveloperNotes(projectPath string) []string {
	notes := []string{
		"Use npm-check-updates to update dependencies: ncu -u",
		"Debug with Chrome DevTools: node --inspect server.js",
		"Generate dependency graph: npm install -g npm-license && npm-license",
	}

	if _, err := os.Stat(filepath.Join(projectPath, "tsconfig.json")); err == nil {
		notes = append(notes, "TypeScript project - compile with tsc or use ts-node for development")
	}

	return notes
}

func (p *NodeJSParser) ParseDependencies(projectPath string, project *models.Project) error {
	packagePath := filepath.Join(projectPath, "package.json")
	data, err := ioutil.ReadFile(packagePath)
	if err != nil {
		return err
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
		PeerDependencies map[string]string `json:"peerDependencies"`
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}

	for name, version := range pkg.Dependencies {
		project.Dependencies["production"] = append(project.Dependencies["production"],
			models.Dependency{Name: name, Version: version, Type: "production"})
	}

	for name, version := range pkg.DevDependencies {
		project.Dependencies["development"] = append(project.Dependencies["development"],
			models.Dependency{Name: name, Version: version, Type: "development"})
	}

	for name, version := range pkg.PeerDependencies {
		project.Dependencies["peer"] = append(project.Dependencies["peer"],
			models.Dependency{Name: name, Version: version, Type: "peer"})
	}

	return nil
}

func (p *NodeJSParser) GenerateSampleOutput(projectPath string) string {
	sample := map[string]interface{}{
		"project": map[string]interface{}{
			"name":        "example-node-app",
			"type":        "Node.js",
			"description": "Sample Node.js API application",
		},
		"endpoints": []map[string]interface{}{
			{
				"method":      "GET",
				"path":        "/api/users",
				"description": "Get list of users",
				"parameters": []map[string]interface{}{
					{
						"name":        "limit",
						"type":        "integer",
						"required":    false,
						"description": "Number of items to return",
					},
				},
			},
		},
		"dependencies": map[string]interface{}{
			"express":    "^4.17.1",
			"mongoose":   "^5.12.0",
			"typescript": "^4.2.0",
		},
	}

	jsonData, _ := json.MarshalIndent(sample, "", "  ")
	return string(jsonData)
}
