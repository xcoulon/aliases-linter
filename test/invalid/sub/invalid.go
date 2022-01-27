// +build:never
package invalid

import (
	alias1 "errors"

	alias2 "github.com/codeready-toolchain/registration-service/pkg/errors"
	alias3 "k8s.io/api/core/v1"
	alias4 "k8s.io/apimachinery/pkg/api/errors"
	alias5 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	_ = alias1.New("an error")
	_ = alias4.IsAlreadyExists(nil)
	_ = alias2.NewBadRequest("", "")
	_ = alias5.ObjectMeta{}
	_ = alias3.ConditionTrue
}
