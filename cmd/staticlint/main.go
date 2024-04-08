// Package main - Мултичекер, содержащий следующие анализаторы -
// "github.com/golangci/gofmt/goimports" - проверка/исправление пакетного импорта,
// ищет и добавляет недостающие пакеты в import и удаляет из импорта пакеты, на которые никто не ссылается
// CheckOsExitCall - проверка на прямой запуск os.Exit в пакете main
// "github.com/kisielk/errcheck/errcheck" - проверка на использование возвращаемых ошибок
// "honnef.co/go/tools/staticcheck" - все проверки пакета SA, а также "ST1006" и "S1031"
// проврки стандартного анализатора 	"golang.org/x/tools/go/analysis/passes/printf",
// "golang.org/x/tools/go/analysis/passes/shadow", "golang.org/x/tools/go/analysis/passes/structtag"
// Запуск - ./mycheck.exe "pkg", список флагов и команд - ./mycheck.exe -help
package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"os"
	"path/filepath"

	"github.com/golangci/gofmt/goimports"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	// импортируем дополнительный анализатор
)

type checkPackageExpr struct {
	packageName string
	packageExpr string
}

// Config — имя файла конфигурации.
const Config = `cmd\staticlint\staticconf.json`

// ConfigData описывает структуру файла конфигурации.
type ConfigData struct {
	Staticcheck []string
}

var CheckOsExitCall = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for direct os exit call ",
	Run:  run,
}

var GoImportsCall = &analysis.Analyzer{
	Name: "goimports",
	Doc:  "running goimports",
	Run:  goimportsRun,
}

var LastPackageName string
var LastFuncDeclared string

func main() {
	appfile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
	if err != nil {
		panic(err)
	}
	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	mychecks := []*analysis.Analyzer{
		GoImportsCall,
		CheckOsExitCall,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		errcheck.Analyzer,
	}
	checks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		checks[v] = true
	}
	// добавляем анализаторы из staticcheck, которые указаны в файле конфигурации
	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	multichecker.Main(
		mychecks...,
	)
}

func run(pass *analysis.Pass) (interface{}, error) {
	checkPackage := checkPackageExpr{
		packageName: "os",
		packageExpr: "Exit",
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// проверяем, какой конкретный тип лежит в узле

			switch x := n.(type) {
			case *ast.File:
				LastPackageName = x.Name.String()
			case *ast.FuncDecl:
				LastFuncDeclared = x.Name.String()
			case *ast.CallExpr:
				// ast.CallExpr представляет вызов функции или метода
				if LastPackageName == "main" && LastFuncDeclared == "main" {
					switch node := x.Fun.(type) {
					case *ast.SelectorExpr:

						packageName := fmt.Sprintf("%s", node.X)
						packageExpr := node.Sel.Name
						if packageName == checkPackage.packageName && packageExpr == checkPackage.packageExpr {
							pass.Reportf(x.Pos(), "using direct call to os.Exit")
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func goimportsRun(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		goimports.Run(file.Name.String())
	}
	return nil, nil
}
