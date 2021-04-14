package berrors

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestBadInputsJSON(t *testing.T) {
	// test for an empty argument
	emptyResultExpectedValue, _ := json.Marshal(BadInput{})
	emptyResult := BadInputsJSON(map[string]string{})
	if bytes.Compare(emptyResult, emptyResultExpectedValue) != 0 {
		t.Errorf("BadInputsJSON return value should be empty. Got: %s\n", string(emptyResult))
	}

	// test for a single data
	singleDataExpectedValue, _ := json.Marshal(BadInput{Fields: []field{{"name", "john"}}})
	singleData := BadInputsJSON(map[string]string{"name": "john"})
	if bytes.Compare(singleData, singleDataExpectedValue) != 0 {
		t.Errorf("BadInputsJSON return value should be equal to:\n%s\ngot:\n%s\n", singleDataExpectedValue, singleData)
	}
}

func TestBadInputJSON(t *testing.T) {
	// test empty data
	emptyResultExpectedValue, _ := json.Marshal(BadInput{Fields: []field{{"", ""}}})
	emptyResult := BadInputJSON("", "")
	if bytes.Compare(emptyResult, emptyResultExpectedValue) != 0 {
		t.Errorf("BadInputsJSON return value should be empty. Got: %s\n", string(emptyResult))
	}

	// test correct data
	dataExpectedValue, _ := json.Marshal(BadInput{Fields: []field{{"name", "john"}}})
	data := BadInputJSON("name", "john")
	if bytes.Compare(data, dataExpectedValue) != 0 {
		t.Errorf("BadInputsJSON return value should be equal to:\n%s\ngot:\n%s\n", dataExpectedValue, data)
	}
}

func TestBadInputsJSONFromType(t *testing.T) {
	// test for an empty argument
	emptyResultExpectedValue, _ := json.Marshal(BadInput{})
	emptyResult := BadInputsJSONFromType(map[string]string{})
	if bytes.Compare(emptyResult, emptyResultExpectedValue) != 0 {
		t.Errorf("BadInputsJSON return value should be empty. Got: %s\n", string(emptyResult))
	}

	// test for an incorrect data
	incorrectDataExpectedValue, _ := json.Marshal(BadInput{Fields: []field{{"username", ""}}})
	incorrectData := BadInputsJSONFromType(map[string]string{"username": "JohnDoe"})
	if bytes.Compare(incorrectData, incorrectDataExpectedValue) != 0 {
		t.Errorf("BadInputsJSON return value should be equal to:\n%s\ngot:\n%s\n", incorrectDataExpectedValue, incorrectData)
	}

	// test for a single data
	singleDataExpectedValue, _ := json.Marshal(BadInput{Fields: []field{{"username", ErrorTypes[Required]}}})
	singleData := BadInputsJSONFromType(map[string]string{"username": string(Required)})
	if bytes.Compare(singleData, singleDataExpectedValue) != 0 {
		t.Errorf("BadInputsJSON return value should be equal to:\n%s\ngot:\n%s\n", singleDataExpectedValue, singleData)
	}
}

func TestBadInputJSONFromType(t *testing.T) {
	// test for empty arguments
	emptyResultExpectedValue, _ := json.Marshal(BadInput{Fields: []field{{"", ""}}})
	emptyResult := BadInputJSONFromType("", "")
	if bytes.Compare(emptyResult, emptyResultExpectedValue) != 0 {
		t.Errorf("BadInputsJSON return value should be empty. Got: %s\n", string(emptyResult))
	}

	// test for incorrect type
	incorrectDataExpectedValue, _ := json.Marshal(BadInput{Fields: []field{{"username", ""}}})
	incorrectData := BadInputJSONFromType("username", "JohnDoe")
	if bytes.Compare(incorrectData, incorrectDataExpectedValue) != 0 {
		t.Errorf("BadInputsJSON return value should be equal to:\n%s\ngot:\n%s\n", incorrectDataExpectedValue, incorrectData)
	}

	// test for correct data
	dataExpectedValue, _ := json.Marshal(BadInput{Fields: []field{{"username", ErrorTypes[Required]}}})
	data := BadInputJSONFromType("username", string(Required))
	if bytes.Compare(data, dataExpectedValue) != 0 {
		t.Errorf("BadInputsJSON return value should be equal to:\n%s\ngot:\n%s\n", dataExpectedValue, data)
	}
}

func TestInternalServerError(t *testing.T) {
	fiberApp := fiber.New()
	ctx := fiberApp.AcquireCtx(&fasthttp.RequestCtx{})
	InternalServerError(ctx, errors.New("test error"))

	if ctx.Context().Response.Header.StatusCode() != http.StatusInternalServerError {
		t.Errorf("Http error should be: %d\ngot: %d\n", http.StatusBadRequest, ctx.Context().Response.Header.StatusCode())
	}
}
