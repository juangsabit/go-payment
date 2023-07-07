package controllers

import (
	"net/http"

	"go-payment/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FindUsers(c *gin.Context) {

	var users []models.User
	results := models.DB.Debug().Find(&users)
	if results.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	var response []models.ResponseUser
	for _, user := range users {
		userResponse := convertToUserResponse(user)
		response = append(response, userResponse)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "count": len(users), "data": response})

}

func FindUserById(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := models.DB.Debug().First(&user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
	}

	response := convertToUserResponse(user)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func CreateUser(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	models.DB.Debug().Create(&user)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if models.DB.Debug().Model(&user).Where("id = ?", id).Updates(&user).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "tidak dapat mengupdate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data berhasil diperbarui"})
}

func DeleteUser(c *gin.Context) {

	var user models.User

	id := c.Param("id")
	if models.DB.Debug().Delete(&user, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data berhasil dihapus"})
}

func convertToUserResponse(user models.User) models.ResponseUser {
	return models.ResponseUser{
		ID:       user.ID,
		Fullname: user.FullName,
		Username: user.Username,
		Email:    user.Email,
	}
}
