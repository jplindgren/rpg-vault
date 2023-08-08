package users

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/jplindgren/rpg-vault/internal/clients"
	"github.com/jplindgren/rpg-vault/internal/validator"
)

// Define constants for the token scope. For now we just define the scope "activation"
// but we'll add additional scopes later in the book.
const (
	ScopeActivation = "activation"
	Authentication  = "authentication"
)

type dbToken struct {
	Hash  []byte `dynamodbav:"hash"`
	Email string `dynamodbav:"email"`
	//Expiry string `dynamodbav:"expiry"`
	Expiry int64  `dynamodbav:"expiry"`
	Scope  string `dynamodbav:"scope"`
}

type TokenService struct {
	Store     *clients.DynamoDbClientWrapper
	TableName string
}

func NewTokenSrv(store *clients.DynamoDbClientWrapper, tableName string) *TokenService {
	return &TokenService{
		Store:     store,
		TableName: tableName,
	}
}

func (s *TokenService) New(email string, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(email, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = s.Insert(token)

	return token, err
}

func (s *TokenService) Insert(token *Token) error {
	item := dbToken{
		Email:  token.Email,
		Hash:   token.Hash,
		Expiry: token.Expiry.Unix(),
		Scope:  token.Scope,
	}

	_, putItemErr := s.Store.PutWrapper(s.TableName, item, nil)
	return putItemErr
}

type TokenKeyBasedStruct struct {
	Hash []byte `dynamodbav:"hash"`
}

func (s *TokenService) Get(tokenPlaintext string) (*Token, error) {
	token := sha256.Sum256([]byte(tokenPlaintext))
	hash := token[:]

	key := &TokenKeyBasedStruct{
		Hash: hash,
	}

	result := &dbToken{}
	_, err := s.Store.GetWrapper(s.TableName, key, result)
	if err != nil {
		return &Token{}, err
	}

	return &Token{
		Hash:   result.Hash,
		Email:  result.Email,
		Expiry: time.Unix(result.Expiry, 0),
		Scope:  result.Scope,
	}, nil
}

// Check that the plaintext token has been provided and is exactly 26 bytes long.
func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "tokenPlaintext", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "tokenPlaintext", "must be 26 bytes long")
}

func generateToken(email string, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		Email:  email,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// Initialize a zero-valued byte slice with a length of 16 bytes.
	randomBytes := make([]byte, 16)
	// Use the Read() function from the crypto/rand package to fill the byte slice with
	// random bytes from your operating system's CSPRNG. This will return an error if
	// the CSPRNG fails to function correctly.
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Encode the byte slice to a base-32-encoded string and assign it to the token
	// Plaintext field. This will be the token string that we send to the user in their
	// welcome email. They will look similar to this:
	//
	// Y3QMGX3PJ3WLRL2YRTQGQ6KRHU
	//
	// Note that by default base-32 strings may be padded at the end with the =
	// character. We don't need this padding character for the purpose of our tokens, so
	// we use the WithPadding(base32.NoPadding) method in the line below to omit them.
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Generate a SHA-256 hash of the plaintext token string. This will be the value
	// that we store in the `hash` field of our database table. Note that the
	// sha256.Sum256() function returns an *array* of length 32, so to make it easier to
	// work with we convert it to a slice using the [:] operator before storing it.
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}
