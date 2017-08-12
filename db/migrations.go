// Copyright (c) 2017 Kuguar <licenses@kuguar.io> Author: Adrian P.K. <apk@kuguar.io>
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package db

import (
	"fmt"
	"log"

	"github.com/DavidHuie/gomigrate"
	_ "github.com/lib/pq" // Import pq without side effects.
)

// Migrate - Run pending database migrations.
func Migrate(migrationsDir string) error {
	return getMigrator(migrationsDir).Migrate()
}

// Rollback - Rollback last database migration.
func Rollback(migrationsDir string) error {
	return getMigrator(migrationsDir).Rollback()
}

// RollbackN - Rollback last n database migrations.
func RollbackN(numMigrations int, migrationsDir string) error {
	return getMigrator(migrationsDir).RollbackN(numMigrations)
}

func getMigrator(migrationsDir string) *gomigrate.Migrator {
	db, err := GetDb()
	if err != nil {
		log.Fatal(err)
		panic(fmt.Sprintf(err.Error()))
	}
	migrator, _ := gomigrate.NewMigrator(db, gomigrate.Postgres{}, migrationsDir)
	return migrator
}
