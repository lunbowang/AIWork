/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package xerr

import (
	"errors"
	"testing"
)

func TestWithMessage(t *testing.T) {
	err := WithMessage(errors.New("测试"), "测试输出")
	t.Log(err)

	type cause interface {
		Cause() error
	}

	e, ok := err.(cause)
	if ok {
		t.Log(e.Cause().Error())
	}
}
