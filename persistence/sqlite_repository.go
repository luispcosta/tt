package persistence

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/luispcosta/go-tt/configuration"
	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/utils"
)

type SqliteRepository struct {
}

func NewSqliteRepository() (*SqliteRepository, error) {

}

func (repo *MongoRepository) Initialize() error {
	return nil
}

func (repo *MongoRepository) Shutdown() error {
	return nil
}
