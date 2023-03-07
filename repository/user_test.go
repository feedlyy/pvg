package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/magiconair/properties/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"pvg/domain"
	"regexp"
	"testing"
)

func TestUserRepository_Fetch(t *testing.T) {
	var (
		query   = `SELECT * FROM "users"`
		returns []domain.Users
	)
	returns = append(returns, domain.Users{
		Username: "fadli",
		Email:    "feedlyy@gmail.com",
	})

	testCases := []struct {
		name        string
		expectedErr bool
		context     context.Context
		doMockDB    func(mock sqlmock.Sqlmock)
		expected    []domain.Users
	}{
		{
			name:        "Success",
			expectedErr: false,
			context:     context.Background(),
			doMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(sqlmock.NewRows([]string{"username", "email"}).
					AddRow("fadli", "feedlyy@gmail.com"))
			},
			expected: returns,
		},
		{
			name:        "Failed",
			expectedErr: true,
			context:     context.Background(),
			doMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("internal server error"))
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err    error
				res    []domain.Users
				db     *sql.DB
				mock   sqlmock.Sqlmock
				mockDB *gorm.DB
			)
			db, mock, err = sqlmock.New()
			if err != nil {
				panic(err)
			}
			defer db.Close()
			tc.doMockDB(mock)

			mockDB, err = gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}))

			repoDB := NewUserRepository(mockDB)
			res, err = repoDB.Fetch(tc.context)
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected, res)
		})
	}
}
