package todos

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

const (
	getPGVersion = "SELECT version();"
)

type PGX struct {
	Conn *pgxpool.Pool
}

func NewPgxDB(dbConnectionString string, maxConnectionsInPool int, log *log.Logger) (Storage, error) {
	var psql PGX
	var successOrFailure = "OK"

	var parsedConfig *pgx.ConnConfig
	var err error
	parsedConfig, err = pgx.ParseConfig(dbConnectionString)
	if err != nil {
		return nil, err
	}

	dbHost := parsedConfig.Host
	dbPort := parsedConfig.Port
	dbUser := parsedConfig.User
	dbPass := parsedConfig.Password
	dbName := parsedConfig.Database

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=%d", dbHost, dbPort, dbUser, dbPass, dbName, maxConnectionsInPool)

	connPool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		successOrFailure = "FAILED"
		log.Printf("Connecting to database %s as user %s : %s \n", dbName, dbUser, successOrFailure)
		log.Fatalf("ERROR TRYING DB CONNECTION : %v ", err)
	} else {
		log.Printf("Connecting to database %s as user %s : %s \n", dbName, dbUser, successOrFailure)
		// golog.Info("Fetching one record to test if db connection is valid...\n")
		var version string
		if errPing := connPool.QueryRow(context.Background(), getPGVersion).Scan(&version); errPing != nil {
			log.Printf("Connection is invalid ! ")
			log.Fatalf("DB ERROR scanning row: %s", errPing)
		}
		log.Printf("SUCCESS Connecting to Postgres version : [%s]", version)
	}

	psql.Conn = connPool
	return &psql, err
}

//Create will store the new task in the store
func (m *PGX) Create(todo NewTodo) (*Todo, error) {
	if len(todo.Task) < 1 {
		return nil, errors.New("todo task cannot be empty")
	}
	if len(todo.Task) < 6 {
		return nil, errors.New("CreateTodo task minLength is 5")
	}
	now := time.Now()
	t := &Todo{
		Completed:   false,
		CompletedAt: nil,
		CreatedAt:   &now,
		Id:          0,
		Task:        todo.Task,
	}
	panic("implement SQL query")
	return t, nil
}

func (m *PGX) List(offset, limit int) ([]Todo, error) {
	var res []Todo
	panic("implement SQL query")
	return res, nil
}

func (m *PGX) Get(id int32) (*Todo, error) {
	if m.Exist(id) {
		panic("implement SQL query")
		return nil, nil
	}
	return nil, errors.New("todo with this id does not exist")
}

// GetMaxId returns the maximum value of todos id existing in store.
func (m *PGX) GetMaxId() (int32, error) {
	existingMaxId := int32(0)
	panic("implement SQL query")
	return existingMaxId, nil
}

// Exist returns true only if a todos with the specified id exists in store.
func (m *PGX) Exist(id int32) bool {
	panic("implement SQL query")
	return true
}

func (m *PGX) Count() (int32, error) {
	panic("implement SQL query")
	return int32(0), nil
}

func (m *PGX) Update(id int32, todo Todo) (*Todo, error) {
	if m.Exist(id) {
		panic("implement SQL query")
		/*
			existingTodo := m.Todos[id]
			now := time.Now()
			// cannot override id field
			todo.Id = existingTodo.Id
			// cannot override CreatedAt field
			todo.CreatedAt = existingTodo.CreatedAt
			switch todo.Completed {
			case true:
				if existingTodo.Completed == false {
					todo.CompletedAt = &now
				}
			case false:
				if existingTodo.Completed == true {
					// task was completed, but user changed it to not completed
					todo.CompletedAt = nil
				}
			// in all other cases the value of CompletedAt should not be changed
			default:
				todo.CompletedAt = existingTodo.CompletedAt
			}
		*/

		return nil, nil
	}
	return nil, errors.New("todo with this id does not exist")
}

func (m *PGX) Delete(id int32) error {
	if m.Exist(id) {
		panic("implement SQL query")
		return nil
	}
	return errors.New("todo with this id does not exist")
}
