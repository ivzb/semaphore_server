package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	a "github.com/ivzb/semaphore_server"
	"github.com/ivzb/semaphore_server/middleware/app"
	"github.com/ivzb/semaphore_server/shared/config"
	"github.com/ivzb/semaphore_server/shared/response"

	dMock "github.com/ivzb/semaphore_server/db/mock"
	tMock "github.com/ivzb/semaphore_server/shared/token/mock"
)

func testHandler(env *a.Env) *response.Message {
	return response.Created("auth_token", "auth token here")
}

func TestAuthHandler_ValidAuthToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/auth", nil)
	req.Header.Add("auth_token", "asdf")

	rr := httptest.NewRecorder()

	env := &a.Env{
		DB: &dMock.DB{
			UserMock: dMock.User{
				ExistsMock: dMock.UserExists{Bool: true, Err: nil},
			},
		},
		Token: &tMock.Tokener{
			DecryptMock: tMock.Decrypt{Dec: "decrypted", Err: nil},
		},
		Config: &config.Config{},
	}

	app := app.App{Env: env, Handler: testHandler}

	var handler http.Handler = Handler(app)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is what we expect.
	expected := `{"message":"auth_token created","results":"auth token here"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: \ngot %v \nwant %v",
			rr.Body.String(), expected)
	}
}

func TestAuthHandler_MissingAuthToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/auth", nil)

	rr := httptest.NewRecorder()

	env := &a.Env{
		DB:     &dMock.DB{},
		Token:  &tMock.Tokener{},
		Config: &config.Config{},
	}

	app := app.App{Env: env, Handler: testHandler}

	var handler http.Handler = Handler(app)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// Check the response body is what we expect.
	expected := `{"message":"missing auth_token"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAuthHandler_InvalidAuthToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/auth", nil)
	req.Header.Add("auth_token", "asdf")

	rec := httptest.NewRecorder()

	env := &a.Env{
		DB: &dMock.DB{},
		Token: &tMock.Tokener{
			DecryptMock: tMock.Decrypt{Dec: "", Err: errors.New("decryption error")},
		},
		Config: &config.Config{},
	}

	app := app.App{Env: env, Handler: testHandler}

	var handler http.Handler = Handler(app)

	handler.ServeHTTP(rec, req)

	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// Check the response body is what we expect.
	expected := `{"message":"invalid auth_token"}`
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rec.Body.String(), expected)
	}
}

func TestAuthHandler_DBError(t *testing.T) {
	req := httptest.NewRequest("GET", "/auth", nil)
	req.Header.Add("auth_token", "asdf")

	rec := httptest.NewRecorder()

	env := &a.Env{
		DB: &dMock.DB{
			UserMock: dMock.User{
				ExistsMock: dMock.UserExists{Bool: false, Err: errors.New("user does not exist")},
			},
		},
		Token: &tMock.Tokener{
			DecryptMock: tMock.Decrypt{Dec: "decrypted", Err: nil},
		},
		Config: &config.Config{},
	}

	app := app.App{Env: env, Handler: testHandler}

	var handler http.Handler = Handler(app)

	handler.ServeHTTP(rec, req)

	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// Check the response body is what we expect.
	expected := `{"message":"an error occurred, please try again later"}`
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rec.Body.String(), expected)
	}
}

func TestAuthHandler_UserDoesNotExist(t *testing.T) {
	req := httptest.NewRequest("GET", "/auth", nil)
	req.Header.Add("auth_token", "asdf")

	rec := httptest.NewRecorder()

	env := &a.Env{
		DB: &dMock.DB{
			UserMock: dMock.User{
				ExistsMock: dMock.UserExists{Bool: false, Err: nil},
			},
		},
		Token: &tMock.Tokener{
			DecryptMock: tMock.Decrypt{Dec: "decrypted", Err: nil},
		},
		Config: &config.Config{},
	}

	app := app.App{Env: env, Handler: testHandler}

	var handler http.Handler = Handler(app)

	handler.ServeHTTP(rec, req)

	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// Check the response body is what we expect.
	expected := `{"message":"invalid auth_token"}`
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rec.Body.String(), expected)
	}
}
