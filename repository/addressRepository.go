package repository

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/helper"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"go.uber.org/zap"
)

type AddressRepository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func NewAddressRepository(db *sql.DB, logger *zap.Logger) AddressRepository {
	return AddressRepository{DB: db, Logger: logger}
}

func (repo AddressRepository) Create(userID string, addressInput model.Address) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction", zap.Error(err), zap.String("Repository", "Address"), zap.String("Function", "Create"))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction", zap.Error(err), zap.String("Repository", "User"), zap.String("Function", "Create"))
			tx.Rollback()
		}
	}()

	defaultAddress, err := repo.GetDefaultAddress(userID)
	if err != nil {
		return err
	}
	if (defaultAddress == model.Address{}) {
		addressInput.IsDefault = true
	} else {
		addressInput.IsDefault = false
	}

	sqlStatement := `INSERT INTO addresses (user_id, name, street, district, city, state, postal_code, country, is_default) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = tx.Exec(sqlStatement, userID, addressInput.Name, addressInput.Street, addressInput.District, addressInput.City, addressInput.State, addressInput.PostalCode, addressInput.Country, addressInput.IsDefault)
	if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.String("Repository", "Address"), zap.String("Function", "Create"))
		return err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction", zap.Error(err), zap.String("Repository", "Address"), zap.String("Function", "Create"))
		return err
	}

	return nil
}

func (repo AddressRepository) GetDefaultAddress(userID string) (model.Address, error) {
	var address model.Address
	sqlStatement := `SELECT id, name, street, district, city, state, postal_code, country, is_default FROM addresses WHERE user_id = $1 AND is_default = true AND status = 'active'`
	err := repo.DB.QueryRow(sqlStatement, userID).Scan(&address.ID, &address.Name, &address.Street, &address.District, &address.City, &address.State, &address.PostalCode, &address.Country, &address.IsDefault)
	if err != nil {
		if err == sql.ErrNoRows {
			repo.Logger.Info("No default address found", zap.String("user_id", userID), zap.String("Repository", "Address"), zap.String("Function", "GetDefaultAddress"))
			return address, nil // No default address found
		} else {
			repo.Logger.Error("Failed to execute query", zap.Error(err), zap.String("Repository", "Address"), zap.String("Function", "GetDefaultAddress"))
			return address, err
		}
	}

	return address, nil
}

func (repo AddressRepository) GetAll(userID string, pagination model.Pagination) ([]model.Address, model.Pagination, error) {
	var addresses []model.Address
	sqlStatement := `SELECT id, name, street, district, city, state, postal_code, country, is_default FROM addresses WHERE user_id = $1 AND status = 'active' ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	limit := pagination.PerPage
	offset := (pagination.Page - 1) / limit

	repo.Logger.Info("Executing query", zap.String("query", sqlStatement), zap.String("Repository", "Address"), zap.String("Function", "GetAll"))
	rows, err := repo.DB.Query(sqlStatement, userID, limit, offset)
	if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.String("Repository", "Address"), zap.String("Function", "GetAll"))
		return nil, pagination, err
	}
	defer rows.Close()

	for rows.Next() {
		var address model.Address
		err = rows.Scan(&address.ID, &address.Name, &address.Street, &address.District, &address.City, &address.State, &address.PostalCode, &address.Country, &address.IsDefault)
		if err != nil {
			repo.Logger.Error("Failed to scan row", zap.Error(err), zap.String("Repository",
				"Address"), zap.String("Function", "GetAll"))
			return nil, pagination, err
		}
		addresses = append(addresses, address)
	}
	return addresses, pagination, nil
}

