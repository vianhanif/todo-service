package errors_test

import (
	"testing"

	"github.com/vianhanif/todo-service/errors"
)

func TestCommonError(t *testing.T) {
	e := errors.NewCommonError("common error")
	if e.Message != "common error" {
		t.Fatal("expected message is 'common error'")
	}
	var ge error
	ge = e
	if ge.Error() != "common error" {
		t.Fatal("expected Error is 'common error'")
	}
}

func TestValidationError(t *testing.T) {
	e := errors.NewValidationError("Validation error")
	// test add field error
	e.FieldError("f1", "field1 error")
	if len(e.Fields) == 0 {
		t.Fatal("expected fields len greater than 0")
	}
	fe := e.GetFieldError("f1")
	if fe == nil {
		t.Fatal("expected field error f1")
	}
	if fe.Message != "field1 error" {
		t.Fatalf("expected field error message to be '%s'", "field1 error")
	}
	// test update field error
	e.FieldError("f1", "field1 error updated")
	fe = e.GetFieldError("f1")
	if fe == nil {
		t.Fatal("expected field error f1")
	}
	if fe.Message != "field1 error updated" {
		t.Fatalf("expected field error message to be '%s'", "field1 error updated")
	}
	// test clear field errors
	e.ClearFieldErrors()
	if len(e.Fields) != 0 {
		t.Fatal("expected fields len 0")
	}

	var ge error
	ge = e
	if ge.Error() != "Validation error" {
		t.Fatal("expected Error is 'Validation error'")
	}
}

func TestAuthError(t *testing.T) {
	e := errors.NewAuthError("auth error")
	if e.Message != "auth error" {
		t.Fatal("expected message is 'auth error'")
	}
	var ge error
	ge = e
	if ge.Error() != "auth error" {
		t.Fatal("expected Error is 'auth error'")
	}
}

func TestServiceError(t *testing.T) {
	e := errors.NewServiceError("service error")
	if e.Message != "service error" {
		t.Fatal("expected message is 'service error'")
	}
	var ge error
	ge = e
	if ge.Error() != "service error" {
		t.Fatal("expected Error is 'service error'")
	}
}

func TestNotFoundError(t *testing.T) {
	e := errors.NewNotFoundError("notfound error")
	if e.Message != "notfound error" {
		t.Fatal("expected message is 'notfound error'")
	}
	var ge error
	ge = e
	if ge.Error() != "notfound error" {
		t.Fatal("expected Error is 'notfound error'")
	}
}
