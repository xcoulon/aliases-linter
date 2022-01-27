package linter_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xcoulon/aliases-linter/pkg/linter"
)

func TestLiner(t *testing.T) {

	t.Run("valid files", func(t *testing.T) {
		// when
		err := linter.Lint("../../test/valid")
		// then
		require.NoError(t, err)
	})

	t.Run("invalid files", func(t *testing.T) {
		// when
		err := linter.Lint("../../test/invalid")
		// then
		require.Error(t, err)
		errs := &multierror.Error{}
		require.True(t, errors.As(err, &errs))
		require.Len(t, errs.Errors, 5)
		t.Log(errs)
		assert.Equal(t, linter.LinterError{
			Filename:      "sub/invalid.go",
			Path:          "errors",
			ExpectedAlias: "",
			ActualAlias:   "alias1",
		}, errs.Errors[0])
		assert.Equal(t, linter.LinterError{
			Filename:      "sub/invalid.go",
			Path:          "github.com/codeready-toolchain/registration-service/pkg/errors",
			ExpectedAlias: "crterrors",
			ActualAlias:   "alias2",
		}, errs.Errors[1])
		assert.Equal(t, linter.LinterError{
			Filename:      "sub/invalid.go",
			Path:          "k8s.io/api/core/v1",
			ExpectedAlias: "corev1",
			ActualAlias:   "alias3",
		}, errs.Errors[2])
		assert.Equal(t, linter.LinterError{
			Filename:      "sub/invalid.go",
			Path:          "k8s.io/apimachinery/pkg/api/errors",
			ExpectedAlias: "apierrors",
			ActualAlias:   "alias4",
		}, errs.Errors[3])
		assert.Equal(t, linter.LinterError{
			Filename:      "sub/invalid.go",
			Path:          "k8s.io/apimachinery/pkg/apis/meta/v1",
			ExpectedAlias: "metav1",
			ActualAlias:   "alias5",
		}, errs.Errors[4])

	})

	t.Run("unknown dir", func(t *testing.T) {
		// when
		err := linter.Lint("../../test/unknown")
		// then
		require.Error(t, err)
	})
}
