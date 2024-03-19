>### Supported commands
>Creating a migration file
>```bash
>migrator create CreateOrdersTable
>```
>```bash
>migrator create --name="createOrdersTable"
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
>```bash
> migrator redo --all
>```
>```bash
> migrator redo --step="1"
>```
>> The **"step"** parameter sets the version of the step to repeat the migration from the end, for example, the current version of the last migration is 10. In step 2, all migrations with versions 10 and 9 will be rolled back and applied.
> 
>Rolling back the last migration
>```bash
> migrator down
>```
>```bash
> migrator down --all
>```
>```bash
> migrator down --step="1"
>```
>> The step option sets the version of the step to roll back from the end, for example, the current version of the last migration is 10. In step 2, it will roll back all migrations with versions 10 and 9.

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