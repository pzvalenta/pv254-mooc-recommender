package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ZaxCZ/docker-mongo/api/internal"

	"github.com/gin-gonic/gin"
)

func main() {
	state, err := internal.NewState("5dceb44288861f034fc60b16")
	if err != nil {
		fmt.Printf("error creating app state: %v", err)
		os.Exit(1)
	}

	router := gin.Default()
	public := router.Group("/api")
	{
		public.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Welcome to sample dockerized golang api")
		})

		public.GET("/random", state.RandomCourse)
		public.GET("/taxonomy", state.TaxonomyCourses)
		public.GET("/overfitting", state.OverfittingCourses)
		public.GET("/getbyid/:id", state.GetCoursebByID)
		public.GET("/getbysubject/:subject/*page",state.GetCoursesBySubjectPaged)
	}

	err = router.Run()
	if err != nil {
		fmt.Printf("router stopped: %v", err)
		os.Exit(1)
	}
}
