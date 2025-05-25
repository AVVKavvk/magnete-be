package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	config "github.com/magnete-library/config"
	"github.com/magnete-library/models"
)

func getCollection(month string) *mongo.Collection {
    return config.MongoClient.Database("magnete").Collection(month)
}

func withTimeoutCtx() (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), 10*time.Second)
}

// HomePage handles GET /
func HomePage(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}

// CreateStudent handles POST /students/:month
func CreateStudent(c echo.Context) error {
    month := c.Param("month")
    var student models.Student

    if err := c.Bind(&student); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
    }
    if err := c.Validate(&student); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    student.ID = primitive.NewObjectID()
    student.RegisterDate = time.Now().Format(time.RFC3339)

    ctx, cancel := withTimeoutCtx()
    defer cancel()

    _, err := getCollection(month).InsertOne(ctx, student)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to insert student"})
    }

    return c.JSON(http.StatusCreated, student)
}

// GetStudents handles GET /students/:month
func GetStudents(c echo.Context) error {
    month := c.Param("month")

    ctx, cancel := withTimeoutCtx()
    defer cancel()

    cursor, err := getCollection(month).Find(ctx, bson.M{})
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch students"})
    }

    var students []models.Student
    if err := cursor.All(ctx, &students); err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error decoding results"})
    }

    return c.JSON(http.StatusOK, students)
}

// SearchStudents handles GET /students/:month/search?q=...
func SearchStudents(c echo.Context) error {
    month := c.Param("month")
    query := c.QueryParam("q")

    filter := bson.M{
        "$or": []bson.M{
            {"name": bson.M{"$regex": query, "$options": "i"}},
            {"phone": bson.M{"$regex": query, "$options": "i"}},
            {"aadhaar": bson.M{"$regex": query, "$options": "i"}},
        },
    }

    ctx, cancel := withTimeoutCtx()
    defer cancel()

    cursor, err := getCollection(month).Find(ctx, filter)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Search failed"})
    }

    var students []models.Student
    if err := cursor.All(ctx, &students); err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error decoding search results"})
    }

    return c.JSON(http.StatusOK, students)
}

// UpdateStudent handles PUT /students/:month/:id
func UpdateStudent(c echo.Context) error {
    month := c.Param("month")
    id := c.Param("id")

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid student ID"})
    }

    var updated models.Student
    if err := c.Bind(&updated); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid data"})
    }

    ctx, cancel := withTimeoutCtx()
    defer cancel()

    update := bson.M{"$set": updated}
    _, err = getCollection(month).UpdateByID(ctx, objectID, update)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Update failed"})
    }

    return c.JSON(http.StatusOK, echo.Map{"message": "Student updated"})
}

// UpdatePayment handles PATCH /students/:month/:id/payment
func UpdatePayment(c echo.Context) error {
    month := c.Param("month")
    id := c.Param("id")

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid student ID"})
    }

    var payload struct {
        AmountPaid bool `json:"amount_paid"`
    }

    if err := c.Bind(&payload); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
    }

    ctx, cancel := withTimeoutCtx()
    defer cancel()

    fmt.Println(month,id)
    _, err = getCollection(month).UpdateByID(ctx, objectID, bson.M{
        "$set": bson.M{"amountpaid": payload.AmountPaid},
    })
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Payment update failed"})
    }

    return c.JSON(http.StatusOK, echo.Map{"message": "Payment status updated"})
}
// UpdateSeatNumber handles PATCH /students/:month/:id/seat
func UpdateSeatNumber(c echo.Context) error {
    month := c.Param("month")
    id := c.Param("id")

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid student ID"})
    }

    var payload struct {
        SeatNumber int `json:"seat_number"`
    }

    if err := c.Bind(&payload); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
    }

    ctx, cancel := withTimeoutCtx()
    defer cancel()

    fmt.Println(month, id)
    _, err = getCollection(month).UpdateByID(ctx, objectID, bson.M{
        "$set": bson.M{"seatnumber": payload.SeatNumber},
    })
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Seat number update failed"})
    }

    return c.JSON(http.StatusOK, echo.Map{"message": "Seat number updated"})
}

func ToggleActiveStatus(c echo.Context) error {
    month := c.Param("month")
    id := c.Param("id")

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid student ID"})
    }

    ctx, cancel := withTimeoutCtx()
    defer cancel()

    _, err = getCollection(month).UpdateByID(ctx, objectID, bson.M{
        "$set": bson.M{"isactive": false},
    })
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Status update failed"})
    }

    return c.JSON(http.StatusOK, echo.Map{"message": "Status updated"})
}


// MigrateMonth handles POST /students/migrate
func MigrateMonth(c echo.Context) error {
    var payload struct {
        FromMonth string `json:"from"`
        ToMonth   string `json:"to"`
    }

    if err := c.Bind(&payload); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
    }

    ctx, cancel := withTimeoutCtx()
    defer cancel()

    cursor, err := getCollection(payload.FromMonth).Find(ctx, bson.M{})
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Migration fetch failed"})
    }

    var students []interface{}
    for cursor.Next(ctx) {
        var student models.Student
        if err := cursor.Decode(&student); err == nil {
            student.ID = primitive.NewObjectID()
            student.RegisterDate = time.Now().Format(time.RFC3339)
            student.AmountPaid = false
            students = append(students, student)
        }
    }

    if len(students) == 0 {
        return c.JSON(http.StatusOK, echo.Map{"message": "No students to migrate"})
    }

    _, err = getCollection(payload.ToMonth).InsertMany(ctx, students)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Migration insert failed"})
    }

    return c.JSON(http.StatusOK, echo.Map{"message": "Migration successful", "count": len(students)})
}

// CreateMonth handles POST /months
func CreateMonth(c echo.Context) error {
    var payload struct {
        Month string `json:"month" validate:"required"`
        Year  int    `json:"year" validate:"required"`
    }

    if err := c.Bind(&payload); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
    }

    payload.Month = strings.Title(strings.ToLower(payload.Month)) // normalize name

    collectionName := fmt.Sprintf("%s-%d", payload.Month, payload.Year)
    doc := bson.M{
        "month":           payload.Month,
        "year":            payload.Year,
        "collection_name": collectionName,
        "created_at":      time.Now(),
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := config.MongoClient.Database("magnete").Collection("months").InsertOne(ctx, doc)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create month"})
    }

    return c.JSON(http.StatusCreated, echo.Map{"message": "Month created", "collection": collectionName})
}

func GetMonths(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := config.MongoClient.Database("magnete").Collection("months").Find(ctx, bson.M{})
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch months"})
    }

    var months []bson.M
    if err := cursor.All(ctx, &months); err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error decoding months"})
    }

    return c.JSON(http.StatusOK, months)
}
func ListCollections(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collections, err := config.MongoClient.Database("magnete").ListCollectionNames(ctx, bson.M{})
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to list collections"})
    }

    return c.JSON(http.StatusOK, echo.Map{"collections": collections})
}