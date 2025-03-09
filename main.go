package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Car represents the structure of a car in our database
type Car struct {
	ID                       primitive.ObjectID `bson:"_id,omitempty"`
	Kenteken                 string             `bson:"kenteken"`
	Merk                     string             `bson:"merk"`
	Handelsbenaming          string             `bson:"handelsbenaming"`
	Variant                  string             `bson:"variant"`
	Uitvoering               string             `bson:"uitvoering"`
	Inrichting               string             `bson:"inrichting"`
	EersteKleur              string             `bson:"eerste_kleur"`
	TweedeKleur              string             `bson:"tweede_kleur"`
	AantalZitplaatsen        int                `bson:"aantal_zitplaatsen"`
	AantalDeuren             int                `bson:"aantal_deuren"`
	AantalWielen             int                `bson:"aantal_wielen"`
	AantalCilinders          int                `bson:"aantal_cilinders"`
	Cilinderinhoud           int                `bson:"cilinderinhoud"`
	Catalogusprijs           int                `bson:"catalogusprijs"`
	Lengte                   int                `bson:"lengte"`
	Wielbasis                int                `bson:"wielbasis"`
	DatumEersteToelating     string             `bson:"datum_eerste_toelating"`
	Zuinigheidsclassificatie string             `bson:"zuinigheidsclassificatie"`
	Picture                  string             `bson:"picture,omitempty"`
	DrivenDate               time.Time          `bson:"driven_date"`
}

var client *mongo.Client
var collection *mongo.Collection

func init() {
	// Initialize MongoDB connection
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

	// Get the collection
	collection = client.Database("car_driven_db").Collection("cars")
}

// saveCarDrivenData saves the car data to the database
func saveCarDrivenData(c Car) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c.DrivenDate = time.Now()

	_, err := collection.InsertOne(ctx, c)
	if err != nil {
		return fmt.Errorf("failed to insert car data: %v", err)
	}

	return nil
}

// getAllCars retrieves all cars from the database
func getAllCars() ([]Car, error) {
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

// getCarByKenteken retrieves a car by its kenteken (license plate)
func getCarByKenteken(kenteken string) (Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var car Car
	err := collection.FindOne(ctx, bson.M{"kenteken": kenteken}).Decode(&car)
	if err != nil {
		return Car{}, fmt.Errorf("failed to find car by kenteken: %v", err)
	}

	return car, nil
}

// addPictureToCar adds a picture to an existing car
func addPictureToCar(kenteken, pictureURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"kenteken": kenteken}, bson.M{"$set": bson.M{"picture": pictureURL}})
	if err != nil {
		return fmt.Errorf("failed to add picture to car: %v", err)
	}

	return nil
}

// getCarDataFromRDWAPI retrieves car data from the RDW API using the given kenteken
func getCarDataFromRDWAPI(kenteken string) (Car, error) {
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

// Helper functions to safely extract values from the API response
func getString(data map[string]interface{}, key string) string {
	if value, ok := data[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

func getInt(data map[string]interface{}, key string) int {
	if value, ok := data[key]; ok {
		if floatValue, ok := value.(float64); ok {
			return int(floatValue)
		}
	}
	return 0
}

func main() {
	godotenv.Load()
	// Initialize Jet views
	views := jet.NewSet(
		jet.NewOSFileSystemLoader("./views"),
		jet.InDevelopmentMode(),
	)

	e := echo.New()

	// Set up routes
	e.GET("/", func(c echo.Context) error {
		cars, err := getAllCars()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to get cars")
		}

		v, err := views.GetTemplate("index.jet")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to load template")
		}

		var buf bytes.Buffer
		if err := v.Execute(&buf, nil, cars); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to execute template")
		}

		return c.HTML(http.StatusOK, buf.String())
	})

	e.POST("/add-car", func(c echo.Context) error {
		kenteken := c.FormValue("kenteken")

		// Get car data from RDW API
		car, err := getCarDataFromRDWAPI(kenteken)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get car data: %v", err))
		}

		// Save the car data
		if err := saveCarDrivenData(car); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to save car data")
		}

		return c.Redirect(http.StatusSeeOther, "/")
	})

	e.POST("/add-picture", func(c echo.Context) error {
		kenteken := c.FormValue("kenteken")
		pictureURL := c.FormValue("picture_url")

		if err := addPictureToCar(kenteken, pictureURL); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to add picture")
		}

		return c.Redirect(http.StatusSeeOther, "/")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
