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

### Get All Working Proxies
- **URL**: `/work-proxies/all`
- **Method**: `GET`
- **Response**: List of all working proxy servers

### Get Single Working Proxy
- **URL**: `/work-proxies/one`
- **Method**: `GET`
- **Response**: Returns a single working proxy server

### Get All Proxies
- **URL**: `/proxies/all`
- **Method**: `GET`
- **Response**: List of all proxy servers (working and non-working)

### Add New Proxy
- **URL**: `/proxies/add`
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

### Get All Sites
- **URL**: `/sites/all`
- **Method**: `GET`
- **Response**: List of all proxy source sites

### Add New Site
- **URL**: `/sites/add`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "url": "https://example.com/proxies"
  }
  ```
- **Response**: Added site URL

### Delete Proxy
- **URL**: `/proxies/delete`
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

### Delete Site
- **URL**: `/sites/delete`
- **Method**: `DELETE`
- **Body**:
  ```json
  {
    "url": "https://example.com/proxies"
  }
  ```
- **Response**: Success message

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