package handlers

import (
	"cost-classifier/backend/db"
	"cost-classifier/backend/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// Handle creating a new label
func CreateNewLabel(c echo.Context) error {
	name := c.QueryParam("name")
	color := c.QueryParam("color")

	label := models.Label{Name: name, Color: color}
	result := db.DB.Create(&label)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create label",
		})
	}

	return c.JSON(http.StatusOK, label)
}

func EditLabel(c echo.Context) error {
	// Get the label name from the query parameter
	oldName := c.QueryParam("old_name")

	// Retrieve the existing label
	var label models.Label
	result := db.DB.Where("name = ?", oldName).First(&label)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Label not found",
		})
	}

	// Parse the new label data from the request body
	var newLabel models.Label
	if err := c.Bind(&newLabel); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	// Update the label with new data
	label.Name = newLabel.Name
	label.Color = newLabel.Color
	db.DB.Save(&label)

	return c.JSON(http.StatusOK, label)
}

func DeleteLabel(c echo.Context) error {
	// Get the label name from the query parameter
	name := c.QueryParam("name")

	// Delete the label
	result := db.DB.Where("name = ?", name).Delete(&models.Label{})
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete label",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Label deleted successfully",
	})
}

// Handle getting a list of labels
func GetLabels(c echo.Context) error {
	var labels []models.Label
	var result *gorm.DB

	countParam := c.QueryParam("count")
	if countParam != "" {
		// If count parameter is provided, retrieve the specified number of labels
		count, err := strconv.Atoi(countParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid count parameter",
			})
		}
		result = db.DB.Limit(count).Find(&labels)
	} else {
		// If count parameter is not provided, retrieve all labels
		result = db.DB.Find(&labels)
	}

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve labels",
		})
	}

	return c.JSON(http.StatusOK, labels)
}
