package handlers

import (
	"strconv"
	"strings"

	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type PropertyInput struct {
	Title           string  `json:"title" form:"title" binding:"required"`
	Description     string  `json:"description" form:"description"`
	PropertyType    string  `json:"property_type" form:"property_type" binding:"required"`
	TransactionType string  `json:"transaction_type" form:"transaction_type" binding:"required"`
	Country         string  `json:"country" form:"country" binding:"required"`
	City            string  `json:"city" form:"city" binding:"required"`
	State           string  `json:"state" form:"state"`
	Neighborhood    string  `json:"neighborhood" form:"neighborhood"`
	ExactAddress    string  `json:"exact_address" form:"exact_address"`
	Latitude        float64 `json:"latitude" form:"latitude"`
	Longitude       float64 `json:"longitude" form:"longitude"`
	Rooms           int     `json:"rooms" form:"rooms" binding:"required"`
	Bathrooms       int     `json:"bathrooms" form:"bathrooms" binding:"required"`
	ShowerType      string  `json:"shower_type" form:"shower_type"`
	Surface         float64 `json:"surface" form:"surface"`
	Furnished       bool    `json:"furnished" form:"furnished"`
	HasWifi         bool    `json:"has_wifi" form:"has_wifi"`
	HasWater        bool    `json:"has_water" form:"has_water"`
	HasElectricity  bool    `json:"has_electricity" form:"has_electricity"`
	HasCourtyard    bool    `json:"has_courtyard" form:"has_courtyard"`
	Price           float64 `json:"price" form:"price" binding:"required"`
	Currency        string  `json:"currency" form:"currency" binding:"required"`
}

func ListProperties(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 10
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.Property{}).Where("is_available = ?", true)

	if q := c.Query("q"); q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where(
			"LOWER(city) LIKE ? OR LOWER(neighborhood) LIKE ? OR LOWER(country) LIKE ? OR LOWER(state) LIKE ? OR LOWER(exact_address) LIKE ? OR LOWER(title) LIKE ?",
			like, like, like, like, like, like,
		)
	}
	if country := c.Query("country"); country != "" {
		query = query.Where("country = ?", country)
	}
	if city := c.Query("city"); city != "" {
		query = query.Where("city = ?", city)
	}
	if pt := c.Query("property_type"); pt != "" {
		query = query.Where("property_type = ?", pt)
	}
	if tt := c.Query("transaction_type"); tt != "" {
		query = query.Where("transaction_type = ?", tt)
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}

	var total int64
	query.Count(&total)

	var props []models.Property
	query.Offset(offset).Limit(pageSize).Preload("Images").Preload("Owner").Find(&props)

	userID, loggedIn := middleware.GetUserID(c)
	for i := range props {
		props[i].ComputeRating()
		if loggedIn {
			var fav models.Favorite
			if err := database.DB.Where("user_id = ? AND property_id = ?", userID, props[i].ID).First(&fav).Error; err == nil {
				props[i].IsFavorited = true
			}
		}
	}

	utils.Paginated(c, props, total, page, pageSize)
}

func GetFeaturedProperties(c *gin.Context) {
	var props []models.Property
	database.DB.Where("is_available = ?", true).
		Preload("Images").
		Preload("Reviews").
		Find(&props)

	type ratedProp struct {
		prop   models.Property
		rating float64
	}
	rated := make([]ratedProp, len(props))
	for i := range props {
		props[i].ComputeRating()
		rated[i] = ratedProp{props[i], props[i].AverageRating}
	}

	for i := 0; i < len(rated); i++ {
		for j := i + 1; j < len(rated); j++ {
			if rated[j].rating > rated[i].rating {
				rated[i], rated[j] = rated[j], rated[i]
			}
		}
	}

	result := make([]models.Property, 0, 4)
	for i := 0; i < len(rated) && i < 4; i++ {
		result = append(result, rated[i].prop)
	}

	utils.OK(c, result)
}

func GetProperty(c *gin.Context) {
	id := c.Param("id")
	var prop models.Property

	if err := database.DB.Preload("Images").Preload("Owner").Preload("Reviews.User").First(&prop, id).Error; err != nil {
		utils.NotFound(c, "Property not found")
		return
	}

	prop.ComputeRating()

	userID, loggedIn := middleware.GetUserID(c)
	if loggedIn {
		var fav models.Favorite
		if err := database.DB.Where("user_id = ? AND property_id = ?", userID, prop.ID).First(&fav).Error; err == nil {
			prop.IsFavorited = true
		}
	}

	utils.OK(c, prop)
}

