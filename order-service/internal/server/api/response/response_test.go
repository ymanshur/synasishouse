package response

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ymanshur/synasishouse/order/internal/typex"
)

func TestResponse_WithError(t *testing.T) {
	err := fmt.Errorf("fn(tx): %w", typex.NewUnprocessableEntityError("unprocessable entity"))
	rsp := New().WithTranslationError(err)
	assert.Equal(t, http.StatusUnprocessableEntity, rsp.GetCode())
}
