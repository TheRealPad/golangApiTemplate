package core

import (
	"fmt"
	"github.com/gorilla/mux"
	"httpServer/src/controller"
	database2 "httpServer/src/database"
	"httpServer/src/initialisation"
	"httpServer/src/middlewares"
	"httpServer/src/models"
	"net/http"
	"strconv"
)

type ApiInterface interface {
	Listen()
	Initialisation()
}

type Api struct {
	Json initialisation.JsonHandler
}

func (a Api) Listen() {
	var configuration models.Configuration
	var dataModel []initialisation.DataModel
	var db database2.DatabaseInterface
	if !a.Initialisation(&configuration, &dataModel) {
		return
	}
	db = &database2.MongoDB{Name: configuration.Db.Name, Url: configuration.Db.Url}
	displayDataTypes(&dataModel)
	r := mux.NewRouter()
	middlewares.GlobalMiddleware(r)
	controller.InitControllers(r, &configuration, &dataModel, db)
	fmt.Println("Server", configuration.Name, "starts listening on port:", configuration.Port)
	http.ListenAndServe(":"+strconv.Itoa(configuration.Port), r)
}

func (a Api) Initialisation(configuration *models.Configuration, dataModel *[]initialisation.DataModel) bool {
	if !a.Json.ReadFile(configuration) {
		return false
	}
	for _, model := range configuration.Models {
		*dataModel = append(*dataModel, initialisation.DataModel{Name: model.Name, Fields: make(initialisation.Field)})
		dataModelPtr := &(*dataModel)[len(*dataModel)-1]
		dataModelPtr.Fields[initialisation.Uuid] = &initialisation.DynamicType{}
		dataModelPtr.Fields[initialisation.Uuid].SetData("", initialisation.Uuid)
		for _, e := range model.Fields {
			dataModelPtr.Fields[e.Name] = &initialisation.DynamicType{}
			dataModelPtr.Fields[e.Name].SetData("", initialisation.Datatype(e.Type))
		}
	}
	displayConfiguration(configuration)
	return true
}

func displayConfiguration(configuration *models.Configuration) {
	fmt.Println("CONFIGURATION:\n")
	fmt.Println("port:", configuration.Port)
	fmt.Println("name:", configuration.Name)
	fmt.Println("Database:")
	fmt.Println("\turl:", configuration.Db.Url)
	fmt.Println("\tname:", configuration.Db.Name)
	fmt.Println("data models:")
	fmt.Println("total:", len(configuration.Models))
	for _, model := range configuration.Models {
		fmt.Println("\tname:", model.Name)
		fmt.Print("\tfields:", len(model.Fields), " ")
		for _, e := range model.Fields {
			fmt.Print(e.Name + " - " + e.Type + " ")
		}
		fmt.Println()
		fmt.Println("\tcreate:", model.Create)
		fmt.Println("\tread one:", model.ReadOne)
		fmt.Println("\tread many:", model.ReadMany)
		fmt.Println("\tupdate:", model.Update)
		fmt.Println("\tdelete:", model.Delete)
		fmt.Println()
	}
}

func displayDataTypes(dataModel *[]initialisation.DataModel) {
	fmt.Println("DATA TYPES:\n")
	for _, elem := range *dataModel {
		fmt.Println(elem.Name)
		for k, f := range elem.Fields {
			fmt.Println("\t", k, ":", f.GetDataType())
		}
	}
}

type ApiService struct {
	Api *Api
}

func (s *ApiService) Listen() {
	s.Api.Listen()
}
