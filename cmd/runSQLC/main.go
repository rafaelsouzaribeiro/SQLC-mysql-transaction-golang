package main

import (
	"context"
	"database/sql"

	"github.com/rafaelsouzaribeiro/SQLC-mysql-transaction-golang/internal/db"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()
	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")

	if err != nil {
		panic(err)
	}

	defer dbConn.Close()

	queries := db.New(dbConn)
	// err = queries.CreateCategry(ctx, db.CreateCategryParams{
	// 	ID:          uuid.New().String(),
	// 	Name:        "Backend",
	// 	Description: sql.NullString{String: "Backend Description", Valid: true},
	// })

	// categories, err := queries.ListCategories(ctx)

	// if err != nil {
	// 	panic(err)
	// }

	// for _, obj := range categories {
	// 	println(obj.ID, obj.Name, obj.Description.String)
	// }

	// queries.UpdateCategory(ctx, db.UpdateCategoryParams{
	// 	ID:          "361772e5-af0c-4e3c-b252-4e4214d8d24d",
	// 	Name:        "Backend updated",
	// 	Description: sql.NullString{String: "Backend Description updated", Valid: true},
	// })

	err = queries.DeleteCategory(ctx, "361772e5-af0c-4e3c-b252-4e4214d8d24d")

	if err != nil {
		panic(err)
	}

	categories, err := queries.ListCategories(ctx)

	if err != nil {
		panic(err)
	}

	for _, obj := range categories {
		println(obj.ID, obj.Name, obj.Description.String)
	}
}
