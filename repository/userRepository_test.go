package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TestUserRepository_Create(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open a mock database connection: %v", err)
	}
	defer db.Close()

	// Create a logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create an instance of UserRepository
	userRepo := NewUserRepository(db, logger)

	tests := []struct {
		name               string
		userInput          model.UserDTO
		mockExpectations   func()
		expectError        bool
		isEmailValid       bool
		isPhoneNumberValid bool
	}{
		{
			name: "Successful Create with valid email",
			userInput: model.UserDTO{
				Name:               "Test User Email",
				Password:           "StrongPassword123!",
				EmailOrPhoneNumber: "test@example.com",
			},
			mockExpectations: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO users").
					WithArgs(sqlmock.AnyArg(), "Test User Email", "test@example.com", "test@example.com", "hashedpassword").
					WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate successful insert
				mock.ExpectCommit()
			},
			expectError:        false,
			isEmailValid:       true,
			isPhoneNumberValid: true,
		},
		{
			name: "Failed Create with valid phone number",
			userInput: model.UserDTO{
				Name:               "Test User Phone Number",
				Password:           "StrongPassword123!",
				EmailOrPhoneNumber: "1234567890",
			},
			mockExpectations: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO users").
					WithArgs(sqlmock.AnyArg(), "Test User Phone Number", "1234567890", "1234567890", "hashedpassword").
					WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate successful insert
				mock.ExpectCommit()
			},
			expectError:        false,
			isEmailValid:       true,
			isPhoneNumberValid: true,
		},
		{
			name: "Failed Transaction with invalid email and phone number",
			userInput: model.UserDTO{
				Name:               "Transaction Fail User",
				Password:           "StrongPassword123!",
				EmailOrPhoneNumber: "test@example.com", // valid email
			},
			mockExpectations: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO users").
					WithArgs(
						sqlmock.AnyArg(),        // user ID
						"Transaction Fail User", // name
						"test@example.com",      // email
						nil,                     // phone_number (set to nil to simulate invalid)
						"hashedpassword",        // password
					).
					WillReturnError(sql.ErrConnDone) // Simulate error during execution
				mock.ExpectRollback() // Expect rollback
			},
			expectError:        true,
			isEmailValid:       true,
			isPhoneNumberValid: false, // simulate phone number not valid because nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpectations() // Set up mock expectations

			// Simulate password hashing
			passwordHashed := "hashedpassword" // Replace with actual password hashing logic if needed

			// Prepare new user input based on the provided structure
			newUserInput := model.User{
				ID:             uuid.NewString(),
				Name:           tt.userInput.Name,
				PasswordHashed: passwordHashed,
				Email:          sql.NullString{String: tt.userInput.EmailOrPhoneNumber, Valid: tt.isEmailValid},
				PhoneNumber:    sql.NullString{String: tt.userInput.EmailOrPhoneNumber, Valid: tt.isPhoneNumberValid},
			}

			// Call the Create function
			err := userRepo.Create(newUserInput)

			// Check for errors
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserRepository_Login(t *testing.T) {
	// Common setup for database and logger
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub database connection: %s", err)
	}
	defer db.Close()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	userRepo := NewUserRepository(db, logger)

	tests := []struct {
		name             string
		userInput        model.UserDTO
		mockExpectations func()
		expectedUser     model.User
		expectError      bool
	}{
		{
			name: "Successful Login",
			userInput: model.UserDTO{
				EmailOrPhoneNumber: "test@example.com",
			},
			mockExpectations: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, password FROM users WHERE (email = $1 OR phone_number = $1) AND status = 'active'")).
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow("1", "hashedpassword"))
			},
			expectedUser: model.User{
				ID:             "1",
				PasswordHashed: "hashedpassword",
			},
			expectError: false,
		},
		{
			name: "Failed Login - User Not Found",
			userInput: model.UserDTO{
				EmailOrPhoneNumber: "notfound@example.com",
			},
			mockExpectations: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, password FROM users WHERE (email = $1 OR phone_number = $1) AND status = 'active'")).
					WithArgs("notfound@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser: model.User{},
			expectError:  true,
		},
		{
			name: "Failed Login - Repository Error",
			userInput: model.UserDTO{
				EmailOrPhoneNumber: "test@example.com",
			},
			mockExpectations: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, password FROM users WHERE (email = $1 OR phone_number = $1) AND status = 'active'")).
					WithArgs("test@example.com").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser: model.User{},
			expectError:  true,
		},
		{
			name: "Failed Login - Invalid Email Format",
			userInput: model.UserDTO{
				EmailOrPhoneNumber: "invalid_email",
			},
			mockExpectations: func() {}, // No DB interaction for invalid email
			expectedUser:     model.User{},
			expectError:      true,
		},
		{
			name: "Failed Login - Empty Email or Phone Number",
			userInput: model.UserDTO{
				EmailOrPhoneNumber: "",
			},
			mockExpectations: func() {}, // No DB interaction for empty input
			expectedUser:     model.User{},
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations for this test case
			tt.mockExpectations()

			// Execute the Login function
			user, err := userRepo.Login(tt.userInput)

			// Validate results
			if (err != nil) != tt.expectError {
				t.Errorf("[%s] expected error: %v, got: %v", tt.name, tt.expectError, err)
			}

			if !usersEqual(user, tt.expectedUser) {
				t.Errorf("[%s] expected user: %+v, got: %+v", tt.name, tt.expectedUser, user)
			}

			// Verify all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("[%s] unfulfilled mock expectations: %s", tt.name, err)
			}
		})
	}
}

func usersEqual(u1, u2 model.User) bool {
	return u1.ID == u2.ID && u1.PasswordHashed == u2.PasswordHashed
}
