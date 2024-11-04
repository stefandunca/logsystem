package logsystem

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const SQLiteDriverID = "sqlite"

type sqliteConfig struct {
	DBPath string   `json:"dbPath"`
	TxAttr []string `json:"txAttr"`
}

// DBDriverFactory implements DriverFactoryInterface
type DBDriverFactory struct {
}

func (f *DBDriverFactory) DriverID() DriverID {
	return DriverID(SQLiteDriverID)
}

func (f *DBDriverFactory) CreateDriver(config json.RawMessage) (drv DriverInterface, err error) {
	var sqliteConfig sqliteConfig
	err = json.Unmarshal(config, &sqliteConfig)
	if err != nil {
		fmt.Printf("Failed to unmarshal SQLite driver config: %v\n", err)
		return nil, errors.Join(errors.New("failed to unmarshal SQLite driver config"), err)
	}

	driver := &SQLiteDriver{
		config: sqliteConfig,
	}

	err = driver.initDB()
	if err != nil {
		return nil, err
	}

	return driver, nil
}

// SQLiteDriver implements DriverInterface
// It logs
type SQLiteDriver struct {
	config sqliteConfig
	db     *sql.DB
}

// TODO: source parameters from available data and config
func (d *SQLiteDriver) initDB() error {
	db, err := sql.Open("sqlite3", d.config.DBPath)
	if err != nil {
		return err
	}

	d.db = db

	attr_columns := ""
	for _, txAttr := range d.config.TxAttr {
		attr_columns += fmt.Sprintf("%s TEXT", txAttr)
	}
	attr_columns += ","

	// TODO: migrate in case of schema changes in the config
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS transactions (
			start_timestamp INTEGER,
			id INTEGER,
			end_timestamp INTEGER,
			%s
			PRIMARY KEY (start_timestamp, id)
		)
	`, attr_columns)

	_, err = d.db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp INTEGER,
			level TEXT,
			message TEXT,
			component TEXT,
			tx_id TEXT,
			FOREIGN KEY(tx_id) REFERENCES transactions(id)
		)
	`)
	return err
}

func (d *SQLiteDriver) Log(data map[Param]string) {
	p := extractKnownParams(data)

	_, err := d.db.Exec(`
		INSERT INTO logs (timestamp, level, message, component, tx_id)
		VALUES (?, ?, ?, ?, ?)
	`, p.Timestamp, p.Level, p.Message, p.Component, p.TxID)

	if err != nil {
		fmt.Printf("Failed to log to SQLite database: %v\n", err)
	}
}

func (d *SQLiteDriver) BeginTx(id TxID, attr map[Param]string) {
	p := extractKnownParams(attr)

	columns := []string{"start_timestamp", "id"}
	values := []interface{}{p.Timestamp, id.String()}

	optionalColumns := []string{}
	optionalValues := []interface{}{}
	for _, txAttr := range d.config.TxAttr {
		if val, ok := attr[Param(txAttr)]; ok {
			optionalColumns = append(optionalColumns, txAttr)
			optionalValues = append(optionalValues, val)
		}
	}

	allColumns := append(columns, optionalColumns...)
	allValues := append(values, optionalValues...)

	placeholders := make([]string, len(allValues))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	insertSQL := fmt.Sprintf(
		"INSERT INTO transactions (%s) VALUES (%s)",
		strings.Join(allColumns, ", "),
		strings.Join(placeholders, ", "),
	)

	// Execute the SQL statement
	_, err := d.db.Exec(insertSQL, allValues...)
	if err != nil {
		fmt.Printf("Failed to log transaction begin to SQLite database: %v\n", err)
		return
	}
}

func (d *SQLiteDriver) EndTx(id TxID) {
	timestamp := time.Now().Unix()

	_, err := d.db.Exec(`
		UPDATE transactions
		SET end_timestamp = ?
		WHERE id = ? AND end_timestamp IS NULL
	`, timestamp, id.String())

	if err != nil {
		fmt.Printf("Failed to log transaction end to SQLite database: %v\n", err)
	}
}

func (d *SQLiteDriver) Stop() {
	d.db.Close()
}
