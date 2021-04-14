package berrors

import (
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type testStruct struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email,omitempty" validate:"required,email"`
}

func TestGetJSONFieldName(t *testing.T) {
	testData := testStruct{}

	// test with an incorrect struct field
	incorrectResult, err := GetJSONFieldName(&testData, "Size")
	if err == nil || incorrectResult != "" {
		t.Errorf("GetJSONFieldName should return an error.")
	}

	// test with an correct struct field
	correctResult, err := GetJSONFieldName(&testData, "Name")
	if err != nil {
		t.Errorf("GetJSONFieldName shouldn't return an error. Err: %s\n", err)
	}
	if correctResult != "name" {
		t.Errorf("GetJSONFieldName should return `name`. Got: %s\n", correctResult)
	}

	// test with omitempty
	omitemptyResult, err := GetJSONFieldName(&testData, "Email")
	if err != nil {
		t.Errorf("GetJSONFieldName shouldn't return an error. Err: %s\n", err)
	}
	if omitemptyResult != "email" {
		t.Errorf("GetJSONFieldName should return `email`. Got: %s\n", omitemptyResult)
	}
}

func TestParseBodyJSONCorrectData(t *testing.T) {
	testData := new(testStruct)

	expectedData := testStruct{
		Name:  "John Doe",
		Email: "test@test.com",
	}

	fiberApp := fiber.New()
	ctx := fiberApp.AcquireCtx(&fasthttp.RequestCtx{})

	ctx.Context().Request.Header.SetContentType("application/json")
	ctx.Context().Request.SetBodyString("{\"name\":\"John Doe\",\"email\":\"test@test.com\"}")
	ok := ParseBodyJSON(ctx, testData)
	if !ok {
		t.Errorf("ParseBodyJSON should return true. Got: %t\n", ok)
	}
	if testData.Name != expectedData.Name || testData.Email != expectedData.Email {
		t.Errorf("testData from ParseBodyJSON does not contain the expected data.\nGot: \n%+v\nExpected: \n%+v\n", *testData, expectedData)
	}
}

func TestParseBodyJSONMissingRequiredData(t *testing.T) {
	testData := new(testStruct)

	expectedError := "{\"fields\":[{\"name\":\"name\",\"reason\":\"The field is required\"}]}"

	fiberApp := fiber.New()
	ctx := fiberApp.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Context().Request.Header.SetContentType("application/json")
	ctx.Context().Request.SetBodyString("{\"email\":\"test@test.com\"}")
	ok := ParseBodyJSON(ctx, testData)
	if ok {
		t.Errorf("ParseBodyJSON should return false. Got: %t\n", ok)
	}
	if ctx.Context().Response.Header.StatusCode() != http.StatusBadRequest ||
		string(ctx.Context().Response.Body()) != expectedError {
		t.Errorf("Http error should be:\n%d: %s\ngot:\n%d: %s\n", http.StatusBadRequest, expectedError, ctx.Context().Response.Header.StatusCode(), string(ctx.Context().Response.Body()))
	}
}

func TestParseBodyJSONIncorrectEmail(t *testing.T) {
	testData := new(testStruct)

	expectedError := "{\"fields\":[{\"name\":\"email\",\"reason\":\"The email given is not correct\"}]}"
	fiberApp := fiber.New()
	ctx := fiberApp.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Context().Request.Header.SetContentType("application/json")
	ctx.Context().Request.SetBodyString("{\"name\":\"John Doe\",\"email\":\"test.com\"}")
	ok := ParseBodyJSON(ctx, testData)
	if ok {
		t.Errorf("ParseBodyJSON should return false. Got: %t\n", ok)
	}
	if ctx.Context().Response.Header.StatusCode() != http.StatusBadRequest ||
		string(ctx.Context().Response.Body()) != expectedError {
		t.Errorf("Http error should be:\n%d: %s\ngot:\n%d: %s\n", http.StatusBadRequest, expectedError, ctx.Context().Response.Header.StatusCode(), string(ctx.Context().Response.Body()))
	}
}

func TestParseBodyJSONIncorrectJSON(t *testing.T) {
	testData := new(testStruct)

	expectedError := "json: cannot unmarshal \"2,\\\"email\\\":\\\"test@test.com\\\"}\" into Go struct field berrors.testStruct.name. of type string"
	fiberApp := fiber.New()
	ctx := fiberApp.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Context().Request.Header.SetContentType("application/json")
	ctx.Context().Request.SetBodyString("{\"name\":2,\"email\":\"test@test.com\"}")
	ok := ParseBodyJSON(ctx, testData)
	if ok {
		t.Errorf("ParseBodyJSON should return false. Got: %t\n", ok)
	}
	if ctx.Context().Response.Header.StatusCode() != http.StatusBadRequest ||
		string(ctx.Context().Response.Body()) != expectedError {
		t.Errorf("Http error should be:\n%d: %s\ngot:\n%d: %s\n", http.StatusBadRequest, expectedError, ctx.Context().Response.Header.StatusCode(), string(ctx.Context().Response.Body()))
	}
}
