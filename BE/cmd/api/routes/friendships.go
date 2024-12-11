package routes

import (
	"BE-chat-app/cmd/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getFriends(context *gin.Context) {
	userID, err := uuid.Parse("870f9f7e-56c6-4a97-864e-712ce0d4967d")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Can't get user id"})
		return
	}

	friends, err := models.GetAllFriends(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Can't fetch friendIDs"})
		return
	}
	context.JSON(http.StatusOK, friends)
}

func sendInvite(context *gin.Context) {
	var friendship models.Friendship

	//Get userID by login
	userID, err := uuid.Parse("897da15a-0e28-4dec-9870-69d677db64b3")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Can't get user id"})
		return
	}

	//Get profileID by click add friend button
	profileID, err := uuid.Parse("870f9f7e-56c6-4a97-864e-712ce0d4967d")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Can't get profile id"})
		return
	}

	//Check invite send before

	//Create invite add friend
	err = friendship.SendInvite(userID, profileID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Can't create invite"})
		return
	}
	context.JSON(http.StatusOK, "Invite friend successfully")
}

func confirmInvite(context *gin.Context) {
	//Get userID who received the invite
	userID, err := uuid.Parse("f4c8f4f8-832d-4766-ac4d-b3f7a9b13654")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Can't get user id"})
		return
	}

	//Get profileID who invited the send
	profileID, err := uuid.Parse("870f9f7e-56c6-4a97-864e-712ce0d4967d")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Can't get profile id"})
		return
	}

	//Get invite from requester, receiver
	invite, err := models.GetInvite(profileID, userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Can't fetch invite"})
		return
	}

	//Get status from body request
	status := "Accepted"

	var friendship models.Friendship
	friendship.ID = invite.ID
	friendship.Created_at = invite.Created_at

	//Accepted or Declined the invite
	err = friendship.ConfirmInvite(status)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Can't update invite"})
		return
	}
	context.JSON(http.StatusOK, "Invite friend is accepted or decline")
}
