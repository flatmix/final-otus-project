package storage

type MigrationDBStruct struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Hash      string `db:"hash"`
	Version   int    `db:"version"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (migrate MigrationDBStruct) IsZero() bool {
	return migrate.ID == 0
}

type MigrationsDBStruct []MigrationDBStruct

func (migrate MigrationsDBStruct) Len() int {
	return len(migrate)
}

func (migrate MigrationsDBStruct) Less(i, j int) bool {
	return migrate[i].ID < migrate[j].ID
}

func (migrate MigrationsDBStruct) Swap(i, j int) {
	migrate[i], migrate[j] = migrate[j], migrate[i]
}
