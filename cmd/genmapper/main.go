package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// FieldMapping описывает как мапить поле
type FieldMapping struct {
	FieldName    string
	DTOType      string
	EntityName   string
	EntityType   string
	IsTimeField  bool
	IsVOField    bool
	IsStringCopy bool
	IsBoolCopy   bool
	IsIDField    bool
	VOPackage    string
}

// MapperConfig содержит конфигурацию для генерации маппера
type MapperConfig struct {
	PackageName string
	Imports     []string
	Mappers     []MapperDefinition
}

// MapperDefinition описывает один маппер DTO <-> Entity
type MapperDefinition struct {
	DTOName      string
	EntityName   string
	EntityPkg    string
	Fields       []FieldMapping
	HasUpdatedAt bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: genmapper <dto-file>...")
		os.Exit(1)
	}

	dtoFiles := os.Args[1:]
	fset := token.NewFileSet()

	config := MapperConfig{
		PackageName: "dto",
		Imports: []string{
			"time",
			"github.com/atumaikin/nexflow/internal/domain/entity",
			"github.com/atumaikin/nexflow/internal/domain/valueobject",
		},
	}

	// Parse all DTO files
	for _, dtoFile := range dtoFiles {
		node, err := parser.ParseFile(fset, dtoFile, nil, parser.ParseComments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing file %s: %v\n", dtoFile, err)
			continue
		}

		// Find all DTO struct types in this file
		mappers := findDTOStructs(node)
		config.Mappers = append(config.Mappers, mappers...)
	}

	// Generate code
	code := generateMapperCode(config)

	// Format code
	formatted, err := format.Source([]byte(code))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting code: %v\n", err)
		os.Exit(1)
	}

	// Write to mapper_gen.go
	outputFile := "mapper_gen.go"
	if err := os.WriteFile(outputFile, formatted, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", outputFile)
}

func findDTOStructs(node *ast.File) []MapperDefinition {
	var mappers []MapperDefinition

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			dtoName := typeSpec.Name.Name

			// Check if this is a DTO (ends with "DTO")
			if !strings.HasSuffix(dtoName, "DTO") {
				continue
			}

			// Determine entity name (remove "DTO" suffix)
			entityName := strings.TrimSuffix(dtoName, "DTO")

			// Extract fields
			fields := extractFields(structType)

			// Check if has UpdatedAt field
			hasUpdatedAt := false
			for _, f := range fields {
				if f.FieldName == "UpdatedAt" {
					hasUpdatedAt = true
					break
				}
			}

			mapper := MapperDefinition{
				DTOName:      dtoName,
				EntityName:   entityName,
				EntityPkg:    "entity",
				Fields:       fields,
				HasUpdatedAt: hasUpdatedAt,
			}

			mappers = append(mappers, mapper)
		}
	}

	return mappers
}

func extractFields(structType *ast.StructType) []FieldMapping {
	var fields []FieldMapping

	if structType.Fields == nil {
		return fields
	}

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		fieldName := field.Names[0].Name

		// Skip tags like json
		// Extract type
		typeName := getFieldType(field.Type)

		// Determine mapping strategy
		mapping := FieldMapping{
			FieldName:  fieldName,
			DTOType:    typeName,
			EntityName: fieldName,
		}

		// Check field type and determine mapping
		switch fieldName {
		case "ID":
			mapping.IsIDField = true
			mapping.EntityType = guessEntityIDType(fieldName)
		case "CreatedAt", "UpdatedAt":
			mapping.IsTimeField = true
			mapping.EntityType = "time.Time"
		case "Channel":
			mapping.IsVOField = true
			mapping.EntityType = "valueobject.Channel"
		case "Role":
			mapping.IsVOField = true
			mapping.EntityType = "valueobject.MessageRole"
		case "Status":
			mapping.IsVOField = true
			mapping.EntityType = "valueobject.TaskStatus"
		case "Version":
			mapping.IsVOField = true
			mapping.EntityType = "valueobject.Version"
		case "CronExpression":
			mapping.IsVOField = true
			mapping.EntityType = "valueobject.CronExpression"
		case "SessionID", "MessageID", "SkillID", "ScheduleID", "TaskID", "UserID":
			mapping.IsVOField = true
			mapping.EntityType = guessEntityIDType(fieldName)
		case "Enabled":
			mapping.IsBoolCopy = true
			mapping.EntityType = "bool"
		default:
			// Default: string copy
			mapping.IsStringCopy = true
			mapping.EntityType = "string"
		}

		fields = append(fields, mapping)
	}

	return fields
}

