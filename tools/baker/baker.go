// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// Output represents the top-level ordered JSON object and metadata for file generation.
type Output struct {
	PackageName   string        `json:"-"`
	Preamble      string        `json:"-"`
	Imports       []*ImportInfo `json:"-"`
	Entries       []*FileEntry
	Uncategorized *FileGroup
}

// ImportInfo captures an import specification from the source file.
type ImportInfo struct {
	Name string
	Path string
	Code string
}

// Ident returns the identifier used to reference this import in Go source code.
func (imp *ImportInfo) Ident() string {
	if imp.Name != "" {
		return imp.Name
	}
	raw := strings.Trim(imp.Path, `"`)
	parts := strings.Split(raw, "/")
	last := parts[len(parts)-1]
	if len(parts) > 1 && strings.HasPrefix(last, "v") && isAllDigits(last[1:]) {
		last = parts[len(parts)-2]
	}
	return last
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// HasLeftovers returns true if there are any uncategorized declarations.
func (o *Output) HasLeftovers() bool {
	if o.Uncategorized == nil {
		return false
	}
	return len(o.Uncategorized.Constants) > 0 ||
		len(o.Uncategorized.Variables) > 0 ||
		len(o.Uncategorized.Structures) > 0 ||
		len(o.Uncategorized.Interfaces) > 0 ||
		len(o.Uncategorized.Functions) > 0
}

// FileEntry represents a target file and its contents.
type FileEntry struct {
	FileName string
	Group    *FileGroup
}

// FileGroup contains the declarations assigned to a potential file.
type FileGroup struct {
	Constants  []*Constant  `json:"constants,omitempty"`
	Variables  []*Variable  `json:"variables,omitempty"`
	Structures []*Structure `json:"structures"`
	Interfaces []*Interface `json:"interfaces"`
	Functions  []*Function  `json:"functions"`
}

// Constant represents a Go package-level constant declaration.
type Constant struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Variable represents a Go package-level variable declaration.
type Variable struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Structure represents a Go struct definition and its receiver methods.
type Structure struct {
	Name    string    `json:"name"`
	Code    string    `json:"code"`
	Methods []*Method `json:"methods"`
}

// Method represents a method associated with a structure.
type Method struct {
	Name     string `json:"name"`
	Receiver string `json:"receiver,omitempty"`
	Code     string `json:"code"`
}

// Interface represents a Go interface definition.
type Interface struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Function represents a standalone function.
type Function struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// MarshalJSON guarantees that file entries are emitted in order and that
// "<uncategorized>" is serialized as the very last key.
func (o *Output) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')

	first := true
	for _, entry := range o.Entries {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		keyBytes, err := json.Marshal(entry.FileName)
		if err != nil {
			return nil, err
		}
		valBytes, err := json.Marshal(entry.Group)
		if err != nil {
			return nil, err
		}

		buf.Write(keyBytes)
		buf.WriteByte(':')
		buf.Write(valBytes)
	}

	// Always emit "<uncategorized>" at the end
	if !first {
		buf.WriteByte(',')
	}
	uncatGroup := o.Uncategorized
	if uncatGroup == nil {
		uncatGroup = &FileGroup{
			Constants:  make([]*Constant, 0),
			Variables:  make([]*Variable, 0),
			Structures: make([]*Structure, 0),
			Interfaces: make([]*Interface, 0),
			Functions:  make([]*Function, 0),
		}
	}
	uncatVal, err := json.Marshal(uncatGroup)
	if err != nil {
		return nil, err
	}

	buf.WriteString(`"<uncategorized>":`)
	buf.Write(uncatVal)
	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func main() {
	var filePath string
	var outDir string
	var compact bool

	flag.StringVar(&filePath, "file", "", "Path to Go source file")
	flag.StringVar(
		&outDir,
		"out",
		"",
		"Output directory for generated Go files",
	)
	flag.BoolVar(
		&compact,
		"compact",
		false,
		"Output compact JSON instead of indented",
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			"Usage: baker [flags] <path-to-go-file>\n\nFlags:\n",
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	if filePath == "" && flag.NArg() > 0 {
		filePath = flag.Arg(0)
	}

	if filePath == "" {
		flag.Usage()
		os.Exit(2)
	}

	output, err := ParseFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", filePath, err)
		os.Exit(2)
	}

	if outDir != "" {
		if err := GenerateFiles(output, outDir, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating files: %v\n", err)
			os.Exit(2)
		}
	} else {
		data, err := EncodeJSON(output, compact)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
			os.Exit(2)
		}
		fmt.Println(string(data))
	}

	if output.HasLeftovers() {
		os.Exit(1)
	}
	os.Exit(0)
}

