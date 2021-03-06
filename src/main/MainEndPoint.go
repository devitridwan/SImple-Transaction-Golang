package main

import (
	"Ridwan/test_kanggo/src/properties"
	"Ridwan/test_kanggo/src/service"
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
)

var logger *zap.Logger

func Ekstration(pathFlac string, nameFlac string, target ...interface{}) error {
	 configPath := flag.String(pathFlac, ".", "Configuration YAML path")
	 configName := flag.String(nameFlac, "config-prod", "Configuration Name { config-prod | config-dev } (Required)")
	 flag.Parse()

	 config file path
	viper.AddConfigPath(configPath)
	// config file name
	viper.SetConfigName(configName)

	err := viper.ReadInConfig()

	if err != nil {
		return err
	}

	for _, element := range target {
		err = viper.Unmarshal(&element)
		if err != nil {
			return err
		}
	}

	viper.OnConfigChange(func(in fsnotify.Event) {
		for _, element := range target {
			err = viper.Unmarshal(&element)
		}
	})

	return nil
}



func main() {
	properties := &properties.EndpointProperties{}
	fmt.Println("trace")
	err := Ekstration("configPath", "configName", &properties)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		Start(properties)
	}
}

func Start(properties *properties.EndpointProperties) {
	prop := properties
	dbCOnf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", prop.Database.Username, prop.Database.Password, prop.Database.Host, prop.Database.Port, prop.Database.Name)
	db, err := sql.Open("mysql", dbCOnf)
	if err != nil {
		panic(err.Error())
	}
	service := service.ServiceKanggo{}
	service.Init()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE}}))

	e.POST(prop.Topic.GetChannelRest("registrasi", "post"), func(context echo.Context) (err error) { return service.Registrasi(context, db) })
	e.POST(prop.Topic.GetChannelRest("login", "post"), func(context echo.Context) (err error) { return service.Login(context, db) })
	e.POST(prop.Topic.GetChannelRest("create-product", "post"), func(context echo.Context) (err error) { return service.CreateProduct(context, db) })
	e.POST(prop.Topic.GetChannelRest("edit-product", "post"), func(context echo.Context) (err error) { return service.EditProduct(context, db) })
	e.POST(prop.Topic.GetChannelRest("delete-product", "post"), func(context echo.Context) (err error) { return service.DeleteProduct(context, db) })
	e.GET(prop.Topic.GetChannelRest("list-product", "get"), func(context echo.Context) (err error) { return service.ListProduct(context, db) })
	e.POST(prop.Topic.GetChannelRest("order-transaction", "post"), func(context echo.Context) (err error) { return service.OrderTransaction(context, db) })
	e.POST(prop.Topic.GetChannelRest("payment-transaction", "post"), func(context echo.Context) (err error) { return service.PaymentTransaction(context, db) })
	// e.POST(this.properties.Topic.GetChannelRest("restore", "post"), func(context echo.Context) (err error) {
	// 	form, err := context.MultipartForm()
	// 	go this.processUploadRestoreRequest(form.File["files"], context.FormValue("id_user"))
	// 	return context.JSON(http.StatusOK, `{"status":"success"}`)
	// })
	// this.logger.Error(e.StartTLS(":"+strconv.Itoa(this.properties.Port), this.properties.CertDir+"cert.pem", this.properties.CertDir+"key.pem").Error())
	logger.Error(e.Start(":" + strconv.Itoa(prop.Port)).Error())

}
