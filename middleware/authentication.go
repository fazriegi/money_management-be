package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/common"
	userModel "github.com/fazriegi/money_management-be/module/master/user/model"
	"github.com/gofiber/fiber/v2"
)

func Authentication(jwt *libs.JWT) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var response = common.Response{}
		header := ctx.Get("Authorization")
		isHasBearer := strings.HasPrefix(header, "Bearer")

		if !isHasBearer {
			response = response.CustomResponse(http.StatusUnauthorized, "sign in to proceed", nil)

			return ctx.Status(response.Code).JSON(response)
		}

		tokenString := strings.Split(header, " ")[1]

		verifiedToken, err := jwt.VerifyJWTTOken(tokenString)
		if err != nil {
			response = response.CustomResponse(http.StatusUnauthorized, err.Error(), nil)

			return ctx.Status(response.Code).JSON(response)
		}

		jsonData, _ := json.Marshal(verifiedToken)

		var user userModel.User
		json.Unmarshal(jsonData, &user)

		ctx.Locals("user", user)

		return ctx.Next()
	}
}
