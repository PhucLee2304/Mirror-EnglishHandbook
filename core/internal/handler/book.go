package handler

import (
	"core/config"
	"core/internal/dto"
	"core/internal/service"
	"core/pkg/gateways/jwt"
	"core/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

func (h *BookHandler) SetupRouter(r *gin.RouterGroup, cfg *config.Config) {
	g := r.Group("/v1/books", jwt.Middleware(cfg))
	{
		g.GET("", h.getList)
		g.GET("/:id", h.getByID)
		g.GET(":id/lessons/lessonID", h.getLessonByID)
	}
}

func (h *BookHandler) getList(c *gin.Context) {
	var query dto.GetBooksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:  response.MessageCodeBadRequest,
			Detail: err.Error(),
		})
		return
	}

	authCtx := jwt.AuthContext{Context: c}
	_, err := authCtx.GetUserID()
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error:  response.MessageCodeUnauthorized,
			Detail: err.Error(),
		})
		return
	}

	resp, httpErr := h.bookService.GetList(c.Request.Context(), query)
	if httpErr != nil {
		c.JSON(httpErr.StatusCode, httpErr.Error)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *BookHandler) getByID(c *gin.Context) {
	var uri dto.BookUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:  response.MessageCodeBadRequest,
			Detail: err.Error(),
		})
		return
	}

	authCtx := jwt.AuthContext{Context: c}
	_, err := authCtx.GetUserID()
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error:  response.MessageCodeUnauthorized,
			Detail: err.Error(),
		})
		return
	}

	var query dto.GetLessonsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:  response.MessageCodeBadRequest,
			Detail: err.Error(),
		})
		return
	}

	resp, httpErr := h.bookService.GetByID(c.Request.Context(), uri.ID, query)
	if httpErr != nil {
		c.JSON(httpErr.StatusCode, httpErr.Error)
	}

	c.JSON(http.StatusOK, resp)
}

func (h *BookHandler) getLessonByID(c *gin.Context) {
	var uri dto.LessonUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:  response.MessageCodeBadRequest,
			Detail: err.Error(),
		})
		return
	}

	authCtx := jwt.AuthContext{Context: c}
	_, err := authCtx.GetUserID()
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error:  response.MessageCodeUnauthorized,
			Detail: err.Error(),
		})
		return
	}

	var query dto.GetQuestionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:  response.MessageCodeBadRequest,
			Detail: err.Error(),
		})
		return
	}

	resp, httpErr := h.bookService.GetLessonByID(c.Request.Context(), uri.LessonID, query)
	if httpErr != nil {
		c.JSON(httpErr.StatusCode, httpErr.Error)
	}

	c.JSON(http.StatusOK, resp)
}
