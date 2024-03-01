package repository

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/common"
	"github.com/dembygenesis/local.tools/internal/persistence/mysql/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew_Success(t *testing.T) {
	db, _, cleanup := helpers.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger: common.GetLogger(context.Background()),
		Db:     db,
	}

	_, err := New(cfg)
	require.NoError(t, err)
}

func TestNew_Fail_Validate(t *testing.T) {
	cfg := &Config{
		Logger: common.GetLogger(context.Background()),
	}

	_, err := New(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "validate:", "expected err to contains 'validate:'")
	require.Contains(t, err.Error(), "validate struct:", "expected err to contains 'validate struct:'")
}

func TestNew_Fail_Ping(t *testing.T) {
	db, _, cleanup := helpers.TestGetMockMariaDB(t)
	defer cleanup(true)

	cfg := &Config{
		Logger: common.GetLogger(context.Background()),
		Db:     db,
	}

	err := db.Close()
	require.NoError(t, err, "expected db to close successfully")

	err = db.Ping()
	require.Error(t, err, "expected to be unable to ping a db")

	_, err = New(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "validate:", "expected err to contains 'validate:'")
	require.Contains(t, err.Error(), "ping:", "expected err to contains 'ping:'")
}
