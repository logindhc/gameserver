package gate

import (
	"fmt"
	clog "gameserver/cherry/logger"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

// 全局注册表
var structRegistry = make(map[string]interface{})

// 注册结构体
func registerStruct(name string, s interface{}) {
	structRegistry[name] = s
}

// 加载指定目录中的结构体
func LoadStructsFromDir(dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 仅处理 .go 文件
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
			return parseFile(path)
		}
		return nil
	})
}

// 解析 Go 文件，提取 S2C 和 C2S 开头的结构体
func parseFile(filePath string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %v", filePath, err)
	}

	for _, decl := range node.Decls {
		// 类型断言为 *ast.GenDecl
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		// 检查是否为类型声明
		if genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			// 类型断言为 *ast.TypeSpec
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structName := typeSpec.Name.Name
			if strings.HasPrefix(structName, "S2C") || strings.HasPrefix(structName, "C2S") {
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					registerStruct(structName, structType)
				}
			}
		}
	}

	return nil
}

func GetStruct(name string) interface{} {
	return structRegistry[name]
}

func PrintStructs() {
	for name, facade := range structRegistry {
		clog.Infof("struct name = %s, facade = %v", name, facade)
	}
}
