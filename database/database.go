package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var client *mongo.Client
var collection *mongo.Collection

func init() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	godotenv.Load()
	mongoDB := os.Getenv("MONGODB_CONNECTION_STRING")
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoDB))
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database("car_driven_db").Collection("cars")
}

func SaveCarDrivenData(c Car) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c.DrivenDate = time.Now()

	//check if car already exists
	res, err := collection.FindOne(ctx, bson.M{"kenteken": c.Kenteken}).DecodeBytes()
	if res != nil {
		return fmt.Errorf("car already exists")
	}

	_, err = collection.InsertOne(ctx, c)
	if err != nil {
		return fmt.Errorf("failed to insert car data: %v", err)
	}

	return nil
}

func GetAllCars() ([]Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find cars: %v", err)
	}
	defer cursor.Close(ctx)

	var cars []Car
	if err := cursor.All(ctx, &cars); err != nil {
		return nil, fmt.Errorf("failed to decode cars: %v", err)
	}

	return cars, nil
}

func GetCarByKenteken(kenteken string) (Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var car Car
	err := collection.FindOne(ctx, bson.M{"kenteken": kenteken}).Decode(&car)
	if err != nil {
		return Car{}, fmt.Errorf("failed to find car by kenteken: %v", err)
	}

	return car, nil
}

func AddPictureToCar(kenteken, pictureURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"kenteken": kenteken}, bson.M{"$set": bson.M{"picture": pictureURL}})
	if err != nil {
		return fmt.Errorf("failed to add picture to car: %v", err)
	}

	return nil
}

func DeletePicterOfCar(kenteken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"kenteken": kenteken}, bson.M{"$unset": bson.M{"picture": ""}})
	if err != nil {
		return fmt.Errorf("failed to delete picture of car: %v", err)
	}

	return nil
}

func GetCarDataFromRDWAPI(kenteken string) (Car, error) {
	// API endpoint
	apiURL := fmt.Sprintf("https://opendata.rdw.nl/resource/m9d7-ebf2.json?$limit=1&kenteken=%s", url.QueryEscape(kenteken))

	// Create a new request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return Car{}, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the app token
	app_token := os.Getenv("RDW_APP_TOKEN")
	req.Header.Set("$$app_token", app_token)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Car{}, fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Car{}, fmt.Errorf("failed to read response body: %v", err)
	}
	fmt.Println("Raw response:", string(body))

	// Parse the JSON response
	var apiResponses []RDWCarResponse
	err = json.Unmarshal(body, &apiResponses)
	if err != nil {
		return Car{}, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Check if we got a response
	if len(apiResponses) == 0 {
		return Car{}, fmt.Errorf("no car data found for kenteken: %s", kenteken)
	}

	apiResponse := apiResponses[0]
	// Convert API response to Car struct
	car := Car{
		Kenteken:                 apiResponse.Kenteken,
		Merk:                     apiResponse.Merk,
		Handelsbenaming:          apiResponse.Handelsbenaming,
		Variant:                  apiResponse.Variant,
		Uitvoering:               apiResponse.Uitvoering,
		Inrichting:               apiResponse.Inrichting,
		EersteKleur:              apiResponse.EersteKleur,
		TweedeKleur:              apiResponse.TweedeKleur,
		AantalZitplaatsen:        convertStringToInt(apiResponse.AantalZitplaatsen),
		AantalDeuren:             convertStringToInt(apiResponse.AantalDeuren),
		AantalWielen:             convertStringToInt(apiResponse.AantalWielen),
		AantalCilinders:          convertStringToInt(apiResponse.AantalCilinders),
		Cilinderinhoud:           convertStringToInt(apiResponse.Cilinderinhoud),
		Catalogusprijs:           convertStringToInt(apiResponse.Catalogusprijs),
		Lengte:                   convertStringToInt(apiResponse.Lengte),
		Wielbasis:                convertStringToInt(apiResponse.Wielbasis),
		DatumEersteToelating:     apiResponse.DatumEersteToelating,
		Zuinigheidsclassificatie: apiResponse.Zuinigheidsclassificatie,
	}
	return car, nil
}

func convertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func UpdateCarData(kenteken string, c Car) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"kenteken": kenteken}, bson.M{"$set": c})
	if err != nil {
		return fmt.Errorf("failed to update car data: %v", err)
	}

	return nil
}

func DeleteCar(kenteken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"kenteken": kenteken})
	if err != nil {
		return fmt.Errorf("failed to delete car: %v", err)
	}

	return nil
}