func getFieldType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", t.X, t.Sel)
	default:
		return "unknown"
	}
}

func guessEntityIDType(fieldName string) string {
	switch fieldName {
	case "ID":
		// Context-specific ID based on DTO name
		return "valueobject.ID" // Will be replaced in template
	case "UserID":
		return "valueobject.UserID"
	case "SessionID":
		return "valueobject.SessionID"
	case "MessageID":
		return "valueobject.MessageID"
	case "SkillID":
		return "valueobject.SkillID"
	case "ScheduleID":
		return "valueobject.ScheduleID"
	case "TaskID":
		return "valueobject.TaskID"
	default:
		return fmt.Sprintf("valueobject.%s", strings.TrimSuffix(fieldName, "ID"))
	}
}

func generateMapperCode(config MapperConfig) string {
	var buf bytes.Buffer

	// Package declaration
	fmt.Fprintf(&buf, "// Code generated by genmapper; DO NOT EDIT.\n\n")
	fmt.Fprintf(&buf, "package %s\n\n", config.PackageName)

	// Imports
	if len(config.Imports) > 0 {
		fmt.Fprint(&buf, "import (\n")
		for _, imp := range config.Imports {
			fmt.Fprintf(&buf, "\t%q\n", imp)
		}
		fmt.Fprint(&buf, ")\n\n")
	}

	// Generate mappers
	for _, mapper := range config.Mappers {
		generateToEntity(&buf, mapper)
		generateFromEntity(&buf, mapper)
	}

	return buf.String()
}

func generateToEntity(buf *bytes.Buffer, mapper MapperDefinition) {
	fmt.Fprintf(buf, "// ToEntity converts %s to entity.%s\n", mapper.DTOName, mapper.EntityName)
	fmt.Fprintf(buf, "func (dto *%s) ToEntity() *entity.%s {\n", mapper.DTOName, mapper.EntityName)

	// Handle CreatedAt/UpdatedAt parsing
	if mapper.HasUpdatedAt {
		fmt.Fprintf(buf, "\tcreatedAt, updatedAt := MustParseTimeFieldsWithUpdatedAt(dto.CreatedAt, dto.UpdatedAt)\n")
	} else {
		fmt.Fprintf(buf, "\tcreatedAt := MustParseTimeFields(dto.CreatedAt)\n")
	}

	// Return statement
	fmt.Fprintf(buf, "\treturn &entity.%s{\n", mapper.EntityName)

	// Fields
	for _, field := range mapper.Fields {
		if field.FieldName == "CreatedAt" || field.FieldName == "UpdatedAt" {
			continue // Already handled
		}

		fmt.Fprintf(buf, "\t\t%s: ", field.FieldName)

		switch {
		case field.IsIDField:
			// Special case for ID - need to determine correct type
			idType := getIDTypeForEntity(mapper.EntityName)
			fmt.Fprintf(buf, "%s(dto.%s),\n", idType, field.FieldName)
		case field.IsVOField:
			fmt.Fprintf(buf, "%s(dto.%s),\n", getVOConstructor(field.EntityType), field.FieldName)
		case field.IsStringCopy:
			fmt.Fprintf(buf, "dto.%s,\n", field.FieldName)
		case field.IsBoolCopy:
			fmt.Fprintf(buf, "dto.%s,\n", field.FieldName)
		default:
			fmt.Fprintf(buf, "dto.%s,\n", field.FieldName)
		}
	}

	// Time fields
	fmt.Fprintf(buf, "\t\tCreatedAt: createdAt,\n")
	if mapper.HasUpdatedAt {
		fmt.Fprintf(buf, "\t\tUpdatedAt: updatedAt,\n")
	}

	fmt.Fprintf(buf, "\t}\n")
	fmt.Fprintf(buf, "}\n\n")
}

