package stores

import (
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMongoStore(t *testing.T) {
	session := NewMyMongoSession(true)

	store := NewMongoRepository(session, session.Collection("lists"))

	gotLists := []models.GetListsResultDto{}
	err := store.Get(&gotLists)
	assert.Equal(t, 0, len(gotLists), "new collection should have zero lists")
	assert.Nil(t, err)

	data := models.SampleList()
	err = store.Add(&data)
	assert.Nil(t, err)

	gotLists = []models.GetListsResultDto{}
	err = store.Get(&gotLists)
	assert.Equal(t, 1, len(gotLists), "after adding a list the new collection should have one list")
	assert.Nil(t, err)

	foundList := models.List{}
	err = store.GetSingle(gotLists[0].ID, &foundList)
	assert.Nil(t, err)
	assert.Equal(t, data.Name, foundList.Name)

	dataToReplace := models.SampleList()
	dataToReplace.Name = "REPLACED"
	err = store.Update(foundList.ID, &dataToReplace)
	assert.Nil(t, err)

	foundList = models.List{}
	err = store.GetSingle(gotLists[0].ID, &foundList)
	assert.Nil(t, err)
	assert.Equal(t, dataToReplace.Name, foundList.Name)

	err = store.Remove(data.ID)
	assert.Nil(t, err)

	foundList = models.List{}
	err = store.GetSingle(gotLists[0].ID, &foundList)
	assert.NotNil(t, err)

	err = store.mongoCollection.DropCollection()
	assert.Nil(t, err)
}
