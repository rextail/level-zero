package getorder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"level-zero/internal/models"
	"level-zero/internal/storage/strgerrs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockOrderGetter struct {
	order string
	err   error
}

func (m *MockOrderGetter) OrderByID(ctx context.Context, ID string) (string, error) {
	return m.order, m.err
}

type MockOrderResponser struct {
	err error
}

func (m MockOrderResponser) OrderResponse(w http.ResponseWriter, order *models.Order) error {
	//TODO implement me
	return m.err
}

func TestNew(t *testing.T) {
	t.Run("case when we found order in storage and have no representation errors", func(t *testing.T) {
		mockGetter := MockOrderGetter{
			order: "some_order",
			err:   nil,
		}
		mockResponser := MockOrderResponser{
			err: nil,
		}
		req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(`{"id":"123"}`))
		rr := httptest.NewRecorder()

		handler := New(context.Background(), &mockGetter, mockResponser)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

	})
	t.Run("case when we didn't found order in storage and have no representation errors", func(t *testing.T) {
		mockGetter := MockOrderGetter{
			order: "some_order",
			err:   strgerrs.ErrZeroRecordsFound,
		}
		mockResponser := MockOrderResponser{
			err: nil,
		}
		req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(`{"id":"123"}`))
		rr := httptest.NewRecorder()

		handler := New(context.Background(), &mockGetter, mockResponser)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
	t.Run("case when we didn't found order in storage and have representation errors", func(t *testing.T) {
		mockGetter := MockOrderGetter{
			order: "some_order",
			err:   strgerrs.ErrZeroRecordsFound,
		}
		mockResponser := MockOrderResponser{
			err: fmt.Errorf("some template error"),
		}
		req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(`{"id":"123"}`))
		rr := httptest.NewRecorder()

		handler := New(context.Background(), &mockGetter, mockResponser)

		handler.ServeHTTP(rr, req)

		var actualResp map[string]interface{}

		json.Unmarshal([]byte(rr.Body.String()), &actualResp)

		expectedResponse := `{"status":"Error","error":"server failed to form response"}`

		var expectedResp map[string]interface{}

		json.Unmarshal([]byte(expectedResponse), &expectedResp)

		assert.Equal(t, actualResp, expectedResp)
	})
}
