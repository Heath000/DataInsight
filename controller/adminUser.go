package controller

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/Heath000/fzuSE2024/model"
	"github.com/gin-gonic/gin"
)

// AdminUserController is the controller for admin user operations
type AdminUserController struct{}

// GetUserList retrieves the list of all users
func (ctrl *AdminUserController) GetUserList(c *gin.Context) {
	var user model.User

	// Retrieve all users from the database
	users, err := user.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetUser retrieves a single user by ID
func (ctrl *AdminUserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User

	// Get user by ID
	if err := user.GetFirstByID(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// CreateUser creates a new user by admin
func (ctrl *AdminUserController) CreateUser(c *gin.Context) {
	var form Signup

	// Bind the form data to the Signup struct
	if err := c.ShouldBind(&form); err == nil {
		// Check if passwords match
		if form.Password != form.Password2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password does not match with confirm password"})
			return
		}

		// Create a new user model
		user := model.User{
			Name:     form.Name,
			Email:    form.Email,
			Password: form.Password,
		}

		// Call the Create method to create the user
		if err := user.Create(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// DeleteUser deletes a user by ID (admin only)
func (ctrl *AdminUserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User

	// Call the DeleteUserByID method to delete the user
	if err := user.DeleteUserByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}



func (ctrl *AdminUserController) UpdateUser(c *gin.Context) {
	var form Signup

	// Bind the form data to the Signup struct
	if err := c.ShouldBind(&form); err == nil {
		// Check if passwords match
		if form.Password != form.Password2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password does not match with confirm password"})
			return
		}

		// Retrieve the user from the database
		id := c.Param("id")
		var user model.User
		if err := user.GetFirstByID(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// Hash the new password before storing it in the database
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Update user data with the hashed password
		user.Name = form.Name
		user.Email = form.Email
		user.Password = string(hashedPassword) // Store the hashed password

		// Call the UpdateUser method to update the user's data
		if err := user.UpdateUser(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
