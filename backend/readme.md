# 用户注册与登录接口说明

## 注册接口

### 请求格式
- **URL**: `/agent/register`
- **Method**: `POST`
- **Content-Type**: `application/json`

### 请求示例
```json 
{
   "name": "yjy01",
   "email": "example@qq.com",
   "password": "12345678" 
 }
 ```
 
### 字段说明
- `name`: 用户名，字符串类型，必填。
- `email`: 用户邮箱，字符串类型，必填。
- `password`: 用户密码，字符串类型，必填。

## 登录接口

### 请求格式
- **URL**: `/agent/login`
- **Method**: `POST`
- **Content-Type**: `application/json`

### 请求示例
```json 
{
  "name": "yjy01",
  "password": "12345678"
}
```
### 字段说明
- `name`: 用户名，字符串类型，必填。
- `password`: 用户密码，字符串类型，必填。

# 列出所有主机信息接口说明

## 接口描述
该接口用于查询当前用户的所有主机信息，支持按时间范围过滤。

## 请求格式
- **URL**: `/agent/list`
- **Method**: `GET`
- **Content-Type**: `application/json`

## 请求参数
| 参数名 | 类型   | 必填 | 说明                                                                 |
|--------|--------|------|--------------------------------------------------------------------|
| from   | string | 否   | 起始时间，格式为 `RFC3339`（如 `2023-01-01T00:00:00Z`），默认为 `1970-01-01T00:00:00Z` |
| to     | string | 否   | 结束时间，格式为 `RFC3339`（如 `2023-12-31T23:59:59Z`），默认为 `9999-12-31T23:59:59Z` |

## 响应格式
- **Content-Type**: `application/json`
- **响应示例**:
```json
[ 
  {
    "id": 2, 
    "hostname": "my-host", 
    "os": "Linux", 
    "platform": "Ubuntu 20.04", 
    "kernel_arch": "x86_64", 
    "host_info_created_at": "2025-03-10T18:17:16.189247Z", 
    "token": ""
  }, 
  { 
    "id": 3, 
    "hostname": "web-server", 
    "os": "Linux", 
    "platform": "CentOS 7", 
    "kernel_arch": "x86_64", 
    "host_info_created_at": "2024-05-15T12:34:56.789012Z", 
    "token": ""
  }, 
  {
    "id": 4, 
    "hostname": "db-server", 
    "os": "Linux", 
    "platform": "Debian 10", 
    "kernel_arch": "x86_64", 
    "host_info_created_at": "2024-08-20T09:45:23.456789Z", 
    "token": ""
  }, 
  { 
    "id": 5, 
    "hostname": "dev-machine", 
    "os": "Windows", 
    "platform": "Windows 10", 
    "kernel_arch": "x86_64", 
    "host_info_created_at": "2024-11-25T15:30:10.123456Z", 
    "token": ""
  }, 
  { 
    "id": 6, 
    "hostname": "test-server", 
    "os": "Linux", 
    "platform": "Fedora 33", 
    "kernel_arch": "x86_64", 
    "host_info_created_at": "2025-01-30T20:10:05.987654Z", 
    "token": ""
  } 
]
```
### 字段说明
| 字段名               | 类型   | 说明                     |
|----------------------|--------|------------------------|
| id                   | int    | 主机唯一标识             |
| hostname             | string | 主机名                   |
| os                   | string | 操作系统                 |
| platform             | string | 操作系统版本             |
| kernel_arch          | string | 内核架构                 |
| host_info_created_at | string | 主机信息创建时间         |
| token                | string | 主机令牌（暂未使用）     |

## 注意事项
1. 请确保在请求头中正确设置 `Content-Type` 为 `application/json`。
2. 时间参数 `from` 和 `to` 必须符合 `RFC3339` 格式。
3. 如果未提供时间参数，默认查询范围为 `1970-01-01T00:00:00Z` 到 `9999-12-31T23:59:59Z`。