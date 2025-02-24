package handlers

import (
	"github.com/KasiditR/forviz-backend-api-test/internal/database"
	"github.com/KasiditR/forviz-backend-api-test/internal/models"
	"github.com/KasiditR/forviz-backend-api-test/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"net/http"
	"time"
)

func AddBook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var book models.Book

		if err := ctx.BindJSON(&book); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if book.BookName == nil || book.Author == nil || book.Category == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid body"})
			return
		}

		var bookFound models.Book
		database.FindOne(bson.M{"book_name": book.BookName, "author": book.Author, "category": book.Category}, &bookFound)

		if bookFound.ID != bson.NilObjectID {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "this book already exits"})
			return
		}

		book.ID = bson.NewObjectID()
		_, err := database.InsertOne(&book)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Add book successfully"})
	}
}

func DeleteBook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bookId := ctx.Param("id")
		if bookId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "this request require name query"})
			return
		}

		result, err := database.DeleteOne(bson.M{"_id": bookId}, "books")
		if err != nil || result.DeletedCount == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "book not found"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Delete book successfully"})
	}
}

func EditBook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var book models.Book

		if err := ctx.BindJSON(&book); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if book.ID == bson.NilObjectID {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "bookId not found"})
			return
		}

		updateFields := bson.M{}
		if book.BookName != nil {
			updateFields["book_name"] = book.BookName
		}
		if book.Author != nil {
			updateFields["author"] = book.Author
		}
		if book.Category != nil {
			updateFields["category"] = book.Category
		}

		if len(updateFields) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "No fields to update"})
			return
		}

		updateQuery := bson.M{"$set": updateFields}

		err := database.FindByIDAndUpdate("books", book.ID, updateQuery)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Edit book successfully"})
	}
}

func GetBook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bookId := ctx.Param("id")
		if bookId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "book id has missing"})
			return
		}

		var book models.Book

		err := database.FindByID(bookId, &book)
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, book)
	}
}

func SearchBook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bookName := ctx.Query("bookName")
		authorName := ctx.Query("authorName")
		categoryName := ctx.Query("categoryName")

		if bookName == "" && authorName == "" && categoryName == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "At least one search query is required"})
			return
		}

		filter := bson.M{}
		if bookName != "" {
			filter["book_name"] = bson.M{"$regex": "^" + bookName, "$options": "i"}
		}
		if authorName != "" {
			filter["author"] = bson.M{"$regex": "^" + authorName, "$options": "i"}
		}
		if categoryName != "" {
			filter["category"] = bson.M{"$regex": "^" + categoryName, "$options": "i"}
		}

		var books []models.Book
		err := database.FindAll(filter, &books)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if len(books) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "No books found matching the search criteria"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"books": books})
	}
}

func BorrowBook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("user_id")

		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		var borrowRequest struct {
			BookID string `json:"book_id"`
		}

		if err := ctx.ShouldBindJSON(&borrowRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var book models.Book
		err := database.FindByID(borrowRequest.BookID, &book)
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		objBookId, err := bson.ObjectIDFromHex(borrowRequest.BookID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		objUserId, err := bson.ObjectIDFromHex(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		var borrowHistory models.BorrowRecord
		database.FindOne(bson.M{"book_id": objBookId, "user_id": objUserId, "returned_at": bson.M{"$exists": false}}, &borrowHistory)

		if borrowHistory.ID != bson.NilObjectID {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "ํYou already borrow this book please return first before borrow "})
			return
		}

		dateNow, err := utils.DataNow()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		dt := bson.DateTime(dateNow.UnixNano() / int64(time.Millisecond))
		borrowRecord := models.BorrowRecord{
			ID:       bson.NewObjectID(),
			BookID:   objBookId,
			UserID:   objUserId,
			Borrowed: &dt,
		}

		_, err = database.InsertOne(&borrowRecord)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Borrow book successfully"})
	}
}

func ReturnBook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("user_id")

		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		var returnRequest struct {
			BookID string `json:"book_id"`
		}

		if err := ctx.ShouldBindJSON(&returnRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		objBookId, err := bson.ObjectIDFromHex(returnRequest.BookID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		objUserId, err := bson.ObjectIDFromHex(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		dateNow, err := utils.DataNow()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		filter := bson.M{"book_id": objBookId, "user_id": objUserId, "returned_at": bson.M{"$exists": false}}
		update := bson.M{"$set": bson.M{"returned_at": dateNow}}
		err = database.FindOneAndUpdate("borrow_records", filter, update)
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Borrow record not found or already returned"})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Returned book successfully"})
	}
}

func GetMostBorrowedBooks() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type MostBorrowedBook struct {
			BookID        bson.ObjectID `json:"book_id" bson:"_id"`
			Title         string        `json:"title" bson:"title"`
			Author        string        `json:"author" bson:"author"`
			BorrowedCount int           `json:"borrowed_count" bson:"borrowed_count"`
		}
		borrowRecords := database.Get().Database.Collection("borrow_records")
		books := database.Get().Database.Collection("books")
		pipeline := mongo.Pipeline{
			{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$book_id"},
				{Key: "borrowed_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			}}},
			{{Key: "$sort", Value: bson.D{{Key: "borrowed_count", Value: -1}}}},
			{{Key: "$limit", Value: 10}},
		}

		cursor, err := borrowRecords.Aggregate(ctx, pipeline)
		if err != nil {
			log.Println("Error in aggregation:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch borrowed books"})
			return
		}

		defer cursor.Close(ctx)
		var mostBorrowed []MostBorrowedBook

		for cursor.Next(ctx) {
			var record struct {
				BookID        bson.ObjectID `bson:"_id"`
				BorrowedCount int           `bson:"borrowed_count"`
			}
			if err := cursor.Decode(&record); err != nil {
				log.Println("Error decoding record:", err)
				continue
			}

			var book struct {
				BookName string `bson:"book_name"`
				Author   string `bson:"author"`
				Category string `bson:"category"`
			}
			err := books.FindOne(ctx, bson.M{"_id": record.BookID}).Decode(&book)
			if err != nil {
				log.Println("Error fetching book details:", err)
				continue
			}

			mostBorrowed = append(mostBorrowed, MostBorrowedBook{
				BookID:        record.BookID,
				Title:         book.BookName,
				Author:        book.Author,
				BorrowedCount: record.BorrowedCount,
			})
		}

		ctx.JSON(http.StatusOK, mostBorrowed)
	}
}
