package stores

import (
	"testing"

	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestMongoStore(t *testing.T) {
	session := NewMyMongoSession(true)

	repository := session.GetRepository("lists")

	gotLists := []models.GetListsResultDto{}
	err := repository.Get(&gotLists, nil, bson.M{"name": 1})
	assert.Equal(t, 0, len(gotLists), "new collection should have zero lists")
	assert.Nil(t, err)

	data := models.SampleList()
	id, err := repository.Add(&data)
	assert.Nil(t, err)
	assert.NotEmpty(t, id)

	foundList := models.List{}
	err = repository.GetOne(&foundList, bson.D{{"_id", id}}, nil)
	assert.Nil(t, err)
	assert.Equal(t, data.Name, foundList.Name)

	gotLists = []models.GetListsResultDto{}
	err = repository.Get(&gotLists, nil, bson.M{"name": 1})
	assert.Equal(t, 1, len(gotLists), "after adding a list the new collection should have one list")
	assert.Nil(t, err)

	foundList = models.List{}
	err = repository.GetOne(&foundList, bson.D{{"_id", gotLists[0].ID}}, nil)
	assert.Nil(t, err)
	assert.Equal(t, data.Name, foundList.Name)

	dataToReplace := models.SampleList()
	dataToReplace.Name = "REPLACED"
	dataToReplace.ID = foundList.ID
	err = repository.Update(bson.D{{"_id", foundList.ID}}, &dataToReplace)
	assert.Nil(t, err)

	foundList = models.List{}
	err = repository.GetOne(&foundList, bson.D{{"_id", gotLists[0].ID}}, nil)
	assert.Nil(t, err)
	assert.Equal(t, dataToReplace.Name, foundList.Name)

	err = repository.Remove(bson.D{{"_id", data.ID}})
	assert.Nil(t, err)

	foundList = models.List{}
	err = repository.GetOne(&foundList, bson.D{{"_id", gotLists[0].ID}}, nil)
	assert.NotNil(t, err)

	err = session.session.DB(session.databaseName).C("lists").DropCollection()
	assert.Nil(t, err)
}
