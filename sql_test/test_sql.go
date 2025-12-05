package main

// import (
// 	"gin_chat/models"

// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func main() {
// 	db, err := gorm.Open(mysql.Open("root:123456a@tcp(127.0.0.1:3306)/gin_chat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect database")
// 	}

// 	// Migrate the schema
// 	// 自动迁移不需要具体数据，只需要结构体字段即可
// 	db.AutoMigrate(&models.User_Basic{})
// 	db.AutoMigrate(&models.Message{})
// 	db.AutoMigrate(&models.Contact{})
// 	db.AutoMigrate(&models.Group{})

// 	// Create
// 	// db.Create(&models.User_Basic{Username: "John", Password: "123456"})

// 	// // Read
// 	// var user models.User_Basic
// 	// // db.First(&user, 1) // find user with integer primary key
// 	// db.First(&user, "username = ?", "John") // find user with username John

// 	// // Update - update user's age to 35
// 	// db.Model(&user).Update("Age", 35)
// 	// // Update - update multiple fields
// 	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
// 	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

// 	// Delete - delete product
// 	// db.Delete(&product, 1)

// }
