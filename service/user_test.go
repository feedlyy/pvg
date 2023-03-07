package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"pvg/domain"
	mock_domain "pvg/domain/mocks"
	"testing"
)

func TestUserService_GetAllUser(t *testing.T) {
	var (
		returns []domain.Users
	)
	returns = append(returns, domain.Users{
		Username: "fadli",
		Email:    "feedlyy@gmail.com",
	})

	testCases := []struct {
		name         string
		expectedErr  bool
		context      context.Context
		mockUserRepo func(mock *mock_domain.MockUserRepository)
		expected     []domain.Users
	}{
		{
			name:        "Success",
			expectedErr: false,
			context:     context.Background(),
			mockUserRepo: func(mock *mock_domain.MockUserRepository) {
				mock.EXPECT().Fetch(gomock.Any()).Return(returns, nil)
			},
			expected: returns,
		},
		{
			name:        "Failed",
			expectedErr: true,
			context:     context.Background(),
			mockUserRepo: func(mock *mock_domain.MockUserRepository) {
				mock.EXPECT().Fetch(gomock.Any()).Return(nil, errors.New("internal server error"))
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err          error
				res          []domain.Users
				mockCtrl     = gomock.NewController(t)
				mockUserRepo = mock_domain.NewMockUserRepository(mockCtrl)
				service      = userService{
					userRepo: mockUserRepo,
					kafka:    domain.KafkaProducer{},
					acRepo:   nil,
				}
			)
			defer mockCtrl.Finish()
			tc.mockUserRepo(mockUserRepo)

			res, err = service.GetAllUser(tc.context)
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected, res)
		})
	}
}
