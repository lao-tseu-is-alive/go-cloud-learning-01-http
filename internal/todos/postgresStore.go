package todos

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

const (
	getPGVersion    = "SELECT version();"
	todosList       = "SELECT id, task, completed, created_at, completed_at FROM todos ORDER BY id;"
	todosGet        = "SELECT id, task, completed, created_at, completed_at FROM todos WHERE id=$1;"
	todosCompleted  = "SELECT completed FROM todos WHERE id=$1"
	todosExist      = "SELECT COUNT(*) FROM todos WHERE id=$1"
	todosCount      = "SELECT COUNT(*) FROM todos"
	todosMaxId      = "SELECT MAX(id) FROM todos"
	todosCreate     = "INSERT INTO todos (task) VALUES($1) RETURNING id;"
	todosUpdate     = "UPDATE todos SET task=$1, completed=$2, completed_at=$3 WHERE id=$4"
	todosUpdateTask = "UPDATE todos SET task=$1 WHERE id=$2"
	todosDelete     = "DELETE FROM todos WHERE id = $1"
)

type PGX struct {
	Conn *pgxpool.Pool
	log  *log.Logger
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
		log.Printf("Connecting to database %s with user %s : %s \n", dbName, dbUser, successOrFailure)
		return nil, errors.New(fmt.Sprintf("error connecting to database. err : %s", err))
	} else {
		log.Printf("Connected to database %s with user %s : %s \n", dbName, dbUser, successOrFailure)
		// let's first check that we can really make a query by querying the postgres version
		var version string
		if errPing := connPool.QueryRow(context.Background(), getPGVersion).Scan(&version); errPing != nil {
			log.Printf("connection is invalid ! ")
			log.Fatalf("DB ERROR scanning row: %s", errPing)
			return nil, errPing
		}
		var numberOfTodos int
		if errTodosTable := connPool.QueryRow(context.Background(), todosCount).Scan(&numberOfTodos); errTodosTable != nil {
			log.Printf("the database does not contain the table «todos»  ! ")
			return nil, errors.New("database does not contain the table «todos»")
		}

		log.Printf("SUCCESS Connected to Postgres DB ver: [%s]", version)
		log.Printf("SUCCESS database contains %d records in todos", numberOfTodos)
	}

	psql.Conn = connPool
	psql.log = log
	return &psql, err
}

// getQueryInt is a postgres helper function for a query expecting an integer result
func (db *PGX) getQueryInt(sql string, arguments ...interface{}) (result int, err error) {
	err = db.Conn.QueryRow(context.Background(), sql, arguments...).Scan(&result)
	if err != nil {
		db.log.Printf("error : getQueryInt(%s) queryRow unexpectedly failed. args : (%v), error : %v\n", sql, arguments, err)
		return 0, err
	}
	return result, err
}

// getQueryBool is a postgres helper function for a query expecting an integer result
func (db *PGX) getQueryBool(sql string, arguments ...interface{}) (result bool, err error) {
	err = db.Conn.QueryRow(context.Background(), sql, arguments...).Scan(&result)
	if err != nil {
		db.log.Printf("error : getQueryBool(%s) queryRow unexpectedly failed. args : (%v), error : %v\n", sql, arguments, err)
		return false, err
	}
	return result, err
}

// execActionQuery is a postgres helper function for an action query, returning the numbers of rows affected
func (db *PGX) execActionQuery(sql string, arguments ...interface{}) (rowsAffected int, err error) {
	commandTag, err := db.Conn.Exec(context.Background(), sql, arguments...)
	if err != nil {
		db.log.Printf("execActionQuery unexpectedly failed with sql: %v . Args(%+v), error : %v", sql, arguments, err)
		return 0, err
	}
	return int(commandTag.RowsAffected()), err
}

func (db *PGX) Close() {
	db.Conn.Close()
	return
}

//Create will store the new task in the store
func (db *PGX) Create(todo NewTodo) (*Todo, error) {
	db.log.Printf("info : Entering Create(%#v)", todo)
	if len(todo.Task) < 1 {
		return nil, errors.New("todo task cannot be empty")
	}
	if len(todo.Task) < 6 {
		return nil, errors.New("CreateTodo task minLength is 5")
	}
	var lastInsertId int = 0
	err := db.Conn.QueryRow(context.Background(), todosCreate, todo.Task).Scan(&lastInsertId)
	if err != nil {
		db.log.Printf("error : Create(%v) unexpectedly failed. error : %v", todo.Task, err)
		return nil, err
	}
	db.log.Printf("info : Create(%v) created with id : %v", todo.Task, lastInsertId)

	// if we get to here all is good, so let's retrieve a fresh copy to send it back
	createdTodo, err := db.Get(int32(lastInsertId))
	if err != nil {
		return nil, GetErrorF("error : todos was created, but can not be retrieved", err)
	}
	return createdTodo, nil
}