// GenerateFiles writes out split Go files to outDir and prints each filename to stdout.
func GenerateFiles(output *Output, outDir string, stdout io.Writer) error {
	cleanDir := filepath.Clean(outDir)
	if err := os.MkdirAll(cleanDir, 0750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	for _, entry := range output.Entries {
		filePath := filepath.Join(cleanDir, entry.FileName)
		content := buildFileContent(output, entry.Group)

		formatted, err := format.Source([]byte(content))
		if err != nil {
			formatted = []byte(content)
		}

		if err := os.WriteFile(
			filepath.Clean(filePath),
			formatted,
			0600,
		); err != nil {
			return fmt.Errorf("write %s: %w", entry.FileName, err)
		}

		_, _ = fmt.Fprintln(stdout, entry.FileName)
	}

	return nil
}

// buildFileContent constructs complete Go source code for a FileGroup.
func buildFileContent(output *Output, group *FileGroup) string {
	var buf strings.Builder

	if output.Preamble != "" {
		buf.WriteString(output.Preamble)
		buf.WriteString("\n\n")
	}

	buf.WriteString("package ")
	buf.WriteString(output.PackageName)
	buf.WriteString("\n\n")

	var bodyParts []string
	for _, c := range group.Constants {
		bodyParts = append(bodyParts, c.Code)
	}
	for _, v := range group.Variables {
		bodyParts = append(bodyParts, v.Code)
	}
	for _, iface := range group.Interfaces {
		bodyParts = append(bodyParts, iface.Code)
	}
	for _, st := range group.Structures {
		if st.Code != "" {
			bodyParts = append(bodyParts, st.Code)
		}
		for _, m := range st.Methods {
			bodyParts = append(bodyParts, m.Code)
		}
	}
	for _, fn := range group.Functions {
		bodyParts = append(bodyParts, fn.Code)
	}
	combinedBody := strings.Join(bodyParts, "\n\n")

	var neededImports []string
	for _, imp := range output.Imports {
		ident := imp.Ident()
		if ident == "_" || isIdentReferenced(combinedBody, ident) {
			neededImports = append(neededImports, imp.Code)
		}
	}

	if len(neededImports) > 0 {
		buf.WriteString("import (\n")
		for _, impCode := range neededImports {
			buf.WriteString("\t")
			buf.WriteString(impCode)
			buf.WriteString("\n")
		}
		buf.WriteString(")\n\n")
	}

	buf.WriteString(combinedBody)
	buf.WriteString("\n")

	return buf.String()
}

func isIdentReferenced(code, ident string) bool {
	if ident == "" {
		return false
	}
	pattern := `\b` + regexp.QuoteMeta(ident) + `\b`
	matched, _ := regexp.MatchString(pattern, code)
	return matched
}

// EncodeJSON serializes the Output struct without HTML escaping, ensuring
// deterministic key ordering and literal "<uncategorized>".
func EncodeJSON(output *Output, compact bool) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if !compact {
		enc.SetIndent("", "  ")
	}
	if err := enc.Encode(output); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

// ParseFile reads and parses a Go source file, grouping declarations by potential file names.
func ParseFile(filePath string) (*Output, error) {
	src, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	fset := token.NewFileSet()
	fileNode, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	base := filepath.Base(filePath)
	inputBase := strings.TrimSuffix(base, filepath.Ext(base))

	var preamble string
	if fileNode.Doc != nil {
		pStart := fset.Position(fileNode.Doc.Pos()).Offset
		pEnd := fset.Position(fileNode.Doc.End()).Offset
		if pStart >= 0 && pEnd <= len(src) && pStart <= pEnd {
			preamble = strings.TrimSpace(string(src[pStart:pEnd]))
		}
	} else if len(fileNode.Comments) > 0 {
		firstComment := fileNode.Comments[0]
		if firstComment.Pos() < fileNode.Package {
			pStart := fset.Position(firstComment.Pos()).Offset
			pEnd := fset.Position(firstComment.End()).Offset
			if pStart >= 0 && pEnd <= len(src) && pStart <= pEnd {
				preamble = strings.TrimSpace(string(src[pStart:pEnd]))
			}
		}
	}

	output := &Output{
		PackageName: fileNode.Name.Name,
		Preamble:    preamble,
		Imports:     make([]*ImportInfo, 0),
		Entries:     make([]*FileEntry, 0),
		Uncategorized: &FileGroup{
			Constants:  make([]*Constant, 0),
			Variables:  make([]*Variable, 0),
			Structures: make([]*Structure, 0),
			Interfaces: make([]*Interface, 0),
			Functions:  make([]*Function, 0),
		},
	}

	// Collect imports
	for _, decl := range fileNode.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}
		for _, spec := range genDecl.Specs {
			impSpec, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			imp := &ImportInfo{
				Path: impSpec.Path.Value,
				Code: extractSource(src, fset, impSpec, impSpec.Doc),
			}
			if impSpec.Name != nil {
				imp.Name = impSpec.Name.Name
			}
			output.Imports = append(output.Imports, imp)
		}
	}

	// Maps for struct lookups
	structMap := make(map[string]*Structure)
	structFileMap := make(map[string]*FileGroup)

	// Collected interfaces to match later
	var interfaces []*Interface

	// First pass: Discover all struct and interface definitions.
	for _, decl := range fileNode.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			switch typeSpec.Type.(type) {
			case *ast.StructType:
				code := extractTypeCode(src, fset, genDecl, typeSpec)
				st := &Structure{
					Name:    typeSpec.Name.Name,
					Code:    code,
					Methods: make([]*Method, 0),
				}
				structMap[typeSpec.Name.Name] = st

				fg := &FileGroup{
					Constants:  make([]*Constant, 0),
					Variables:  make([]*Variable, 0),
					Structures: []*Structure{st},
					Interfaces: make([]*Interface, 0),
					Functions:  make([]*Function, 0),
				}
				structFileMap[typeSpec.Name.Name] = fg

				targetFile := targetFileName(inputBase, typeSpec.Name.Name)
				output.Entries = append(output.Entries, &FileEntry{
					FileName: targetFile,
					Group:    fg,
				})

			case *ast.InterfaceType:
				code := extractTypeCode(src, fset, genDecl, typeSpec)
				interfaces = append(interfaces, &Interface{
					Name: typeSpec.Name.Name,
					Code: code,
				})
			}
		}
	}

	// Match interfaces to structs
	for _, iface := range interfaces {
		matched := false
		if strings.HasPrefix(iface.Name, "I") {
			candidate := iface.Name[1:]
			if fg, ok := structFileMap[candidate]; ok {
				fg.Interfaces = append(fg.Interfaces, iface)
				matched = true
			}
		}

		if !matched {
			trimmed := strings.TrimPrefix(iface.Name, "I")
			for sName, fg := range structFileMap {
				if strings.EqualFold(trimmed, sName) {
					fg.Interfaces = append(fg.Interfaces, iface)
					matched = true
					break
				}
			}
		}

		if !matched {
			output.Uncategorized.Interfaces = append(
				output.Uncategorized.Interfaces,
				iface,
			)
		}
	}

	// Second pass: Process constant and variable declarations.
	for _, decl := range fileNode.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.CONST && genDecl.Tok != token.VAR {
			continue
		}

		allBlank := true
		var firstName string
		for _, spec := range genDecl.Specs {
			valSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, nameIdent := range valSpec.Names {
				if nameIdent.Name != "_" {
					allBlank = false
					if firstName == "" {
						firstName = nameIdent.Name
					}
				}
			}
		}
		if allBlank {
			continue
		}

		code := extractSource(src, fset, genDecl, genDecl.Doc)
		matchedStruct := matchDeclToStruct(firstName, structMap)

		if genDecl.Tok == token.CONST {
			c := &Constant{
				Name: firstName,
				Code: code,
			}
			if matchedStruct != "" {
				if fg, ok := structFileMap[matchedStruct]; ok {
					fg.Constants = append(fg.Constants, c)
					continue
				}
			}
			output.Uncategorized.Constants = append(
				output.Uncategorized.Constants,
				c,
			)
		} else {
			v := &Variable{
				Name: firstName,
				Code: code,
			}
			if matchedStruct != "" {
				if fg, ok := structFileMap[matchedStruct]; ok {
					fg.Variables = append(fg.Variables, v)
					continue
				}
			}
			output.Uncategorized.Variables = append(
				output.Uncategorized.Variables,
				v,
			)
		}
	}

	// Third pass: Process functions and methods.
	for _, decl := range fileNode.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		code := extractSource(src, fset, funcDecl, funcDecl.Doc)

		// Method
		if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
			recvTypeName := getReceiverTypeName(funcDecl.Recv)
			recvStr := extractReceiverString(src, fset, funcDecl.Recv)

			method := &Method{
				Name:     funcDecl.Name.Name,
				Receiver: recvStr,
				Code:     code,
			}

			if st, found := structMap[recvTypeName]; found {
				st.Methods = append(st.Methods, method)
			} else {
				orphanStruct := &Structure{
					Name:    recvTypeName,
					Code:    "",
					Methods: []*Method{method},
				}
				output.Uncategorized.Structures = append(
					output.Uncategorized.Structures,
					orphanStruct,
				)
			}
			continue
		}

		// Standalone Function: Match to a struct's file
		matchedStruct := matchFunctionToStruct(funcDecl, structMap)
		fn := &Function{
			Name: funcDecl.Name.Name,
			Code: code,
		}

		if matchedStruct != "" {
			if fg, ok := structFileMap[matchedStruct]; ok {
				fg.Functions = append(fg.Functions, fn)
				continue
			}
		}

		output.Uncategorized.Functions = append(
			output.Uncategorized.Functions,
			fn,
		)
	}

	return output, nil
}

