package stores

import (
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMongoStore(t *testing.T) {
	session := NewMyMongoSession(true)

	store := NewMongoRepository(session)

	gotLists, err := store.GetLists()
	assert.Equal(t, 0, len(gotLists), "new collection should have zero lists")
	assert.Nil(t, err)

	data := models.SampleList()
	err = store.AddList(&data)
	assert.Nil(t, err)

	gotLists, err = store.GetLists()
	assert.Equal(t, 1, len(gotLists), "after adding a list the new collection should have one list")
	assert.Nil(t, err)

	gotList, err := store.GetSingleList(gotLists[0].ID)
	assert.Nil(t, err)
	assert.Equal(t, data.Name, gotList.Name)

	dataToReplace := models.SampleList()
	dataToReplace.Name = "REPLACED"
	err = store.UpdateList(gotList.ID, &dataToReplace)
	assert.Nil(t, err)

	gotList, err = store.GetSingleList(gotLists[0].ID)
	assert.Nil(t, err)
	assert.Equal(t, dataToReplace.Name, gotList.Name)

	err = store.RemoveList(data.ID)
	assert.Nil(t, err)

	_, err = store.GetSingleList(gotLists[0].ID)
	assert.NotNil(t, err)

	err = store.mongoSession.Collection(listsCollectionName).DropCollection()
	assert.Nil(t, err)
}
