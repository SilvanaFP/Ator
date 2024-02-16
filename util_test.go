package main

import (
	"database/sql"
	"testing"
)

type testHarness struct {
	db *sql.DB
}

func newTestHarness(t *testing.T) testHarness {
	db, err := InitDB(DBConfig{User: "ator", Password: "ator", Name: "test_ator", Host: "localhost", Port: "5432"})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := deleteTestData(db); err != nil {
			t.Fatal(err)
		}
		db.Close()
	})
	//TODO: implement use of migrations
	if err := populateTestDB(db); err != nil {
		t.Fatal(err)
	}

	return testHarness{db: db}
}

func populateTestDB(pg *sql.DB) error {
	_, err := pg.Exec(`
	INSERT INTO users (name, password) VALUES ('TestUser','1234');

	INSERT INTO rewards (name, required_points, created_at) VALUES ('cinema', 20, NOW());
	INSERT INTO rewards (name, required_points) VALUES ('weekend trip', 10);
	
	INSERT INTO tasks (name, reward_id) VALUES ('laundry', 1);
	INSERT INTO tasks (name, reward_id) VALUES ('take out trash', 1);
	
	INSERT INTO tasks (name, reward_id) VALUES ('yoga', 2);
	INSERT INTO tasks (name, reward_id) VALUES ('go for a run', 2);
	
	INSERT INTO tasks_done (user_id, task_id , verified) VALUES (1, 1, 'f');
	INSERT INTO tasks_done (user_id, task_id , verified) VALUES (1, 4, 'f');
	`,
	)

	return err
}

func deleteTestData(pg *sql.DB) error {
	_, err := pg.Exec(`
	TRUNCATE tasks_done RESTART IDENTITY CASCADE;
	TRUNCATE tasks RESTART IDENTITY CASCADE;
	TRUNCATE rewards RESTART IDENTITY CASCADE;
	TRUNCATE users RESTART IDENTITY CASCADE;
	`,
	)

	return err
}
