package handlers

import (
	"strconv"

	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/repositories"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type PropertyHandler struct {
	service services.PropertyService
}

func NewPropertyHandler(svc services.PropertyService) *PropertyHandler {
	return &PropertyHandler{service: svc}
}

func (h *PropertyHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	userID, loggedIn := middleware.GetUserID(c)

	props, total, err := h.service.ListProperties(repositories.PropertyFilters{
		Q:               c.Query("q"),
		Country:         c.Query("country"),
		City:            c.Query("city"),
		PropertyType:    c.Query("property_type"),
		TransactionType: c.Query("transaction_type"),
		MinPrice:        c.Query("min_price"),
		MaxPrice:        c.Query("max_price"),
	}, page, userID, loggedIn)
	if handleErr(c, err) {
		return
	}
	pageSize := 10
	if page < 1 {
		page = 1
	}
	utils.Paginated(c, props, total, page, pageSize)
}

func (h *PropertyHandler) GetFeatured(c *gin.Context) {
	props, err := h.service.GetFeatured()
	if handleErr(c, err) {
		return
	}
	utils.OK(c, props)
}

func (h *PropertyHandler) Get(c *gin.Context) {
	id := paramUint(c, "id")
	userID, loggedIn := middleware.GetUserID(c)

	prop, err := h.service.GetProperty(id, userID, loggedIn)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, prop)
}

func (h *PropertyHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input struct {
		Title           string  `form:"title" binding:"required"`
		Description     string  `form:"description"`
		PropertyType    string  `form:"property_type" binding:"required"`
		TransactionType string  `form:"transaction_type" binding:"required"`
		Country         string  `form:"country" binding:"required"`
		City            string  `form:"city" binding:"required"`
		State           string  `form:"state"`
		Neighborhood    string  `form:"neighborhood"`
		ExactAddress    string  `form:"exact_address"`
		Latitude        float64 `form:"latitude"`
		Longitude       float64 `form:"longitude"`
		Rooms           int     `form:"rooms" binding:"required"`
		Bathrooms       int     `form:"bathrooms" binding:"required"`
		ShowerType      string  `form:"shower_type"`
		Surface         float64 `form:"surface"`
		Furnished       bool    `form:"furnished"`
		HasWifi         bool    `form:"has_wifi"`
		HasWater        bool    `form:"has_water"`
		HasElectricity  bool    `form:"has_electricity"`
		HasCourtyard    bool    `form:"has_courtyard"`
		Price           float64 `form:"price" binding:"required"`
		Currency        string  `form:"currency" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	form, _ := c.MultipartForm()
	prop, err := h.service.CreateProperty(services.CreatePropertyInput{
		Title: input.Title, Description: input.Description,
		PropertyType: input.PropertyType, TransactionType: input.TransactionType,
		Country: input.Country, City: input.City, State: input.State,
		Neighborhood: input.Neighborhood, ExactAddress: input.ExactAddress,
		Latitude: input.Latitude, Longitude: input.Longitude,
		Rooms: input.Rooms, Bathrooms: input.Bathrooms,
		ShowerType: input.ShowerType, Surface: input.Surface,
		Furnished: input.Furnished, HasWifi: input.HasWifi,
		HasWater: input.HasWater, HasElectricity: input.HasElectricity,
		HasCourtyard: input.HasCourtyard, Price: input.Price, Currency: input.Currency,
		Images: form.File["images"],
	}, userID)
	if handleErr(c, err) {
		return
	}
	utils.Created(c, prop)
}

func (h *PropertyHandler) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id := paramUint(c, "id")

	var input struct {
		Title           string  `form:"title" binding:"required"`
		Description     string  `form:"description"`
		PropertyType    string  `form:"property_type" binding:"required"`
		TransactionType string  `form:"transaction_type" binding:"required"`
		Country         string  `form:"country" binding:"required"`
		City            string  `form:"city" binding:"required"`
		State           string  `form:"state"`
		Neighborhood    string  `form:"neighborhood"`
		ExactAddress    string  `form:"exact_address"`
		Latitude        float64 `form:"latitude"`
		Longitude       float64 `form:"longitude"`
		Rooms           int     `form:"rooms" binding:"required"`
		Bathrooms       int     `form:"bathrooms" binding:"required"`
		ShowerType      string  `form:"shower_type"`
		Surface         float64 `form:"surface"`
		Furnished       bool    `form:"furnished"`
		HasWifi         bool    `form:"has_wifi"`
		HasWater        bool    `form:"has_water"`
		HasElectricity  bool    `form:"has_electricity"`
		HasCourtyard    bool    `form:"has_courtyard"`
		Price           float64 `form:"price" binding:"required"`
		Currency        string  `form:"currency" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	prop, err := h.service.UpdateProperty(id, services.UpdatePropertyInput{
		Title: input.Title, Description: input.Description,
		PropertyType: input.PropertyType, TransactionType: input.TransactionType,
		Country: input.Country, City: input.City, State: input.State,
		Neighborhood: input.Neighborhood, ExactAddress: input.ExactAddress,
		Latitude: input.Latitude, Longitude: input.Longitude,
		Rooms: input.Rooms, Bathrooms: input.Bathrooms,
		ShowerType: input.ShowerType, Surface: input.Surface,
		Furnished: input.Furnished, HasWifi: input.HasWifi,
		HasWater: input.HasWater, HasElectricity: input.HasElectricity,
		HasCourtyard: input.HasCourtyard, Price: input.Price, Currency: input.Currency,
	}, userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, prop)
}

func (h *PropertyHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id := paramUint(c, "id")

	if handleErr(c, h.service.DeleteProperty(id, userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Property deleted"})
}

func (h *PropertyHandler) ToggleAvailability(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id := paramUint(c, "id")

	newVal, err := h.service.ToggleAvailability(id, userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, gin.H{"is_available": newVal})
}

func (h *PropertyHandler) MyProperties(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	props, err := h.service.GetMyProperties(userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, props)
}

func (h *PropertyHandler) AddImage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	propID := paramUint(c, "id")

	file, err := c.FormFile("image")
	if err != nil {
		utils.BadRequest(c, "No image provided")
		return
	}

	img, err := h.service.AddImage(propID, userID, file)
	if handleErr(c, err) {
		return
	}
	utils.Created(c, img)
}

func (h *PropertyHandler) DeleteImage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	imgID := paramUint(c, "id")

	if handleErr(c, h.service.DeleteImage(imgID, userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Image deleted"})
}

func (h *PropertyHandler) SetMainImage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	imgID := paramUint(c, "id")

	if handleErr(c, h.service.SetMainImage(imgID, userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Main image updated"})
}
