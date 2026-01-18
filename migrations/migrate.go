package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 数据库连接字符串
	dsn := "ielts_alex:Yx180236@tcp(172.25.138.133:3306)/ielts_database"

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database successfully!")

	// 执行 migration
	migrations := []string{
		// Add user_id column to listening_part_list table
		`ALTER TABLE listening_part_list ADD COLUMN user_id VARCHAR(255) NULL COMMENT '创建者用户ID'`,

		// Add user_id column to reading_part_list table
		`ALTER TABLE reading_part_list ADD COLUMN user_id VARCHAR(255) NULL COMMENT '创建者用户ID'`,

		// Add user_id column to writing_part_list table
		`ALTER TABLE writing_part_list ADD COLUMN user_id VARCHAR(255) NULL COMMENT '创建者用户ID'`,

		// Add indexes for better query performance
		`CREATE INDEX idx_listening_part_user_id ON listening_part_list(user_id)`,
		`CREATE INDEX idx_reading_part_user_id ON reading_part_list(user_id)`,
		`CREATE INDEX idx_writing_part_user_id ON writing_part_list(user_id)`,
	}

	// 执行每个 migration
	for i, migration := range migrations {
		fmt.Printf("Executing migration %d/%d...\n", i+1, len(migrations))
		_, err := db.Exec(migration)
		if err != nil {
			// 检查是否是字段已存在或索引已存在的错误
			if isAlreadyExistsError(err) {
				fmt.Printf("  ⚠️  Already exists, skipping: %v\n", err)
				continue
			}
			log.Fatalf("Migration %d failed: %v\nSQL: %s", i+1, err, migration)
		}
		fmt.Printf("  ✓ Migration %d completed successfully\n", i+1)
	}

	fmt.Println("\n✅ All migrations completed successfully!")
}

func isAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	// Check for MySQL duplicate column error (1060) or duplicate key error (1061)
	return contains(errMsg, "Duplicate column name") ||
		   contains(errMsg, "Duplicate key name") ||
		   contains(errMsg, "Error 1060") ||
		   contains(errMsg, "Error 1061")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
