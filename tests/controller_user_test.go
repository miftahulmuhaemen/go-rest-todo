package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go-rest-todo/modules/api/common"
	controller "go-rest-todo/modules/api/controller/user"

	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestControllerUserRegister(t *testing.T) {

	e := echo.New()
	e.POST("/user", controller.Register)

	server := httptest.NewServer(e)
	defer server.Close()

	httpTest := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  server.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	tests := []struct {
		name        string
		retActual   interface{}
		retExpected map[string]interface{}
		errExpected common.ResponseCode
	}{
		{
			name: "Parameters empty",
			retActual: httpTest.Builder(func(req *httpexpect.Request) {}).
				POST("/user").
				WithJSON(map[string]interface{}{}).
				Expect().
				Status(http.StatusBadRequest).JSON().Raw(),
			errExpected: common.ResponseCode{
				RC:       "52C",
				Message:  "request validation mismatch",
				Messages: "Empty",
			},
		},
		{
			name: "Existed Username",
			retActual: httpTest.Builder(func(req *httpexpect.Request) {}).
				POST("/user").
				WithJSON(map[string]interface{}{
					"role_id":  "3f82fa1f-42a4-401b-8607-0674b94b6dab",
					"username": "existuser",
					"password": "password",
					"name":     "name"}).
				Expect().
				Status(http.StatusBadRequest).JSON().Raw(),
			errExpected: common.ResponseCode{
				RC:      "51C",
				Message: "username already exist",
			},
		},
		{
			name: "Only some parameter not empty",
			retActual: httpTest.Builder(func(req *httpexpect.Request) {}).
				POST("/user").
				WithJSON(map[string]interface{}{
					"role_id": "3f82fa1f-42a4-401b-8607-0674b94b6dab",
					"name":    "name"}).
				Expect().
				Status(http.StatusBadRequest).JSON().Raw(),
			errExpected: common.ResponseCode{
				RC:       "52C",
				Message:  "request validation mismatch",
				Messages: "Username,Password,",
			},
		},
		{
			name: "Sucess;",
			retActual: httpTest.Builder(func(req *httpexpect.Request) {}).
				POST("/user").
				WithJSON(map[string]interface{}{
					"role_id":  "3f82fa1f-42a4-401b-8607-0674b94b6dab",
					"username": "user",
					"password": "password",
					"name":     "name",
				}).
				Expect().
				Status(http.StatusCreated).JSON().Raw(),
			retExpected: map[string]interface{}{
				"role_id":  "3f82fa1f-42a4-401b-8607-0674b94b6dab",
				"username": "user",
				"password": "password",
				"name":     "name",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {

			t.Parallel()

			if !reflect.DeepEqual(test.errExpected, common.ResponseCode{}) {
				var rc common.ResponseCode
				raw, _ := json.Marshal(test.retActual)
				_ = json.Unmarshal(raw, &rc)

				assert.Equal(t, test.errExpected, rc)
			} else {
				assert.Equal(t, test.retExpected, test.retActual)
			}

		})
	}

}
