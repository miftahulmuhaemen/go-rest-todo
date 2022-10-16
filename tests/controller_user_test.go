package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"go-rest-todo/business"
	core "go-rest-todo/core/user"
	"go-rest-todo/modules/api/common"
	controller "go-rest-todo/modules/api/controller/user"

	mocks "go-rest-todo/business/mocks/user"

	"github.com/jonboulle/clockwork"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestControllerUserRegister(t *testing.T) {

	e := echo.New()
	clock := clockwork.NewFakeClock()

	type mock struct {
		service *mocks.Service
	}

	tests := []struct {
		name        string
		prep        func(f *mock)
		reqBody     map[string]interface{}
		status      int
		retExpected core.User
		errExpected common.ResponseCode
	}{
		{
			name: "Parameters empty",
			prep: func(f *mock) {
				f.service.EXPECT().Register(core.User{}).Return(core.User{}, business.ErrInvalidSpec)
			},
			reqBody: map[string]interface{}{},
			status:  http.StatusBadRequest,
			errExpected: common.ResponseCode{
				RC:      "52C",
				Message: "request validation mismatch",
			},
		},
		{
			name: "Existed Username",
			prep: func(f *mock) {
				f.service.EXPECT().Register(core.User{
					RoleID:   "3f82fa1f-42a4-401b-8607-0674b94b6dab",
					Username: "existuser",
					Password: "password",
					Name:     "name",
				}).Return(core.User{}, business.ErrInvalidExistingUsername)
			},
			reqBody: map[string]interface{}{
				"role_id":  "3f82fa1f-42a4-401b-8607-0674b94b6dab",
				"username": "existuser",
				"password": "password",
				"name":     "name",
			},
			status: http.StatusBadRequest,
			errExpected: common.ResponseCode{
				RC:      "51C",
				Message: "username already exist",
			},
		},
		{
			name: "Only some parameter not empty",
			prep: func(f *mock) {
				f.service.EXPECT().Register(core.User{
					RoleID: "3f82fa1f-42a4-401b-8607-0674b94b6dab",
					Name:   "name",
				}).Return(core.User{}, business.ErrInvalidSpec)
			},
			reqBody: map[string]interface{}{
				"role_id": "3f82fa1f-42a4-401b-8607-0674b94b6dab",
				"name":    "name",
			},
			status: http.StatusBadRequest,
			errExpected: common.ResponseCode{
				RC:      "52C",
				Message: "request validation mismatch",
			},
		},
		{
			name: "Success;",
			prep: func(f *mock) {
				f.service.EXPECT().Register(core.User{
					RoleID:   "3f82fa1f-42a4-401b-8607-0674b94b6dab",
					Username: "user",
					Password: "password",
					Name:     "name",
				}).Return(core.User{
					ID:        "testID",
					RoleID:    "3f82fa1f-42a4-401b-8607-0674b94b6dab",
					Username:  "user",
					Password:  "password",
					Name:      "name",
					CreatedAt: clock.Now(),
					Updatet:   clock.Now(),
				}, nil)
			},
			reqBody: map[string]interface{}{
				"role_id":  "3f82fa1f-42a4-401b-8607-0674b94b6dab",
				"username": "user",
				"password": "password",
				"name":     "name",
			},
			status: http.StatusOK,
			retExpected: core.User{
				ID:        "testID",
				RoleID:    "3f82fa1f-42a4-401b-8607-0674b94b6dab",
				Username:  "user",
				Password:  "password",
				Name:      "name",
				CreatedAt: clock.Now(),
				Updatet:   clock.Now(),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {

			t.Parallel()

			service := new(mocks.Service)
			if test.prep != nil {
				test.prep(&mock{
					service: service,
				})
			}

			inputArg, _ := json.Marshal(test.reqBody)

			req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(string(inputArg)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			ctrl := controller.NewController(service)
			if assert.NoError(t, ctrl.Register(c)) {
				if !reflect.DeepEqual(test.errExpected, common.ResponseCode{}) {
					var rc common.ResponseCode
					json.Unmarshal(rec.Body.Bytes(), &rc)
					assert.Equal(t, test.errExpected, rc)
				} else {
					var retActual core.User
					json.Unmarshal(rec.Body.Bytes(), &retActual)
					assert.Equal(t, test.retExpected, retActual)
				}
			}
		})
	}
}
