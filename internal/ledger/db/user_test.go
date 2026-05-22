package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Shihasz/go-fintech-ledger/internal/common/util"
	"github.com/stretchr/testify/require"
)

// createRandomUser inserts a unique user into the database and returns it
func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword("secret123")
	require.NoError(t, err)

	// Use UnixNano to guarantee unique usernames and emails for every test run
	uniqueSuffix := time.Now().UnixNano()

	arg := CreateUserParams{
		Username:       fmt.Sprintf("user_%d", uniqueSuffix),
		HashedPassword: hashedPassword,
		FullName:       "John Doe",
		Email:          fmt.Sprintf("john_%d@example.com", uniqueSuffix),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)

	return user

}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
