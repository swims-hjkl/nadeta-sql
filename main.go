package main

// What should the usage look like?
// sqlmigrate init --db "connection string"
// sqlmigrate create --name
// sqlmigrate delete --name
// sqlmigrate up
// sqlmigrate up --steps 2
// sqlmigrate up --steps 2 --dryrun
// sqlmigrate up --name migrationname --dryrun
// sqlmigrate down
// sqlmigrate down --steps 2
// sqlmigrate down --steps 2 --dryrun
// sqlmigrate status
// sqlmigrate status --pending
// sqlmigrate status --applied
// sqlmigrate list

import (
	"fmt"
	"log"
	"os"
	_ "time"

	"github.com/swims/nadeta-sql/dbutil"
	"github.com/swims/nadeta-sql/helpers"
	"github.com/swims/nadeta-sql/migrations"
	"github.com/swims/nadeta-sql/processors"
	"github.com/swims/nadeta-sql/types"
)

func main() {

	flagSetInit := getFlagSetInit()
	flagSetCreate := getFlagSetCreate()
	flagSetDelete := getFlagSetDelete()
	flagSetUp := getFlagSetUp()
	flagSetDown := getFlagSetDown()
	flagSetStatus := getFlagSetStatus()

	if len(os.Args) < 2 {
		log.Fatal("at least 1 argument required")
	}

	if (os.Args[1] == "--help") || (os.Args[1] == "-help") {
		os.Args[1] = "help"
	}

	var config *types.Config
	var err error
	var dbUtilObj *dbutil.DBUtil
	var migrationStore *migrations.MigrationStore

	fmt.Printf("\n\n")

	if os.Args[1] != "init" && os.Args[1] != "help" {
		config, err = helpers.GetConfig()
		if err != nil && err.Error() == "Config not present" {
			log.Fatal("did not find configuration file, did you forget to run init?")
		}
		if err != nil {
			log.Fatal(err)
		}
		dbUtilObj, err = dbutil.NewDBUtil(config.ConnectionString)
		if err != nil {
			log.Fatal(err)
		}
		defer dbUtilObj.CloseDatabase()
		migrationStore = migrations.NewMigrationStore(dbUtilObj, config)
		if err != nil {
			log.Fatal(err)
		}
	}

	switch os.Args[1] {
	case "init":
		flagSetInit.FlagSet.Parse(os.Args[2:])
		err := processors.RunInit(flagSetInit)
		if err != nil {
			log.Fatal(err)
		}
	case "create":
		flagSetCreate.FlagSet.Parse(os.Args[2:])
		err := processors.RunCreate(flagSetCreate, migrationStore)
		if err != nil {
			log.Fatal(err)
		}
	case "delete":
		flagSetDelete.FlagSet.Parse(os.Args[2:])
		err := processors.RunDelete(flagSetDelete, migrationStore)
		if err != nil {
			log.Fatal(err)
		}
	case "up":
		flagSetUp.FlagSet.Parse(os.Args[2:])
		err := processors.RunUp(flagSetUp, migrationStore)
		if err != nil {
			log.Fatal(err)
		}
	case "down":
		flagSetDown.FlagSet.Parse(os.Args[2:])
		err := processors.RunDown(flagSetDown, migrationStore)
		if err != nil {
			log.Fatal(err)
		}
	case "status":
		flagSetStatus.FlagSet.Parse(os.Args[2:])
		err := processors.RunStatus(flagSetStatus, migrationStore)
		if err != nil {
			log.Fatal(err)
		}
	case "list":
		err := processors.RunList(migrationStore)
		if err != nil {
			log.Fatal(err)
		}
	case "purge":
		err := processors.RunPurge(migrationStore)
		if err != nil {
			log.Fatal(err)
		}
	case "help":
		fmt.Print("\n\ninit: used to initialize sqlmigrate (mandatory before using the tool)\n")
		flagSetInit.FlagSet.PrintDefaults()
		fmt.Print("\n\ncreate: used to create migration\n")
		flagSetCreate.FlagSet.PrintDefaults()
		fmt.Print("\n\ndelete: used to delete migration\n")
		flagSetDelete.FlagSet.PrintDefaults()
		fmt.Print("\n\nup: used to apply migrations\n")
		flagSetUp.FlagSet.PrintDefaults()
		fmt.Print("\n\ndown: used to rollback migrations\n")
		flagSetDown.FlagSet.PrintDefaults()
		fmt.Print("\n\nstatus: displays what migrations have been applied and what are pending\n")
		flagSetStatus.FlagSet.PrintDefaults()
		fmt.Print("\n\nlist: lists all the migrations known to the system\n")
		fmt.Print("\n\npurge: Caution! Purges all data related to migrations\n")
		fmt.Print("\n\n")
	default:
		fmt.Println("Please provide valid arguments\n")
		os.Exit(0)
	}
	fmt.Printf("\n\n")
}
