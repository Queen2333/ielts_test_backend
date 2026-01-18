# Database Migrations

## 运行 Migration

### 方式一：使用 MySQL 命令行

```bash
# 连接到数据库
mysql -h 172.25.138.133 -u ielts_alex -p ielts_database

# 执行 migration
source /Users/alexye/Documents/ielts_test_backend/migrations/001_add_user_id_to_part_lists.sql
```

### 方式二：使用 mysql 命令直接执行

```bash
mysql -h 172.25.138.133 -u ielts_alex -pYx180236 ielts_database < /Users/alexye/Documents/ielts_test_backend/migrations/001_add_user_id_to_part_lists.sql
```

## 回滚 Migration

如果需要回滚这个 migration：

```bash
mysql -h 172.25.138.133 -u ielts_alex -pYx180236 ielts_database < /Users/alexye/Documents/ielts_test_backend/migrations/001_add_user_id_to_part_lists_down.sql
```

## Migration 说明

### 001_add_user_id_to_part_lists.sql

为以下三个表添加 `user_id` 字段：
- `listening_part_list`
- `reading_part_list`
- `writing_part_list`

**字段说明：**
- 类型：VARCHAR(255)
- 允许为空：是
- 用途：记录创建该记录的用户ID

**索引：**
为每个表的 `user_id` 字段添加索引以优化查询性能。

**业务逻辑：**
当 `type=3` 时，查询接口会自动筛选出属于当前用户的数据。

## 已完成的代码修改

### 模型更新
- ✅ `models.ListeningPartItem` - 添加 `UserID` 字段
- ✅ `models.ReadingPartItem` - 添加 `UserID` 字段
- ✅ `models.WritingPartItem` - 添加 `UserID` 字段

### 控制器更新
- ✅ `controllers.AddListeningPart` - 创建时自动填充 `user_id`
- ✅ `controllers.AddReadingPart` - 创建时自动填充 `user_id`
- ✅ `controllers.AddWritingPart` - 创建时自动填充 `user_id`

### 查询逻辑
- ✅ `utils.ProcessRequest` - 当 `type=3` 时自动添加 `user_id` 到查询条件

## 注意事项

1. 执行 migration 前请确保已连接到正确的数据库
2. 建议先在测试环境执行，确认无误后再在生产环境执行
3. 执行前建议先备份相关表的数据
4. 现有数据的 `user_id` 将为 NULL，这是正常的
5. 新创建的数据将自动填充创建者的 `user_id`
