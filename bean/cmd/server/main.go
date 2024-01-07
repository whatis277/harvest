package main

import (
	"fmt"

	"harvest/bean/internal/entity"

	envAdapter "harvest/bean/internal/adapter/env"

	"harvest/bean/internal/driver/database"
	tokenDS "harvest/bean/internal/driver/datasource/token"
	userDS "harvest/bean/internal/driver/datasource/user"
)

func main() {
	env, err := envAdapter.New()
	if err != nil {
		panic(
			fmt.Errorf("error reading env: %v", err),
		)
	}

	db, err := database.New(&database.DSNBuilder{
		Host:        env.DB.Host,
		Name:        env.DB.Name,
		Username:    env.DB.Username,
		Password:    env.DB.Password,
		Tls:         true,
		Interpolate: true,
		ParseTime:   true,
	})
	if err != nil {
		panic(
			fmt.Errorf("error connecting db: %v", err),
		)
	}
	defer db.Close()

	u := createUserIfNotExists(db)
	if u == nil {
		fmt.Println("user not found")
		return
	}

	testLoginTokens(db, u)
}

func createUserIfNotExists(db *database.DB) *entity.User {
	users := userDS.New(db)

	u, _ := users.FindByEmail("test@gmail.com")
	if u != nil {
		fmt.Println("user found by email: ", u.ID, u.Email, u.CreatedAt, u.UpdatedAt)
		return u
	}

	u, err := users.Create(&entity.User{Email: "test@gmail.com"})
	if u != nil {
		fmt.Println("user created: ", u.ID, u.Email, u.CreatedAt, u.UpdatedAt)
	} else {
		fmt.Println("user not created: ", err)
		return nil
	}

	u, err = users.FindById(u.ID)
	if u != nil {
		fmt.Println("user found by id: ", u.ID, u.Email, u.CreatedAt, u.UpdatedAt)
	} else {
		fmt.Println("user not found by id: ", err)
	}

	return u
}

func testLoginTokens(db *database.DB, u *entity.User) {
	tokens := tokenDS.New(db)

	t1 := [60]byte{[]byte("123")[0], []byte("123")[1], []byte("123")[2]}
	t2 := [60]byte{[]byte("456")[0], []byte("456")[1], []byte("456")[2]}

	err := tokens.Create(&entity.LoginToken{Email: u.Email, HashedToken: t1})
	if err != nil {
		fmt.Println("error creating token: ", err)
		return
	}

	t, err := tokens.FindUnexpired(&entity.LoginToken{Email: u.Email, HashedToken: t1})
	if t != nil {
		fmt.Println("token found: ", t.ID, t.Email, string(t.HashedToken[:]), t.CreatedAt, t.ExpiresAt)
	} else {
		fmt.Println("token not found: ", err)
	}

	err = tokens.Create(&entity.LoginToken{Email: u.Email, HashedToken: t2})
	if err != nil {
		fmt.Println("error overwriting token: ", err)
	}

	t, err = tokens.FindUnexpired(&entity.LoginToken{Email: u.Email, HashedToken: t2})
	if t != nil {
		fmt.Println("token found: ", t.ID, t.Email, string(t.HashedToken[:]), t.CreatedAt, t.ExpiresAt)
	} else {
		fmt.Println("token not found: ", err)
		return
	}

	err = tokens.Delete(t)
	if err != nil {
		fmt.Println("error deleting token: ", err)
	}
}