// matchDeclToStruct finds a matching struct name for a constant or variable declaration.
func matchDeclToStruct(name string, structMap map[string]*Structure) string {
	lowerName := strings.ToLower(name)
	var bestMatch string
	for sName := range structMap {
		lowerStruct := strings.ToLower(sName)
		if strings.HasPrefix(lowerName, lowerStruct) {
			if len(sName) > len(bestMatch) {
				bestMatch = sName
			}
		}
	}
	return bestMatch
}

// matchFunctionToStruct finds a matching struct name for a standalone function.
func matchFunctionToStruct(
	fn *ast.FuncDecl,
	structMap map[string]*Structure,
) string {
	if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
		retBase := getBaseTypeName(fn.Type.Results.List[0].Type)
		if _, ok := structMap[retBase]; ok {
			return retBase
		}
	}

	if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
		paramBase := getBaseTypeName(fn.Type.Params.List[0].Type)
		if _, ok := structMap[paramBase]; ok {
			return paramBase
		}
	}

	fnName := fn.Name.Name
	for sName := range structMap {
		if strings.EqualFold(fnName, "New"+sName) ||
			strings.EqualFold(fnName, "NewEmpty"+sName) ||
			strings.EqualFold(fnName, "Init"+sName) ||
			strings.EqualFold(fnName, "InitEmpty"+sName) ||
			strings.EqualFold(fnName, sName+"Init") {
			return sName
		}
	}

	return ""
}

