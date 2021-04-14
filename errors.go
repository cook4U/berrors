// license that can be found in the LICENSE file.

// backend-errors is a package developed and used by alexandrio.
// It is made so that every micro service can use the same errors methods.
//
package berrors

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// BadInputsJSON creates the error JSON using the struct BadInput.
// The key of the map given correspond to the Name and the value to the Reason.
// It returns the JSON in []byte.
func BadInputsJSON(fields map[string]string) []byte {
	badInputData := BadInput{}
	for key, element := range fields {
		badInputData.Fields = append(badInputData.Fields, field{
			Name:   key,
			Reason: element,
		})
	}

	data, _ := json.Marshal(badInputData)
	return data
}

// BadInputJSON is simply a call to BadInputsJSON to create a single bad input error.
// It returns the JSON of the struct BadInput in []byte.
func BadInputJSON(name string, reason string) []byte {
	return BadInputsJSON(map[string]string{name: reason})
}

// BadInputsJSONFromType create a BadInput JSON from a key and a value corresponding to an ErrorType.
// It replace the Value with the defined string corresponding to the ErrorType.
// It returns the JSON in []byte.
func BadInputsJSONFromType(fields map[string]string) []byte {
	newFields := make(map[string]string)
	for key, element := range fields {
		newFields[key] = ErrorTypes[ErrorType(element)]
	}
	return BadInputsJSON(newFields)
}

// BadInputJSONFromType is simply a call to BadInputsJSONFromType to create a single bad input error.
// It returns the JSON of the struct BadInput in []byte.
func BadInputJSONFromType(name string, errorType string) []byte {
	return BadInputsJSONFromType(map[string]string{name: errorType})
}

// InternalServerError set a 500 http error and log the error.
func InternalServerError(ctx *fiber.Ctx, err error) {
	_ = ctx.SendStatus(http.StatusInternalServerError)
	log.Println(err)
}
