package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	userModel "github.com/fazriegi/money_management-be/module/master/user/model"
	"github.com/gofiber/fiber/v2"
)

func Authentication(jwt *libs.JWT) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var response = model.Response{}
		header := ctx.Get("Authorization")
		isHasBearer := strings.HasPrefix(header, "Bearer")

		if !isHasBearer {
			status := libs.CustomResponse(http.StatusUnauthorized, "sign in to proceed")
			response.Status = status

			return ctx.Status(status.Code).JSON(response)
		}

		tokenString := strings.Split(header, " ")[1]

		verifiedToken, err := jwt.VerifyJWTTOken(tokenString)
		if err != nil {
			status := libs.CustomResponse(http.StatusUnauthorized, err.Error())
			response.Status = status

			return ctx.Status(status.Code).JSON(response)
		}

		jsonData, _ := json.Marshal(verifiedToken)

		var user userModel.User
		json.Unmarshal(jsonData, &user)

		ctx.Locals("user", user)

		return ctx.Next()
	}
}
