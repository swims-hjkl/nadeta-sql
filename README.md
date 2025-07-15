# nadeta-sql


## Building using Make

Available targets:<br />
  build   - Build the app<br />
  test    - Run unit tests<br />
  clean   - Remove build artifacts<br />
  fmt     - Format code<br />
  help    - Show this help<br />
<br />
<br />
## Usage:

init: used to initialize sqlmigrate (mandatory before using the tool)<br />
  -connection-string string<br />
    	connection string to the database<br />
  -directory-name string<br />
    	directory name to use for storing migrations (default "sqlmigratior_migration_file")<br />
<br />
<br />
create: used to create migration<br />
  -name string<br />
    	name of the migration<br />
<br />
<br />
delete: used to delete migration<br />
  -name string
    	name of the migration
<br />
<br />
up: used to apply migrations<br />
  -dryrun<br />
    	peek into migrations which will run without actually running them<br />
  -name string<br />
    	name of a specific migration<br />
  -steps int<br />
    	determines how many migrations to run, 0 by default runs all<br />
<br />
<br />
down: used to rollback migrations<br />
  -dryrun<br />
    	peek into migrations which will run without actually running them<br />
  -name string<br />
    	name of a specific migration<br />
  -steps int<br />
    	determines how many migrations to run, 0 by default runs all<br />
<br />
<br />
status: displays what migrations have been applied and what are pending<br />
  -applied<br />
    	show only applied migrations<br />
  -pending<br />
    	show only pending migrations<br />
<br />
<br />
list: lists all the migrations known to the system<br />
