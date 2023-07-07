package controllers

import (
	"go-payment/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FindProducts(c *gin.Context) {
	var products []models.Product
	results := models.DB.Debug().Find(&products)
	if results.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "count": len(products), "data": products})
}

type Activity struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Username    string `gorm:"type:varchar(255)" json:"username"`
	Information string `gorm:"type:varchar(255)" json:"information"`
	CreatedAt   *time.Time
}

func FindProductById(c *gin.Context) {
	var product models.Product
	var activity []models.ResponseActivity

	id := c.Param("id")

	if err := models.DB.Debug().Unscoped().First(&product, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
	}

	// models.DB.Debug().Unscoped() for ignore where deleteAt is null
	if err := models.DB.Debug().Unscoped().Select("activities.id, username, information, activities.created_at").Joins("JOIN users on users.id = activities.user_id").Where("table_name = ? AND row_id = ?", "product", id).Table("activities").Find(&activity).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
	}

	var productResponse models.ResponseProduct

	productResponse.ID = product.ID
	productResponse.ProductName = product.ProductName
	productResponse.Description = product.Description
	productResponse.Price = product.Price
	productResponse.CreatedAt = product.CreatedAt

	var response []models.ResponseActivity
	for _, activity := range activity {
		activityResponse := convertToActivityResponse(activity)
		response = append(response, activityResponse)
	}

	productResponse.Activity = response

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": productResponse})
}

func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	models.DB.Debug().Create(&product)

	str_id, _ := c.Get("userID")
	userID := str_id.(uint)

	models.SaveActivity(userID, "create product", "product", product.ID)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Product has been created"})
}

func UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if models.DB.Debug().Model(&product).Where("id = ?", id).Updates(&product).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Product not found"})
		return
	}

	str_id, _ := c.Get("userID")
	userID := str_id.(uint)

	models.SaveActivity(userID, "update product", "product", models.StrToUint(id))
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Product has been updated"})
}

func DeleteProduct(c *gin.Context) {

	var product models.Product

	id := c.Param("id")
	if models.DB.Debug().Delete(&product, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Data not found"})
		return
	}

	str_id, _ := c.Get("userID")
	userID := str_id.(uint)

	models.SaveActivity(userID, "deleted product", "product", models.StrToUint(id))
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Success deleted product"})
}

func convertToActivityResponse(activity models.ResponseActivity) models.ResponseActivity {
	return models.ResponseActivity{
		ID:          activity.ID,
		Username:    activity.Username,
		Information: activity.Information,
		CreatedAt:   activity.CreatedAt,
	}
}
