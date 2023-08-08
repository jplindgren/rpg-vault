package users

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jplindgren/rpg-vault/internal/clients"
	"github.com/jplindgren/rpg-vault/internal/validator"
)

// Define a custom ErrorDuplicateEmail error.
var (
	ErrorDuplicateEmail = errors.New("duplicate email")
)

//const usersTable = "rpg_users"

type dbUserItem struct {
	Email         string `dynamodbav:"email"`
	CreatedAt     string `dynamodbav:"created_at"`
	Name          string `dynamodbav:"name"`
	Password_hash []byte `dynamodbav:"password_hash"`
	Activated     bool   `dynamodbav:"activated"`
	Version       int    `dynamodbav:"version"`
}

func New(store *clients.DynamoDbClientWrapper, tableName string) *UserService {
	return &UserService{
		Store:     store,
		TableName: tableName,
	}
}

type UserService struct {
	Store     *clients.DynamoDbClientWrapper
	TableName string
}

func (s UserService) Insert(user *User) error {
	item := dbUserItem{
		Email:         user.Email,
		CreatedAt:     user.CreatedAt,
		Activated:     user.Activated,
		Password_hash: user.Password.hash,
		Name:          user.Name,
		Version:       user.Version,
	}

	_, putItemErr := s.Store.PutWrapper(s.TableName, item, aws.String("attribute_not_exists(email)"))
	return putItemErr
}

type UserKeyBasedStruct struct {
	Email string `dynamodbav:"email"`
}

// Retrieve the User details from the database based on the user's email address.
// Because we have a UNIQUE constraint on the email column, this SQL query will only
// return one record (or none at all, in which case we return a ErrorRecordNotFound error).
func (s UserService) GetByEmail(email string) (*User, error) {
	key := &UserKeyBasedStruct{
		Email: email,
	}
	result := &dbUserItem{}
	_, err := s.Store.GetWrapper(s.TableName, key, result)
	if err != nil {
		return &User{}, err
	}

	return NewUser(result.Email,
		result.Name,
		result.CreatedAt,
		result.Password_hash,
		result.Version,
		result.Activated), nil
}

// Update the details for a specific user. Notice that we check against the version
// field to help prevent any race conditions during the request cycle, just like we did
// when updating a movie. And we also check for a violation of the "users_email_key"
// constraint when performing the update, just like we did when inserting the user
// record originally.
func (s UserService) Update(user *User) error {
	panic("missing Update")
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlainText(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	// Call the standalone ValidateEmail() helper.
	ValidateEmail(v, user.Email)

	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlainText() helper.
	if user.Password.plaintext != nil {
		ValidatePasswordPlainText(v, *user.Password.plaintext)
	}

	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