func generateFromEntity(buf *bytes.Buffer, mapper MapperDefinition) {
	fmt.Fprintf(buf, "// FromEntity converts entity.%s to %s\n", mapper.EntityName, mapper.DTOName)
	fmt.Fprintf(buf, "func %sFromEntity(%s *entity.%s) *%s {\n",
		mapper.DTOName, strings.ToLower(mapper.EntityName[:1])+mapper.EntityName[1:], mapper.EntityName, mapper.DTOName)

	fmt.Fprintf(buf, "\treturn &%s{\n", mapper.DTOName)

	// Fields
	for _, field := range mapper.Fields {
		fmt.Fprintf(buf, "\t\t%s: ", field.FieldName)

		switch {
		case field.IsIDField || field.IsVOField:
			fmt.Fprintf(buf, "string(%s.%s),\n", strings.ToLower(mapper.EntityName[:1])+mapper.EntityName[1:], field.FieldName)
		case field.FieldName == "CreatedAt":
			fmt.Fprintf(buf, "%s.CreatedAt.Format(time.RFC3339),\n", strings.ToLower(mapper.EntityName[:1])+mapper.EntityName[1:])
		case field.FieldName == "UpdatedAt":
			fmt.Fprintf(buf, "%s.UpdatedAt.Format(time.RFC3339),\n", strings.ToLower(mapper.EntityName[:1])+mapper.EntityName[1:])
		case field.IsStringCopy:
			fmt.Fprintf(buf, "%s.%s,\n", strings.ToLower(mapper.EntityName[:1])+mapper.EntityName[1:], field.FieldName)
		case field.IsBoolCopy:
			fmt.Fprintf(buf, "%s.%s,\n", strings.ToLower(mapper.EntityName[:1])+mapper.EntityName[1:], field.FieldName)
		default:
			fmt.Fprintf(buf, "%s.%s,\n", strings.ToLower(mapper.EntityName[:1])+mapper.EntityName[1:], field.FieldName)
		}
	}

	fmt.Fprintf(buf, "\t}\n")
	fmt.Fprintf(buf, "}\n\n")
}

func getIDTypeForEntity(entityName string) string {
	switch entityName {
	case "User":
		return "valueobject.UserID"
	case "Session":
		return "valueobject.SessionID"
	case "Message":
		return "valueobject.MessageID"
	case "Task":
		return "valueobject.TaskID"
	case "Skill":
		return "valueobject.SkillID"
	case "Schedule":
		return "valueobject.ScheduleID"
	default:
		return fmt.Sprintf("valueobject.%sID", entityName)
	}
}

func getVOConstructor(entityType string) string {
	switch entityType {
	case "valueobject.Channel":
		return "valueobject.MustNewChannel"
	case "valueobject.MessageRole":
		return "valueobject.MustNewMessageRole"
	case "valueobject.TaskStatus":
		return "valueobject.MustNewTaskStatus"
	case "valueobject.Version":
		return "valueobject.MustNewVersion"
	case "valueobject.CronExpression":
		return "valueobject.MustNewCronExpression"
	case "valueobject.UserID":
		return "valueobject.MustNewUserID"
	case "valueobject.SessionID":
		return "valueobject.MustNewSessionID"
	case "valueobject.MessageID":
		return "valueobject.MustNewMessageID"
	case "valueobject.SkillID":
		return "valueobject.MustNewSkillID"
	case "valueobject.ScheduleID":
		return "valueobject.MustNewScheduleID"
	case "valueobject.TaskID":
		return "valueobject.MustNewTaskID"
	case "valueobject.Role":
		return "valueobject.MustNewMessageRole"
	case "valueobject.Status":
		return "valueobject.MustNewTaskStatus"
	default:
		parts := strings.Split(string(entityType), ".")
		if len(parts) > 1 {
			return fmt.Sprintf("valueobject.MustNew%s", parts[1])
		}
		return fmt.Sprintf("valueobject.MustNew%s", entityType)
	}
}
