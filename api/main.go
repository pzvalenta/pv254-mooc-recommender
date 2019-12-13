package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ZaxCZ/docker-mongo/api/internal"
	cors "github.com/rs/cors/wrapper/gin"

	"github.com/gin-gonic/gin"
)

func main() {
	state, err := internal.NewState()
	if err != nil {
		fmt.Printf("error creating app state: %v", err)
		os.Exit(1)
	}

	router := gin.Default()
	router.Use(cors.AllowAll())
	public := router.Group("/api")
	{
		public.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Welcome to sample dockerized golang api")
		})

		public.GET("/randomRecommending", state.RandomCourse)
		public.GET("/taxonomyRecommending", state.TaxonomyCourses)
		public.GET("/overfittingRecommending", state.OverfittingCourses)
		public.GET("/categoryRecommending", state.CategoryRecommending)
		public.GET("/getbyid/:id", state.GetCoursebByID)
		public.GET("/courses", state.GetCoursesByQuery)
		public.GET("/subjects", state.GetAllSubjects)
		public.GET("/categories", state.GetAllCategories)
		public.GET("/generalRecommending", state.GeneralModelCourses)
		public.POST("/enrollUser/:authId/:courseId",state.EnrollUser)
		public.POST("/removeUserEnrollement/:authId/:courseId",state.RemoveUserEnrollment)
		public.POST("/createUser/:name/:authId", state.CreateUser)
		public.GET("/getUserById/:id", state.GetUserByID)
		public.GET("/getUserByAuthId/:authId", state.GetUserByAuthID)
		public.GET("/getUserCoursesByAuth/:authId", state.GetUserCourses)
		public.GET("/getUserCourses/:id", state.GetUserCourses)
	}

	err = router.Run()
	if err != nil {
		fmt.Printf("router stopped: %v", err)
		os.Exit(1)
	}
}
