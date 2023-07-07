package controllers

import (
	"errors"
	"go-payment/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTransaction(c *gin.Context) {
	// define varible from models
	var transaction models.Transaction

	// bind json from request body
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// get userID from session, the session created in middlewares/jwt.go
	strID, _ := c.Get("userID")
	userID := strID.(uint)
	userBalance := models.CheckBalanaceUserByID(userID)

	// check if balance user < trasanction amount
	if userBalance < int(transaction.Amount) {
		strAmount := models.IntToString(userBalance)
		message := "Not enough balance, your balance is " + strAmount
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": message})
		return
	}

	// assign sender user id & status
	transaction.SenderUserID = userID
	transaction.Status = "pending"

	// insert transaction
	models.DB.Debug().Create(&transaction)

	transactionResponse, err := getCustomFieldTransaction(transaction.ID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err})
		return
	}

	// insert log
	models.SaveActivity(userID, "create transaction", "transaction", transaction.ID)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Transaction has been created", "data": transactionResponse})
}

func ApproveTransaction(c *gin.Context) {
	// define varible from models
	var user models.User
	var transaction models.Transaction
	var transactionRequest models.ApproveTransactionRequest

	// get session user id
	strID, _ := c.Get("userID")
	userID := strID.(uint)

	// get session user detail by id
	userLogin, _ := models.GetUserByID(userID)

	// make sure user id is approval role
	if userLogin.RoleID != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not user approval"})
		return
	}

	// bind json from request body
	if err := c.ShouldBindJSON(&transactionRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// get transaction by id
	trx, _ := models.GetTransactionById(transactionRequest.ID)

	// check transaction already updated
	if trx.Status != "pending" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Transaction have already been updated"})
		return
	}

	// update transaction approve/reject
	models.DB.Debug().Model(&transaction).Where("ID = ?", transactionRequest.ID).Update("status", transactionRequest.Status)

	// do transaction approve
	if transactionRequest.Status == "approve" {

		// get senderID & receiverID
		senderID := trx.SenderUserID
		receiverID := trx.ReceiverUserID

		// get sender detail by id
		senderIdDetails, _ := models.GetUserByID(senderID)
		// get receiver detail by id
		receiverIdDetails, _ := models.GetUserByID(receiverID)

		// update balance sender
		updateBalanceSender := senderIdDetails.Balance - trx.Amount
		// update balance receiver
		updateBalanceReceiver := receiverIdDetails.Balance + trx.Amount

		// minus balance to sender
		models.DB.Debug().Model(&user).Where("ID = ?", senderID).Update("balance", updateBalanceSender)
		// add balance to receiver
		models.DB.Debug().Model(&user).Where("ID = ?", receiverID).Update("balance", updateBalanceReceiver)
	}

	message := transactionRequest.Status + " transaction"

	// insert log
	models.SaveActivity(userID, message, "transaction", uint(transactionRequest.ID))
	message = "Transaction has been " + transactionRequest.Status

	// get response custom field
	transactionResponse, _ := getCustomFieldTransaction(transactionRequest.ID)

	// response ok
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": message, "data": transactionResponse})
}

func FindTransactions(c *gin.Context) {
	// define variable array object
	var transaction []models.Transaction

	// query get  all transaction
	results := models.DB.Debug().Find(&transaction)

	if results.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "count": len(transaction), "data": transaction})
}

func FindTransactionByID(c *gin.Context) {
	// get id from segment url
	id := c.Param("id")

	// query get transaction by id
	results, err := models.GetTransactionById(models.StrToUint(id))

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": results})
}

func convertToTransactionResponse(trx models.Transaction) models.ResponseTransaction {
	return models.ResponseTransaction{
		ID:       trx.ID,
		Sender:   trx.Status,
		Receiver: trx.Status,
		Amount:   trx.Amount,
		Status:   trx.Status,
	}
}

func getCustomFieldTransaction(trx_id uint) (models.ResponseTransaction, error) {

	var transactionResponse models.ResponseTransaction

	// get detail transaction
	if err := models.DB.Debug().Unscoped().Select("transactions.id, u1.username sender, u2.username receiver, amount, status").Joins("JOIN users u1 on u1.id = transactions.sender_user_id").Joins("JOIN users u2 on u2.id = transactions.receiver_user_id").Where("transactions.id = ?", trx_id).Table("transactions").Find(&transactionResponse).Error; err != nil {
		return transactionResponse, errors.New("Transaction not found!")
	}

	return transactionResponse, nil
}
