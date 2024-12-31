# proxy_parser_checker

A Go program that automatically parses and checks proxy servers from various sources.

## Build

```bash
make build
./proxy_parser_checker
```

### Static build 

```bash
make build_static
./proxy_parser_checker_static
```

## API Endpoints

All endpoints are prefixed with `/api/v1`

### Get All Working Proxies
- **URL**: `/proxies/working`
- **Method**: `GET`
- **Response**: List of all working proxy servers

### Get First Working Proxy
- **URL**: `/proxies/working/first`
- **Method**: `GET`
- **Response**: Returns a single working proxy server

### Get All Proxies
- **URL**: `/proxies`
- **Method**: `GET`
- **Response**: List of all proxy servers (working and non-working)

### Add New Proxy
- **URL**: `/proxies`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "ip": "192.168.1.1",
    "port": "8080",
    "protocol": "http"
  }
  ```
- **Response**: Added proxy details

### Delete Proxy
- **URL**: `/proxies`
- **Method**: `DELETE`
- **Body**:
  ```json
  {
    "ip": "192.168.1.1",
    "port": "8080",
    "protocol": "http"
  }
  ```
- **Response**: Success message

### Get All Sites
- **URL**: `/sites`
- **Method**: `GET`
- **Response**: List of all proxy source sites

### Add New Site
- **URL**: `/sites`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "url": "https://example.com/proxies"
  }
  ```
- **Response**: Added site URL

### Delete Site
- **URL**: `/sites`
- **Method**: `DELETE`
- **Body**:
  ```json
  {
    "url": "https://example.com/proxies"
  }
  ```
- **Response**: Success message

### Get Proxies Stats
- **URL**: `/stats`
- **Method**: `GET`
- **Response**: Returns statistics about proxies and checking process
  ```json
  {
    "success": true,
    "data": {
      "total_proxies": 1000,
      "worked_proxies": 400,
      "blocked_proxies": 50,
      "not_checked_proxies": 550,
      "check_rate": 2.5,
      "estimated_minutes": 220.0,
      "estimated_time": "3h 40m"
    }
  }
  ```

### Response Format
All endpoints return JSON responses in the following format:
```json
{
  "success": true|false,
  "data": <response_data>,
  "error": "error message if any"
}
```

![](https://asdertasd.site/counter/proxy_parser_checker?a=1)