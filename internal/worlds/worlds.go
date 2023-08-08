package worlds

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	common "github.com/jplindgren/rpg-vault/internal"
	"github.com/jplindgren/rpg-vault/internal/clients"
	"github.com/jplindgren/rpg-vault/internal/uploader"
	"github.com/jplindgren/rpg-vault/internal/validator"
)

//const worldTableName = "rpg_worlds"

type WorldService struct {
	db        *clients.DynamoDbClientWrapper
	s3        *clients.S3ClientWrapper
	tableName string
}

func New(db *clients.DynamoDbClientWrapper, s3 *clients.S3ClientWrapper, tableName string) *WorldService {
	return &WorldService{
		db:        db,
		s3:        s3,
		tableName: tableName,
	}
}

const coverImageDestination = "%s/world/cover.png"

func (ws *WorldService) Insert(world *World) error {
	world.Id = common.GenerateToken()
	coverUrl, err := uploader.UploadCoverImage(ws.s3,
		world.CoverImage,
		fmt.Sprintf(coverImageDestination, world.Id),
	)
	if err != nil {
		return err
	}

	world.CreatedAt = common.GetIsoString()
	world.CoverImage = coverUrl

	_, err = ws.db.PutWrapper(ws.tableName, world, nil)
	return err
}

type WorldKey struct {
	UserId string `dynamodbav:"userId"`
	Id     string `dynamodbav:"id"`
}

func (ws *WorldService) Get(userId, id string) (*World, error) {
	key := &WorldKey{
		UserId: userId,
		Id:     id,
	}
	var result World
	_, err := ws.db.GetWrapper(ws.tableName, key, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ws *WorldService) Update(userId, id string, world *World, imageUpdated bool) error {
	key := &WorldKey{
		UserId: userId,
		Id:     id,
	}

	if imageUpdated {
		coverUrl, err := uploader.UploadCoverImage(ws.s3,
			world.CoverImage,
			fmt.Sprintf(coverImageDestination, world.Id),
		)
		if err != nil {
			return err
		}
		world.CoverImage = coverUrl
	}

	world.UpdatedAt = common.GetIsoString()

	update := expression.Set(
		expression.Name("name"), expression.Value(world.Name),
	).Set(
		expression.Name("intro"), expression.Value(world.Intro),
	).Set(
		expression.Name("genres"), expression.Value(world.Genres),
	).Set(
		expression.Name("coverImage"), expression.Value(world.CoverImage),
	).Set(
		expression.Name("updatedAt"), expression.Value(world.UpdatedAt),
	)

	_, err := ws.db.UpdateWrapper(ws.tableName, key, update)
	return err
}

func (ws *WorldService) List(userId string) (*[]World, error) {
	keyEx := expression.Key("userId").Equal(expression.Value(userId))
	var resultArr []World
	_, err := ws.db.QueryWrapper(ws.tableName, keyEx, &resultArr)
	if err != nil {
		return nil, err
	}

	return &resultArr, nil
}

func (ws *WorldService) Delete(userId, id string) error {
	key := WorldKey{
		UserId: userId,
		Id:     id,
	}

	_, err := ws.db.DeleteWrapper(ws.tableName, key)
	return err
}

func ValidateWorld(v *validator.Validator, world *World) {
	v.Check(world.Name != "", "name", "must be provided")
	v.Check(len(world.Name) < 200, "name", "must not be more than 200 characteres long")

	v.Check(validator.Unique(world.Genres), "genres", "must not contain duplicate values")
	v.Check(len(world.Genres) <= 5, "genres", "must not contain more than 5 genres")
}
