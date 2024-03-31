package usecase

const createMigrationsTable = `CREATE TABLE IF NOT EXISTS public.%s
	(
		id   serial PRIMARY KEY   NOT NULL ,
		name  varchar             NOT NULL,
		hash varchar              NOT NULL,
		version integer default 0,
		created_at timestamp default now(),
		updated_at timestamp default now()
	);`

const selectExistsTable = `SELECT EXISTS (
   SELECT FROM information_schema.tables 
   WHERE  table_schema = '%s'
   AND    table_name   = '%s'
   );`

const insertMigrationsTable = `INSERT INTO public.%s (name, hash, version, updated_at) VALUES ($1, $2, $3, $4)`

const deleteRowFromMigrationsTable = `DELETE FROM public.%s where name=$1`

const getMaxVersionFromMigrationsTable = `SELECT max(version) FROM public.%s`

const migrationFileTemplate = `--migrate:UP
--write your sql for migration...


--migrate:DOWN
--write your sql for rollback migration...

`

const regexpUpTemplate = `--migrate:UP([\s\S]+?)--migrate:DOWN`

const regexpDownTemplate = `--migrate:DOWN([\s\S]+?)$`

const selectMigrationTable = `SELECT id, name, hash, version, created_at, updated_at FROM public.%s
   WHERE  name = $1 ORDER BY id DESC LIMIT 1`

const selectAllMigrationsTable = `SELECT id, name, hash, version, created_at, updated_at
FROM public.%s WHERE %s.version >= $1 ORDER BY version DESC, id DESC`

type FilesMap map[string]FileStruct

type FileInfo interface {
	Name() string // base name of the file
	Size() int64  // length in bytes for regular files; system-dependent for others
	//	Mode() fs.FileMode  // file mode bits
	//	ModTime() time.Time // modification time
	//	IsDir() bool        // abbreviation for Mode().IsDir()
	//	Sys() any           // underlying data source (can return nil)
}

type FileStruct struct {
	File FileInfo
	Hash string
}

type Outs []*Out

type Out struct {
	Name        string
	Status      string
	Version     string
	TimeMigrate string
}
