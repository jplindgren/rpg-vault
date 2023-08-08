package services

import (
	"github.com/jplindgren/rpg-vault/internal/characters"
	"github.com/jplindgren/rpg-vault/internal/clients"
	"github.com/jplindgren/rpg-vault/internal/users"
	"github.com/jplindgren/rpg-vault/internal/worlds"
)

type Services struct {
	Users      *users.UserService
	Tokens     *users.TokenService
	Worlds     *worlds.WorldService
	Characters *characters.CharacterService
}

// Interface to mock models and help unit tests
// type Services struct {
// 	Users interface {
// 		Insert(user *User) error
// 		Get(id int64) (*User, error)
// 		Update(user *User) error
// 		Delete(id int64) error
// 	}
// }

func NewServices(dynClientWrapper *clients.DynamoDbClientWrapper, s3ClientWrapper *clients.S3ClientWrapper) Services {
	return Services{
		Users:      users.New(dynClientWrapper, "rpg_users"),
		Tokens:     users.NewTokenSrv(dynClientWrapper, "rpg_usertokens"),
		Worlds:     worlds.New(dynClientWrapper, s3ClientWrapper, "rpg_worlds"),
		Characters: characters.New(dynClientWrapper, "rpg_characters"),
	}
}

// Create a helper function which returns a Models instance containing the mock models
// only.
// func NewServices() Services {
//     return Services{
//         Users: MockUserService{},
//     }
// }
