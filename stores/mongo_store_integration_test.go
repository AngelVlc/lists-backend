package stores

import (
	"fmt"
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

	fmt.Println(gotLists[0], gotLists[0].ID)

	gotList, err := store.GetSingleList(gotLists[0].ID)
	assert.Nil(t, err)
	assert.Equal(t, gotList.Name, data.Name)

	dataToReplace := models.SampleList()
	dataToReplace.Name = "REPLACED"
	err = store.UpdateList(gotList.ID, &dataToReplace)
	assert.Nil(t, err)

	gotList, err = store.GetSingleList(gotLists[0].ID)
	assert.Nil(t, err)
	assert.Equal(t, gotList.Name, dataToReplace.Name)

	err = store.RemoveList(data.ID)
	assert.Nil(t, err)

	_, err = store.GetSingleList(gotLists[0].ID)
	assert.NotNil(t, err)

	err = store.mongoSession.Collection(listsCollectionName).DropCollection()
	assert.Nil(t, err)
}
