<h1>Database Migration ReadMe</h1>

<h2>Run database migrations on code instead of docker</h2>

This is so helpful because it can simplify the docker file, removing the logic to run the database migrations. Doing this on the code instead.

Add the path to the migration file in to our `app.env` file

```env
    MIGRATION_URL=file://db/migration
```

Add these modules so we can leverage them to make the migrations since Go code.

```go
    "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
```

Create the functions to run the database migration

```go
    func runDBMigration(migrationURL string, dbSource string) {
        migration, err := migrate.New(migrationURL, dbSource)
        if err != nil {
            log.Fatal("cannot create new migrate instance")
        }

        if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
            log.Fatal("iled to run migrate up")
        }
    }
```