// getBaseTypeName extracts the base type identifier from an expression.
func getBaseTypeName(expr ast.Expr) string {
	for expr != nil {
		switch v := expr.(type) {
		case *ast.StarExpr:
			expr = v.X
		case *ast.IndexExpr:
			expr = v.X
		case *ast.IndexListExpr:
			expr = v.X
		case *ast.ParenExpr:
			expr = v.X
		case *ast.Ident:
			return v.Name
		default:
			return ""
		}
	}
	return ""
}

// toSnakeCase converts a CamelCase or PascalCase identifier to snake_case.
func toSnakeCase(s string) string {
	var buf strings.Builder
	runes := []rune(s)
	for i := range runes {
		r := runes[i]
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := runes[i-1]
				if unicode.IsLower(prev) || unicode.IsDigit(prev) {
					if buf.Len() > 0 && !strings.HasSuffix(buf.String(), "_") {
						buf.WriteRune('_')
					}
				} else if unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					if buf.Len() > 0 && !strings.HasSuffix(buf.String(), "_") {
						buf.WriteRune('_')
					}
				}
			}
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(unicode.ToLower(r))
		}
	}

	res := buf.String()
	for strings.Contains(res, "__") {
		res = strings.ReplaceAll(res, "__", "_")
	}
	return strings.Trim(res, "_")
}

