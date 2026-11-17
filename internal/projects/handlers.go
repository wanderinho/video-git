package projects

import (
	"errors"
	"net/http"
	"strconv"
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
			return c.JSON(http.StatusConflict, echo.Map{"message": err.Error()})

		case errors.Is(err, ErrInvalidBranchName),
			errors.Is(err, ErrInvalidProjectName),
			errors.Is(err, ErrSizeName):
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})

		default:
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
		}
	}

	return c.JSON(http.StatusCreated, ProjectResponse{
		ID:            createdProject.id,
		Name:          createdProject.name,
		Desc:          createdProject.desc,
		DefaultBranch: createdProject.defaultBranch,
	})
}

func (ph *ProjectHandler) GetProjects(c echo.Context) error {
	userID, ok := auth.UserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthorized"})
	}

	projects, err := ph.service.ListByUser(c.Request().Context(), userID)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
	}

	projectsResponse := ProjectsList{
		UserID:   userID,
		Projects: make([]ProjectResponse, 0, len(projects.projectsList)),
	}
	for _, v := range projects.projectsList {
		projectResponse := ProjectResponse{
			ID:            v.id,
			Name:          v.name,
			Desc:          v.desc,
			DefaultBranch: v.defaultBranch,
		}

		projectsResponse.Projects = append(projectsResponse.Projects, projectResponse)
	}

	return c.JSON(http.StatusOK, projectsResponse)
}

func (ph *ProjectHandler) GetProjectByID(c echo.Context) error {
	userID, ok := auth.UserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthorized"})
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid project id"})
	}

	project, err := ph.service.GetProjectByID(c.Request().Context(), projectID, userID)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "project not found"})
		}
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
	}

	return c.JSON(http.StatusOK, FullProjectResponse{
		ID: project.id,
		OwnerID: project.ownerID,
		Name: project.name,
		Desc: project.desc,
		CreatedAt: project.createdAt,
		DefaultBranchID: project.defaultBranchID,
		DefaultBranchName: project.defaultBranchName,
	})
}
