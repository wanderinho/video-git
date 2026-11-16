package projects

import (
	"errors"
	"net/http"
	"video-git/internal/auth"

	"github.com/labstack/echo/v4"
)

type ProjectHandler struct {
	service *ProjectService
}

func NewProjectHandler(service *ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (ph *ProjectHandler) CreateProject(c echo.Context) error {
	userID, ok := auth.UserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthorized"})
	}

	var createProjectReq CreateProjectRequest
	if err := c.Bind(&createProjectReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request body"})
	}

	createdProject, err := ph.service.Create(c.Request().Context(), userID, createProjectReq.Desc, createProjectReq.Name, createProjectReq.BranchName)
	if err != nil {
		switch {
		case errors.Is(err, ErrAlreadyExists):
			return c.JSON(http.StatusConflict, echo.Map{
				"message": err.Error(),
			})

		case errors.Is(err, ErrInvalidBranchName),
			errors.Is(err, ErrInvalidProjectName),
			errors.Is(err, ErrSizeName):
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})

		default:
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "internal server error",
			})
		}
	}

	return c.JSON(http.StatusCreated, CreateProjectResponse{
		ID: createdProject.ID,
		Name: createdProject.Name,
		Desc: createdProject.Desc,
		DefaultBranch: createdProject.DefaultBranch,
	})
}
