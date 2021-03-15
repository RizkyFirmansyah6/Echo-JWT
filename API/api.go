package API

import (
	"EchoAPI/Helper"
	"EchoAPI/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"time"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type User struct {
	ID       string `json:"id" form:"id"`
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Nama     string `json:"nama" form:"nama"`
	Foto     string `json:"foto" form:"foto"`
}

func StartServer(address string) {
	r := echo.New()

	// Middleware
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())

	r.Validator = &CustomValidator{validator: validator.New()}
	r.HTTPErrorHandler = func(err error, c echo.Context) {
		report, ok := err.(*echo.HTTPError)
		if !ok {
			report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if castedObject, ok := err.(validator.ValidationErrors); ok {
			for _, err := range castedObject {
				switch err.Tag() {
				case "required":
					report.Message = fmt.Sprintf("%s is required",
						err.Field())
				case "email":
					report.Message = fmt.Sprintf("%s is not valid email",
						err.Field())
				case "gte":
					report.Message = fmt.Sprintf("%s value must be greater than %s",
						err.Field(), err.Param())
				case "lte":
					report.Message = fmt.Sprintf("%s value must be lower than %s",
						err.Field(), err.Param())
				}

				break
			}
		}

		c.Logger().Error(report)
		c.JSON(report.Code, report)
	}

	// Login route
	r.POST("/login", login)

	r.GET("/", func(ctx echo.Context) error {
		data := "Hello from  /index"
		return ctx.String(http.StatusOK, data)
	})

	// User route
	e := r.Group("/users")
	e.Use(middleware.JWT([]byte("secret")))
	e.GET("", func(ctx echo.Context) error {
		data := model.GetUsers()
		return ctx.JSON(http.StatusOK, data)
	})

	e.POST("", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		if err := c.Validate(u); err != nil {
			return err
		}
		imageName, err := Helper.FileUpload(c)
		if err != nil {
			return c.String(http.StatusBadRequest, "Error upload image")
			//checking whether any error occurred retrieving image
		}
		u.Foto = imageName
		u.Password = Helper.EncryptPassword(u.Password)
		model.PostUsers(u.Username, u.Password, u.Nama, u.Foto)

		return c.JSON(http.StatusOK, "success")
	})

	e.PUT("", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		if err := c.Validate(u); err != nil {
			return err
		}
		imageName, err := Helper.FileUpload(c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, "Error upload image")
			//checking whether any error occurred retrieving image
		}
		u.Foto = imageName
		u.Password = Helper.EncryptPassword(u.Password)
		model.PutUser(u.ID, u.Username, u.Password, u.Nama, u.Foto)

		return c.JSON(http.StatusOK, "success")
	})

	e.DELETE("", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		model.DelUser(u.ID)

		return c.JSON(http.StatusOK, "success")
	})

	r.Start(":8080")
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	password = Helper.EncryptPassword(password)

	// Throws unauthorized error
	data := model.Login(username, password)
	if data.ID == "" {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Jon Snow"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
