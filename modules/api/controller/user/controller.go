package user

import (
	"net/http"
	"reflect"

	"go-rest-todo/modules/api/common"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {

	type Rets struct {
		RoleID   string `json:"role_id" validate:"required"`
		Name     string `json:"name" validate:"required"`
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	ret := new(Rets)
	c.Bind(ret)

	if reflect.DeepEqual(*ret, Rets{}) {
		return c.JSON(http.StatusBadRequest, common.GetErrorMessage("52C", "Empty"))
	} else if ret.Username == "existuser" {
		return c.JSON(http.StatusBadRequest, common.GetErrorMessage("51C", ""))
	}

	vald := *validator.New()
	err := vald.Struct(ret)
	if err != nil {
		var errStr string
		for _, er := range err.(validator.ValidationErrors) {
			errStr += er.Field() + ","
		}
		return c.JSON(http.StatusBadRequest, common.GetErrorMessage("52C", errStr))
	}

	return c.JSON(http.StatusCreated, ret)
}
