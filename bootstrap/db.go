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
//

package bootstrap

import (
	"github.com/adrianpk/fundacja/db"
	"github.com/adrianpk/fundacja/logger"
	"path"
)

const (
	migrate  = true
	rollback = true
)

func initMigrationOrRollback() {
	dbConfig()
	migrationSwitches()
	dir := migrationsDir()
	if doMigration {
		logger.Debug("Migrating...")
		initMigration(dir)
	} else if doRollback {
		if rollbackAll {
			logger.Debug("Doing rollback all...")
			initRollbackN(migrationsNum, dir)
		} else {
			logger.Debug("Doing rollback...")
			initRollback(dir)
		}
	}
}

func migrationsDir() string {
	migrationsDir := path.Join(BaseDir, "resources/migrations")
	//logger.Info("Migration's dir is %s", migrationsDir)
	return migrationsDir
}

func initMigration(migrationsDir string) {
	if migrate {
		logger.Debug("Migration init...")
		db.Migrate(migrationsDir)
		logger.Debug("Migration completed.")
	}
}

func initRollback(migrationsDir string) {
	if rollback {
		logger.Debug("Rollback init...")
		db.Rollback(migrationsDir)
		logger.Debug("Rollback completed.")
	}
}

func initRollbackN(numRollback int, migrationsDir string) {
	if rollback {
		logger.Debugf("Rollback %d migrations init...", numRollback)
		db.RollbackN(numRollback, migrationsDir)
		logger.Debug("Rollback completed.")
	}
}

func migrationSwitches() {
	if mgrSw == "m" {
		doMigration = true
		doRollback = false
	} else if mgrSw == "r" {
		doMigration = false
		doRollback = true
	} else {
		doMigration = false
		doRollback = false
	}
}

func dbConfig() {
	db.DBConfig.Host = AppConfig.DBHost
	db.DBConfig.DB = AppConfig.Database
	db.DBConfig.User = AppConfig.DBUser
	db.DBConfig.Pass = AppConfig.DBPass
	db.DBConfig.SSL = AppConfig.DBSSL
}
