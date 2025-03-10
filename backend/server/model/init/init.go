package init

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

const createTableSQL = `
-- roles 表
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
);

-- users 表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    isverified BOOLEAN DEFAULT FALSE,
    role_id INT ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- host表
CREATE TABLE IF NOT EXISTS host_info (
	id SERIAL PRIMARY KEY,
    user_name VARCHAR ,
	hostname VARCHAR(255)  UNIQUE,
	os TEXT NOT NULL,
	platform TEXT NOT NULL,
	kernel_arch TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- TIMESTAMP WITH TIME ZONE 加上时区
);

-- system_info表
CREATE TABLE IF NOT EXISTS system_info (
	id SERIAL PRIMARY KEY,
	host_info_id INT ,
	host_name VARCHAR(255) ,
	cpu_info JSONB,
	memory_info JSONB,
	process_info JSONB,
	network_info JSONB,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- token表
CREATE TABLE IF NOT EXISTS hostandtoken (
	id SERIAL PRIMARY KEY,
	host_name VARCHAR(255) ,
	token TEXT NOT NULL,
	last_heartbeat TIMESTAMP DEFAULT NOW(),
	status VARCHAR(10) DEFAULT 'offline'
);

-- 在system_info表的host_info_id字段上创建索引，加速通过主机ID查找系统信息
CREATE INDEX IF NOT EXISTS idx_system_info_host_info_id ON system_info(host_info_id);

-- 对于system_info表中的JSONB字段(cpu_info, memory_info等)，如果需要根据某些键值进行查询，
-- 可以考虑创建GIN (Generalized Inverted Index) 索引，例如：
-- 假设经常需要基于cpu_info内的某个键（如percent）来查询
CREATE INDEX IF NOT EXISTS idx_system_info_cpu_percent ON system_info USING GIN ((cpu_info->'percent') jsonb_path_ops);

-- 在hostandtoken表的host_name字段上创建索引，加速主机名查找
CREATE INDEX IF NOT EXISTS idx_hostandtoken_host_name ON hostandtoken(host_name);

-- 如果经常按last_heartbeat查询或排序，可以在此字段上创建索引
CREATE INDEX IF NOT EXISTS idx_hostandtoken_last_heartbeat ON hostandtoken(last_heartbeat);
`

// cpu_info示例，每次一新的数据就追加进json里面，这样可以保存多个时间戳的数据
// {
//     "cpu_info": [
//         {
//             "time": "2023-10-10T12:34:56Z",
//             "data": {
//                 "id": 1,
//                 "model_name": "Intel Xeon E5-2678 v3",
//                 "cores_num": 12,
//                 "percent": 45.7,
//                 "cpu_info_created_at": "2023-10-10T12:34:56Z",
//                 "updated_at": "2023-10-10T12:34:56Z"
//             }
//         },
//         {
//             "time": "2023-10-10T12:35:56Z",
//             "data": {
//                 "id": 1,
//                 "model_name": "Intel Xeon E5-2678 v3",
//                 "cores_num": 12,
//                 "percent": 50.2,
//                 "cpu_info_created_at": "2023-10-10T12:35:56Z",
//                 "updated_at": "2023-10-10T12:35:56Z"
//             }
//         }
//     ]
// }

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

// -- cpu表
// CREATE TABLE IF NOT EXISTS cpu_info (
// 	id SERIAL PRIMARY KEY,
// 	host_id INT REFERENCES host_info(id),
// 	model_name TEXT NOT NULL,
// 	cores_num INT NOT NULL,
// 	percent NUMERIC(5,2) NOT NULL,
// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

// -- memory 表
// CREATE TABLE IF NOT EXISTS memory_info (
// 	id SERIAL PRIMARY KEY,
// 	host_id INT REFERENCES host_info(id),
// 	total NUMERIC(10,2) NOT NULL,
// 	available NUMERIC(10,2) NOT NULL,
// 	used NUMERIC(10,2) NOT NULL,
// 	free NUMERIC(10,2) NOT NULL,
// 	user_percent NUMERIC(5,2) NOT NULL,
// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

// -- process 表
// CREATE TABLE IF NOT EXISTS process_info (
// 	id SERIAL PRIMARY KEY,
// 	host_id INT REFERENCES host_info(id),
// 	pid INT NOT NULL,
// 	cpu_percent NUMERIC(5,2) NOT NULL,
// 	mem_percent NUMERIC(5,2) NOT NULL,
// 	cmdline TEXT NOT NULL,
// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

// -- net_info表
// CREATE TABLE IF NOT EXISTS network_info (
// 	id SERIAL PRIMARY KEY,
// 	host_id INT REFERENCES host_info(id),
// 	bytesrecv BIGINT NOT NULL,
// 	bytessent BIGINT NOT NULL,
// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//);
