package controllers

import (
	"TestHeroBackendGo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Controller structure that holds the GORM DB instance
type TestTopicDataController struct {
	DB *gorm.DB
}

// CreateTestTopicData handles creating a new TestTopicData entry
func (controller *TestTopicDataController) CreateTestTopicData(c *gin.Context) {
	var testTopicData models.TestTopicData
	if err := c.ShouldBindJSON(&testTopicData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := controller.DB.Create(&testTopicData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, testTopicData)
}

// GetTestTopicData handles fetching a single TestTopicData by ID
func (controller *TestTopicDataController) GetTestTopicData(c *gin.Context) {
	id := c.Param("id")
	var testTopicData models.TestTopicData
	if err := controller.DB.First(&testTopicData, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TestTopicData not found"})
		return
	}
	c.JSON(http.StatusOK, testTopicData)
}

// UpdateTestTopicData handles updating a TestTopicData entry by ID
func (controller *TestTopicDataController) UpdateTestTopicData(c *gin.Context) {
	id := c.Param("id")
	var testTopicData models.TestTopicData
	if err := controller.DB.First(&testTopicData, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TestTopicData not found"})
		return
	}

	// Bind updated data
	if err := c.ShouldBindJSON(&testTopicData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := controller.DB.Save(&testTopicData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, testTopicData)
}

// DeleteTestTopicData handles deleting a TestTopicData entry by ID
func (controller *TestTopicDataController) DeleteTestTopicData(c *gin.Context) {
	id := c.Param("id")
	if err := controller.DB.Delete(&models.TestTopicData{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TestTopicData not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "TestTopicData deleted"})
}

// ListTestTopicData handles listing all TestTopicData entries
func (controller *TestTopicDataController) ListTestTopicData(c *gin.Context) {
	var testTopicData []models.TestTopicData
	if err := controller.DB.Find(&testTopicData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, testTopicData)
}

// ListByTestType handles listing all TestTopicData entries matching a specific TestType
func (controller *TestTopicDataController) ListByTestType(c *gin.Context) {
	testType := c.Param("test_type")
	var testTopicData []models.TestTopicData
	if err := controller.DB.Where("test_type = ?", testType).Find(&testTopicData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, testTopicData)
}

// ListBySubject handles listing all TestTopicData entries matching a specific Subject
func (controller *TestTopicDataController) ListBySubject(c *gin.Context) {
	subject := c.Param("subject")
	var testTopicData []models.TestTopicData
	if err := controller.DB.Where("subject = ?", subject).Find(&testTopicData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, testTopicData)
}

// ListByTopic handles listing all TestTopicData entries matching a specific Topic
func (controller *TestTopicDataController) ListByTopic(c *gin.Context) {
	topic := c.Param("topic")
	var testTopicData []models.TestTopicData
	if err := controller.DB.Where("topic = ?", topic).Find(&testTopicData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, testTopicData)
}

// ListBySubtopic handles listing all TestTopicData entries matching a specific Subtopic
func (controller *TestTopicDataController) ListBySubtopic(c *gin.Context) {
	subtopic := c.Param("subtopic")
	var testTopicData []models.TestTopicData
	if err := controller.DB.Where("subtopic = ?", subtopic).Find(&testTopicData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, testTopicData)
}

// ListBySpecificTopic handles listing all TestTopicData entries matching a specific SpecificTopic
func (controller *TestTopicDataController) ListBySpecificTopic(c *gin.Context) {
	specificTopic := c.Param("specific_topic")
	var testTopicData []models.TestTopicData
	if err := controller.DB.Where("specific_topic = ?", specificTopic).Find(&testTopicData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, testTopicData)
}
