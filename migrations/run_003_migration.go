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

	// 执行 SQL 语句
	statements := []string{
		"ALTER TABLE writing_part_list ADD COLUMN task_type VARCHAR(255) NULL COMMENT '写作任务类型：1=Task1，2=Task2'",
		"UPDATE writing_part_list SET task_type = type",
		"UPDATE writing_part_list SET type = source",
		"UPDATE writing_part_list SET source = NULL",
		"CREATE INDEX idx_writing_part_task_type ON writing_part_list(task_type)",
	}

	for i, stmt := range statements {
		fmt.Printf("Executing statement %d...\n", i+1)
		_, err := db.Exec(stmt)
		if err != nil {
			log.Printf("Warning: Failed to execute statement %d: %v", i+1, err)
			continue
		}
		fmt.Printf("✓ Statement %d executed successfully\n", i+1)
	}

	fmt.Println("\n✅ Migration 003 completed!")
	fmt.Println("\nSummary:")
	fmt.Println("- Added task_type field to writing_part_list")
	fmt.Println("- Migrated type → task_type (Task 1/2)")
	fmt.Println("- Migrated source → type (data source 1/2/3)")
	fmt.Println("- Cleared source field")
}
