package routes

import (
	"BE-chat-app/cmd/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getProfiles(context *gin.Context) {
	profiles, err := models.GetAllProfile()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Can't fetch list profiles"})
		return
	}
	context.JSON(http.StatusOK, profiles)
}

func getProfileByID(context *gin.Context) {
	profileID, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Can't get profile id"})
		return
	}

	profile, err := models.GetProfileByID(profileID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Can't fetch profile"})
		return
	}
	context.JSON(http.StatusOK, profile)
}

func createProfile(context *gin.Context) {
	var profile models.Profile
	err := context.ShouldBindJSON(&profile)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Can't parse profile data"})
		return
	}

	err = profile.CreateProfile()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Can't create profile"})
		return
	}
	context.JSON(http.StatusOK, "Create successfully")
}

func updateProfile(context *gin.Context) {
	profileID, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Can't get profile id"})
		return
	}

	profile, err := models.GetProfileByID(profileID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Can't fetch profile"})
		return
	}

	var updateProfile models.Profile
	err = context.ShouldBindJSON(&updateProfile)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Can't parse JSON profile"})
		return
	}

	updateProfile.ID = profile.ID
	updateProfile.Created_at = profile.Created_at
	err = updateProfile.UpdateProfile()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Can't update profile"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Update successfully"})
}

func deleteProfile(context *gin.Context) {
	profileID, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Can't get profile id"})
		return
	}

	profile, err := models.GetProfileByID(profileID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Can't fetch profile"})
		return
	}

	err = profile.DeleteProfile()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Can't delete profile"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Delete successfully"})
}
