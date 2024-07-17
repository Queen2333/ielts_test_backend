package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitializeDB 初始化数据库连接池
func InitializeDB(connectionString string) error {
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// GetDB 返回数据库连接
func GetDB() *sql.DB {
	return db
}

// IsNoRowsError 检查错误是否是 sql.ErrNoRows
func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

// PaginationQuery 封装分页查询列表的方法
func PaginationQuery(tableName string, pageNo, pageLimit int, conditions map[string]interface{}) ([]map[string]interface{}, int, error) {
	// 构建查询条件
	var args []interface{}
	var conditionsStr string
	for key, value := range conditions {
		conditionsStr += fmt.Sprintf(" AND %s = ?", key)
		args = append(args, value)
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	if conditionsStr != "" {
		countQuery += " WHERE 1 = 1" + conditionsStr
	}
	var total int
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	if conditionsStr != "" {
		query += " WHERE 1 = 1" + conditionsStr
	}
	if pageLimit > 0 {
		offset := (pageNo - 1) * pageLimit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageLimit, offset)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 解析查询结果
	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, 0, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			rawValue := values[i]
			if rawValue != nil {
				entry[col] = convertType(rawValue)
			} else {
				entry[col] = nil
			}
		}
		results = append(results, entry)
	}

	return results, total, nil
}

// convertType converts database types to appropriate Go types
func convertType(value interface{}) interface{} {
	switch v := value.(type) {
	case int64:
		// 将int64类型转换为int类型
		return int(v)
	case int32:
		return int(v)
	case []byte:
		str := string(v)
		if isJSONArray(str) {
            // 如果是 JSON 数组，则解析为 []interface{} 类型
            var jsonArray []interface{}
			if err := json.Unmarshal(v, &jsonArray); err == nil {
				for i, val := range jsonArray {
					jsonArray[i] = convertType(val)
				}
				return jsonArray
			}
        } else if isJSONObject(str) {
			// 如果是JSON对象，则解析为interface{}类型
			var jsonObject map[string]interface{}
			if err := json.Unmarshal(v, &jsonObject); err == nil {
				for key, val := range jsonObject {
					jsonObject[key] = convertType(val)
				}
				return jsonObject
			}
		}
		return str
	case nil:
		return nil
	default:
		return v
	}
}


func isJSONArray(s string) bool {
	// JSON 数组应该以 [ 开头，以 ] 结尾
    if len(s) >= 2 && s[0] == '[' && s[len(s)-1] == ']' {
        var js []interface{}
        return json.Unmarshal([]byte(s), &js) == nil
    }
    return false
}

func isJSONObject(s string) bool {
	// JSON 对象应该以 { 开头，以 } 结尾
    if len(s) >= 2 && s[0] == '{' && s[len(s)-1] == '}' {
        var js map[string]interface{}
        return json.Unmarshal([]byte(s), &js) == nil
    }
    return false
}

// GetListeningPartsByIds 根据ID查询表中的数据
func GetPartsByIds(tableName string, ids []int) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id IN (?", tableName) + strings.Repeat(",?", len(ids)-1) + ")"
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			rawValue := values[i]
			if rawValue != nil {
				entry[col] = convertType(rawValue)
			} else {
				entry[col] = nil
			}
		}
		results = append(results, entry)
	}

	return results, nil
}

// InsertData 插入数据到指定的表
func InsertData(tableName string, data interface{}, handleType string) (int, error) {
	v := reflect.ValueOf(data).Elem()
    t := v.Type()

    var columns []string
    var placeholders []string
    var values []interface{}
    var idValue interface{}
    idIndex := -1

    // Collect columns, placeholders, and values
    for i := 0; i < v.NumField(); i++ {
        field := t.Field(i)
        column := field.Tag.Get("json")
        if column == "" {
            column = strings.ToLower(field.Name)
        } else {
            column = strings.Split(column, ",")[0] // Handle omitempty and other tags
        }
        if column == "-" { // Skip fields with `json:"-"`
            continue
        }

        columns = append(columns, column)
        placeholders = append(placeholders, "?")
        value := v.Field(i).Interface()
        if isStructOrSlice(v.Field(i)) {
            jsonValue, _ := json.Marshal(value)
            values = append(values, string(jsonValue))
        } else {
            values = append(values, value)
        }

        // Check if the field is the ID field
        if field.Name == "ID" {
            idIndex = len(columns) - 1
            idValue = value
        }
    }

    var query string
	var isInsert bool
	var newID int

	if handleType == "update" { // 修改
		if idIndex == -1 || isEmptyID(idValue) {
			return 0, fmt.Errorf("wrong data: %v", "id cannot be empty!")
		}
		exists, err := checkIDExists(tableName, idValue)
        if err != nil {
            return 0, fmt.Errorf("failed to check if ID exists: %w", err)
        }

		if !exists {
			return 0, fmt.Errorf("wrong data: %v", "wrong id!")
		}

		query = generateUpdateQuery(tableName, columns, idIndex)
		values = append(values[:idIndex], values[idIndex+1:]...)
		values = append(values, idValue)
		// fmt.Println(query, values, "query")
	} else { // 新增
		newID = generateRandomID()
		values[idIndex] = newID
		query = generateInsertQuery(tableName, columns, placeholders)
		isInsert = true
	}

    // Execute the query
    _, err := db.Exec(query, values...)
    if err != nil {
        return 0, fmt.Errorf("failed to execute query: %v, values: %v, error: %w", query, values, err)
    }

	if isInsert {
        return newID, nil // Return generated ID for insert operation
    }

    return 0, nil
}

func generateUpdateQuery(tableName string, columns []string, idValue interface{}) string {
	var setClauses []string
	for _, column := range columns {
		if column == "id" { // Skip ID in update statement
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s=?", column))
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=?", tableName, strings.Join(setClauses, ","))
	return query
}

func generateInsertQuery(tableName string, columns []string, placeholders []string) string {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ","), strings.Join(placeholders, ","))
	return query
}

// Check if the ID value is empty
func isEmptyID(idValue interface{}) bool {
    v := reflect.ValueOf(idValue)
    switch v.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return v.Int() == 0
    case reflect.String:
        return v.String() == ""
    default:
        return false
    }
}

// Generate a random ID (example implementation)
func generateRandomID() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1000000) // You can customize this range
}

// Check if the ID exists in the table
func checkIDExists(tableName string, idValue interface{}) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=?)", tableName)
	err := db.QueryRow(query, idValue).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func isStructOrSlice(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Struct, reflect.Slice:
		return true
	default:
		return false
	}
}