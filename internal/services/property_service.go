package services

import (
	"fmt"
	"mime/multipart"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/repositories"
	"github.com/badersalis/gidana_backend/internal/storage"
	"github.com/badersalis/gidana_backend/internal/utils"
	appws "github.com/badersalis/gidana_backend/internal/ws"
)

type CreatePropertyInput struct {
	Title           string
	Description     string
	PropertyType    string
	TransactionType string
	Country         string
	City            string
	State           string
	Neighborhood    string
	ExactAddress    string
	Latitude        float64
	Longitude       float64
	Rooms           int
	Bathrooms       int
	ShowerType      string
	Surface         float64
	Furnished       bool
	HasWifi         bool
	HasWater        bool
	HasElectricity  bool
	HasCourtyard    bool
	Price           float64
	Currency        string
	Images          []*multipart.FileHeader
}

type UpdatePropertyInput struct {
	Title           string
	Description     string
	PropertyType    string
	TransactionType string
	Country         string
	City            string
	State           string
	Neighborhood    string
	ExactAddress    string
	Latitude        float64
	Longitude       float64
	Rooms           int
	Bathrooms       int
	ShowerType      string
	Surface         float64
	Furnished       bool
	HasWifi         bool
	HasWater        bool
	HasElectricity  bool
	HasCourtyard    bool
	Price           float64
	Currency        string
}

type PropertyService interface {
	ListProperties(filters repositories.PropertyFilters, page int, userID uint, loggedIn bool) ([]models.Property, int64, error)
	GetProperty(id uint, userID uint, loggedIn bool) (*models.Property, error)
	GetFeatured() ([]models.Property, error)
	GetMyProperties(userID uint) ([]models.Property, error)
	CreateProperty(input CreatePropertyInput, userID uint) (*models.Property, error)
	UpdateProperty(id uint, input UpdatePropertyInput, userID uint) (*models.Property, error)
	DeleteProperty(id uint, userID uint) error
	ToggleAvailability(id uint, userID uint) (bool, error)
	AddImage(propID uint, userID uint, file *multipart.FileHeader) (*models.PropertyImage, error)
	DeleteImage(imgID uint, userID uint) error
	SetMainImage(imgID uint, userID uint) error
}

type propertyService struct {
	propRepo  repositories.PropertyRepository
	imageRepo repositories.PropertyImageRepository
	favRepo   repositories.FavoriteRepository
	alertRepo repositories.AlertRepository
	userRepo  repositories.UserRepository
	fileStore storage.FileStorage
	hub       *appws.Hub
}

func NewPropertyService(
	propRepo repositories.PropertyRepository,
	imageRepo repositories.PropertyImageRepository,
	favRepo repositories.FavoriteRepository,
	alertRepo repositories.AlertRepository,
	userRepo repositories.UserRepository,
	fileStore storage.FileStorage,
	hub *appws.Hub,
) PropertyService {
	return &propertyService{
		propRepo:  propRepo,
		imageRepo: imageRepo,
		favRepo:   favRepo,
		alertRepo: alertRepo,
		userRepo:  userRepo,
		fileStore: fileStore,
		hub:       hub,
	}
}

