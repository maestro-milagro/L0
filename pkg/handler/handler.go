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
	return &Handler{
		services: services,
	}
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
		"order_uid":          message1.OrderUid,
		"track_number":       message1.TrackNumber,
		"entry":              message1.Entry,
		"name":               message1.Delivery.Name,
		"phone":              message1.Delivery.Phone,
		"zip":                message1.Delivery.Zip,
		"city":               message1.Delivery.City,
		"address":            message1.Delivery.Address,
		"region":             message1.Delivery.Region,
		"email":              message1.Delivery.Email,
		"transaction":        message1.Payment.Transaction,
		"request_id":         message1.Payment.RequestId,
		"currency":           message1.Payment.Currency,
		"provider":           message1.Payment.Provider,
		"amount":             message1.Payment.Amount,
		"payment_dt":         message1.Payment.PaymentDt,
		"bank":               message1.Payment.Bank,
		"delivery_cost":      message1.Payment.DeliveryCost,
		"goods_total":        message1.Payment.GoodsTotal,
		"custom_fee":         message1.Payment.CustomFee,
		"chrt_id":            message1.Items[0].ChrtId,
		"track_numberI":      message1.Items[0].TrackNumber,
		"price":              message1.Items[0].Price,
		"rid":                message1.Items[0].Rid,
		"nameI":              message1.Items[0].Name,
		"sale":               message1.Items[0].Sale,
		"size":               message1.Items[0].Size,
		"total_price":        message1.Items[0].TotalPrice,
		"nm_id":              message1.Items[0].NmId,
		"brand":              message1.Items[0].Brand,
		"status":             message1.Items[0].Status,
		"locale":             message1.Locale,
		"internal_signature": message1.InternalSignature,
		"customer_id":        message1.CustomerId,
		"delivery_service":   message1.DeliveryService,
		"shardkey":           message1.Shardkey,
		"sm_id":              message1.SmId,
		"date_created":       message1.DateCreated,
		"oof_shard":          message1.OofShard,
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
