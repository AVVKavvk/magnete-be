package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name         string             `json:"name" validate:"required"`
    Phone        string             `json:"phone" validate:"required"`
    Aadhaar      string             `json:"aadhaar" validate:"required"`
    Photo        string             `json:"photo" validate:"required"`
    AadhaarPhoto string             `json:"aadhaar_photo" validate:"required"`
    FatherName   string             `json:"father_name" validate:"required"`
    Address      string             `json:"address" validate:"required"`
    RegisterDate string             `json:"register_date"`
    SeatNumber   int                `json:"seat_number"`
    AmountPaid   bool               `json:"amount_paid"`
    IsActive     bool               `json:"is_active"`
}
