// Package team_action_status_api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package team_action_status_api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// TeamActionStatus defines model for TeamActionStatus.
type TeamActionStatus struct {
	CompletedAt    time.Time `json:"completed_at"`
	Notes          string    `json:"notes"`
	ResolutionLink string    `json:"resolution_link"`
	ResultValue    int       `json:"result_value"`
}

// TeamActionStatusResponse defines model for TeamActionStatusResponse.
type TeamActionStatusResponse struct {
	CompletedAt    time.Time `json:"completed_at"`
	Notes          string    `json:"notes"`
	ResolutionLink string    `json:"resolution_link"`
	ResultValue    int       `json:"result_value"`
	TimelineId     *int      `json:"timeline_id,omitempty"`
	TrackTeamId    *int      `json:"track_team_id,omitempty"`
}

// TeamActionStatusUpdate defines model for TeamActionStatusUpdate.
type TeamActionStatusUpdate struct {
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	ResolutionLink *string    `json:"resolution_link,omitempty"`
	ResultValue    *int       `json:"result_value,omitempty"`
}

// TeamId defines model for TeamId.
type TeamId = int

// TimelineId defines model for TimelineId.
type TimelineId = int

// XTeamId defines model for XTeamId.
type XTeamId = int

// XTimelineId defines model for XTimelineId.
type XTimelineId = int

// GetTeamActionStatusParams defines parameters for GetTeamActionStatus.
type GetTeamActionStatusParams struct {
	// XTeamId Header идентификатор команды
	XTeamId *XTeamId `json:"XTeamId,omitempty"`

	// XTimelineId Header идентификатор таймлайна
	XTimelineId *XTimelineId `json:"XTimelineId,omitempty"`
}

// PostTeamActionStatusJSONRequestBody defines body for PostTeamActionStatus for application/json ContentType.
type PostTeamActionStatusJSONRequestBody = TeamActionStatus

// PutTeamActionStatusTimelineIdTeamIdJSONRequestBody defines body for PutTeamActionStatusTimelineIdTeamId for application/json ContentType.
type PutTeamActionStatusTimelineIdTeamIdJSONRequestBody = TeamActionStatusUpdate

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Получить все результаты команд по ID команды или ID таймлайна
	// (GET /team-action-status)
	GetTeamActionStatus(ctx echo.Context, params GetTeamActionStatusParams) error
	// Создать результат команды
	// (POST /team-action-status)
	PostTeamActionStatus(ctx echo.Context) error
	// Удаление конкретного результата в конкретный таймлайн
	// (DELETE /team-action-status/{TimelineId}/{TeamId})
	DeleteTeamActionStatusTimelineIdTeamId(ctx echo.Context, timelineId TimelineId, teamId TeamId) error
	// Получить результат конкретной команды в конкретный таймлайн
	// (GET /team-action-status/{TimelineId}/{TeamId})
	GetTeamActionStatusTimelineIdTeamId(ctx echo.Context, timelineId TimelineId, teamId TeamId) error
	// Изменить результат конкретной команды в конкретный таймлайн
	// (PUT /team-action-status/{TimelineId}/{TeamId})
	PutTeamActionStatusTimelineIdTeamId(ctx echo.Context, timelineId TimelineId, teamId TeamId) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetTeamActionStatus converts echo context to params.
func (w *ServerInterfaceWrapper) GetTeamActionStatus(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetTeamActionStatusParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "XTeamId" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("XTeamId")]; found {
		var XTeamId XTeamId
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for XTeamId, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "XTeamId", runtime.ParamLocationHeader, valueList[0], &XTeamId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter XTeamId: %s", err))
		}

		params.XTeamId = &XTeamId
	}
	// ------------- Optional header parameter "XTimelineId" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("XTimelineId")]; found {
		var XTimelineId XTimelineId
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for XTimelineId, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "XTimelineId", runtime.ParamLocationHeader, valueList[0], &XTimelineId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter XTimelineId: %s", err))
		}

		params.XTimelineId = &XTimelineId
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetTeamActionStatus(ctx, params)
	return err
}

// PostTeamActionStatus converts echo context to params.
func (w *ServerInterfaceWrapper) PostTeamActionStatus(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostTeamActionStatus(ctx)
	return err
}

// DeleteTeamActionStatusTimelineIdTeamId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteTeamActionStatusTimelineIdTeamId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "TimelineId" -------------
	var timelineId TimelineId

	err = runtime.BindStyledParameterWithLocation("simple", false, "TimelineId", runtime.ParamLocationPath, ctx.Param("TimelineId"), &timelineId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter TimelineId: %s", err))
	}

	// ------------- Path parameter "TeamId" -------------
	var teamId TeamId

	err = runtime.BindStyledParameterWithLocation("simple", false, "TeamId", runtime.ParamLocationPath, ctx.Param("TeamId"), &teamId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter TeamId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteTeamActionStatusTimelineIdTeamId(ctx, timelineId, teamId)
	return err
}

// GetTeamActionStatusTimelineIdTeamId converts echo context to params.
func (w *ServerInterfaceWrapper) GetTeamActionStatusTimelineIdTeamId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "TimelineId" -------------
	var timelineId TimelineId

	err = runtime.BindStyledParameterWithLocation("simple", false, "TimelineId", runtime.ParamLocationPath, ctx.Param("TimelineId"), &timelineId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter TimelineId: %s", err))
	}

	// ------------- Path parameter "TeamId" -------------
	var teamId TeamId

	err = runtime.BindStyledParameterWithLocation("simple", false, "TeamId", runtime.ParamLocationPath, ctx.Param("TeamId"), &teamId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter TeamId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetTeamActionStatusTimelineIdTeamId(ctx, timelineId, teamId)
	return err
}

// PutTeamActionStatusTimelineIdTeamId converts echo context to params.
func (w *ServerInterfaceWrapper) PutTeamActionStatusTimelineIdTeamId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "TimelineId" -------------
	var timelineId TimelineId

	err = runtime.BindStyledParameterWithLocation("simple", false, "TimelineId", runtime.ParamLocationPath, ctx.Param("TimelineId"), &timelineId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter TimelineId: %s", err))
	}

	// ------------- Path parameter "TeamId" -------------
	var teamId TeamId

	err = runtime.BindStyledParameterWithLocation("simple", false, "TeamId", runtime.ParamLocationPath, ctx.Param("TeamId"), &teamId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter TeamId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutTeamActionStatusTimelineIdTeamId(ctx, timelineId, teamId)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/team-action-status", wrapper.GetTeamActionStatus)
	router.POST(baseURL+"/team-action-status", wrapper.PostTeamActionStatus)
	router.DELETE(baseURL+"/team-action-status/:TimelineId/:TeamId", wrapper.DeleteTeamActionStatusTimelineIdTeamId)
	router.GET(baseURL+"/team-action-status/:TimelineId/:TeamId", wrapper.GetTeamActionStatusTimelineIdTeamId)
	router.PUT(baseURL+"/team-action-status/:TimelineId/:TeamId", wrapper.PutTeamActionStatusTimelineIdTeamId)

}
