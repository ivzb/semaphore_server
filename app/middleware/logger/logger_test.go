package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	a "github.com/ivzb/semaphore_server"
	"github.com/ivzb/semaphore_server/middleware/app"
	"github.com/ivzb/semaphore_server/shared/config"
	"github.com/ivzb/semaphore_server/shared/logger/mock"
	"github.com/ivzb/semaphore_server/shared/response"
)

func testHandler(env *a.Env) *response.Message {
	return response.Ok("ok", 1, "OK")
}

func TestLoggerHandler_Log(t *testing.T) {
	req := httptest.NewRequest("GET", "/logger", nil)

	rec := httptest.NewRecorder()

	env := &a.Env{
		Log:    &mock.Logger{},
		Config: &config.Config{},
	}

	app := app.App{Env: env, Handler: testHandler}

	var handler http.Handler = Handler(app)

	handler.ServeHTTP(rec, req)

	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"message":"ok found","results":"OK"}`
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rec.Body.String(), expected)
	}
}
