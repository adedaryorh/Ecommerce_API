package db_test

import (
	"context"
	db "github.com/adedaryorh/ecommerceapi/db/sqlc"
	"github.com/adedaryorh/ecommerceapi/utils"
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
	"time"
)

func clean_up() {
	err := testQuery.DeleteAllUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func TestDeleteUser(t *testing.T) {
	defer clean_up()
	user := createRandomUser(t)

	err := testQuery.DeleteUser(context.Background(), user.ID)

	assert.NoError(t, err)

	newUser, err := testQuery.GetUserByID(context.Background(), user.ID)
	assert.Error(t, err)
	assert.Empty(t, newUser)

}

func TestListUser(t *testing.T) {
	defer clean_up()
	//go_routine call
	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		//concurrency
		go func() {
			defer wg.Done()
			createRandomUser(t)
		}()
	}
	wg.Wait()
	arg := db.ListUserParams{
		Offset: 0,
		Limit:  30,
	}
	users, err := testQuery.ListUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, users)
	assert.Equal(t, len(users), 30)
}

func TestGetUserByEmail(t *testing.T) {
	defer clean_up()
	user := createRandomUser(t)

	newUser, err := testQuery.GetUserByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.NotEmpty(t, newUser)

	assert.Equal(t, newUser.HashedPassword, user.HashedPassword)
	assert.Equal(t, user.Email, newUser.Email)
}

func createRandomUser(t *testing.T) db.User {
	hashedPass, err := utils.GenerateHashedPassword(utils.RandomString(8))
	if err != nil {
		log.Fatal("Unable to generate Pass", err)
	}

	arg := db.CreateUserParams{
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPass,
		Username:       utils.RandomString(6),
	}

	user, err := testQuery.CreateUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, user.Email, arg.Email)
	assert.Equal(t, user.HashedPassword, arg.HashedPassword)
	assert.Equal(t, user.Username, arg.Username)
	assert.WithinDuration(t, user.CreatedAt, time.Now(), 2*time.Second)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 2*time.Second)

	return user
}

func TestCreateUser(t *testing.T) {
	defer clean_up()

	user1 := createRandomUser(t)

	user2, err := testQuery.CreateUser(context.Background(), db.CreateUserParams{
		Email:          user1.Email,
		HashedPassword: user1.HashedPassword,
		Username:       user1.Username,
	})
	assert.Error(t, err)
	assert.Empty(t, user2)
}

func TestUpdateUser(t *testing.T) {
	defer clean_up()
	user := createRandomUser(t)

	newPassword, err := utils.GenerateHashedPassword(utils.RandomString(8))
	if err != nil {
		log.Fatal("Unable to generate Pass", err)
	}
	arg := db.UpdateUserPasswordParams{
		HashedPassword: newPassword,
		UpdatedAt:      time.Now(),
		ID:             user.ID,
	}
	newUser, err := testQuery.UpdateUserPassword(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, newUser)
	assert.Equal(t, newUser.HashedPassword, arg.HashedPassword)
	assert.Equal(t, user.Email, newUser.Email)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 2*time.Second)
}
