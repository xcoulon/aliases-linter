// +build:never
package valid

import (
	"errors"

	crterrors "github.com/codeready-toolchain/registration-service/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	_ = errors.New("an error")
	_ = apierrors.IsAlreadyExists(nil)
	_ = crterrors.NewBadRequest("", "")
	_ = metav1.ObjectMeta{}
	_ = corev1.ConditionTrue
}
