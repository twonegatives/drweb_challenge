package drweb_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/mocks"
)

func TestWithCallbacks(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	before := mocks.NewMockCallback(mockCtrl)
	after := mocks.NewMockCallback(mockCtrl)

	rr := httptest.NewRecorder()
	handler := func(http.ResponseWriter, *http.Request) {}

	gomock.InOrder(
		before.EXPECT().Invoke(gomock.Any(), gomock.Any()),
		after.EXPECT().Invoke(gomock.Any(), gomock.Any()),
	)

	function := drweb.WithCallbacks(handler, before, after)
	function(rr, &http.Request{})
}
