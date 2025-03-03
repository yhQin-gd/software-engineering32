package model

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 连接数据库并创建表
func InitDB() (*sql.DB, error) {
	connStr := "host=192.168.31.251 port=5432 user=postgres password=cCyjKKMyweCer8f3 dbname=monitor sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// cpu表
	createCPUSQL := `
	ALTER TABLE hostandtoken
ADD COLUMN last_heartbeat TIMESTAMP DEFAULT NOW(),
ADD COLUMN status VARCHAR(10) DEFAULT 'offline';

	CREATE TABLE IF NOT EXISTS cpu_info (
		id SERIAL PRIMARY KEY,
		host_id INT REFERENCES host_info(id),
		model_name TEXT NOT NULL,
		cores_num INT NOT NULL,
		percent FLOAT NOT NULL,
		cpu_info_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
`
	// memory表
	createMEMSQL := `
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
	
`
	// host表
	createHOSTSQL := `
	CREATE TABLE IF NOT EXISTS host_info (
		id SERIAL PRIMARY KEY,
		hostname TEXT  UNIQUE,
		os TEXT NOT NULL,
		platform TEXT NOT NULL,
		kernel_arch TEXT NOT NULL,
		host_info_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	// process表
	createPROSQL := `
	CREATE TABLE IF NOT EXISTS process_info (
		id SERIAL PRIMARY KEY,
		host_id INT REFERENCES host_info(id),
		pid INT NOT NULL,
		cpu_percent FLOAT NOT NULL,
		mem_percent FLOAT NOT NULL,
		cmdline TEXT NOT NULL,
		pro_info_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
`
	// net_info表
	createNetSQL := `
	CREATE TABLE IF NOT EXISTS network_info (
		id SERIAL PRIMARY KEY,
		host_id INT REFERENCES host_info(id),
		bytesrecv INT NOT NULL,
		bytessent INT NOT NULL,
		net_info_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// system_info表
	createTableSQL := `
	ALTER TABLE system_info
ADD COLUMN network_info_id INT REFERENCES network_info(id);
	CREATE TABLE IF NOT EXISTS system_info (
		id SERIAL PRIMARY KEY,
		cpu_info_id INT REFERENCES cpu_info(id),
		memory_info_id INT REFERENCES memory_info(id),
		host_info_id INT REFERENCES host_info(id),
		process_info_id INT REFERENCES process_info(id),
		system_info_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	//token表
	createTokenSQL := `
	CREATE TABLE IF NOT EXISTS hostandtoken (
		id SERIAL PRIMARY KEY,
		host_name TEXT NOT NULL,
		token TEXT NOT NULL
	);`

	_, err = db.Exec(createCPUSQL)
	if err != nil {
		fmt.Printf("failed to create cpu_info table: %v", err)
		return nil, err
	}
	_, err = db.Exec(createMEMSQL)
	if err != nil {
		fmt.Printf("failed to create memory_info table: %v", err)
		return nil, err
	}
	_, err = db.Exec(createHOSTSQL)
	if err != nil {
		fmt.Printf("failed to create host_info table: %v", err)
		return nil, err
	}
	_, err = db.Exec(createPROSQL)
	if err != nil {
		fmt.Printf("failed to create process_info table: %v", err)
		return nil, err
	}
	_, err = db.Exec(createNetSQL)
	if err != nil {
		fmt.Printf("failed to create net_info table: %v", err)
		return nil, err
	}
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Printf("failed to create system_info table: %v", err)
		return nil, err
	}
	_, err = db.Exec(createTokenSQL)
	if err != nil {
		fmt.Printf("failed to create hostandtoken table: %v", err)
		return nil, err
	}
	return db, nil
}

type HostInfo struct {
	ID         int       `json:"id"` // 添加 ID 字段
	Hostname   string    `json:"hostname"`
	OS         string    `json:"os"`
	Platform   string    `json:"platform"`
	KernelArch string    `json:"kernel_arch"`
	CreatedAt  time.Time `json:"host_info_created_at"` // 添加 CreatedAt 字段
	Token      string    `json:"token"`
}

func InsertHostInfo(db *sql.DB, hostInfo HostInfo, username string) (int, string, error) {
	var hostInfoID int
	var hostname string
	var exists bool

	// 检查主机记录是否存在
	querySQL := `
    SELECT id, hostname, EXISTS (SELECT 1 FROM host_info WHERE hostname = $1 AND os = $2 AND platform = $3 AND kernel_arch = $4)
    FROM host_info WHERE hostname = $1 AND os = $2 AND platform = $3 AND kernel_arch = $4`

	err := db.QueryRow(querySQL, hostInfo.Hostname, hostInfo.OS, hostInfo.Platform, hostInfo.KernelArch).Scan(&hostInfoID, &hostname, &exists)
	if err == sql.ErrNoRows {
		fmt.Println("No matching host info found.")
		exists = false
	} else if err != nil {
		fmt.Printf("Failed to query host info: %v\n", err)
		return 0, "", err
	}

	if exists {
		// 更新已存在的主机记录
		updateSQL := `
        UPDATE host_info
        SET host_info_created_at = CURRENT_TIMESTAMP
        WHERE id = $1`
		_, err = db.Exec(updateSQL, hostInfoID)
		if err != nil {
			fmt.Printf("Failed to update host_info_created_at: %v\n", err)
			return 0, "", err
		}
		fmt.Printf("Updated existing host_info with ID: %d\n", hostInfoID)
	} else {
		// 插入新的主机记录
		insertSQL := `
        INSERT INTO host_info (hostname, os, platform, kernel_arch, host_info_created_at, user_name)
        VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, $5)
        RETURNING id, hostname`
		err = db.QueryRow(insertSQL, hostInfo.Hostname, hostInfo.OS, hostInfo.Platform, hostInfo.KernelArch, username).Scan(&hostInfoID, &hostname)
		if err != nil {
			fmt.Printf("Failed to insert host_info: %v\n", err)
			return 0, "", err
		}
		fmt.Printf("Inserted new host_info with ID and Name: %d and %v\n", hostInfoID, hostname)
	}

	return hostInfoID, hostname, nil
}

type CPUInfo struct {
	ID        int       `json:"id"` // 添加 ID 字段
	ModelName string    `json:"model_name"`
	CoresNum  int       `json:"cores_num"`
	Percent   float64   `json:"percent"`
	CreatedAt time.Time `json:"cpu_info_created_at"` // 添加 CreatedAt 字段
}

func InsertCpuInfo(db *sql.DB, cpuInfo CPUInfo, hostInfoID int, hostname string) (int, error) {
	cpuSQL := ` 
    INSERT INTO cpu_info (model_name, cores_num, percent, host_id, hostname, cpu_info_created_at)
    VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
    RETURNING id`
	var cpuInfoID int
	err := db.QueryRow(cpuSQL, cpuInfo.ModelName, cpuInfo.CoresNum, cpuInfo.Percent, hostInfoID, hostname).Scan(&cpuInfoID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert CPU info: %v", err)
	}
	fmt.Printf("Inserted cpu_info with ID: %d\n", cpuInfoID)
	return cpuInfoID, nil
}

type ProcessInfo struct {
	ID         int       `json:"id"` // 添加 ID 字段
	PID        int       `json:"pid"`
	CPUPercent float64   `json:"cpu_percent"`
	MemPercent float64   `json:"mem_percent"`
	Cmdline    string    `json:"cmdline"`
	CreatedAt  time.Time `json:"pro_info_created_at"` // 添加 CreatedAt 字段
}

func InsertProcessInfo(db *sql.DB, processInfo ProcessInfo, hostInfoID int, hostname string) (int, error) {
	processSQL := ` 
    INSERT INTO process_info (pid, cpu_percent, mem_percent, cmdline, host_id, hostname, pro_info_created_at)
    VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
    RETURNING id`
	var processInfoID int
	err := db.QueryRow(processSQL, processInfo.PID, processInfo.CPUPercent, processInfo.MemPercent, processInfo.Cmdline, hostInfoID, hostname).Scan(&processInfoID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert process info: %v", err)
	}
	fmt.Printf("Inserted process_info with ID: %d\n", processInfoID)
	return processInfoID, nil
}

type MemoryInfo struct {
	ID          int       `json:"id"` // 添加 ID 字段
	Total       string    `json:"total"`
	Available   string    `json:"available"`
	Used        string    `json:"used"`
	Free        string    `json:"free"`
	UserPercent float64   `json:"user_percent"`
	CreatedAt   time.Time `json:"mem_info_created_at"` // 添加 CreatedAt 字段
}

func InsertMemoryInfo(db *sql.DB, memoryInfo MemoryInfo, hostInfoID int, hostname string) (int, error) {
	memSQL := `
    INSERT INTO memory_info (total, available, used, free, user_percent, host_id, hostname, mem_info_created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
    RETURNING id`
	var memoryInfoID int
	err := db.QueryRow(memSQL, memoryInfo.Total, memoryInfo.Available, memoryInfo.Used, memoryInfo.Free, memoryInfo.UserPercent, hostInfoID, hostname).Scan(&memoryInfoID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert memory info: %v", err)
	}
	fmt.Printf("Inserted memory_info with ID: %d\n", memoryInfoID)
	return memoryInfoID, nil
}

// 定义网络信息结构体
type NetworkInfo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	BytesRecv uint64    `json:"bytes_recv"` // 接收字节数
	BytesSent uint64    `json:"bytes_sent"` // 发送字节数
	CreatedAt time.Time `json:"net_info_created_at"`
}

func InsertNetworkInfo(db *sql.DB, networkInfo NetworkInfo, hostInfoID int, hostname string) (int, error) {
	netSQL := `
    INSERT INTO network_info (name, bytesrecv, bytessent, net_info_created_at)
    VALUES ($1, $2,$3, CURRENT_TIMESTAMP)
    RETURNING id`
	var netInfoID int
	err := db.QueryRow(netSQL, networkInfo.Name, networkInfo.BytesRecv, networkInfo.BytesSent, hostInfoID, hostname).Scan(&netInfoID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert memory info: %v", err)
	}
	log.Printf("Inserted net_info with ID: %d\n", netInfoID)
	return netInfoID, nil
}
func InsertSystemInfo(db *sql.DB, cpuInfo []CPUInfo, memoryInfo MemoryInfo, hostInfo HostInfo, processInfo []ProcessInfo, networkInfo NetworkInfo, username string) error {
	// 插入或更新主机信息
	hostInfoID, hostname, err := InsertHostInfo(db, hostInfo, username)
	if err != nil {
		return err
	}

	// 插入内存信息
	memoryInfoID, err := InsertMemoryInfo(db, memoryInfo, hostInfoID, hostname)
	if err != nil {
		return err
	}

	// 插入 CPU 信息
	var cpuInfoIDs []int
	for _, cpu := range cpuInfo {
		id, err := InsertCpuInfo(db, cpu, hostInfoID, hostname)
		if err != nil {
			return err
		}
		cpuInfoIDs = append(cpuInfoIDs, id)
	}
	// 插入网卡信息
	netInfoID, err := InsertNetworkInfo(db, networkInfo, hostInfoID, hostname)
	if err != nil {
		return err
	}
	// 插入进程信息
	var processInfoIDs []int
	for _, proc := range processInfo {
		id, err := InsertProcessInfo(db, proc, hostInfoID, hostname)
		if err != nil {
			return err
		}
		processInfoIDs = append(processInfoIDs, id)
	}

	// 插入系统信息，使用最近插入的 CPU、内存、主机和进程信息的 ID
	systemInfoSQL := `
    INSERT INTO system_info (cpu_info_id, memory_info_id, host_info_id, process_info_id,network_info_id, system_info_created_at)
    VALUES ($1, $2, $3, $4,$5, CURRENT_TIMESTAMP)`
	_, err = db.Exec(systemInfoSQL, cpuInfoIDs[0], memoryInfoID, hostInfoID, processInfoIDs[0], netInfoID)
	if err != nil {
		fmt.Printf("Failed to insert system info: %v", err)
		return err
	}
	fmt.Println("Inserted system_info successfully")

	return nil
}

func InsertHostandToken(db *sql.DB, UserName string, Token string) error {

	// 插入新的记录
	fmt.Println("Inserting new host")
	insertSQL := `
	INSERT INTO hostandtoken (host_name, token)
	VALUES ($1, $2) RETURNING token`
	var token string
	err := db.QueryRow(insertSQL, UserName, Token).Scan(&token)
	if err != nil {
		log.Fatalf("Failed to query host info: %v\n", err)
		return err
	}
	log.Println("Insert successfully")

	return nil
}
func ReadDB(db *sql.DB, queryType, from, to string, hostname string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 查询主机信息
	if queryType == "host" || queryType == "all" {
		row := db.QueryRow("SELECT id, hostname, os, platform, kernel_arch, host_info_created_at FROM host_info WHERE hostname = $1", hostname)
		var id int
		var os, platform, kernelArch string
		var createdAt time.Time
		err := row.Scan(&id, &hostname, &os, &platform, &kernelArch, &createdAt)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("未找到指定的主机记录")
			}
			return nil, fmt.Errorf("查询主机信息时发生错误: %v", err)
		}
		result["host"] = map[string]interface{}{
			"id":                   id,
			"hostname":             hostname,
			"os":                   os,
			"platform":             platform,
			"kernel_arch":          kernelArch,
			"host_info_created_at": createdAt,
		}
	}

	// 查询内存信息
	if queryType == "memory" || queryType == "all" {
		row := db.QueryRow("SELECT id, total, available, used, free, user_percent, mem_info_created_at FROM memory_info WHERE hostname = $1 AND mem_info_created_at BETWEEN $2 AND $3", hostname, from, to)
		var id int
		var total, available, used, free, userPercent string
		var createdAt time.Time
		err := row.Scan(&id, &total, &available, &used, &free, &userPercent, &createdAt)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("未找到指定的内存记录")
			}
			return nil, fmt.Errorf("查询内存信息时发生错误: %v", err)
		}
		result["memory"] = map[string]interface{}{
			"id":                  id,
			"total":               total,
			"available":           available,
			"used":                used,
			"free":                free,
			"user_percent":        userPercent,
			"hostname":            hostname,
			"mem_info_created_at": createdAt,
		}
	}
	// 查询网卡信息
	if queryType == "net" || queryType == "all" {
		row := db.QueryRow("SELECT id, name, bytesrecv, bytessent, net_info_created_at FROM memory_info WHERE hostname = $1 AND net_info_created_at BETWEEN $2 AND $3", hostname, from, to)
		var id int
		var name, bytesrecv, bytes_sent string
		var createdAt time.Time
		err := row.Scan(&id, &name, &bytesrecv, &bytes_sent, &createdAt)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("未找到指定的内存记录")
			}
			return nil, fmt.Errorf("查询内存信息时发生错误: %v", err)
		}
		result["net"] = map[string]interface{}{
			"id":                  id,
			"name":                name,
			"bytesrecv":           bytes_sent,
			"bytessent":           bytes_sent,
			"net_info_created_at": createdAt,
		}
	}
	// 查询 CPU 信息
	if queryType == "cpu" || queryType == "all" {
		rows, err := db.Query("SELECT id, model_name, cores_num, percent, cpu_info_created_at FROM cpu_info WHERE hostname = $1 AND cpu_info_created_at BETWEEN $2 AND $3", hostname, from, to)
		if err != nil {
			return nil, fmt.Errorf("查询 CPU 信息时发生错误: %v", err)
		}
		defer rows.Close()

		var cpuInfos []map[string]interface{}
		for rows.Next() {
			var id int
			var modelName string
			var coresNum int
			var percent float64
			var createdAt time.Time
			err := rows.Scan(&id, &modelName, &coresNum, &percent, &createdAt)
			if err != nil {
				return nil, fmt.Errorf("扫描 CPU 信息记录时发生错误: %v", err)
			}
			cpuInfos = append(cpuInfos, map[string]interface{}{
				"id":                  id,
				"model_name":          modelName,
				"cores_num":           coresNum,
				"percent":             percent,
				"hostname":            hostname,
				"cpu_info_created_at": createdAt,
			})
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("处理 CPU 信息记录时发生错误: %v", err)
		}
		result["cpu"] = cpuInfos
	}

	// 查询进程信息
	if queryType == "process" || queryType == "all" {
		rows, err := db.Query("SELECT id, pid, cpu_percent, mem_percent, cmdline, pro_info_created_at FROM process_info WHERE hostname = $1 AND pro_info_created_at BETWEEN $2 AND $3", hostname, from, to)
		if err != nil {
			return nil, fmt.Errorf("查询进程信息时发生错误: %v", err)
		}
		defer rows.Close()

		var processInfos []map[string]interface{}
		for rows.Next() {
			var id, pid int
			var cpuPercent, memPercent float64
			var cmdline string
			var createdAt time.Time
			err := rows.Scan(&id, &pid, &cpuPercent, &memPercent, &cmdline, &createdAt)
			if err != nil {
				return nil, fmt.Errorf("扫描进程信息记录时发生错误: %v", err)
			}
			processInfos = append(processInfos, map[string]interface{}{
				"id":                  id,
				"pid":                 pid,
				"cpu_percent":         cpuPercent,
				"mem_percent":         memPercent,
				"cmdline":             cmdline,
				"hostname":            hostname,
				"pro_info_created_at": createdAt,
			})
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("处理进程信息记录时发生错误: %v", err)
		}
		result["process"] = processInfos
	}

	return result, nil
}

func UpdateDB(db *sql.DB, host_id int, new_cpu_info []map[string]string, new_memory_info map[string]string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 更新CPU信息
	for _, cpu_info := range new_cpu_info {
		_, err = tx.Exec(
			"UPDATE cpu_info SET model_name = $1, cores_num = $2, percent = $3, updated_at = $4 WHERE host_id = $5",
			cpu_info["ModelName"], cpu_info["CoresNum"], cpu_info["Percent"], time.Now(), host_id,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 更新内存信息
	_, err = tx.Exec(
		"UPDATE memory_info SET total = $1, available = $2, used = $3, free = $4, user_percent = $5, updated_at = $6 WHERE host_id = $7",
		new_memory_info["Total"], new_memory_info["Available"], new_memory_info["Used"], new_memory_info["Free"], new_memory_info["UserPercent"], time.Now(), host_id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return err
	}

	fmt.Printf("Updated CPU and Memory info for host_id: %d\n", host_id)
	return nil
}

func DeleteDB(db *sql.DB, host_id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 删除CPU信息
	_, err = tx.Exec("DELETE FROM cpu_info WHERE host_id = $1", host_id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除内存信息
	_, err = tx.Exec("DELETE FROM memory_info WHERE host_id = $1", host_id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除进程信息
	_, err = tx.Exec("DELETE FROM process_info WHERE host_id = $1", host_id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除主机信息
	_, err = tx.Exec("DELETE FROM host_info WHERE host_id = $1", host_id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
