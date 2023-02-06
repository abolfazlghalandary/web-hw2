// package main

// import (
// 	"fmt"
// 	"time"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// type user_account struct {
// 	gorm.Model
// 	ID        uint           `gorm:"primaryKey"`
// 	Email         string
// 	Phone_number  string
// 	Gender        string
// 	First_name    string
// 	Last_name     string
// 	Password_hash string
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
// }

// type unauthorized_token struct {
// 	gorm.Model
// 	ID        uint           `gorm:"primaryKey"`
// 	token      string
// 	expiration time.Time
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
// }

// func main() {
// 	dsn := "host=localhost user=postgres password=dev dbname=dev port=5432 sslmode=disable TimeZone=Asia/Shanghai"
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// 	user1 := &user_account{
// 		Email:        "asfd",
// 		Phone_number: "sfdsfg",
// 	}

	
// 	token1 := &unauthorized_token{
// 		token: "sdfs",
// 		expiration: time.Now(),
// 	}


// 	// Migrate the schema
// 	db.AutoMigrate(&user_account{})

// 	fmt.Println(user1)

// 	// var result = db.Create(&user1)
// 	var result = db.Create(&token1)

// 	fmt.Println(err)
// 	fmt.Println(result.Error)
// }