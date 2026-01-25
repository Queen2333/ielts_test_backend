package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 连接数据库
	dsn := "ielts_alex:Yx180236@tcp(172.25.138.133:3306)/ielts_database?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database successfully")

	// 执行删除 source 列的 SQL
	stmt := "ALTER TABLE writing_part_list DROP COLUMN source"

	fmt.Println("Executing: DROP COLUMN source from writing_part_list...")
	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatalf("Failed to execute statement: %v", err)
	}

	fmt.Println("✓ Successfully dropped source column")
	fmt.Println("\n✅ Migration 004 completed!")
	fmt.Println("\nSummary:")
	fmt.Println("- Removed source field from writing_part_list")
	fmt.Println("- writing_part_list now has: type (data source), task_type (Task 1/2)")
}
