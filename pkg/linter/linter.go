package linter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableQuote: true, // see https://github.com/sirupsen/logrus/issues/608#issuecomment-745137306
	})
}

func Lint(basepath string) error {
	var lintErrs error
	if err := filepath.WalkDir(basepath, func(filename string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() && strings.HasSuffix(d.Name(), ".go") {
			if lintErr := lint(basepath, filename); lintErr != nil {
				lintErrs = multierror.Append(lintErrs, lintErr)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return lintErrs
}

var aliases = map[string]string{
	"errors":                             "",
	"http://github.com/pkg/errors":       "errs",
	"k8s.io/apimachinery/pkg/api/errors": "apierrors",
	"github.com/codeready-toolchain/registration-service/pkg/errors": "crterrors",
	"k8s.io/apimachinery/pkg/apis/meta/v1":                           "metav1",
	"k8s.io/api/core/v1":                                             "corev1",
}

func lint(basepath, filename string) error {
	fset := token.NewFileSet()
	log.Infof("checking %s", filename)
	file, err := parser.ParseFile(fset, filename, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}
	var lintErrs error
	ast.Inspect(file, func(n ast.Node) bool {
		if imp, ok := n.(*ast.ImportSpec); ok {
			path := strings.Trim(imp.Path.Value, `"`)
			var actualAlias string
			if imp.Name != nil {
				actualAlias = imp.Name.Name
			} else {
				actualAlias = ""
			}
			log.Infof("found import of '%s' with alias '%s'", path, actualAlias)
			if expectedAlias, found := aliases[path]; found && expectedAlias != actualAlias {
				f, _ := filepath.Rel(basepath, filename)
				lintErrs = multierror.Append(lintErrs, LinterError{
					Filename:      f,
					Path:          path,
					ActualAlias:   actualAlias,
					ExpectedAlias: expectedAlias,
				})
			}
		}

		return true
	})
	return lintErrs
}

type LinterError struct {
	Filename      string
	Path          string
	ActualAlias   string
	ExpectedAlias string
}

var _ error = LinterError{}

func (e LinterError) Error() string {
	return fmt.Sprintf("invalid alias for package '%s' in %s: got '%s', expected '%s'", e.Path, e.Filename, e.ActualAlias, e.ExpectedAlias)
}
