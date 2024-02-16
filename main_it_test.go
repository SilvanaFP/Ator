package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_GetAllRewards(t *testing.T) {
	h := newTestHarness(t)
	srv := server{
		db: h.db,
	}
	actual, err := srv.getAllRewards()
	assert.Nil(t, err)

	expected := []reward{
		{
			Name:           "cinema",
			RequiredPoints: 20,
		},
		{
			Name:           "weekend trip",
			RequiredPoints: 10,
		},
	}

	assert.Equal(t, expected, actual)
}

func TestServer_GetAllTasks(t *testing.T) {
	h := newTestHarness(t)
	srv := server{
		db: h.db,
	}
	actual, err := srv.getAllTasks()
	assert.Nil(t, err)

	expected := []task{
		{
			Name: "laundry",
			ID:   1,
		},
		{
			Name: "take out trash",
			ID:   2,
		},
		{
			Name: "yoga",
			ID:   3,
		},
		{
			Name: "go for a run",
			ID:   4,
		},
	}

	assert.Equal(t, expected, actual)
}

func TestServer_GetAllUSers(t *testing.T) {
	h := newTestHarness(t)
	srv := server{
		db: h.db,
	}
	actual, err := srv.getAllUsers()
	assert.Nil(t, err)

	expected := []user{
		{
			Name: "TestUser",
		},
	}

	assert.Equal(t, expected, actual)
}

func TestServer_GetAllTasksDoneByUser(t *testing.T) {
	h := newTestHarness(t)
	srv := server{
		db: h.db,
	}
	actual, err := srv.getAllTasksDoneByUser(1)
	assert.Nil(t, err)

	expected := []int{1, 4}

	assert.Equal(t, expected, actual)
}

func TestServer_InsertTasksDone(t *testing.T) {
	h := newTestHarness(t)
	srv := server{
		db: h.db,
	}

	//first we check how many tasks done the user 1 has
	actual, err := srv.getAllTasksDoneByUser(1)
	assert.Nil(t, err)

	expected := []int{1, 4}

	assert.Equal(t, expected, actual)

	//now we insert a new entry
	err = srv.insertTasksDone(1, 3)
	assert.Nil(t, err)

	//now we confirm that the task with id 3 has been added
	actual, err = srv.getAllTasksDoneByUser(1)
	assert.Nil(t, err)

	expected = []int{1, 4, 3}

	assert.Equal(t, expected, actual)

}
