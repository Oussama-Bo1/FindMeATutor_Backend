package MongoDB

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                   primitive.ObjectID `json:"_id" bson:"_id"`
	AboutMe              string             `json:"aboutMe" bson:"aboutMe"`
	Address              Address            `json:"address" bson:"address"`
	BirthDate            primitive.DateTime `json:"birthDate" bson:"birthDate"`
	CurrentStudents      []string           `json:"currentStudents" bson:"currentStudents"`
	CurrentTeachers      []string           `json:"currentTeachers" bson:"currentTeachers"`
	Email                string             `json:"email" bson:"email"`
	FirstName            string             `json:"firstName" bson:"firstName"`
	IdentityVerified     bool               `json:"identity_verified" bson:"identity_verified"`
	LastName             string             `json:"lastName" bson:"lastName"`
	MemberSince          primitive.DateTime `json:"memberSince" bson:"memberSince"`
	Password             string             `json:"password" bson:"password"`
	Phone                string             `json:"phone" bson:"phone"`
	PictureUrl           string             `json:"picture_url" bson:"picture_url"`
	PricePerHour         int                `json:"price_per_hour" bson:"price_per_hour"`
	Reviews              []Review           `json:"reviews" bson:"reviews"`
	Role                 string             `json:"role" bson:"role"`
	Searching            bool               `json:"searching" bson:"searching"`
	SearchingFor         []string           `json:"searchingFor" bson:"searchingFor"`
	SearchingForSubjects []string           `json:"searchingForSubjects" bson:"searchingForSubjects"`
	Skills               []string           `json:"skills" bson:"skills"`
	Username             string             `json:"username" bson:"username"`
}

type Address struct {
	City   string `json:"city" bson:"city"`
	State  string `json:"state" bson:"state"`
	Street string `json:"street" bson:"street"`
	Zip    string `json:"zip" bson:"zip"`
}

type Review struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Date       primitive.DateTime `json:"date" bson:"date"`
	Rating     int                `json:"rating" bson:"rating"`
	Review     string             `json:"review" bson:"review"`
	Reviewee   string             `json:"reviewee" bson:"reviewee"`
	RevieweeID primitive.ObjectID `json:"reviewee_id" bson:"reviewee_id"`
	Reviewer   string             `json:"reviewer" bson:"reviewer"`
	ReviewerID primitive.ObjectID `json:"reviewer_id" bson:"reviewer_id"`
}
