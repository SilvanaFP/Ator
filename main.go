package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	dbConfig := DBConfig{
		User:     "ator",
		Password: "ator",
		Name:     "ator",
		Host:     "localhost",
		Port:     "5432",
	}

	db, err := InitDB(dbConfig)
	if err != nil {
		panic(err)
	}

	srv := server{
		db: db,
	}
	r := gin.Default()
	r.GET("/hello", sayHelloHandler)
	r.GET("/time", tellTimeHandler)
	r.GET("/rewards", srv.getAllRewardsHandler)
	r.GET("/users", srv.getAllUsersHandler)
	r.GET("/tasks", srv.getAllTasksHandler)
	r.GET("/tasksdone", srv.getAllTasksDoneByUserHandler)
	r.POST("/tasksdone", srv.insertTasksDoneHandler)
	r.Run()

}

func sayHelloHandler(c *gin.Context) {
	c.JSON(200, gin.H{"hello": "world"})

}

func tellTimeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"now": time.Now().Format(time.Kitchen)})

}

type DBConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}

func InitDB(c DBConfig) (*sql.DB, error) {
	psqlconn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", c.User, c.Password, c.Name, c.Host, c.Port)
	var err error
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

type server struct {
	db *sql.DB
}

type reward struct {
	Name           string
	RequiredPoints int
}

func (s *server) getAllRewards() ([]reward, error) {
	rows, err := s.db.Query("SELECT name, required_points FROM rewards")
	if err != nil {
		return []reward{}, err
	}
	rewards := []reward{}
	for rows.Next() {
		var r reward
		err := rows.Scan(&r.Name, &r.RequiredPoints)
		if err != nil {
			return rewards, err
		}
		rewards = append(rewards, r)
	}
	return rewards, nil
}

func (s *server) getAllRewardsHandler(c *gin.Context) {
	rewards, err := s.getAllRewards()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, gin.H{"rewards": rewards})

}

type user struct {
	Name string
}

type task struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func (s *server) getAllUsers() ([]user, error) {
	rows, err := s.db.Query("SELECT name FROM users")
	if err != nil {
		return []user{}, err
	}
	users := []user{}
	for rows.Next() {
		var u user
		err := rows.Scan(&u.Name)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *server) getAllUsersHandler(c *gin.Context) {
	users, err := s.getAllUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, gin.H{"users": users})

}

func (s *server) getAllTasks() ([]task, error) {
	rows, err := s.db.Query("SELECT name, id FROM tasks")
	if err != nil {
		return []task{}, err
	}
	tasks := []task{}
	for rows.Next() {
		var t task
		err := rows.Scan(&t.Name, &t.ID)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *server) getAllTasksHandler(c *gin.Context) {
	tasks, err := s.getAllTasks()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, gin.H{"tasks": tasks})

}

func (s *server) getAllTasksDoneByUser(userID int) ([]int, error) {
	rows, err := s.db.Query("SELECT task_id FROM tasks_done WHERE user_id = $1", userID)
	if err != nil {
		return []int{}, err
	}
	taskIDs := []int{}
	for rows.Next() {
		var taskID int
		err := rows.Scan(&taskID)
		if err != nil {
			return taskIDs, err
		}
		taskIDs = append(taskIDs, taskID)
	}
	return taskIDs, nil
}

func (s *server) getAllTasksDoneByUserHandler(c *gin.Context) {
	userIDstr := c.Query("userID")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "user ID is not an int")
		return
	}
	taskIDs, err := s.getAllTasksDoneByUser(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, gin.H{"tasks done": taskIDs})

}

//TODO: ensenyar-me altra sintaxi
func (s *server) insertTasksDone(userID int, taskID int) error {
	_, err := s.db.Exec("INSERT INTO tasks_done (user_id, task_id) VALUES ($1, $2);", userID, taskID)

	return err

}

func (s *server) insertTasksDoneHandler(c *gin.Context) {
	userIDstr := c.Query("userID")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "user ID is not an int")
		return
	}

	taskIDstr := c.Query("taskID")
	taskID, err := strconv.Atoi(taskIDstr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "task ID is not an int")
		return
	}

	err = s.insertTasksDone(userID, taskID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, nil)

}
