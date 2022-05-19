package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email         string             `json:"email" bson:"email"`
	Password      string             `json:"password" bson:"password"`
	RefreshTokens primitive.A        `json:"tokens" bson:"tokens"`
	//RegisteredAt time.Time            `json:"registeredAt" bson:"registeredAt"`
	//LastVisitAt  time.Time            `json:"lastVisitAt" bson:"lastVisitAt"`
	//Verification Verification         `json:"verification" bson:"verification"`
}
