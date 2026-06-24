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

// List godoc
// @Summary      List properties with optional filters
// @Tags         properties
// @Produce      json
// @Param        q                query  string  false  "Free-text search"
// @Param        country          query  string  false  "Country"
// @Param        city             query  string  false  "City"
// @Param        property_type    query  string  false  "Property type (Apartment, Studio, Bedsitter…)"
// @Param        transaction_type query  string  false  "rent or sale"
// @Param        min_price        query  number  false  "Minimum price"
// @Param        max_price        query  number  false  "Maximum price"
// @Param        page             query  int     false  "Page number (default 1)"
// @Success      200  {object}  PropertyListResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /properties [get]
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

// GetFeatured godoc
// @Summary      Get featured properties
// @Tags         properties
// @Produce      json
// @Success      200  {object}  PropertyListResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /properties/featured [get]
func (h *PropertyHandler) GetFeatured(c *gin.Context) {
	props, err := h.service.GetFeatured()
	if handleErr(c, err) {
		return
	}
	utils.OK(c, props)
}

// Get godoc
// @Summary      Get a single property by ID
// @Tags         properties
// @Produce      json
// @Param        id  path  int  true  "Property ID"
// @Success      200  {object}  PropertyResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /properties/{id} [get]
func (h *PropertyHandler) Get(c *gin.Context) {
	id := paramUint(c, "id")
	userID, loggedIn := middleware.GetUserID(c)

	prop, err := h.service.GetProperty(id, userID, loggedIn)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, prop)
}

// Create godoc
// @Summary      Create a new property listing
// @Tags         properties
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        title            formData  string   true   "Title"
// @Param        description      formData  string   false  "Description"
// @Param        property_type    formData  string   true   "Bedsitter | Studio | Apartment | Maisonette | Bungalow | Townhouse | Villa | Commercial"
// @Param        transaction_type formData  string   true   "rent | sale"
// @Param        country          formData  string   true   "Country"
// @Param        city             formData  string   true   "City"
// @Param        state            formData  string   false  "State / region"
// @Param        neighborhood     formData  string   false  "Neighborhood"
// @Param        exact_address    formData  string   false  "Exact address"
// @Param        latitude         formData  number   false  "Latitude"
// @Param        longitude        formData  number   false  "Longitude"
// @Param        rooms            formData  integer  true   "Number of rooms"
// @Param        bathrooms        formData  integer  true   "Number of bathrooms"
// @Param        shower_type      formData  string   false  "en_suite | shared"
// @Param        surface          formData  number   false  "Surface area in m²"
// @Param        furnished        formData  boolean  false  "Is furnished"
// @Param        has_wifi         formData  boolean  false  "Has Wi-Fi"
// @Param        has_water        formData  boolean  false  "Has water"
// @Param        has_electricity  formData  boolean  false  "Has electricity"
// @Param        has_courtyard    formData  boolean  false  "Has courtyard"
// @Param        price            formData  number   true   "Price"
// @Param        currency         formData  string   true   "ISO 4217 currency code (e.g. XAF, USD)"
// @Param        images           formData  file     false  "Property images (multiple)"
// @Success      201  {object}  PropertyResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /properties [post]
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

// Update godoc
// @Summary      Update a property listing
// @Tags         properties
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id               path      int      true   "Property ID"
// @Param        title            formData  string   true   "Title"
// @Param        description      formData  string   false  "Description"
// @Param        property_type    formData  string   true   "Property type"
// @Param        transaction_type formData  string   true   "rent | sale"
// @Param        country          formData  string   true   "Country"
// @Param        city             formData  string   true   "City"
// @Param        state            formData  string   false  "State / region"
// @Param        neighborhood     formData  string   false  "Neighborhood"
// @Param        exact_address    formData  string   false  "Exact address"
// @Param        latitude         formData  number   false  "Latitude"
// @Param        longitude        formData  number   false  "Longitude"
// @Param        rooms            formData  integer  true   "Number of rooms"
// @Param        bathrooms        formData  integer  true   "Number of bathrooms"
// @Param        shower_type      formData  string   false  "en_suite | shared"
// @Param        surface          formData  number   false  "Surface area in m²"
// @Param        furnished        formData  boolean  false  "Is furnished"
// @Param        has_wifi         formData  boolean  false  "Has Wi-Fi"
// @Param        has_water        formData  boolean  false  "Has water"
// @Param        has_electricity  formData  boolean  false  "Has electricity"
// @Param        has_courtyard    formData  boolean  false  "Has courtyard"
// @Param        price            formData  number   true   "Price"
// @Param        currency         formData  string   true   "ISO 4217 currency code"
// @Success      200  {object}  PropertyResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /properties/{id} [put]
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

// Delete godoc
// @Summary      Delete a property listing
// @Tags         properties
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Property ID"
// @Success      200  {object}  MessageResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /properties/{id} [delete]
func (h *PropertyHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id := paramUint(c, "id")

	if handleErr(c, h.service.DeleteProperty(id, userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Property deleted"})
}

// ToggleAvailability godoc
// @Summary      Toggle a property's availability
// @Tags         properties
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Property ID"
// @Success      200  {object}  AvailabilityResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /properties/{id}/availability [patch]
func (h *PropertyHandler) ToggleAvailability(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id := paramUint(c, "id")

	newVal, err := h.service.ToggleAvailability(id, userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, gin.H{"is_available": newVal})
}

// MyProperties godoc
// @Summary      List the current user's property listings
// @Tags         properties
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  PropertyListResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /properties/my/listings [get]
func (h *PropertyHandler) MyProperties(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	props, err := h.service.GetMyProperties(userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, props)
}

// AddImage godoc
// @Summary      Add an image to a property
// @Tags         properties
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int   true  "Property ID"
// @Param        image  formData  file  true  "Image file"
// @Success      201  {object}  PropertyImageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Router       /properties/{id}/images [post]
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

// DeleteImage godoc
// @Summary      Delete a property image
// @Tags         images
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Image ID"
// @Success      200  {object}  MessageResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /images/{id} [delete]
func (h *PropertyHandler) DeleteImage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	imgID := paramUint(c, "id")

	if handleErr(c, h.service.DeleteImage(imgID, userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Image deleted"})
}

// SetMainImage godoc
// @Summary      Set an image as the main property photo
// @Tags         images
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Image ID"
// @Success      200  {object}  MessageResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /images/{id}/main [patch]
func (h *PropertyHandler) SetMainImage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	imgID := paramUint(c, "id")

	if handleErr(c, h.service.SetMainImage(imgID, userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Main image updated"})
}
