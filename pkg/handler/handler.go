package handler

import (
	message "awesomeProject"
	"awesomeProject/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("templates/*")
	messages := router.Group("/message")
	{
		messages.POST("/", h.createMessage)
		messages.GET("/", h.getAllMessages)
		messages.GET("/:MessageId", h.getMessageById)
		messages.DELETE("/:MessageId", h.deleteMessage)

	}

	return router
}
func (h *Handler) createMessage(c *gin.Context) {
	var input message.Message
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.MessagesS.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllListsResponse struct {
	Messages []message.Message `json:"messages"`
}

func (h *Handler) getAllMessages(c *gin.Context) {
	list, err := h.services.MessagesS.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllListsResponse{
		Messages: list,
	})
}
func (h *Handler) getMessageById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("MessageId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	message1, err := h.services.MessagesS.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"message": message1,
	})
	//	c.JSON(http.StatusOK, message1)
}
func (h *Handler) deleteMessage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("MessageId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	err = h.services.MessagesS.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
