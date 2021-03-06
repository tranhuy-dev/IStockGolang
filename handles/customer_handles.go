package handles

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/tranhuy-dev/IStockGolang/core/constant"
	"github.com/tranhuy-dev/IStockGolang/database"
	models "github.com/tranhuy-dev/IStockGolang/models"
)

// Get customer
func GetCustomer(c echo.Context) error {
	customer := database.RetrieveAllCustomer()
	return c.JSON(http.StatusOK, models.SuccessReponse{
		Code:    200,
		Message: "Retrieve customers success",
		Data:    customer})
}

// Create customer
func CreateCustomer(c echo.Context) error {
	var req models.CustomerReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    constant.BadRequest,
			Message: "Bad parameter"})
	}
	customer := database.InsertCustomer(req)
	return c.JSON(http.StatusOK, models.SuccessReponse{
		Code:    constant.Success,
		Message: "Create success",
		Data:    customer})
}

func UpdateCustomer(c echo.Context) error {
	var req models.CustomerReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    constant.BadRequest,
			Message: "Bad request"})
	}
	token := GetTokenHeader(c)
	updateResult, err := database.UpdateCustomer(req, token)

	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    constant.NotFound,
			Message: err.Error()})
	}
	return c.JSON(http.StatusOK, models.SuccessReponse{
		Code:    constant.Success,
		Message: "update success",
		Data:    updateResult})
}

func DeleteCustomer(c echo.Context) error {
	token := GetTokenHeader(c)
	deleteResult, err := database.DeleteCustomer(token)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    constant.NotFound,
			Message: err.Error()})
	}
	return c.JSON(http.StatusOK, models.SuccessReponse{
		Code:    constant.Success,
		Message: "Delete success",
		Data:    deleteResult})
}

func GetCustomerByEmail(c echo.Context) error {
	idCustomer := c.Param("email")
	token := GetTokenHeader(c)
	customerData, err := database.FindUserByEmail(token , idCustomer)
	if err != nil {
		return c.JSON(http.StatusOK, models.ErrorResponse{Code: constant.NotFound, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, customerData)
}

func ChangePassword(c echo.Context) error {
	var req models.ChangePasswordReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    constant.BadRequest,
			Message: "Bad request"})
	}

	token := GetTokenHeader(c)
	dataChangePassword, err := database.ChangePassword(req, token)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    401,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.SuccessReponse{
		Code:    201,
		Message: "Change password success",
		Data:    dataChangePassword})
}

func FindUserByFilter(c echo.Context) error {
	var filterBody models.FilterUser
	filterBody.Age, _ = strconv.Atoi(c.FormValue("age"))
	filterBody.Address = c.FormValue("address")
	filterBody.Name = c.FormValue("name")
	customers, err := database.RetrieveCustomerByFilter(filterBody)
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.ErrorResponse{
			Code:    constant.ExpectedError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.SuccessReponse{
		Code:    200,
		Message: "Retrieve customer filter",
		Data:    customers})
}

func LoginAccount(c echo.Context) error {
	var loginBody models.LoginBody
	err := c.Bind(&loginBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    constant.BadRequest,
			Message: "Bad parameter"})
	}

	dataLogin , err := database.LoginAccount(loginBody)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    constant.NotFound,
			Message: err.Error()})
	}
	return c.JSON(http.StatusOK, models.SuccessReponse{
		Code:    constant.Success,
		Message: "Login success",
		Data:    dataLogin})
}
