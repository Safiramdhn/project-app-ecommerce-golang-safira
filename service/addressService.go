package service

import (
	"database/sql"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"go.uber.org/zap"
)

type AddressService struct {
	Repo   repository.MainRepository
	Logger *zap.Logger
}

func NewAddressService(repo repository.MainRepository, log *zap.Logger) AddressService {
	return AddressService{Repo: repo, Logger: log}
}

func (s *AddressService) GetAddressById(addressId int) (*model.Address, error) {
	return s.Repo.AddressRepository.GetByID(addressId)
}

func (s *AddressService) AddAddress(userId string, addressInput model.AddressDTO) error {
	var district = sql.NullString{String: "", Valid: false}
	var city = sql.NullString{String: "", Valid: false}
	var state = sql.NullString{String: "", Valid: false}

	if addressInput.District != "" {
		district = sql.NullString{String: addressInput.District, Valid: true}
	}
	if addressInput.City != "" {
		city = sql.NullString{String: addressInput.City, Valid: true}
	}
	if addressInput.State != "" {
		state = sql.NullString{String: addressInput.State, Valid: true}
	}

	newAddressInput := model.Address{
		Name:       addressInput.Name,
		Street:     addressInput.Street,
		District:   district,
		City:       city,
		State:      state,
		Country:    addressInput.Country,
		PostalCode: addressInput.PostalCode,
	}
	return s.Repo.AddressRepository.Create(userId, newAddressInput)
}

func (s *AddressService) RemoveAddress(addressId int, userId string) error {
	return s.Repo.AddressRepository.Delete(addressId, userId)
}

func (s *AddressService) SetDeafultAddress(addressId int, userId string, setAsDefault bool) error {
	s.Logger.Info("Getting current default address", zap.String("userId", userId), zap.String("service", "Address"), zap.String("function", "setAsDefault"))
	if setAsDefault {
		currentDefaultAddress, err := s.Repo.AddressRepository.GetDefaultAddress(userId)
		if err != nil {
			return err
		}

		if currentDefaultAddress.ID != 0 {
			err := s.Repo.AddressRepository.UpdateDefaultAddress(currentDefaultAddress.ID, userId, false)
			if err != nil {
				return err
			}
		}
	}
	return s.Repo.AddressRepository.UpdateDefaultAddress(addressId, userId, setAsDefault)
}

func (s *AddressService) GetAllAddresses(userID string, paginationInput model.Pagination) ([]model.Address, model.Pagination, error) {
	if paginationInput.Page == 0 {
		paginationInput.Page = 1
	}
	if paginationInput.PerPage == 0 {
		paginationInput.PerPage = 5
	}

	return s.Repo.AddressRepository.GetAll(userID, paginationInput)
}

func (s *AddressService) UpdateAdress(addressId int, userId string, addressInput model.AddressDTO) error {
	var district = sql.NullString{String: "", Valid: false}
	var city = sql.NullString{String: "", Valid: false}
	var state = sql.NullString{String: "", Valid: false}

	if addressInput.District != "" {
		district = sql.NullString{String: addressInput.District, Valid: true}
	}
	if addressInput.City != "" {
		city = sql.NullString{String: addressInput.City, Valid: true}
	}
	if addressInput.State != "" {
		state = sql.NullString{String: addressInput.State, Valid: true}
	}

	newAddressInput := model.Address{
		Name:       addressInput.Name,
		Street:     addressInput.Street,
		District:   district,
		City:       city,
		State:      state,
		Country:    addressInput.Country,
		PostalCode: addressInput.PostalCode,
	}
	return s.Repo.AddressRepository.Update(addressId, userId, newAddressInput)
}
