package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.GET("/profiles", getProfiles)
	server.GET("/profiles/:id", getProfileByID)
	server.POST("/profiles", createProfile)
	server.PUT("/profiles/:id", updateProfile)
	server.DELETE("/profiles/:id", deleteProfile)

	server.GET("/friends", getFriends)
	server.POST("/friends", sendInvite)
	server.PUT("/friends", confirmInvite)
}
