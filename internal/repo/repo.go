package repo

import (
	"context"

	"github.com/autocorrectoff/banners/internal/dto"
	mdb "github.com/minus5/svckit/pkg/mdb2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dbName     = "matches"
	colName    = "matches"
	idAttr     = "_id"
)

type Db struct {
	mdb.Mdb
}

func New() (*Db, error) {
	return NewWith(mdb.DefaultConnStr())
}

func NewWith(connStr string) (*Db, error) {
	db := &Db{}
	if err := db.Init(connStr,
		mdb.MajoritySafe(),
		mdb.Name(dbName)); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *Db) InsertOrUpdate(match *dto.Match) error {
	err := db.Insert(colName, match)
	if err != nil {
		if err == mdb.ErrDuplicate {
			db.Use(colName, "update_match", func(c *mongo.Collection) error {
				ctx := context.Background()
				filter := bson.M{idAttr: bson.M{"$eq": match.BaseId}}
				_, err := c.ReplaceOne(ctx, filter, match)
				return err
			})
		} else {
			return err
		}
	}
	return nil
}

func (db *Db) FindByIdIn(ids []int) ([]dto.Match, error) {
	var matches []dto.Match
	err := db.Use(colName, "find_match", func(c *mongo.Collection) error {
		ctx := context.Background()
		filter := bson.M{idAttr: bson.M{"$in": ids}}
		cursor, err := c.Find(ctx, filter)
		if err != nil {
			return err
		}
		return cursor.All(ctx, &matches)
	})
	return matches, err
}
