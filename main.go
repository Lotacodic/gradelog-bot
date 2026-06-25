package main

// import (
// 	"fmt"

// 	"gradelog-sys/services"
// 	"gradelog-sys/storage"
// )

import (
	"os"

	"gradelog-sys/bot"
)

func main() {
	// // Test registration
	// err := services.RegisterStudent("John Adebayo", "john@example.com", "password123")
	// if err != nil {
	// 	fmt.Println("Registration failed:", err)
	// } else {
	// 	fmt.Println("Registration successful!")
	// }

	// // Test login with correct password
	// student, err := services.LoginStudent("john@example.com", "password123")
	// if err != nil {
	// 	fmt.Println("Login failed:", err)
	// } else {
	// 	fmt.Println("Login sucessful! Welcome,", student.FullName)
	// }

	// // Test login with wrong password
	// _, err = services.LoginStudent("john@example.com", "wrongpassword")
	// if err != nil {
	// 	fmt.Println("Login failed:", err)
	// } else {
	// 	fmt.Print("Login successful!")
	// }

	// // Load the database (or create a fresh one)
	// db, err := storage.LoadDB()
	// if err != nil {
	// 	fmt.Println("Error loading database:", err)
	// 	return
	// }
	// if len(db.Students) > 0 {
	// 	fmt.Println("Student name:", db.Students[0].FullName)
	// 	fmt.Println("Hashed password:", db.Students[0].Password)
	// } else {
	// 	fmt.Println("No students registered yet.")
	// }
	// cli.Start()

	token := os.Getenv("BOT_TOKEN")

	if token == "" {
		panic("BOT_TOKEN environment variable not set!")
	}
	bot.Start(token)
}