func (repo AddressRepository) GetByID(id int) (*model.Address, error) {
	var address model.Address
	sqlStatement := `SELECT id, name, street, district, city, state, postal_code, country, is_default FROM addresses WHERE id = $1 AND status = 'active'`

	repo.Logger.Info("Running query", zap.String("query", sqlStatement), zap.String("Repository", "Address"), zap.String("Function", "GetByID"))
	err := repo.DB.QueryRow(sqlStatement, id).Scan(&address.ID, &address.Name, &address.Street, &address.District, &address.City, &address.State, &address.PostalCode, &address.Country, &address.IsDefault)
	if err != nil {
		if err == sql.ErrNoRows {
			repo.Logger.Info("Address not found", zap.Int("id", id), zap.String("Repository", "Address"), zap.String("Function", "GetByID"))
			return nil, nil
		}
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.String("Repository", "Address"), zap.String("Function", "GetByID"))
		return nil, err
	}

	return &address, nil
}

func (repo AddressRepository) Update(id int, userID string, addressInput model.Address) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction", zap.Error(err), zap.String("Repository",
			"Address"), zap.String("Function", "Update"))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction", zap.Error(err), zap.String("Repository", "User"), zap.String("Function", "Create"))
			tx.Rollback()
		}
	}()

	// Create a map to hold the fields to update
	fields := map[string]interface{}{}

	if addressInput.Name != "" {
		fields["name"] = addressInput.Name
	}
	if addressInput.Street != "" {
		fields["street"] = addressInput.Street
	}
	if addressInput.District.Valid && addressInput.District.String != "" {
		fields["district"] = addressInput.District.String
	}
	if addressInput.City.Valid && addressInput.City.String != "" {
		fields["city"] = addressInput.City.String
	}
	if addressInput.State.Valid && addressInput.State.String != "" {
		fields["state"] = addressInput.State.String
	}
	if addressInput.PostalCode != "" {
		fields["postal_code"] = addressInput.PostalCode
	}
	if addressInput.Country != "" {
		fields["country"] = addressInput.Country
	}

	// Add the updated_at field
	fields["updated_at"] = time.Now()

	// Build the SET clause dynamically
	setClauses := []string{}
	values := []interface{}{}
	index := 1
	for field, value := range fields {
		setClauses = append(setClauses, field+"=$"+strconv.Itoa(index))
		values = append(values, value)
		index++
	}

	// Check if there are fields to update
	if len(setClauses) == 0 {
		return errors.New("no fields to update")
	}

	// Build the final query
	queryStatement := `
		UPDATE addresses
		SET ` + helper.JoinStrings(setClauses, ", ") + `
		WHERE id = $` + strconv.Itoa(index) +
		` AND user_id = $` + strconv.Itoa(index+1) +
		` AND status = 'active'`
	values = append(values, id, userID)

	// Execute the query
	repo.Logger.Info("Executing query", zap.Int("address_id", id),
		zap.String("query", queryStatement), zap.String("repository", "Address"),
		zap.String("function", "Update"))
	_, err = tx.Exec(queryStatement, values...)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (repo AddressRepository) UpdateDefaultAddress(id int, userID string, setAsDefault bool) error {
	// Check if the address exists
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction", zap.Error(err), zap.String("Repository", "User"), zap.String("Function", "Create"))
			tx.Rollback()
		}
	}()

	sqlStatement := `UPDATE addresses SET is_default = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`

	repo.Logger.Info("Executing query", zap.Bool("boolean", setAsDefault),
		zap.String("query", sqlStatement), zap.String("repository", "Address"),
		zap.String("function", "UpdateDefaultAddress"))

	_, err = tx.Exec(sqlStatement, setAsDefault, id, userID)
	if err != nil {

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (repo AddressRepository) Delete(id int, userID string) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction", zap.Error(err), zap.String("Repository", "User"),
				zap.String("Function", "Delete"))
			tx.Rollback()
		}
	}()

	sqlStatement := `UPDATE addresses SET status = 'deleted', deleted_at = NOW() WHERE id = $1 AND user_id = $2`
	repo.Logger.Info("Executing query", zap.Int("address_id", id),
		zap.String("query", sqlStatement), zap.String("repository", "Address"),
		zap.String("function", "Delete"))
	_, err = tx.Exec(sqlStatement, id, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