func (s *propertyService) ListProperties(filters repositories.PropertyFilters, page int, userID uint, loggedIn bool) ([]models.Property, int64, error) {
	if page < 1 {
		page = 1
	}
	pageSize := 10
	offset := (page - 1) * pageSize

	props, total, err := s.propRepo.List(filters, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	for i := range props {
		props[i].ComputeRating()
		if loggedIn {
			if exists, _ := s.favRepo.ExistsForUserAndProperty(userID, props[i].ID); exists {
				props[i].IsFavorited = true
			}
		}
	}
	return props, total, nil
}

func (s *propertyService) GetProperty(id uint, userID uint, loggedIn bool) (*models.Property, error) {
	prop, err := s.propRepo.GetByID(id)
	if err != nil {
		return nil, ErrNotFound("Property not found")
	}
	prop.ComputeRating()
	if loggedIn {
		if exists, _ := s.favRepo.ExistsForUserAndProperty(userID, prop.ID); exists {
			prop.IsFavorited = true
		}
		// Expose contact info only to the owner or seeker with "pro" plan
		viewer, _ := s.userRepo.GetByID(userID)
		if viewer != nil && (prop.OwnerID == userID || viewer.SubscriptionPlan == "pro") {
			phone := prop.PhoneContact
			whatsapp := prop.WhatsappContact
			prop.ContactPhone = &phone
			prop.ContactWhatsapp = &whatsapp
		}
	}
	return prop, nil
}

func (s *propertyService) GetFeatured() ([]models.Property, error) {
	props, err := s.propRepo.GetFeatured()
	if err != nil {
		return nil, err
	}

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
	return result, nil
}

func (s *propertyService) GetMyProperties(userID uint) ([]models.Property, error) {
	props, err := s.propRepo.GetByOwnerID(userID)
	if err != nil {
		return nil, err
	}
	for i := range props {
		props[i].ComputeRating()
	}
	return props, nil
}

func (s *propertyService) CreateProperty(input CreatePropertyInput, userID uint) (*models.Property, error) {
	if len(input.Images) < 3 {
		return nil, ErrBadRequest("At least 3 images are required")
	}

	// Enforce per-plan listing limit
	owner, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrInternal("Failed to fetch user")
	}
	limit := config.GetListingLimit(owner.LandlordPlan)
	if limit >= 0 {
		count, _ := s.propRepo.CountActiveByOwner(userID)
		if int(count) >= limit {
			return nil, ErrForbidden("listing_limit_reached")
		}
	}

	prop := &models.Property{
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

	if err := s.propRepo.Create(prop); err != nil {
		return nil, ErrInternal("Failed to create property")
	}

	for i, fh := range input.Images {
		url, err := s.fileStore.Save(fh)
		if err != nil {
			s.propRepo.Delete(prop)
			return nil, ErrInternal("Failed to upload image: " + err.Error())
		}
		img := &models.PropertyImage{
			Filename:   url,
			PropertyID: prop.ID,
			IsMain:     i == 0,
		}
		if err := s.imageRepo.Create(img); err != nil {
			s.propRepo.Delete(prop)
			return nil, ErrInternal("Failed to save image record")
		}
	}

	created, err := s.propRepo.GetByID(prop.ID)
	if err != nil {
		return nil, ErrInternal("Failed to reload property")
	}

	go s.notifyMatchingAlerts(*created)
	return created, nil
}

func (s *propertyService) UpdateProperty(id uint, input UpdatePropertyInput, userID uint) (*models.Property, error) {
	prop, err := s.propRepo.GetByID(id)
	if err != nil {
		return nil, ErrNotFound("Property not found")
	}
	if prop.OwnerID != userID {
		return nil, ErrForbidden("Not authorized")
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
	if err := s.propRepo.Update(prop, updates); err != nil {
		return nil, ErrInternal("Failed to update property")
	}
	return s.propRepo.GetByID(id)
}

func (s *propertyService) DeleteProperty(id uint, userID uint) error {
	prop, err := s.propRepo.GetByID(id)
	if err != nil {
		return ErrNotFound("Property not found")
	}
	if prop.OwnerID != userID {
		return ErrForbidden("Not authorized")
	}
	for _, img := range prop.Images {
		s.fileStore.Delete(img.Filename)
	}
	return s.propRepo.Delete(prop)
}

func (s *propertyService) ToggleAvailability(id uint, userID uint) (bool, error) {
	prop, err := s.propRepo.GetByID(id)
	if err != nil {
		return false, ErrNotFound("Property not found")
	}
	if prop.OwnerID != userID {
		return false, ErrForbidden("Not authorized")
	}
	newVal := !prop.IsAvailable
	s.propRepo.ToggleAvailability(prop, newVal)
	return newVal, nil
}

func (s *propertyService) AddImage(propID uint, userID uint, file *multipart.FileHeader) (*models.PropertyImage, error) {
	prop, err := s.propRepo.GetByID(propID)
	if err != nil {
		return nil, ErrNotFound("Property not found")
	}
	if prop.OwnerID != userID {
		return nil, ErrForbidden("Not authorized")
	}

	url, err := s.fileStore.Save(file)
	if err != nil {
		return nil, ErrBadRequest(err.Error())
	}

	count, _ := s.imageRepo.CountByPropertyID(prop.ID)
	img := &models.PropertyImage{
		Filename:   url,
		PropertyID: prop.ID,
		IsMain:     count == 0,
	}
	if err := s.imageRepo.Create(img); err != nil {
		return nil, ErrInternal("Failed to save image record")
	}
	return img, nil
}

func (s *propertyService) DeleteImage(imgID uint, userID uint) error {
	img, err := s.imageRepo.GetByID(imgID)
	if err != nil {
		return ErrNotFound("Image not found")
	}
	prop, err := s.propRepo.GetByID(img.PropertyID)
	if err != nil {
		return ErrNotFound("Property not found")
	}
	if prop.OwnerID != userID {
		return ErrForbidden("Not authorized")
	}
	s.fileStore.Delete(img.Filename)
	return s.imageRepo.Delete(img)
}

func (s *propertyService) SetMainImage(imgID uint, userID uint) error {
	img, err := s.imageRepo.GetByID(imgID)
	if err != nil {
		return ErrNotFound("Image not found")
	}
	prop, err := s.propRepo.GetByID(img.PropertyID)
	if err != nil {
		return ErrNotFound("Property not found")
	}
	if prop.OwnerID != userID {
		return ErrForbidden("Not authorized")
	}
	s.imageRepo.UnsetMainForProperty(img.PropertyID)
	return s.imageRepo.SetMain(img)
}

func (s *propertyService) notifyMatchingAlerts(prop models.Property) {
	alerts, err := s.alertRepo.FindMatchingAlerts(prop)
	if err != nil || len(alerts) == 0 {
		return
	}

	title := "New property available"
	body := fmt.Sprintf("%s in %s, %s — %s (%.0f %s)",
		prop.PropertyType, prop.Neighborhood, prop.City, prop.TransactionType, prop.Price, prop.Currency)
	data := map[string]any{"property_id": prop.ID}

	event := appws.Event{Type: "property_alert", Data: map[string]any{
		"property_id":      prop.ID,
		"title":            prop.Title,
		"city":             prop.City,
		"neighborhood":     prop.Neighborhood,
		"property_type":    prop.PropertyType,
		"transaction_type": prop.TransactionType,
		"price":            prop.Price,
		"currency":         prop.Currency,
	}}

	for _, alert := range alerts {
		s.hub.Emit(alert.UserID, event)
		if !s.hub.IsOnline(alert.UserID) {
			if user, err := s.userRepo.GetByIDWithToken(alert.UserID); err == nil {
				utils.SendExpoPush(user.ExpoPushToken, title, body, data)
			}
		}
	}
}
