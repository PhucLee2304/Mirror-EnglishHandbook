package handler

import (
	"core/internal/dto"
	"core/internal/service"
	"core/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) SetupRouter(r *gin.RouterGroup) *gin.RouterGroup {
	g := r.Group("/v1/auth")
	{
		g.POST("/firebase", h.login)
	}

	return g
}

// login
// @Summary login with Firebase ID Token
// @Description Authenticate user using Firebase ID token and receive JWT tokens
// @Tags Auth
// @Accept JSON
// @Produce JSON
// @Param request body dto.LoginBody true "login Request"
// @Success 200 {object} dto.LoginResponse "login successful"
// @Failure 400 {object} core.ErrorResponse "Invalid request payload"
// @Failure 401 {object} core.ErrorResponse "Unauthorized: Invalid or expired ID token"
// @Failure 500 {object} core.ErrorResponse "Internal server error during authentication"
// @Router /api/v1/auth/firebase [post]
func (h *AuthHandler) login(c *gin.Context) {
	var body dto.LoginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:  response.MessageCodeBadRequest,
			Detail: err.Error(),
		})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), body)
	if err != nil {
		c.JSON(err.StatusCode, err.Error)
		return
	}

	c.JSON(http.StatusOK, resp)
}
