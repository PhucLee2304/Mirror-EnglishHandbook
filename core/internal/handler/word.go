package handler

import (
	"core/config"
	"core/internal/dto"
	"core/internal/service"
	"core/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WordHandler struct {
	wordService *service.WordService
}

func NewWordHandler(wordService *service.WordService) *WordHandler {
	return &WordHandler{wordService: wordService}
}

func (h *WordHandler) SetupRouter(r *gin.RouterGroup, cfg *config.Config) {
	g := r.Group("/v1/words")
	{
		g.GET("/:id", h.getByID)
		g.GET("", h.getList)
	}
}

// getByID
// @Summary Get word details by ID
// @Description Retrieve detailed information about a specific word using its unique ID (phonetics, meanings, definitions).
// @Tags Words
// @Produce JSON
// @Param id path int true "Word ID"
// @Success 200 {object} dto.WordBase "Word details retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid ID format"
// @Failure 404 {object} response.ErrorResponse "Word not found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/words/{id} [get]
func (h *WordHandler) getByID(c *gin.Context) {
	var uri dto.WordUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:  response.MessageCodeBadRequest,
			Detail: err.Error(),
		})
		return
	}

	word, httpErr := h.wordService.GetByID(c.Request.Context(), uri.ID)
	if httpErr != nil {
		c.JSON(httpErr.StatusCode, httpErr.Error)
		return
	}
	c.JSON(http.StatusOK, word)
}

// getList
// @Summary Get a list of words
// @Description Retrieve a paginated list of words from the dictionary. Supports text search and filtering by part of speech (noun, verb, etc.).
// @Tags Words
// @Produce json
// @Param query query dto.GetWordsQuery false "Pagination and Filter parameters"
// @Success 200 {object} dto.GetWordsResponse "Words retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid query parameters (e.g., invalid limit/offset)"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/words [get]
func (h *WordHandler) getList(c *gin.Context) {
	var query dto.GetWordsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:  response.MessageCodeBadRequest,
			Detail: err.Error(),
		})
		return
	}

	response, httpErr := h.wordService.GetList(c.Request.Context(), query)
	if httpErr != nil {
		c.JSON(httpErr.StatusCode, httpErr.Error)
		return
	}
	c.JSON(http.StatusOK, response)
}
