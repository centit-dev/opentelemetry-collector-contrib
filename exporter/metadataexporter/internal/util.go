package internal

import (
	"database/sql"
	"fmt"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent"
)

func valueHash(name, value string) string {
	return fmt.Sprintf("%s:%s", name, value)
}

func commit(tx *ent.Tx) error {
	commitErr := tx.Commit()
	if commitErr == nil || commitErr == sql.ErrTxDone {
		return nil
	}

	rollbackErr := tx.Rollback()
	if rollbackErr == nil || rollbackErr == sql.ErrTxDone {
		return commitErr
	}

	return fmt.Errorf("commit fails: %v\nrollback fails: %w", commitErr, rollbackErr)
}

func rollback(tx *ent.Tx, err error) error {
	if err == nil {
		return nil
	}

	rollbackErr := tx.Rollback()
	if rollbackErr == nil || rollbackErr == sql.ErrTxDone {
		return err
	}
	return fmt.Errorf("original error: %v\nrollback fails: %w", err, rollbackErr)
}
