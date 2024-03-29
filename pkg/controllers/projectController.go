package controllers

import (
	"auth-service/pkg/database/mongodb/models"
	"auth-service/pkg/database/mongodb/repository"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProject(c *gin.Context) (err error) {
	var project models.Project

	if err = c.BindJSON(&project); err != nil {
		return
	}

	project.Id = primitive.NewObjectID()

	// Create the project in the database
	err = repository.CreateProject(project)
	if err != nil {
		return fmt.Errorf("error occurred while creating project: %v", err)
	}

	return
}

// Function that updates a project
func UpdateProject(c *gin.Context) (err error) {
	var project models.Project

	if err = c.BindJSON(&project); err != nil {
		return
	}

	// Extract project ID from the request
	projectID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return
	}
	// Update the project in the database
	err = repository.UpdateProject(projectID, project)
	if err != nil {
		return errors.New("error occurred while updating project")
	}

	return
}

// Function that returns all projects
func GetAllProjects(c *gin.Context) (projects []models.Project, err error) {
	projects, err = repository.GetAllProjects()
	return
}
