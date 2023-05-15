package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/athiramjayaprasad/ecommerse-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("can't find the product")
	ErrUserIdIsNotValid = errors.New("this user is not valid")
	ErrCantUpdateUser = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove the product from the cart")
	ErrCantGetItem =  errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem = errors.New("cannot update the purchase")
	
)
func AddProductToCart(ctx context.Context, productCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error{
	search_from_db, err := productCollection.Find(ctx, bson.M{"_id":productId})
	if err != nil {
		log.Panicln(err)
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	err = search_from_db.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value:  id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "user_cart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil

}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, user_id string) error {
	id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"user_cart":bson.M{"_id": productId}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil
	      
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, user_id string) error {
	id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var get_cart_items models.User
	var order_cart models.Order            

	 order_cart.Order_Id = primitive.NewObjectID()
	 order_cart.Orderd_At = time.Now()
	 order_cart.Order_Cart = make([]models.ProductUser, 0)
	 order_cart.Payment_Method.COD =true

	 unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$user_cart"}}}}
	 grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$user_cart.price"}} }}}}

	 current_results, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	 ctx.Done()
	 if err != nil {
		panic(err)
	 }

	 var get_user_cart []bson.M
	 if err = current_results.All(ctx, &get_user_cart); err != nil {
		panic(err)
	 }

	 var total_price int32

	 for _, user_item := range get_user_cart {
		price := user_item["total"]
		total_price = price.(int32)
	 }
	 order_cart.Price = int(total_price)
	 filter := bson.D{primitive.E{Key: "_id", Value: id}}
	 update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order_cart}}}}
	 _, err = userCollection.UpdateMany(ctx, filter, update)
	 if err != nil {
		log.Println(err)
	 }

	 err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&get_cart_items)
	 if err != nil {
		log.Println(err)
	 }

	 filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	 update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each":get_cart_items.User_Cart}}}
	 _, err = userCollection.UpdateOne(ctx, filter2, update2)
	 if err != nil {
		log.Println(err)
	 }
	 user_cart_empty := make([]models.ProductUser, 0)
	 filter3:= bson.D{primitive.E{Key: "_id", Value: id}}
	 update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "user_cart", Value: user_cart_empty}}}}
	 _, err = userCollection.UpdateOne(ctx, filter3, update3)
	 if err != nil {
		return ErrCantBuyCartItem
	 }
	 return nil                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      
}

func InstanceBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var product_details models.ProductUser
	var orders_detail models.Order

	orders_detail.Order_Id = primitive.NewObjectID()
	orders_detail.Orderd_At =  time.Now()
	orders_detail.Order_Cart = make([]models.ProductUser, 0)
	orders_detail.Payment_Method.COD = true

	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productId}}).Decode(&product_details)
	if err != nil {
		log.Println(err)
	}
	orders_detail.Price = int(product_details.Price)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orders_detail}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push":bson.M{"orders.$[].orders_list":product_details}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	return nil



	
}