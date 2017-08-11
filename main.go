package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	restful "github.com/emicklei/go-restful"
	swagger "github.com/emicklei/go-restful-swagger12"
	_ "github.com/go-sql-driver/mysql"

	"github.com/dohko-io/sysreqs/mysql"
	"github.com/dohko-io/sysreqs/pkgmnt"
)

func newPackageService() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/package").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("A package knowledge base").
		ApiVersion("1.0")

	service.Route(service.GET("/{name}").To(func(request *restful.Request, response *restful.Response) {
		name := request.PathParameter("name")
		repository, err := mysql.NewPackageRepository(os.Getenv("SYSREQS_DATASOURCE"))

		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
		}

		pkg, err := repository.GetWithName(name)

		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
		}

		response.WriteEntity(pkg)
	}).
		Doc("Get a package with a given name").
		Param(service.PathParameter("name", "Package's name to be returned").DataType("string")).
		Do(returns200, returns500))

	service.Route(service.PUT("").To(func(request *restful.Request, response *restful.Response) {
		var pkg pkgmnt.Package
		err := request.ReadEntity(&pkg)

		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
		}

		repository, err := mysql.NewPackageRepository(os.Getenv("SYSREQS_DATASOURCE"))
		repository.Store(pkg)

		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
		}

		response.WriteHeader(http.StatusCreated)
	}).
		Doc("Add a new package").
		Reads(pkgmnt.Package{}))

	return service
}

func returns200(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", pkgmnt.Package{})
}

func returns500(b *restful.RouteBuilder) {
	b.Returns(http.StatusInternalServerError, "Sorry, something went wrong", nil)
}

func main() {
	restful.Add(newPackageService())

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	config := swagger.Config{
		WebServices:    restful.DefaultContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: fmt.Sprintf("http://localhost:%s", port),
		ApiPath:        "/apidocs.json",
		// Optionally, specify where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: os.Getenv("SYSREQS_SWAGGER_UI_FILE_PATH")}

	swagger.RegisterSwaggerService(config, restful.DefaultContainer)

	log.Printf("Starting listening on port: %s\n", port)
	// server := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: restful.DefaultContainer}
	// log.Fatal(server.ListenAndServe())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), restful.DefaultContainer))
}

// func main() {
// 	r := mux.NewRouter()
// 	port := os.Getenv("PORT")
//
// 	if port == "" {
// 		log.Fatal("$PORT must be set")
// 	}
//
// 	r.HandleFunc("/package/{name}", FindPackageByName).
// 		Methods("GET")
//
// 	r.HandleFunc("/package", CreatePackage).
// 		Methods("PUT").
// 		Headers("Content-Type", "application/json")
//
// 	log.Fatal(http.ListenAndServe(":"+port, r))
// }

// FindPackageByName returns the data of a given package
// func FindPackageByName(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	repository, err := mysql.NewPackageRepository(os.Getenv("SYSREQS_DATASOURCE"))
//
// 	//
// 	if err != nil {
// 		panic(err)
// 	}
// 	//
// 	pkg, err := repository.GetWithName(vars["name"])
//
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	if len(pkg.Name) > 0 {
// 		result, _ := json.Marshal(pkg)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(result)
// 	} else {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusNotFound)
// 	}
// }
//
// // CreatePackage handles new packages to be included in the system
// func CreatePackage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Adding a new package")
// 	repository, err := mysql.NewPackageRepository(os.Getenv("SYSREQS_DATASOURCE"))
//
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	defer r.Body.Close()
// 	decoder := json.NewDecoder(r.Body)
// 	var pkg pkgmnt.Package
// 	err = decoder.Decode(&pkg)
//
// 	if err != nil {
// 		w.Write([]byte(err.Error()))
// 	}
//
// 	err = repository.Store(pkg)
//
// 	if err != nil {
// 		w.Write([]byte(err.Error()))
// 	}
//
// 	w.WriteHeader(http.StatusCreated)
// }
