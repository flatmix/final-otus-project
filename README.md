>### Supported commands
>Creating a migration file
>```bash
>migrator create CreateOrders
>```
>Applying all migrations
>
>```bash
> migrator up
>```
>Repeat the last migration (rollback + up)
>```bash
> migrator redo
>```
>Rolling back the last migration
>```bash
> migrator down
>```

>### Set db config examples
>```bash
>migrator up --user="user" --pass="password" --db="database" --host="localhost" --port=5432 --sslMode=true
>```
>```bash
>migrator up --dbDSN="postgres://user:password@localhost:5432/database?sslmode=disable"
>```
>```bash
>export MIGRATOR_DB_DSN=postgres://postgres:password@localhost:5432/postgres?sslmode=disable
>migrator up
>```