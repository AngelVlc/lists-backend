package stores

import (
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMongoStore(t *testing.T) {
	session := NewMyMongoSession(true)

	store := NewMongoStore(session)

	gotLists := store.GetLists()
	assert.Equal(t, 0, len(gotLists), "new collection should have zero lists")

	data := models.SampleList()
	err := store.AddList(&data)
	assert.Nil(t, err)

	gotLists = store.GetLists()
	assert.Equal(t, 1, len(gotLists), "new collection should have zero lists")

	gotList, err := store.GetSingleList(gotLists[0].ID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, gotList.Name, data.Name)

	dataToReplace := models.SampleList()
	dataToReplace.Name = "REPLACED"
	err = store.UpdateList(gotList.ID.Hex(), &dataToReplace)
	assert.Nil(t, err)

	gotList, err = store.GetSingleList(gotLists[0].ID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, gotList.Name, dataToReplace.Name)

	err = store.RemoveList(data.ID.Hex())
	assert.Nil(t, err)

	_, err = store.GetSingleList(gotLists[0].ID.Hex())
	assert.NotNil(t, err)

	err = store.mongoSession.Collection(ListsCollectionName).DropCollection()
	assert.Nil(t, err)
}
