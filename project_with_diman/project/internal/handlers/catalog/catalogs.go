package catalog

import (
	"Project/internal/manager"
	"Project/internal/models"
	"Project/internal/repository/user/handler"
	"Project/internal/token/jwt"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Middleware(project *manager.Project) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			value, err := jwt.DecodeJwt(c.Request().Context(), token)
			if err != nil {
				project.Logger.Println("Error decoding token ", err)
				return c.JSON(401, "You're not authorized")
			}

			val, err := project.CacheUser.GetData(c.Request().Context(), strconv.Itoa(value.UserId))
			if err != nil {
				project.Logger.Println("Cache error:", err)
				return c.JSON(401, "Not bv")
			}

			if val != strconv.Itoa(value.Version) {
				project.Logger.Println("Version userId not looks like ", value.UserId)
				return c.JSON(401, "")
			}
			c.Set("userId", value.UserId)
			return next(c)
		}
	}
}

func HandlerCatalogs(e *echo.Echo, project *manager.Project) {
	catalogGroup := e.Group("/catalogs", Middleware(project))

	catalogGroup.GET("", func(c echo.Context) error {
		return GetCatalogs(c, project)
	})
	catalogGroup.POST("", func(c echo.Context) error {
		return CreateCatalogs(c, project)
	})
	catalogGroup.PATCH("", func(c echo.Context) error {
		return UpdateCatalogs(c, project)
	})
	catalogGroup.GET("/bucket", func(c echo.Context) error {
		return GetBucket(c, project)
	})

	catalogGroup.POST("/bucket", func(c echo.Context) error {
		return AddUpdateBucket(c, project)
	})
	//catalogGroup.POST("", GetCatalogsHandler(project))
	//catalogGroup.PATCH("/:id", DeleteCatalogHandler(project))
	//catalogGroup.DELETE("",Delete)
}

func GetBucket(c echo.Context, project *manager.Project) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	userId := c.Get("userId").(int)
	value, err := handler.GetBucketHandler(ctx, project, userId)
	if err != nil {
		return c.JSON(500, "Please try again later.")
	}

	return c.JSON(200, value)
}

func AddUpdateBucket(c echo.Context, project *manager.Project) error {
	var user models.AddBucketRequest
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	userId := c.Get("userId").(int)
	err := handler.CreateBucketHandler(ctx, project, user.Id, userId)
	if err != nil {
		if errors.Is(err, handler.ErrAlreadyExists) {
			return c.JSON(400, "Вы уже добавили такой товар")
		}
		return c.JSON(500, "Please try again later.")
	}

	return c.JSON(200, "Success")

}

func CreateCatalogs(c echo.Context, project *manager.Project) error {
	var catalog models.CreateCatalogRequest
	if err := c.Bind(&catalog); err != nil {
		project.Logger.Println("Error parsing catalog ", err)
		return c.JSON(500, "Please try later")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	if catalog.Filter == "" {
		return c.JSON(400, "Нужно указать фильтр")
	}
	value, ok := handler.Validate(ctx, catalog.Filter)
	if !ok {
		return c.JSON(400, "Не валидный фильтр")
	}
	catalog.Filter = value

	userId := c.Get("userId").(int)
	ok, err := handler.CreateCatalogHandler(ctx, catalog, userId, project)
	if err != nil {
		return c.JSON(500, "Oops something went wrong")
	}
	if !ok {
		return c.JSON(500, "Oops something went wrong")
	}
	return c.JSON(204, nil)
}

func UpdateCatalogs(c echo.Context, project *manager.Project) error {
	var catalog models.UpdateCatalogRequest
	if err := c.Bind(&catalog); err != nil {
		project.Logger.Println("Error parsing catalog ", err)
		return c.JSON(500, "Please try later")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	valid, ok := handler.Validate(ctx, catalog.Filter)
	if !ok {
		return c.JSON(400, "Ошибка параметра фильтр")
	}
	catalog.Filter = valid
	userId := c.Get("userId").(int)

	err := handler.UpdateCatalogHandler(ctx, catalog, userId, project)
	if err != nil {
		switch err {
		case handler.ErrClient:
			return c.JSON(400, "Wrong data")
		default:
			return c.JSON(500, "Oops something went wrong")
		}
	}
	return c.JSON(204, nil)
}

func GetCatalogs(c echo.Context, project *manager.Project) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	filter := c.QueryParam("filter")
	if filter == "" {
		log.Println("Filter query:", filter)
		return c.JSON(500, "Pleas writy query param")
	}

	value, ok := handler.Validate(ctx, filter)
	if !ok {
		project.Logger.Println("Error validating request")
		return c.JSON(500, "Фронт ошибка") //Просто так
	}

	data, err := handler.GetCatalogHandler(ctx, value, project)
	if err != nil {
		project.Logger.Println("Error getting catalog", err)
		return c.JSON(500, "Error serve")
	}
	return c.JSON(http.StatusOK, data)
}
