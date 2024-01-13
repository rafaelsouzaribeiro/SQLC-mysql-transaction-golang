package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rafaelsouzaribeiro/SQLC-mysql-transaction-golang/internal/db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type CourseDB struct {
	dbConn *sql.DB
	*db.Queries
}

type CourseParams struct {
	ID          string
	Name        string
	Description sql.NullString
	Price       float64
}

type CategoryParams struct {
	ID          string
	Name        string
	Description sql.NullString
}

func (c *CourseDB) CallTx(ctx context.Context, fn func(*db.Queries) error) error {
	// Inicia uma transação
	tx, err := c.dbConn.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := db.New(tx)

	err = fn(q)

	// se der erro na fn anônima não executa
	if err != nil {
		// Se der erro ele executa o rollback verficando se deu erro nele mesmo
		if errb := tx.Rollback(); errb != nil {
			return fmt.Errorf("error rollback %v, original error %w", errb, err)
		}
		return err
	}

	return tx.Commit()
}

func (c *CourseDB) CreateCourseAndCatgeory(ctx context.Context, argsCaCategory CategoryParams, argsCourse CourseParams) error {

	// Somente executa o commit se não der erro quanto no criar a category e o course
	err := c.CallTx(ctx, func(q *db.Queries) error {
		var err error
		err = q.CreateCategry(ctx, db.CreateCategryParams{
			ID:          argsCaCategory.ID,
			Name:        argsCaCategory.Name,
			Description: argsCaCategory.Description,
		})

		if err != nil {
			return err
		}

		err = q.CreateCourse(ctx, db.CreateCourseParams{
			ID:          argsCourse.ID,
			Name:        argsCourse.Name,
			Description: argsCourse.Description,
			Price:       argsCourse.Price,
			CategoryID:  argsCaCategory.ID,
		})

		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return err
	}

	return nil
}

func NewCourseDb(dbConn *sql.DB) *CourseDB {
	return &CourseDB{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

func main() {
	ctx := context.Background()
	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")

	if err != nil {
		panic(err)
	}

	defer dbConn.Close()

	CreateCourse(ctx, dbConn)
	ListCourses(ctx, dbConn)

}

func ListCourses(ctx context.Context, dbConn db.DBTX) {
	queries := db.New(dbConn)
	couses, err := queries.ListCourses(ctx)

	if err != nil {
		panic(err)
	}

	for _, course := range couses {
		fmt.Printf("Category: %s, Course ID: %s, Course Name: %s, Course Description: %s, Course Price: %f\n",
			course.CategoryName, course.ID, course.Name, course.Description.String, course.Price)
	}
}

func CreateCourse(ctx context.Context, dbConn *sql.DB) {

	courseArgs := CourseParams{
		ID:          uuid.New().String(),
		Name:        "GO 3",
		Description: sql.NullString{String: "GO description 3", Valid: true},
		Price:       10,
	}

	categoryArgs := CategoryParams{
		ID:          uuid.New().String(),
		Name:        "Backend 3",
		Description: sql.NullString{String: "Backend description 3", Valid: true},
	}

	CourseDB := NewCourseDb(dbConn)
	err := CourseDB.CreateCourseAndCatgeory(ctx, categoryArgs, courseArgs)

	if err != nil {
		panic(err)
	}
}
