package user

import (
	"gfly/app/constants"
	"gfly/app/dto"
	"gfly/app/http"
	"gfly/app/http/request"
	"gfly/app/http/transformers"
	"gfly/app/services"
	"gfly/app/utils"
	"github.com/gflydev/core"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

type UpdateUserStatusApi struct {
	core.Api
}

func NewUpdateUserStatusApi() *UpdateUserStatusApi {
	return &UpdateUserStatusApi{}
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

func (h UpdateUserStatusApi) Validate(c *core.Ctx) error {
	// Receive path parameter ID
	itemID, errData := http.PathID(c)
	if errData != nil {
		return c.Error(errData)
	}

	// Receive request data.
	var requestBody request.UpdateUserStatus
	if errData := http.Parse(c, &requestBody); errData != nil {
		return c.Error(errData)
	}

	// Update updateUserStatus struct with the given userID
	requestBody.ID = itemID

	// Convert to DTO
	requestDto := requestBody.ToDto()

	// Validate DTO
	if errData := utils.Validate(requestDto); errData != nil {
		return c.Error(errData)
	}

	// Store data into context.
	c.SetData(constants.Data, requestDto)

	return nil
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle Process main logic for API.
// @Summary Update user's status by ID
// @Description Update user's status by ID. <b>Administrator privilege required</b>
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body request.UpdateUserStatus true "Update user status data"
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Unauthorized
// @Success 200 {object} response.User
// @Security ApiKeyAuth
// @Router /users/{id}/status [put]
func (h UpdateUserStatusApi) Handle(c *core.Ctx) error {
	updateUserStatusDto := c.GetData(constants.Data).(dto.UpdateUserStatus)

	// Bind data to service
	user, err := services.UpdateUserStatus(updateUserStatusDto)
	if err != nil {
		c.Status(core.StatusBadRequest)
		return err
	}

	// Transform response data
	userResponse := transformers.ToUserResponse(*user)

	return c.Success(userResponse)
}