func CreateProperty(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input PropertyInput
	if err := c.ShouldBind(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["images"]
	if len(files) < 3 {
		utils.BadRequest(c, "At least 3 images are required")
		return
	}

	prop := models.Property{
		Title:           input.Title,
		Description:     input.Description,
		PropertyType:    input.PropertyType,
		TransactionType: input.TransactionType,
		Country:         input.Country,
		City:            input.City,
		State:           input.State,
		Neighborhood:    input.Neighborhood,
		ExactAddress:    input.ExactAddress,
		Latitude:        input.Latitude,
		Longitude:       input.Longitude,
		Rooms:           input.Rooms,
		Bathrooms:       input.Bathrooms,
		ShowerType:      input.ShowerType,
		Surface:         input.Surface,
		Furnished:       input.Furnished,
		HasWifi:         input.HasWifi,
		HasWater:        input.HasWater,
		HasElectricity:  input.HasElectricity,
		HasCourtyard:    input.HasCourtyard,
		Price:           input.Price,
		Currency:        input.Currency,
		OwnerID:         userID,
		IsAvailable:     true,
	}

	if err := database.DB.Create(&prop).Error; err != nil {
		utils.InternalError(c, "Failed to create property")
		return
	}

	savedCount := 0
	for _, fh := range files {
		url, err := saveFile(c, fh)
		if err != nil {
			database.DB.Delete(&prop)
			utils.InternalError(c, "Failed to upload image: "+err.Error())
			return
		}
		img := models.PropertyImage{
			Filename:   url,
			PropertyID: prop.ID,
			IsMain:     savedCount == 0,
		}
		if err := database.DB.Create(&img).Error; err != nil {
			database.DB.Delete(&prop)
			utils.InternalError(c, "Failed to save image record")
			return
		}
		savedCount++
	}

	database.DB.Preload("Images").First(&prop, prop.ID)
	utils.Created(c, prop)

	go notifyMatchingAlerts(prop)
}

func UpdateProperty(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id := c.Param("id")

	var prop models.Property
	if err := database.DB.First(&prop, id).Error; err != nil {
		utils.NotFound(c, "Property not found")
		return
	}
	if prop.OwnerID != userID {
		utils.Forbidden(c, "Not authorized")
		return
	}

	var input PropertyInput
	if err := c.ShouldBind(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{
		"title":            input.Title,
		"description":      input.Description,
		"property_type":    input.PropertyType,
		"transaction_type": input.TransactionType,
		"country":          input.Country,
		"city":             input.City,
		"state":            input.State,
		"neighborhood":     input.Neighborhood,
		"exact_address":    input.ExactAddress,
		"latitude":         input.Latitude,
		"longitude":        input.Longitude,
		"rooms":            input.Rooms,
		"bathrooms":        input.Bathrooms,
		"shower_type":      input.ShowerType,
		"surface":          input.Surface,
		"furnished":        input.Furnished,
		"has_wifi":         input.HasWifi,
		"has_water":        input.HasWater,
		"has_electricity":  input.HasElectricity,
		"has_courtyard":    input.HasCourtyard,
		"price":            input.Price,
		"currency":         input.Currency,
	}

	database.DB.Model(&prop).Updates(updates)
	database.DB.Preload("Images").First(&prop, prop.ID)
	utils.OK(c, prop)
}

func DeleteProperty(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id := c.Param("id")

	var prop models.Property
	if err := database.DB.Preload("Images").First(&prop, id).Error; err != nil {
		utils.NotFound(c, "Property not found")
		return
	}
	if prop.OwnerID != userID {
		utils.Forbidden(c, "Not authorized")
		return
	}

	for _, img := range prop.Images {
		deleteStorageFile(img.Filename)
	}

	database.DB.Select("Images", "Rentals", "Reviews", "Favorites").Delete(&prop)
	utils.OK(c, gin.H{"message": "Property deleted"})
}

func AddPropertyImage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	propID := c.Param("id")

	var prop models.Property
	if err := database.DB.First(&prop, propID).Error; err != nil {
		utils.NotFound(c, "Property not found")
		return
	}
	if prop.OwnerID != userID {
		utils.Forbidden(c, "Not authorized")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		utils.BadRequest(c, "No image provided")
		return
	}

	url, err := saveFile(c, file)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var imageCount int64
	database.DB.Model(&models.PropertyImage{}).Where("property_id = ?", prop.ID).Count(&imageCount)

	img := models.PropertyImage{
		Filename:   url,
		PropertyID: prop.ID,
		IsMain:     imageCount == 0,
	}
	if err := database.DB.Create(&img).Error; err != nil {
		utils.InternalError(c, "Failed to save image record")
		return
	}
	utils.Created(c, img)
}

func DeletePropertyImage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	imgID := c.Param("id")

	var img models.PropertyImage
	if err := database.DB.First(&img, imgID).Error; err != nil {
		utils.NotFound(c, "Image not found")
		return
	}

	var prop models.Property
	database.DB.First(&prop, img.PropertyID)
	if prop.OwnerID != userID {
		utils.Forbidden(c, "Not authorized")
		return
	}

	deleteStorageFile(img.Filename)
	database.DB.Delete(&img)
	utils.OK(c, gin.H{"message": "Image deleted"})
}

func SetMainImage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	imgID := c.Param("id")

	var img models.PropertyImage
	if err := database.DB.First(&img, imgID).Error; err != nil {
		utils.NotFound(c, "Image not found")
		return
	}

	var prop models.Property
	database.DB.First(&prop, img.PropertyID)
	if prop.OwnerID != userID {
		utils.Forbidden(c, "Not authorized")
		return
	}

	database.DB.Model(&models.PropertyImage{}).Where("property_id = ?", img.PropertyID).Update("is_main", false)
	database.DB.Model(&img).Update("is_main", true)
	utils.OK(c, gin.H{"message": "Main image updated"})
}

func ToggleAvailability(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id := c.Param("id")

	var prop models.Property
	if err := database.DB.First(&prop, id).Error; err != nil {
		utils.NotFound(c, "Property not found")
		return
	}
	if prop.OwnerID != userID {
		utils.Forbidden(c, "Not authorized")
		return
	}

	database.DB.Model(&prop).Update("is_available", !prop.IsAvailable)
	utils.OK(c, gin.H{"is_available": !prop.IsAvailable})
}

func MyProperties(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var props []models.Property
	database.DB.Where("owner_id = ?", userID).Preload("Images").Find(&props)
	for i := range props {
		props[i].ComputeRating()
	}
	utils.OK(c, props)
}
