# nadeta-sql


## Building using Make

Available targets:
  build   - Build the app
  test    - Run unit tests
  clean   - Remove build artifacts
  fmt     - Format code
  help    - Show this help


## Usage:

init: used to initialize sqlmigrate (mandatory before using the tool)
  -connection-string string
    	connection string to the database
  -directory-name string
    	directory name to use for storing migrations (default "sqlmigratior_migration_file")


create: used to create migration
  -name string
    	name of the migration


delete: used to delete migration
  -name string
    	name of the migration


up: used to apply migrations
  -dryrun
    	peek into migrations which will run without actually running them
  -name string
    	name of a specific migration
  -steps int
    	determines how many migrations to run, 0 by default runs all


down: used to rollback migrations
  -dryrun
    	peek into migrations which will run without actually running them
  -name string
    	name of a specific migration
  -steps int
    	determines how many migrations to run, 0 by default runs all


status: displays what migrations have been applied and what are pending
  -applied
    	show only applied migrations
  -pending
    	show only pending migrations


list: lists all the migrations known to the system
