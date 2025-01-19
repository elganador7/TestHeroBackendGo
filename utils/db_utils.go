package utils

import (
	"TestHeroBackendGo/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetQuestionByID(questionId string, question models.Question, db *gorm.DB, c *gin.Context) (models.Question, *gin.Context, bool, error) {
	if err := db.First(&question, "id = ?", questionId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return question, c, false, nil
		} else {
			return question, c, false, err
		}
	}

	return question, c, true, nil
}
