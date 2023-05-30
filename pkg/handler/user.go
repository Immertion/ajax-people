package handler

import (
	user "backend_ajax-people"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) createUser(c *gin.Context) {
	var input user.User

	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("Failed to create a user: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, "Failed to create a user")
		return
	}

	id, err := h.services.UserAction.CreateUser(input)
	if err != nil {
		fmt.Printf("Failed data: %s\n", err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllUsers(c *gin.Context) {
	userList, err := h.services.GetAllUsers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, userList)
}

func (h *Handler) getUserById(c *gin.Context) {
	getId, err := getUserId(c)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.services.GetUserById(getId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) updateUser(c *gin.Context) {
	getId, err := getUserId(c)
	var input user.UpdateUserInput

	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("Failed to update a user: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, "Failed to update a user")
		return
	}

	err = h.services.UpdateUser(getId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Ok")
}

func (h *Handler) deleteUser(c *gin.Context) {
	userId, err := getUserId(c)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.DeleteUser(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Ok")
}

type Message struct {
	Code string `json:"code" binding:"required"`
}

func (h *Handler) checkActivationUser(c *gin.Context) {
	//userId, _ := getUserId(c)
	//
	//var code Message
	//
	//if err := c.BindJSON(&code); err != nil {
	//	c.JSON(http.StatusBadRequest, "Wrong content")
	//	return
	//}
	//
	//verified, err := h.services.CheckCodeActivation(userId, code.Code)
	//if err != nil {
	//	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//if !verified {
	//	c.JSON(http.StatusBadRequest, "Code is not correct")
	//	return
	//}

	c.JSON(http.StatusOK, "Successfully")
}

func (h *Handler) selectUsers(c *gin.Context) {
	var input user.User

	idUser, _, _, _ := getJWT(h, c)

	var inputSelect user.UpdateUserInput

	items, _ := strconv.Atoi(c.Query("items"))
	page, _ := strconv.Atoi(c.Query("page"))
	input.Id, _ = strconv.Atoi(c.Query("id"))
	input.Age, _ = strconv.Atoi(c.Query("age"))
	input.FirstName = c.Query("firstName")
	input.LastName = c.Query("lastName")
	input.StatusUserId, _ = strconv.Atoi(c.Query("statusUserId"))
	input.StudyProgramId, _ = strconv.Atoi(c.Query("studyProgramId"))
	input.EducationLevelId, _ = strconv.Atoi(c.Query("educationLevelId"))
	input.SchoolId, _ = strconv.Atoi(c.Query("schoolId"))
	input.AdmissionYear, _ = strconv.Atoi(c.Query("admissionYear"))
	input.GraduationYear, _ = strconv.Atoi(c.Query("graduationYear"))
	interests := strings.Split(c.Query("interests"), `,`)
	for _, s := range interests {
		num, err := strconv.Atoi(s)
		if err == nil {
			inputSelect.IdsInterests = append(inputSelect.IdsInterests, num)
		}
	}

	inputSelect.Id = &input.Id
	inputSelect.SchoolId = &input.SchoolId
	inputSelect.Age = &input.Age
	if input.LastName != "" {
		inputSelect.LastName = &input.FirstName
	}
	if input.FirstName != "" {
		inputSelect.FirstName = &input.FirstName
	}
	inputSelect.StudyProgramId = &input.StudyProgramId
	inputSelect.StatusUser = &input.StatusUserId
	inputSelect.EducationLevel = &input.EducationLevelId
	inputSelect.AdmissionYear = &input.AdmissionYear
	inputSelect.GraduationYear = &input.GraduationYear

	userList, err := h.services.SelectedDataUser(inputSelect, idUser, page, items)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, userList)
}

func (h *Handler) coincidenceSend(c *gin.Context) {
	idSender, _, _, err := getJWT(h, c)
	var input Message

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "Failed to get Mail")
		return
	}

	idCoincidence, err := h.services.RequestСorrespondence(idSender, input.Code)
	if idCoincidence == -1 {
		c.JSON(http.StatusBadRequest, "Request exists")
		return
	}
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": idCoincidence,
	})
}

func (h *Handler) coincidenceAccept(c *gin.Context) {
	var reqId int

	reqId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.services.AcceptMessageRequest(reqId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "Accept")
}

func (h *Handler) ChangeUserOnAdmin(c *gin.Context) {
	var id int

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.services.ChangeUserOnAdmin(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, "Ok")
}
