package shortcut_helper

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

var (
	namedKeys     = []string{KeyEnter, KeyEsc, KeySpace, KeyDelete, KeyPgUp, KeyPgDn, KeyTab}
	functionKey   = regexp.MustCompile(`^F([1-9]|1[0-2])$`)
	modifiedKey   = regexp.MustCompile(`^(ctrl|shift|alt)\+(.+)$`)
	wrongModifier = regexp.MustCompile(`(?i)^(ctrl|shift|alt)\+`)
)

func checkKeyName(name string) string {
	switch {
	case utf8.RuneCountInString(name) == 1:
		return ""
	case slices.Contains(namedKeys, name) || functionKey.MatchString(name):
		return ""
	case modifiedKey.MatchString(name):
		return checkKeyName(modifiedKey.FindStringSubmatch(name)[2])
	case wrongModifier.MatchString(name):
		return "modifiers are lowercase, e.g. ctrl+f"
	}
	for _, key := range namedKeys {
		if strings.EqualFold(name, key) {
			return "named keys are capitalized: " + key
		}
	}
	return "unknown key name"
}

func TestCheckKeyName(t *testing.T) {
	for _, valid := range []string{"h", "H", "?", "*", "↑", "Enter", "Esc", "⭾", "F1", "F2", "ctrl+q", "shift+⭾", "shift+↑", "ctrl+Enter"} {
		if problem := checkKeyName(valid); problem != "" {
			t.Errorf("expected %q to be valid, got problem: %s", valid, problem)
		}
	}
	for _, invalid := range []string{"Ctrl+q", "Shift+↑", "enter", "esc", "space", "pgup", "F13", "ctrl+enter", "Return"} {
		if problem := checkKeyName(invalid); problem == "" {
			t.Errorf("expected %q to be invalid", invalid)
		}
	}
}

func TestShortcutStyle(t *testing.T) {
	uiDir := ".."
	checked := 0
	err := filepath.WalkDir(uiDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, path, nil, 0)
		if err != nil {
			return err
		}
		check := func(literal *ast.BasicLit) {
			name, err := strconv.Unquote(literal.Value)
			if err != nil {
				return
			}
			checked++
			position := fileSet.Position(literal.Pos())
			if problem := checkKeyName(name); problem != "" {
				t.Errorf("%s: key %q: %s", position, name, problem)
			} else if slices.Contains(namedKeys, name) {
				t.Errorf("%s: key %q: use the constant from key_names.go", position, name)
			}
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.KeyValueExpr:
				key, ok := node.Key.(*ast.Ident)
				if !ok {
					return true
				}
				switch key.Name {
				case "KeyCombo":
					if list, ok := node.Value.(*ast.CompositeLit); ok {
						for _, element := range list.Elts {
							if literal, ok := element.(*ast.BasicLit); ok && literal.Kind == token.STRING {
								check(literal)
							}
						}
					}
				}
			case *ast.CallExpr:
				if selector, ok := node.Fun.(*ast.SelectorExpr); ok && (selector.Sel.Name == "Ctrl" || selector.Sel.Name == "Shift" || selector.Sel.Name == "Alt") {
					for _, argument := range node.Args {
						if literal, ok := argument.(*ast.BasicLit); ok && literal.Kind == token.STRING {
							check(literal)
						}
					}
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error walking dir: %v", err)
	}
	if checked == 0 {
		t.Errorf("expected at least 1 key check, got %d", checked)
	}
}
