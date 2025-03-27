package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Bio         string `json:"bio,omitempty"`
	Interests   string `json:"interests,omitempty"`
	WorkHistory string `json:"workHistory,omitempty"`
	Education   string `json:"education,omitempty"`
	Certificates string `json:"certificates,omitempty"`
}

var users = make(map[string]User)

func main() {
	r := gin.Default()

	r.GET("/users", getUsers)
	r.GET("/users/:id", getUserByID)
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	r.Run(":8080")
}

func getUsers(c *gin.Context) {
	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	c.JSON(http.StatusOK, userList)
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	if user, exists := users[id]; exists {
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = uuid.New().String()
	if user.Certificates == "" {
		user.Certificates = "[]" // Default value
	}
	users[user.ID] = user
	c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	if _, exists := users[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedUser.ID = id // Keep existing ID
	users[id] = updatedUser
	c.JSON(http.StatusOK, updatedUser)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	if _, exists := users[id]; exists {
		delete(users, id)
		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
}
