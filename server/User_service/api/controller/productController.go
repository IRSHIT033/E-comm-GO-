package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/IRSHIT033/E-comm-GO-/server/User_service/domain_user"
	"github.com/IRSHIT033/E-comm-GO-/server/User_service/grpc_client"
	"github.com/IRSHIT033/E-comm-GO-/server/User_service/kafka_producer"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	ProductUsecase domain_user.ProductUseCase
}

func (pc *ProductController) AddProductToCustomers(c *gin.Context) {

	var Productidreq domain_user.ProductIdRequest
	err := c.ShouldBind(&Productidreq)

	product := grpc_client.GetProductViaGRPC(Productidreq.ProductId)

	log.Println(product)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain_user.ErrorResponse{Message: err.Error()})
		return
	}

	userID := c.GetString("x-user-id")
	userid, err := strconv.ParseUint(userID, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain_user.ErrorResponse{Message: err.Error()})
		return
	}

	//send the "adding product to cart" message to kafka broker
	var cart domain_user.KafkaMessagePayload
	cart.CustomerId = uint(userid)
	cart.Product = product
	cart.Operation = "Add"
	go kafka_producer.ProduceCart(cart)
	////////////////////////////////////////////////////////////

	err = pc.ProductUsecase.Add(c, uint(userid), &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain_user.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain_user.SuccessResponse{
		Message: "product added to cart successfully",
	})

}

func (pc *ProductController) GetproductOfUser(c *gin.Context) {

	userid, err := strconv.ParseUint(c.GetString("x-user-id"), 10, 32)
	userID := uint(userid)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain_user.ErrorResponse{Message: err.Error()})
		return
	}

	cart, err := pc.ProductUsecase.FetchByUserID(c, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain_user.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)

}

func (pc *ProductController) RemoveProductFromCart(c *gin.Context) {

	productid, err := strconv.ParseUint(c.Param("id"), 10, 32)
	productID := uint(productid)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain_user.ErrorResponse{Message: err.Error()})
		return
	}

	userid, err := strconv.ParseUint(c.GetString("x-user-id"), 10, 32)
	userID := uint(userid)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain_user.ErrorResponse{Message: err.Error()})
		return
	}

	err = pc.ProductUsecase.Remove(c, productID, userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain_user.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain_user.SuccessResponse{
		Message: "product deleted to cart successfully",
	})
}

func (pc *ProductController) GetCartWithUser(c *gin.Context) {

	userid, err := strconv.ParseUint(c.GetString("x-user-id"), 10, 32)
	userID := uint(userid)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain_user.ErrorResponse{Message: err.Error()})
		return
	}

	var user domain_user.User
	user, err = pc.ProductUsecase.FetchUserCart(c, userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain_user.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)

}
