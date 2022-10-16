package user

import (
	"net/http"

	businessErr "go-rest-todo/business"
	business "go-rest-todo/business/user"
	core "go-rest-todo/core/user"
	"go-rest-todo/modules/api/common"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	service business.Service
}

func NewController(service business.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (controller *Controller) Register(c echo.Context) error {

	bindValue := new(core.User)
	c.Bind(bindValue)

	ret, err := controller.service.Register(*bindValue)
	if err != nil {
		if err == businessErr.ErrInvalidSpec {
			return c.JSON(http.StatusBadRequest, common.GetErrorMessage("52C", ""))
		} else if err == businessErr.ErrInvalidExistingUsername {
			return c.JSON(http.StatusBadRequest, common.GetErrorMessage("51C", ""))
		}
		return c.JSON(http.StatusInternalServerError, common.GetErrorMessage("53S", ""))
	}

	return c.JSON(http.StatusCreated, ret)
}
