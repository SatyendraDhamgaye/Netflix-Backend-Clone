package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/SatyendraDhamgaye/mongoDbApi/helpers"
	"github.com/SatyendraDhamgaye/mongoDbApi/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connection string from created database from mongodb
const connectionString = "mongodb+srv://satye:Lavra@golangprac.3cu4y7i.mongodb.net/?retryWrites=true&w=majority"

// database details
const dbname = "netflix"
const colName = "watchlist"

// importent part
var collection *mongo.Collection

// connection with mongo DB
func init() {
	//client option list like jdbs in java
	clientOption := options.Client().ApplyURI(connectionString)

	//connect to mongoDb
	//background works in background anutomatically on repeat
	//TODO works in foreground
	client, err := mongo.Connect(context.TODO(), clientOption)

	helpers.ErrorCatcher(err)

	fmt.Println("MongoDB connection success")

	collection = client.Database(dbname).Collection(colName)

	//collection instance
	fmt.Println("Collection instance is ready")
}

//mongodb helper methods should go in seperate file

func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)

	helpers.ErrorCatcher(err)
	fmt.Println("Inserted 1 movie in the database with Id: ", inserted.InsertedID)
}

func updateOneMovie(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)
	helpers.ErrorCatcher(err)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	helpers.ErrorCatcher(err)

	fmt.Println("Modified count: ", result.ModifiedCount)
}

func deleteOneMovie(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)
	helpers.ErrorCatcher(err)
	filter := bson.M{"_id": id}
	deletecount, err := collection.DeleteOne(context.Background(), filter)
	helpers.ErrorCatcher(err)
	fmt.Println("Movie got deleted with delete count: ", deletecount)
}

func deleteAllMovie() int64 {
	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("NUmber of movies delete: ", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

func getAllMovies() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	helpers.ErrorCatcher(err)

	var movies []primitive.M

	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		helpers.ErrorCatcher(err)
		movies = append(movies, movie)
	}
	defer cur.Close(context.Background())
	return movies
}

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix

	err := json.NewDecoder(r.Body).Decode(&movie)
	helpers.ErrorCatcher(err)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)
}
