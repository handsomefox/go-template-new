package user_test

import (
	"context"
	"testing"

	"project-template/database/sqlc"
	"project-template/services/user"
	"project-template/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_Create(t *testing.T) {
	// Setup
	testDB, err := testhelper.SetupTestDatabase()
	require.NoError(t, err)
	defer testDB.TearDown()

	// Arrange
	ctx := context.Background()
	service := user.New(testDB.DB)

	// Act
	createdUser, err := service.Create(ctx, "John Doe", "john@example.com")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, "John Doe", createdUser.Name)
	assert.Equal(t, "john@example.com", createdUser.Email)

	// Verify the user was actually created in the database
	var dbUser sqlc.Users
	err = testDB.DB.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = $1", createdUser.ID).
		Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, dbUser.ID)
	assert.Equal(t, createdUser.Name, dbUser.Name)
	assert.Equal(t, createdUser.Email, dbUser.Email)
}
