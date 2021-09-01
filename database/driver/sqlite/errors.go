package driver

import "errors"

var (
    // ErrCreateDatabaseNotImplemented is
    ErrCreateDatabaseNotImplemented = errors.New("SQLite doesn't support CREATE DATABASE statement")
    // ErrDropDatabaseNotImplemented is
    ErrDropDatabaseNotImplemented = errors.New("SQLite doesn't support DROP DATABASE statementl")
    // ErrPrimaryKeyNotImplemented is
    ErrPrimaryKeyNotImplemented = errors.New("SQLite doesn't support ALTER TABLE statement to add or drop primary key")
    // ErrForeignKeyNotImplemented is
    ErrForeignKeyNotImplemented = errors.New("SQLite doesn't support ALTER TABLE statement to add or drop foreign key")
    // ErrAlertTableNotImplemented is
    ErrAlertTableNotImplemented = errors.New("SQLite doesn't support ALTER TABLE statement to modify、drop、rename a column, more details https://www.techonthenet.com/sqlite/tables/alter_table.php")
    // ErrAlertTableMultipleAddNotImplemented is
    ErrAlertTableMultipleAddNotImplemented = errors.New("SQLite doesn't support ALTER TABLE statement to adding multiple columns to a table using a single statement, more details https://www.techonthenet.com/sqlite/tables/alter_table.php")
    // ErrConstraintsNotImplemented is
    ErrConstraintsNotImplemented = errors.New("Constraints not implemented on sqlite, more details Url")
)
