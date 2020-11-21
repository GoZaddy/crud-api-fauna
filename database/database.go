package database

import (
	"fmt"
	f "github.com/fauna/faunadb-go/v3/faunadb"
	"github.com/gozaddy/crud-api-fauna/models"
)

const (
	ReadingItemCollection = "reading_items"
	ReadingItemByTypeIndex = "reading_items_by_type"
)

type faunaDB struct {
	client *f.FaunaClient
	secret string
}

type FaunaDB interface {
	Init() error
	NewID() (string, error)
	FaunaClient() *f.FaunaClient
	GetDocument(collection string, documentID string) (f.Value, error)
	AddDocument(collection string, object models.Model) (f.Value, error)
	DeleteDocument(collection string, documentID string) error
	UpdateDocument(collection string, documentID string, update interface{}) error
}

func NewFaunaDB(secret string) FaunaDB {
	return &faunaDB{nil, secret}
}

func (fdb *faunaDB) Init() error{
	fdb.client = f.NewFaunaClient(fdb.secret)
	res, err := fdb.client.Query(
		f.If(
			f.Exists(f.Collection(ReadingItemCollection)),
			"Exists!",
			f.Do(
				f.CreateCollection(
					f.Obj{
						"name": "reading_items",
					},
				),
				f.CreateIndex(
					f.Obj{
						"name": ReadingItemByTypeIndex,
						"source": f.Collection(ReadingItemCollection),
						"terms": f.Arr{f.Obj{"field": f.Arr{"data", "type"}}},
					},
				),
			),
		),
	)
	if err != nil{
		return err
	}
	fmt.Println(res)
	return nil
}

func (fdb *faunaDB) NewID() (string, error) {
	res, err := fdb.client.Query(f.NewId())
	if err != nil{
		return "", err
	}
	var id string
	if err = res.Get(&id); err != nil{
		return "", err
	}
	return id, nil
}

func (fdb *faunaDB) FaunaClient() *f.FaunaClient {
	return fdb.client
}

func (fdb *faunaDB) GetDocument(collection string, documentID string) (f.Value, error) {
	return fdb.client.Query(f.Get(f.RefCollection(f.Collection(collection), documentID)))
}


func (fdb *faunaDB) AddDocument(collection string, object models.Model) (f.Value, error) {
	if err := object.Validate(); err != nil{
		return nil, err
	}

	return fdb.client.Query(f.Create(f.RefCollection(f.Collection(collection), object.UniqueID()), f.Obj{"data": object}))
}

func (fdb *faunaDB) DeleteDocument(collection string, documentID string) error {
	_, err := fdb.client.Query(f.Delete(f.RefCollection(f.Collection(collection), documentID)))
	return err
}

func (fdb *faunaDB) UpdateDocument(collection string, documentID string, update interface{}) error {
	_, err := fdb.client.Query(
		f.Update(
			f.RefCollection(f.Collection(collection), documentID),
			f.Obj{"data": update},
		),
	)
	return err
}
