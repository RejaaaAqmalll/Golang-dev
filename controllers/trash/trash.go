package trash

import (
	"net/http"
	"nyoba/configg"
	"nyoba/models"

	"github.com/gin-gonic/gin"
)

type respon struct {
	EmployeeId     string
	DivisionId     int
	Division       models.Division
	JobPositionId  int
	JobPosition    models.JobPosition
	JobLevelId     int
	JobLevel       models.JobLevel
	EmployeeStatus string
	PersonalId     int
	PersonalInfo   models.PersonalInfo `gorm:"foreignkey:PersonalId"`
}

func (respon) TableName() string {
	return "employees"
}

func DetailEmp(c *gin.Context) {
	var detail respon

	id := c.Param("id")

	if err := configg.KoneksiData().Preload("Division").Preload("JobPosition").
		Preload("JobLevel").Preload("PersonalInfo").
		Where("employee_id = ?", id).
		Find(&detail).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Failed Query",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": detail,
	})
}
