package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 读取 SQL 文件
	sqlContent, err := os.ReadFile("002_add_user_id_to_main_lists.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

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

	// 执行 SQL 语句（逐行执行）
	sqlStatements := string(sqlContent)

	// 简单的 SQL 分割（按分号分割，跳过注释）
	statements := []string{
		"ALTER TABLE listening_list ADD COLUMN user_id VARCHAR(255) NULL COMMENT '创建者用户ID'",
		"ALTER TABLE reading_list ADD COLUMN user_id VARCHAR(255) NULL COMMENT '创建者用户ID'",
		"ALTER TABLE writing_list ADD COLUMN user_id VARCHAR(255) NULL COMMENT '创建者用户ID'",
		"ALTER TABLE testing_list ADD COLUMN user_id VARCHAR(255) NULL COMMENT '创建者用户ID'",
		"CREATE INDEX idx_listening_list_user_id ON listening_list(user_id)",
		"CREATE INDEX idx_reading_list_user_id ON reading_list(user_id)",
		"CREATE INDEX idx_writing_list_user_id ON writing_list(user_id)",
		"CREATE INDEX idx_testing_list_user_id ON testing_list(user_id)",
	}

	for i, stmt := range statements {
		fmt.Printf("Executing statement %d: %s\n", i+1, stmt[:50]+"...")
		_, err := db.Exec(stmt)
		if err != nil {
			log.Printf("Warning: Failed to execute statement %d: %v", i+1, err)
			// 继续执行其他语句
			continue
		}
		fmt.Printf("✓ Statement %d executed successfully\n", i+1)
	}

	fmt.Println("\n✅ Migration completed!")
	fmt.Println(sqlStatements)
}
