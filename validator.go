package berrors

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/fatih/structtag"
	"github.com/gofiber/fiber/v2"

	"gopkg.in/go-playground/validator.v9"
)

// GetJSONFieldName is used to get the json tag of a given field in a struct
func GetJSONFieldName(object interface{}, fieldName string) (string, error) {
	field, ok := reflect.TypeOf(object).Elem().FieldByName(fieldName)
	if !ok {
		return "", errors.New("Field name `" + fieldName + "` not found within the given object `" + reflect.TypeOf(object).Elem().Name() + "`")
	}
	tags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		log.Println(err)
		return "", err
	}
	jsonTag, err := tags.Get("json")
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Remove `,omitempty` from the tag
	jsonTagName := strings.Replace(jsonTag.Name, ",omitempty", "", -1)

	return jsonTagName, nil
}

// ParseBodyJSON parse and validate a body contained in the fiber context to the given object
// If an error occur, the correct http error is called and false is returned
// The validator errors messages are using BadInputsJSONFromType
func ParseBodyJSON(ctx *fiber.Ctx, object interface{}) bool {
	if err := ctx.BodyParser(object); err != nil {
		_ = ctx.Status(http.StatusBadRequest).SendString(err.Error())
		log.Println(err)
		return false
	}

	v := validator.New()
	if err := v.Struct(object); err != nil {
		log.Println(err)
		errorMap := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			jsonTagName, err := GetJSONFieldName(object, e.Field())
			if err != nil {
				InternalServerError(ctx, err)
				return false
			}
			errorMap[jsonTagName] = e.Tag()
		}
		_ = ctx.Status(http.StatusBadRequest).Send(BadInputsJSONFromType(errorMap))
		return false
	}
	return true
}
