package characters

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	common "github.com/jplindgren/rpg-vault/internal"
	"github.com/jplindgren/rpg-vault/internal/clients"
)

//const characterTable = "rpg_characters"

type CharacterKey struct {
	WorldId string `dynamodbav:"worldId"`
	Id      string `dynamodbav:"id"`
}

type CharacterService struct {
	db        *clients.DynamoDbClientWrapper
	tableName string
}

func New(db *clients.DynamoDbClientWrapper, tableName string) *CharacterService {
	return &CharacterService{
		db:        db,
		tableName: tableName,
	}
}

func (cs *CharacterService) Insert(character *Character) error {
	character.Id = common.GenerateToken()
	character.CreatedAt = common.GetIsoString()
	character.UpdatedAt = ""

	_, err := cs.db.PutWrapper(cs.tableName, &character, nil)
	return err
}

func (cs *CharacterService) Get(worldId, id string) (*Character, error) {
	key := CharacterKey{
		WorldId: worldId,
		Id:      id,
	}

	var result Character
	_, err := cs.db.GetWrapper(cs.tableName, key, &result)
	if err != nil {
		return nil, err
	}

	if result.AttributesJSON != "" {
		err = json.Unmarshal([]byte(result.AttributesJSON), &result.Attributes)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (cs *CharacterService) List(worldId string) (*[]Character, error) {
	keyEx := expression.Key("worldId").Equal(expression.Value(worldId))

	var resultArr []Character
	_, err := cs.db.QueryWrapper(cs.tableName, keyEx, &resultArr)
	if err != nil {
		return nil, err
	}

	return &resultArr, nil
}

func (cs *CharacterService) ListKeys(worldId string) ([]map[string]string, error) {
	keyEx := expression.Key("worldId").Equal(expression.Value(worldId))
	proj := expression.NamesList(expression.Name("id"), expression.Name("worldId"))

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyEx).
		WithProjection(proj).
		Build()
	if err != nil {
		return nil, err
	}

	var resultArr []map[string]string
	_, err = cs.db.QueryWithExpressionWrapper(cs.tableName, expr, &resultArr)
	if err != nil {
		return nil, err
	}

	return resultArr, nil
}

func (cs *CharacterService) Update(worldId, id string, uc *Character) error {
	key := CharacterKey{
		WorldId: worldId,
		Id:      id,
	}

	uc.UpdatedAt = common.GetIsoString()

	update := expression.Set(
		expression.Name("attributesJSON"),
		expression.Value(uc.AttributesJSON),
	).Set(
		expression.Name("name"),
		expression.Value(uc.Name),
	).Set(
		expression.Name("intro"),
		expression.Value(uc.Intro),
	).Set(
		expression.Name("updatedAt"),
		expression.Value(uc.UpdatedAt),
	).Set(
		expression.Name("coverImage"),
		expression.Value(uc.CoverImage),
	)

	_, err := cs.db.UpdateWrapper(cs.tableName, key, update)
	return err
}

func (cs *CharacterService) Delete(worldId, id string) error {
	key := CharacterKey{
		WorldId: worldId,
		Id:      id,
	}
	_, err := cs.db.DeleteWrapper(cs.tableName, key)
	return err
}

func (cs *CharacterService) DeleteByKeys(keys []map[string]string) error {
	if len(keys) == 0 {
		return nil
	}
	_, err := cs.db.BatchDeleteWrapper(cs.tableName, keys)
	return err
}
