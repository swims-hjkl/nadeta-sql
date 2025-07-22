package main

import (
	"flag"

	"github.com/swims/nadeta-sql/types"
)

const FOLDER_NAME string = "nadeta-migration-files"

func getFlagSetInit() *types.InitFlagData {
	flagSet := flag.NewFlagSet("init", flag.ExitOnError)
	db := flagSet.String("connection-string", "", "connection string to the database")
	directoryName := flagSet.String("directory-name", FOLDER_NAME, "directory name to use for storing migrations")
	return &types.InitFlagData{
		FlagSet:          *flagSet,
		ConnectionString: db,
		DirectoryName:    directoryName,
	}
}

func getFlagSetCreate() *types.CreateFlagData {
	flagSet := flag.NewFlagSet("create", flag.ExitOnError)
	name := flagSet.String("name", "", "name of the migration")
	return &types.CreateFlagData{
		FlagSet: *flagSet,
		Name:    name,
	}
}

func getFlagSetDelete() *types.DeleteFlagData {
	flagSet := flag.NewFlagSet("delete", flag.ExitOnError)
	name := flagSet.String("name", "", "name of the migration")
	return &types.DeleteFlagData{
		FlagSet: *flagSet,
		Name:    name,
	}
}

func getFlagSetUp() *types.UpFlagData {
	flagSet := flag.NewFlagSet("up", flag.ExitOnError)
	steps := flagSet.Int("steps", 0, "determines how many migrations to run, 0 by default runs all")
	dryrun := flagSet.Bool("dryrun", false, "peek into migrations which will run without actually running them")
	name := flagSet.String("name", "", "name of a specific migration")
	return &types.UpFlagData{
		FlagSet: *flagSet,
		Steps:   steps,
		Dryrun:  dryrun,
		Name:    name,
	}
}

func getFlagSetDown() *types.DownFlagData {
	flagSet := flag.NewFlagSet("down", flag.ExitOnError)
	steps := flagSet.Int("steps", 0, "determines how many migrations to run, 0 by default runs all")
	dryrun := flagSet.Bool("dryrun", false, "peek into migrations which will run without actually running them")
	name := flagSet.String("name", "", "name of a specific migration")
	return &types.DownFlagData{
		FlagSet: *flagSet,
		Steps:   steps,
		Dryrun:  dryrun,
		Name:    name,
	}
}

func getFlagSetStatus() *types.StatusFlagData {
	flagSet := flag.NewFlagSet("status", flag.ExitOnError)
	pending := flagSet.Bool("pending", false, "show only pending migrations")
	applied := flagSet.Bool("applied", false, "show only applied migrations")
	return &types.StatusFlagData{
		FlagSet: *flagSet,
		Pending: pending,
		Applied: applied,
	}
}
