package authentication

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func encryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func isCorrectPassword(password string, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

var validate *validator.Validate

func main() {
	validate = validator.New()
}

func ValidateIsEmail(id string) bool {
	err := validate.Var(id, "required,email")
	if err != nil {
		return false
	}
	return true
}

type MongoRepository struct {
	mongoDbClient *mongo.Client
}

func NewMongoRepository(mongoDbClient *mongo.Client) *MongoRepository {
	return &MongoRepository{
		mongoDbClient: mongoDbClient,
	}
}

func (m *MongoRepository) CreateNewUser(ctx context.Context, userDetail *UserCredential) error {
	// check if user already exists
	// if yes, return error
	// if no, create user

	res := m.mongoDbClient.Database("test").Collection("users").FindOne(ctx, bson.D{{"_id", userDetail.Id}})
	if res.Err() != nil {
		return res.Err()
	}

	// encrypt password
	password, err := encryptPassword(userDetail.Password)
	if err != nil {
		return err
	}

	userDetail.Password = password

	if err != nil {
		return err
	}

	_, err = m.mongoDbClient.Database("test").Collection("users").InsertOne(ctx, userDetail)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoRepository) GetUserCredential(ctx context.Context, userId UserID) (userCredential *UserCredential, err error) {

	res := m.mongoDbClient.Database("test").Collection("users").FindOne(ctx, bson.D{{"_id", string(userId)}})
	if res.Err() != nil {
		return nil, res.Err()
	}

	err = res.Decode(&userCredential)
	if err != nil {
		return nil, err
	}

	return
}

func (m *MongoRepository) GetUserDetail(ctx context.Context, userId UserID) (*UserDetail, error) {
	res := m.mongoDbClient.Database("test").Collection("users").FindOne(ctx, bson.D{{"_id", string(userId)}})

	if res.Err() != nil {
		return nil, res.Err()
	}

	var userDetail *UserDetail
	err := res.Decode(&userDetail)
	if err != nil {
		return nil, err
	}

	return userDetail, nil
}

func (m *MongoRepository) DeleteUser(ctx context.Context, userId UserID) error {
	res, err := m.mongoDbClient.Database("test").Collection("users").UpdateOne(ctx, bson.D{{"_id", string(userId)}}, bson.D{{"$set", bson.D{{"deleted", true}}}})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (m *MongoRepository) isValidPassword(ctx context.Context, password string, encryptedPassword string) bool {
	return isCorrectPassword(password, encryptedPassword)
}
