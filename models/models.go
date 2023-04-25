package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID  primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name *string	`json:"first_name" validate:"required, min=2, max=30"`
	Last_Name *string	`json:"last_name" validate:"required, min=2, max=30"`
	Password *string	`json:"password" validate:"required, min=6"`
	Email *string	`json:"email" validate:"required, email"`
	Phone *string	`json:"phone" validate:"required"`
	Token *string	`json:"token"`
	Refresh_Token *string 	`json:"refresh_token"`
	Created_At time.Time	`json:"created_at"`
	Updated_At time.Time	`json:"updated_at"`
	User_Id string	`json:"user_id"`
	User_Cart []ProductUser	`json:"user_cart" bson:"user_cart"`
	Address_Details []Address	`json:"address_details" bson:"address_details"`
	Order_Status []Order	`json:"order_status" bson:"order_status"`
}

type Product struct {
	Product_Id	primitive.ObjectID	`json:"_id" bson:"_id"`
	Product_Name *string	`json:"project_name"`
	Price	uint64	`json:"price"`
	Rating 	uint8	`json:"rating"`
	Image	*string	`json:"image"`
}

type ProductUser struct {
	Product_Id	primitive.ObjectID `json:"_id" bson:"_id"`
	Product_Name	*string	`json:"product_name" bson:"product_name"`
	Price		uint64	`json:"price"`
	Rating		uint8	`json:"rating"`
	Image		*string	`json:"image"`

}

type Address struct {
	Address_Id	primitive.ObjectID	`json:"_id" bson:"_id"`
	House	*string	`json:"house"`
	Street	*string	`json:"street"`
	City 	*string	`json:"city"`
	Pincode	*string	`json:"pincode"`

}

type Order struct {
	Order_Id	primitive.ObjectID	`json:"_id" bson:"_id"`
	Order_Cart	[]ProductUser	`json:"order_cart" bson:"order_cart"`
	Orderd_At	time.Time	`json:"ordered_at" bson:"ordered_at"`
	Price		int	`json:"price" bson:"price"`
	Discount	*int	`json:"discount" bson:"discount"`
	Payment_Method Payment	`json:"payment" bson:"payment"`
}

type Payment struct {
	Digital bool 
	COD bool       

}                                                                                                                                                        