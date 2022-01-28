package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/cfloress-gb-cl/final-project-bootcamp/repository/user"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type config struct {
	Port        int    `env:"MONGO_PORT" envDefault:"27017"`
	Host        string `env:"MONGO_HOST" envDefault:"172.17.0.2"`
	DefaultDB   string `env:"MONGO_DEFAULTDB" envDefault:"globant"`
	DefaultCOLL string `env:"MONGO_DEFAULTCOLL" envDefault:"user"`
}

//MongoUserRepository - is a mongoDB implementation of users repository
type singleMongo struct {
	db *mongo.Client
}

var collection *mongo.Collection

func initMongoRepository() (*mongo.Client, error) {

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port))
	dbLocalRef, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		fmt.Println("Error en dbConnMongo!!!:", err.Error())
		return nil, err
	}

	if err := dbLocalRef.Ping(context.Background(), nil); err != nil {
		return nil, err
	}
	fmt.Println("initMongoUserRepository...")
	collection = dbLocalRef.Database("globant").Collection("user")

	return dbLocalRef, nil
}

//NewMongoUserRepository - returns a Mongoepository type pointer
func NewMongoUserRepository() (*singleMongo, error) {

	db, err := initMongoRepository()

	if err != nil {
		fmt.Println("Error en dbConnMongo!!!:", err.Error())
		return nil, err
	}
	fmt.Println("NewMongoUserRepository...")

	return &singleMongo{
		db: db,
	}, nil
}

//Add - adds a user to the repository
func (r *singleMongo) Add(ctx context.Context, usr user.User) (int, error) {

	res, err := collection.InsertOne(ctx, usr)
	if err != nil {
		fmt.Println("Internal error:", err)
		return 0, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		fmt.Println("Cannot find User ID", id)
		return 0, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot find User ID"),
		)
	}

	filter := bson.M{"_id": bson.M{"$eq": id}}

	update := bson.M{"$inc": bson.M{"id": 1}}

	collection.UpdateOne(context.Background(), filter, update)

	fmt.Println("mongoDB inserted id-->", id.Hex())

	return usr.ID, nil
}

//GetByID - retrieves a user from the repository based on the object id from mongo!
func (r *singleMongo) GetByID(ctx context.Context, userID int) (user.User, error) {

	usr := user.User{}

	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	cursor := collection.FindOne(ctx, bson.M{"id": userID})

	if err := cursor.Decode(&usr); err != nil {
		return usr, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("no user found: %v", err),
		)
	}

	return usr, nil

}

//GetByEmail - retrieves a user from the repository based on the email address
func (r *singleMongo) GetByEmail(ctx context.Context, email string) (user.User, error) {

	var usr []user.User
	var u user.User

	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	cursor := collection.FindOne(ctx, bson.M{"email": email})

	if err := cursor.Decode(&u); err != nil {
		fmt.Println("no user found-->:", err)
		return u, nil
	}

	regData, err := json.Marshal(usr)
	if err != nil {
		fmt.Println("Error marshal to json GetByEmail:", err.Error()) // proper error handling instead of panic in your app
		return u, err

	}
	errUn := json.Unmarshal(regData, &u)
	if errUn != nil {
		fmt.Println("Error unmarshal json GetByEmail:", errUn)
		return u, errUn
	}

	return u, nil
}

//GetAll - retrieves all the users from the repository
func (r *singleMongo) GetAll(ctx context.Context) ([]user.User, error) {

	usrs := []user.User{}

	queryOptions := options.Find()
	queryOptions.SetSort(bson.M{"email": 1})
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{}, queryOptions)

	if err != nil {
		fmt.Println("Error on execute GetAll mongo query:", err.Error()) // proper error handling instead of panic in your app
		return nil, err

	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var userss user.User
		cursor.Decode(&userss)
		usrs = append(usrs, userss)
	}
	if err := cursor.Err(); err != nil {
		fmt.Println("Error on cursor.Err() GetAll mongo query:", err.Error()) // proper error handling instead of panic in your app
		return nil, err
	}
	if len(usrs) <= 0 {
		fmt.Println("no records found") // proper error handling instead of panic in your app
		return nil, errors.New("no records found")
	}
	return usrs, nil
}

//Update -  updates the information of a user
func (r *singleMongo) Update(ctx context.Context, usr user.User) error {

	filter := bson.M{"email": usr.Email}

	cursor, updateErr := collection.ReplaceOne(context.Background(), filter, usr)
	if updateErr != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update object in MongoDB: %v", updateErr),
		)
	}

	if cursor.ModifiedCount == 0 || updateErr != nil {
		fmt.Println("cursor.ModifiedCount-->", cursor.ModifiedCount)
		return errors.New("no records were updated")
	}

	return nil
}

//Delete - deletes a user from the repository
func (r *singleMongo) Delete(ctx context.Context, userID int) error {

	filter := bson.M{"id": bson.M{"$eq": userID}}
	update := bson.M{"$set": bson.M{"status": 0}}

	cursor, updateErr := collection.UpdateOne(context.Background(), filter, update)

	if updateErr != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete object in MongoDB: %v", updateErr),
		)
	}

	if cursor.ModifiedCount == 0 || updateErr != nil {
		return errors.New("no records were updated")
	}

	return nil
}

/*func main() {
	_, err := NewMongoUserRepository()
	if err != nil {
		panic(fmt.Sprintf("mongoDB connection failed: %s", err))
	}
}*/
