package usecase

import (
	"context"
	"errors"
	"intern-project-v2/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) GetAll(ctx context.Context) ([]*domain.Customer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *domain.CustomerRequest) (*domain.Customer, error) {
	args := m.Called(ctx, customer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Update(ctx context.Context, id string, customerReq *domain.CustomerRequest) (*domain.Customer, error) {
	args := m.Called(ctx, id, customerReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id string) (*domain.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func TestCustomerUsecase_GetAll(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(*MockCustomerRepository)
		expectedResult []*domain.Customer
		expectedError  error
	}{
		{
			name: "Success - Get all customers",
			mockSetup: func(mockRepo *MockCustomerRepository) {
				customers := []*domain.Customer{
					{
						Id:    bson.NewObjectID(),
						Name:  "John Doe",
						Email: "john@example.com",
						Phone: "123456789",
					},
					{
						Id:    bson.NewObjectID(),
						Name:  "Jane Smith",
						Email: "jane@example.com",
						Phone: "987654321",
					},
				}
				mockRepo.On("GetAll", mock.Anything).Return(customers, nil)
			},
			expectedResult: []*domain.Customer{
				{
					Id:    bson.NewObjectID(),
					Name:  "John Doe",
					Email: "john@example.com",
					Phone: "123456789",
				},
				{
					Id:    bson.NewObjectID(),
					Name:  "Jane Smith",
					Email: "jane@example.com",
					Phone: "987654321",
				},
			},
			expectedError: nil,
		},
		{
			name: "Error - Repository fails",
			mockSetup: func(mockRepo *MockCustomerRepository) {
				mockRepo.On("GetAll", mock.Anything).Return([]*domain.Customer(nil), errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
		{
			name: "Success - Empty result",
			mockSetup: func(mockRepo *MockCustomerRepository) {
				mockRepo.On("GetAll", mock.Anything).Return([]*domain.Customer{}, nil)
			},
			expectedResult: []*domain.Customer{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockCustomerRepository)
			tt.mockSetup(mockRepo)

			usecase := NewCustomerUsecase(mockRepo)
			ctx := context.Background()

			// Act
			result, err := usecase.GetAll(ctx)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, len(tt.expectedResult))
			}

			// Verify mock was called
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCustomerUsecase_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		customerID     string
		mockSetup      func(*MockCustomerRepository)
		expectedResult *domain.Customer
		expectedError  error
	}{
		{
			name:       "Success - Customer found",
			customerID: "64f1a2b3c4d5e6f7a8b9c0d1",
			mockSetup: func(mockRepo *MockCustomerRepository) {
				customer := &domain.Customer{
					Id:    bson.NewObjectID(),
					Name:  "John Doe",
					Email: "john@example.com",
					Phone: "123456789",
				}
				mockRepo.On("GetByID", mock.Anything, "64f1a2b3c4d5e6f7a8b9c0d1").Return(customer, nil)
			},
			expectedResult: &domain.Customer{
				Id:    bson.NewObjectID(),
				Name:  "John Doe",
				Email: "john@example.com",
				Phone: "123456789",
			},
			expectedError: nil,
		},
		{
			name:       "Error - Customer not found",
			customerID: "64f1a2b3c4d5e6f7a8b9c0d1",
			mockSetup: func(mockRepo *MockCustomerRepository) {
				mockRepo.On("GetByID", mock.Anything, "64f1a2b3c4d5e6f7a8b9c0d1").Return(nil, errors.New("customer not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("customer not found"),
		},
		{
			name:       "Error - Invalid ID format",
			customerID: "invalid-id",
			mockSetup: func(mockRepo *MockCustomerRepository) {
				mockRepo.On("GetByID", mock.Anything, "invalid-id").Return(nil, errors.New("invalid ID format"))
			},
			expectedResult: nil,
			expectedError:  errors.New("invalid ID format"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockCustomerRepository)
			tt.mockSetup(mockRepo)

			usecase := NewCustomerUsecase(mockRepo)
			ctx := context.Background()

			// Act
			result, err := usecase.GetByID(ctx, tt.customerID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Name, result.Name)
				assert.Equal(t, tt.expectedResult.Email, result.Email)
				assert.Equal(t, tt.expectedResult.Phone, result.Phone)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCustomerUsecase_Create(t *testing.T) {
	tests := []struct {
		name           string
		customerReq    *domain.CustomerRequest
		mockSetup      func(*MockCustomerRepository)
		expectedResult *domain.Customer
		expectedError  error
	}{
		{
			name: "Success - Create customer",
			customerReq: &domain.CustomerRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Phone: "123456789",
			},
			mockSetup: func(mockRepo *MockCustomerRepository) {
				customerReq := &domain.CustomerRequest{
					Name:  "John Doe",
					Email: "john@example.com",
					Phone: "123456789",
				}
				createdCustomer := &domain.Customer{
					Id:    bson.NewObjectID(),
					Name:  "John Doe",
					Email: "john@example.com",
					Phone: "123456789",
				}
				mockRepo.On("Create", mock.Anything, customerReq).Return(createdCustomer, nil)
			},
			expectedResult: &domain.Customer{
				Id:    bson.NewObjectID(),
				Name:  "John Doe",
				Email: "john@example.com",
				Phone: "123456789",
			},
			expectedError: nil,
		},
		{
			name: "Error - Repository fails",
			customerReq: &domain.CustomerRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Phone: "123456789",
			},
			mockSetup: func(mockRepo *MockCustomerRepository) {
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockCustomerRepository)
			tt.mockSetup(mockRepo)

			usecase := NewCustomerUsecase(mockRepo)
			ctx := context.Background()

			// Act
			result, err := usecase.Create(ctx, tt.customerReq)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Name, result.Name)
				assert.Equal(t, tt.expectedResult.Email, result.Email)
				assert.Equal(t, tt.expectedResult.Phone, result.Phone)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
