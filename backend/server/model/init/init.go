package init

import (
	"fmt"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

)


const createTableSQL = `
-- ALTER TABLE hostandtoken
-- ADD COLUMN last_heartbeat TIMESTAMP DEFAULT NOW(),
-- ADD COLUMN status VARCHAR(10) DEFAULT 'offline';

-- user 表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    isverified BOOLEAN DEFAULT FALSE
);

-- host表
CREATE TABLE IF NOT EXISTS host_info (
	id SERIAL PRIMARY KEY,
	hostname TEXT  UNIQUE,
	os TEXT NOT NULL,
	platform TEXT NOT NULL,
	kernel_arch TEXT NOT NULL,
	host_info_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- cpu表
CREATE TABLE IF NOT EXISTS cpu_info (
	id SERIAL PRIMARY KEY,
	host_id INT REFERENCES host_info(id),
	model_name TEXT NOT NULL,
	cores_num INT NOT NULL,
	percent FLOAT NOT NULL,
	cpu_info_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- memory 表
CREATE TABLE IF NOT EXISTS memory_info (
	id SERIAL PRIMARY KEY,
	host_id INT REFERENCES host_info(id),
	total TEXT NOT NULL,
	available TEXT NOT NULL,
	used TEXT NOT NULL,
	free TEXT NOT NULL,
	user_percent TEXT NOT NULL,
	mem_info_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- process 表
CREATE TABLE IF NOT EXISTS process_info (
	id SERIAL PRIMARY KEY,
	host_id INT REFERENCES host_info(id),
	pid INT NOT NULL,
	cpu_percent FLOAT NOT NULL,
	mem_percent FLOAT NOT NULL,
	cmdline TEXT NOT NULL,
	pro_info_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- net_info表
CREATE TABLE IF NOT EXISTS network_info (
	id SERIAL PRIMARY KEY,
	host_id INT REFERENCES host_info(id),
	bytesrecv INT NOT NULL,
	bytessent INT NOT NULL,
	net_info_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- system_info表
CREATE TABLE IF NOT EXISTS system_info (
	id SERIAL PRIMARY KEY,
	cpu_info_id INT REFERENCES cpu_info(id),
	memory_info_id INT REFERENCES memory_info(id),
	host_info_id INT REFERENCES host_info(id),
	process_info_id INT REFERENCES process_info(id),
	system_info_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	network_info_id INT,
    FOREIGN KEY (network_info_id) REFERENCES network_info(id)
);

-- token表
CREATE TABLE IF NOT EXISTS hostandtoken (
	id SERIAL PRIMARY KEY,
	host_name TEXT NOT NULL,
	token TEXT NOT NULL,
	last_heartbeat TIMESTAMP DEFAULT NOW(),
	status VARCHAR(10) DEFAULT 'offline'
);`

var DB *gorm.DB

// ConnectDatabase 连接到数据库
func ConnectDatabase() error {
	var err error

	// 获取数据库连接信息
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	// 使用gorm打开数据库连接
	DB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err // 返回连接错误
	}
	return nil
}

// InitDB 初始化数据库，创建所需的表
func InitDB() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized") // 检查数据库连接是否已初始化
	}

	tx := DB.Begin() // 开始事务
	if tx.Error != nil {
		return tx.Error // 返回事务错误
	}

	if err := tx.Exec(createTableSQL).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err    // 返回创建表时的错误
	}

	if err := tx.Commit().Error; err != nil {
		return err // 返回提交事务时的错误
	}

	return nil
}