package models

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MONGODB = "mongodb://localhost:27017"

type DB struct {
	ctx context.Context
	db  *mongo.Database
}

func NewDB(name string) *DB {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(MONGODB))
	db := client.Database("test")

	return &DB{
		ctx: ctx,
		db:  db,
	}
}

func (d *DB) Collection(name string) *mongo.Collection {
	return d.db.Collection(name)
}

func (d *DB) DropCollection(name string) {
	d.db.Collection(name).Drop(d.ctx)
}

type SortPostSlice []Post

// Len is the number of elements in the collection.
func (s SortPostSlice) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s SortPostSlice) Less(i int, j int) bool {
	iRune, _ := utf8.DecodeRuneInString(s[i].Title)
	jRune, _ := utf8.DecodeRuneInString(s[j].Title)
	return int32(iRune) < int32(jRune)
}

// Swap swaps the elements with indexes i and j.
func (s SortPostSlice) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

func TestPost(t *testing.T) {
	db := NewDB("test")

	t.Run("Test GetMongoCollectionName", func(t *testing.T) {
		want := "posts"

		p := Post{}
		get := p.GetMongoCollectionName()

		if get != want {
			t.Errorf("Get: %s, want: %s\n", get, want)
		}
	})

	dataSet := []Post{
		Post{Title: "", Content: "", Created: "", Author: Author{}},
		Post{Title: "Test", Content: "Test", Created: "Test", Author: Author{}},
		Post{Title: "Title1", Content: "Content1", Created: "some text", Author: Author{"User", "nohoby"}},
	}

	p := Post{}
	db.DropCollection(p.GetMongoCollectionName())

	t.Run("Test insert and select post", func(t *testing.T) {
		for idx, _ := range dataSet {
			item := &dataSet[idx]
			err := item.Insert(db.ctx, db.db)
			if err != nil {
				t.Error(err)
			}

			get, err := GetPost(db.ctx, db.db, item.ID)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(get, item) {
				t.Errorf("Get %v, want %v", get, item)
			}
		}
	})

	t.Run("Test update", func(t *testing.T) {
		for idx, _ := range dataSet {
			item := &dataSet[idx]
			item.Title = item.Title + "update"
			if err := item.Update(db.ctx, db.db); err != nil {
				t.Error(err)
			}

			get, err := GetPost(db.ctx, db.db, item.ID)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(get, item) {
				t.Errorf("Get %v, want %v", get, item)
			}
		}
	})

	t.Run("Test GetPosts", func(t *testing.T) {
		getSlice, err := GetPosts(db.ctx, db.db)
		sort.Sort(SortPostSlice(getSlice))
		sort.Sort(SortPostSlice(dataSet))
		if err != nil {
			t.Error(err)
		}
		if len(getSlice) != len(dataSet) {
			t.Errorf("Get %d, want %d", len(getSlice), len(dataSet))
		}
		for i, item := range dataSet {
			if !reflect.DeepEqual(getSlice[i], dataSet[i]) {
				t.Errorf("Get %v, want %v", getSlice[i], item)
			}
		}
	})
}
