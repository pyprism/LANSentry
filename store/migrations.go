package store

// migrate runs database migrations.
func (s *Store) migrate() error {
	// Create migrations table if not exists
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Get current version
	var currentVersion int
	row := s.db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_migrations`)
	if err := row.Scan(&currentVersion); err != nil {
		currentVersion = 0
	}

	migrations := []struct {
		version int
		sql     string
	}{
		{
			version: 1,
			sql: `
				CREATE TABLE IF NOT EXISTS devices (
					mac TEXT PRIMARY KEY NOT NULL,
					ip TEXT NOT NULL,
					hostname TEXT DEFAULT '',
					manufacturer TEXT DEFAULT '',
					os_hint TEXT DEFAULT '',
					first_seen INTEGER NOT NULL,
					last_seen INTEGER NOT NULL,
					times_seen INTEGER DEFAULT 1
				);
				CREATE INDEX IF NOT EXISTS idx_devices_last_seen ON devices(last_seen);
				CREATE INDEX IF NOT EXISTS idx_devices_first_seen ON devices(first_seen);
			`,
		},
		{
			version: 2,
			sql: `
				CREATE TABLE IF NOT EXISTS config (
					key TEXT PRIMARY KEY NOT NULL,
					value TEXT NOT NULL,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
			`,
		},
		{
			version: 3,
			sql: `
				CREATE TABLE IF NOT EXISTS scan_history (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					scanned_at INTEGER NOT NULL,
					devices_found INTEGER NOT NULL,
					new_devices INTEGER DEFAULT 0,
					rejoin_devices INTEGER DEFAULT 0
				);
				CREATE INDEX IF NOT EXISTS idx_scan_history_scanned_at ON scan_history(scanned_at);
			`,
		},
	}

	for _, m := range migrations {
		if m.version <= currentVersion {
			continue
		}

		tx, err := s.db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(m.sql); err != nil {
			tx.Rollback()
			return err
		}

		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, m.version); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// RecordScan records a scan in the history.
func (s *Store) RecordScan(devicesFound, newDevices, rejoinDevices int) error {
	_, err := s.db.Exec(`
		INSERT INTO scan_history (scanned_at, devices_found, new_devices, rejoin_devices)
		VALUES (strftime('%s', 'now'), ?, ?, ?)
	`, devicesFound, newDevices, rejoinDevices)
	return err
}