func (db *PGX) List(offset, limit int) ([]*Todo, error) {
	var res []*Todo

	err := pgxscan.Select(context.Background(), db.Conn, &res, todosList)
	if err != nil {
		db.log.Printf("error : List pgxscan.Select unexpectedly failed, error : %v", err)
		return nil, err
	}
	if res == nil {
		db.log.Println("info : List returned no results ")
		return nil, errors.New("records not found")
	}

	return res, nil
}

func (db *PGX) Get(id int32) (*Todo, error) {
	db.log.Printf("info : Get(%d) entering...", id)
	if db.Exist(id) == true {
		res := &Todo{
			Completed:   false,
			CompletedAt: nil,
			CreatedAt:   nil,
			Id:          0,
			Task:        "",
		}
		err := pgxscan.Get(context.Background(), db.Conn, res, todosGet, id)
		if err != nil {
			db.log.Printf("error : Get(%d) pgxscan.Select unexpectedly failed, error : %v", id, err)
			return nil, err
		}
		if res == nil {
			db.log.Printf("info : Get(%d) returned no results ", id)
			return nil, errors.New("records not found")
		}
		return res, nil
	}
	db.log.Printf("info : Get(%d) id does not exist", id)
	return nil, errors.New("todo with this id does not exist")
}

// GetMaxId returns the maximum value of todos id existing in store.
func (db *PGX) GetMaxId() (int32, error) {
	existingMaxId, err := db.getQueryInt(todosMaxId)
	if err != nil {
		db.log.Printf("getMaxId() could not be retrieved from DB. failed db.Query err: %v", err)
		return 0, err
	}
	return int32(existingMaxId), nil
}

// Exist returns true only if a todos with the specified id exists in store.
func (db *PGX) Exist(id int32) bool {
	count, err := db.getQueryInt(todosExist, id)
	if err != nil {
		db.log.Printf("exist(%d) could not be retrieved from DB. failed db.Query err: %v", id, err)
		return false
	}
	if count > 0 {
		db.log.Printf("info : Exist(%d) id does exist  count:%v", id, count)
		return true
	} else {
		db.log.Printf("info : Exist(%d) id does not exist count:%v", id, count)
		return false
	}
}

// Count returns the number of todos stored in DB
func (db *PGX) Count() (int32, error) {
	count, err := db.getQueryInt(todosCount)
	if err != nil {
		db.log.Printf("count(*) could not be retrieved from DB. failed db.Query err: %v", err)
		return 0, err
	}
	return int32(count), nil
}

// Update the todos stored in DB with given id and other information in struct
func (db *PGX) Update(id int32, todo Todo) (*Todo, error) {
	if db.Exist(id) {
		// first check business rules for task field
		if len(todo.Task) < 1 {
			return nil, errors.New("todo task cannot be empty")
		}
		if len(todo.Task) < 6 {
			return nil, errors.New("CreateTodo task minLength is 5")
		}
		updateAll := true
		var rowsAffected int = 0
		var err error
		now := time.Now()
		// implements basic Business Rules !
		// let's first check if task was already completed in DB
		alreadyCompleted, _ := db.getQueryBool(todosCompleted, id)
		switch todo.Completed {
		case true:
			if alreadyCompleted == false {
				// this task was not completed, now it is, so update CompletedAt
				todo.CompletedAt = &now
			}
		case false:
			if alreadyCompleted == true {
				// this task was completed, but user changed it to not completed so update with nil
				todo.CompletedAt = nil
			}
		default:
			// in all other cases the values of Completed and CompletedAt fields should not be changed in DB
			// so here let's update only the Task field
			rowsAffected, err = db.execActionQuery(todosUpdateTask, todo.Task, id)
			updateAll = false
		}
		if updateAll {
			rowsAffected, err = db.execActionQuery(todosUpdate, todo.Task, todo.Completed, todo.CompletedAt, id)
		}
		if err != nil {
			return nil, GetErrorF("error : todos could not be updated", err)
		}
		if rowsAffected < 1 {
			return nil, GetErrorF("error : todos was not updated", err)
		}
		// if we get to here all is good, so let's retrieve a fresh copy to send it back
		updatedTodo, err := db.Get(id)
		if err != nil {
			return nil, GetErrorF("error : todos was updated, but can not be retrieved", err)
		}
		return updatedTodo, nil
	}
	db.log.Printf("info : Update(%d) id does not exist", id)
	return nil, errors.New("todo with this id does not exist")
}

// Delete the todos stored in DB with given id
func (db *PGX) Delete(id int32) error {
	if db.Exist(id) {
		rowsAffected, err := db.execActionQuery(todosDelete, id)
		if err != nil {
			return GetErrorF("error : todos could not be deleted", err)
		}
		if rowsAffected < 1 {
			return GetErrorF("error : todos was not deleted", err)
		}
		// if we get to here all is good
		return nil
	}
	db.log.Printf("info : Delete(%d) id does not exist", id)
	return errors.New("todo with this id does not exist")
}