// targetFileName constructs <input_basename>_<snake_case_name>.go.
func targetFileName(inputBase, structName string) string {
	snake := toSnakeCase(structName)
	return fmt.Sprintf("%s_%s.go", inputBase, snake)
}

// extractTypeCode extracts the code for a type specification, preserving comments and formatting.
func extractTypeCode(
	src []byte,
	fset *token.FileSet,
	genDecl *ast.GenDecl,
	typeSpec *ast.TypeSpec,
) string {
	if len(genDecl.Specs) == 1 {
		return extractSource(src, fset, genDecl, genDecl.Doc)
	}

	doc := typeSpec.Doc
	if doc == nil && genDecl.Doc != nil {
		doc = genDecl.Doc
	}
	specCode := extractSource(src, fset, typeSpec, doc)
	if !strings.HasPrefix(strings.TrimSpace(specCode), "type ") {
		specCode = "type " + specCode
	}
	return specCode
}

// extractSource slices the exact bytes from src for a given AST node, including its doc comment.
func extractSource(
	src []byte,
	fset *token.FileSet,
	node ast.Node,
	doc *ast.CommentGroup,
) string {
	startPos := node.Pos()
	if doc != nil && doc.Pos() < startPos {
		startPos = doc.Pos()
	}

	startOffset := fset.Position(startPos).Offset
	endOffset := fset.Position(node.End()).Offset

	if startOffset < 0 {
		startOffset = 0
	}
	if endOffset > len(src) {
		endOffset = len(src)
	}
	if startOffset > endOffset {
		return ""
	}

	return strings.TrimSpace(string(src[startOffset:endOffset]))
}

// getReceiverTypeName extracts the base type name of a method receiver.
func getReceiverTypeName(recv *ast.FieldList) string {
	if recv == nil || len(recv.List) == 0 {
		return ""
	}

	t := recv.List[0].Type
	for {
		switch v := t.(type) {
		case *ast.StarExpr:
			t = v.X
		case *ast.IndexExpr:
			t = v.X
		case *ast.IndexListExpr:
			t = v.X
		case *ast.ParenExpr:
			t = v.X
		case *ast.Ident:
			return v.Name
		default:
			return ""
		}
	}
}

// extractReceiverString returns the string representation of the receiver, e.g. "(s *MyStruct)".
func extractReceiverString(
	src []byte,
	fset *token.FileSet,
	recv *ast.FieldList,
) string {
	if recv == nil {
		return ""
	}
	startOffset := fset.Position(recv.Opening).Offset
	endOffset := fset.Position(recv.Closing + 1).Offset
	if startOffset >= 0 && endOffset <= len(src) && startOffset < endOffset {
		return string(src[startOffset:endOffset])
	}
	return ""
}
