package main

import (
	_ "database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"gitlab.com/garrettcoleallen/happy/factors"
	"gitlab.com/garrettcoleallen/happy/helpers"
	"gitlab.com/garrettcoleallen/happy/notifications"
	"gitlab.com/garrettcoleallen/happy/ratings"
	"gitlab.com/garrettcoleallen/happy/reports"
	"gitlab.com/garrettcoleallen/happy/users"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)

	//setup the config reader
	configFilePath := flag.String("config", "config.json", "Path to config file")
	flag.Parse()
	helpers.SetConfigFilePath(*configFilePath)

	//Make sure db is migrated
	// err := tryMigrateToLatestDBVersion()
	// if err != nil {
	// 	log.Panicf("Failed to run migrations!\n%v", err)
	// }

	//support routes
	supportRouter := r.PathPrefix("/support").Subrouter()
	supportRouter.HandleFunc("/sendfeedback", users.AuthorizeRequest(users.SendFeedbackEmail)).Methods("POST")
	supportRouter.HandleFunc("/getuserchangelogseen", users.AuthorizeRequest(users.GetUserChangelogSeenVersion)).Methods("POST")
	supportRouter.HandleFunc("/setuserchangelogseen", users.AuthorizeRequest(users.SetUserChangelogSeenVersion)).Methods("POST")
	supportRouter.HandleFunc("/addinterestedpeopleemail", users.AddInterestedPeopleEmail).Methods("POST")

	//user route
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/authenticate", users.AuthenticateUser).Methods("POST")
	userRouter.HandleFunc("/logout", users.LogoutUser).Methods("POST")
	userRouter.HandleFunc("/newsession", users.CreateNewSessionForUser).Methods("POST")
	userRouter.HandleFunc("/doesuserexist", users.CheckIfUsernameExists).Methods("POST")
	userRouter.HandleFunc("/doesemailexist", users.CheckIfEmailExists).Methods("POST")
	userRouter.HandleFunc("/createuser", users.CreateUser).Methods("POST")
	userRouter.HandleFunc("/forgotpassword", users.UserForgotPassword).Methods("POST")
	userRouter.HandleFunc("/resetpassword", users.UserResetPassword).Methods("POST")
	userRouter.HandleFunc("/verifyguid", users.AuthorizeRequest(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})).Methods("POST")
	userRouter.HandleFunc("/getuser", users.AuthorizeRequest(users.GetUser)).Methods("POST")
	userRouter.HandleFunc("/setphone", users.AuthorizeRequest(users.UserSetPhoneNumber)).Methods("POST")

	//device logging
	userRouter.HandleFunc("/logdevice", users.AuthorizeRequest(users.LogDeviceToUser)).Methods("POST")

	//notifications routes
	notificationsRouter := r.PathPrefix("/notifications").Subrouter()
	notificationsRouter.HandleFunc("/getall", users.AuthorizeRequest(notifications.GetAllUserReminderSettings)).Methods("POST")
	notificationsRouter.HandleFunc("/updateall", users.AuthorizeRequest(notifications.UpdateUserReminderSettings)).Methods("POST")
	notificationsRouter.HandleFunc("/deleteall", users.AuthorizeRequest(notifications.DeleteUserReminderSettings)).Methods("POST")

	//Ratings route
	ratingsRouter := r.PathPrefix("/ratings").Subrouter()
	ratingsRouter.HandleFunc("/getall", users.AuthorizeRequest(ratings.GetAllRatingsForUser))
	ratingsRouter.HandleFunc("/get", users.AuthorizeRequest(ratings.GetRatingByID))

	ratingsRouter.HandleFunc("/newrating", users.AuthorizeRequest(ratings.NewRating)).Methods("POST")
	ratingsRouter.HandleFunc("/updaterating", users.AuthorizeRequest(ratings.UpdateRating)).Methods("POST")
	ratingsRouter.HandleFunc("/deleterating", users.AuthorizeRequest(ratings.DeleteRating)).Methods("POST")

	ratingsRouter.HandleFunc("/newratingfactor", users.AuthorizeRequest(ratings.NewRatingFactor)).Methods("POST")
	ratingsRouter.HandleFunc("/newratingfactors", users.AuthorizeRequest(ratings.NewRatingFactors)).Methods("POST")
	ratingsRouter.HandleFunc("/deleteratingfactors", users.AuthorizeRequest(ratings.DeleteRatingFactors)).Methods("POST")
	ratingsRouter.HandleFunc("/updateratingfactors", users.AuthorizeRequest(ratings.UpdateRatingFactors)).Methods("POST")

	ratingsRouter.HandleFunc("/getratingswithfactors", users.AuthorizeRequest(ratings.GetRatingsWithFactorsForUserByDateRange)).Methods("POST")
	ratingsRouter.HandleFunc("/getallratingswithfactors", users.AuthorizeRequest(ratings.GetAllRatingsWithFactorsForUser)).Methods("POST")

	//Factors route
	factorsRouter := r.PathPrefix("/factors").Subrouter()
	factorsRouter.HandleFunc("/getall", users.AuthorizeRequest(factors.GetAllFactors))
	factorsRouter.HandleFunc("/getallfactortypes", users.AuthorizeRequest(factors.GetAllFactorTypes))
	factorsRouter.HandleFunc("/newfactor", users.AuthorizeRequest(factors.NewFactor)).Methods("POST")
	factorsRouter.HandleFunc("/setfactorarchive", users.AuthorizeRequest(factors.SetFactorArchive)).Methods("POST")
	factorsRouter.HandleFunc("/renamefactor", users.AuthorizeRequest(factors.RenameFactor)).Methods("POST")

	//factor aspects route
	factorAspectsRouter := r.PathPrefix("/factoraspects").Subrouter()
	factorAspectsRouter.HandleFunc("/getall", users.AuthorizeRequest(factors.GetAllFactorAspects))
	factorAspectsRouter.HandleFunc("/newfactoraspect", users.AuthorizeRequest(factors.NewFactorAspect)).Methods("POST")
	factorAspectsRouter.HandleFunc("/setfactoraspectarchive", users.AuthorizeRequest(factors.SetFactorAspectArchive)).Methods("POST")
	factorAspectsRouter.HandleFunc("/renamefactoraspect", users.AuthorizeRequest(factors.RenameFactorAspect)).Methods("POST")

	//reports route
	reportsRouter := r.PathPrefix("/reports").Subrouter()
	reportsRouter.HandleFunc("/factortoratingoccurrences", users.AuthorizeRequest(reports.GetReportFactorToRatingOccurrences))
	reportsRouter.HandleFunc("/request", users.AuthorizeRequest(reports.RequestReport))
	reportsRouter.HandleFunc("/results", users.AuthorizeRequest(reports.GetReportResultsForReportType))
	reportsRouter.HandleFunc("/getreportpending", users.AuthorizeRequest(reports.GetPendingRequestsForReport))
	reportsRouter.HandleFunc("/getfulldataexport", users.AuthorizeRequest(reports.GetFullDataExport))

	http.Handle("/", cors.Default().Handler(r))

	port, err := helpers.GetConfigValue("web.port")
	if err != nil {
		log.Fatal(err)
	}

	certPath, err := helpers.GetConfigValue("web.ssl.cert")
	if err != nil {
		log.Fatal(err)
	}

	keyPath, err := helpers.GetConfigValue("web.ssl.key")
	if err != nil {
		log.Fatal(err)
	}

	// err = tryMigrateToLatestDBVersion()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Printf("Listening on %v...\n", port)
	err = http.ListenAndServeTLS(fmt.Sprintf(":%v", port), certPath, keyPath, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func tryMigrateToLatestDBVersion() error {
	dbMigPath, err := helpers.GetConfigValue("db.migrations_folder")
	if err != nil {
		return err
	}

	dbMigPath = fmt.Sprintf("file://%s", dbMigPath)

	dbMigVersionString, err := helpers.GetConfigValue("db.version")
	if err != nil {
		return err
	}

	dbMigVersion, err := strconv.ParseUint(dbMigVersionString, 10, 64)
	if err != nil {
		return err
	}

	db := helpers.NewDatabaseConnection()

	driver, err := postgres.WithInstance(db, &postgres.Config{DatabaseName: "happy", MigrationsTable: "database_migrations"})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		dbMigPath,
		"postgres",
		driver,
	)

	defer m.Close()

	if err != nil {
		return err
	}

	dbVersion, dirty, _ := m.Version()
	//this throws an error if no migrations have been run so we ignore it
	// if err != nil {
	// 	return err
	// }

	//if we are in a dirty state then error out and tell the fam
	if dirty {
		return fmt.Errorf("Our database is dirty on version %v! Config specifics version %v but it will not be migrated until the dirty status is resolved manually!", dbVersion, dbMigVersion)
	}

	if dbVersion != uint(dbMigVersion) {
		fmt.Printf("\nMigrating from %v to %v\n", dbVersion, dbMigVersion)
		return m.Migrate(uint(dbMigVersion))
	}

	return nil
}
