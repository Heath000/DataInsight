# SE2024
2024软工实践-时空数据管理系统
使用golang-gin框架 采用skeleton脚手架 restful风格

## 在postman中测试chart
### 饼图
1. Set Request Type and URL:  
* Request Type: `POST`
* URL: `http://localhost:8080/chart/pie`
2. Set Request Headers:  
* Add a `Content-Type` header with the value `application/json`
3. Set Request Body:  
* Select the `raw` option and set the type to `JSON`
Enter the JSON data in the request body, for example:
```json
{
  "title": "Sample Pie Chart",
  "data": [
    {"name": "Category A", "value": 30},
    {"name": "Category B", "value": 70}
  ]
}
```
4. Send Request:  
* Click the Send button to send the request

### 折线图
1. Set Request Type and URL:
* Request Type: `POST`
* URL: `http://localhost:8080/chart/line`
2. Set Request Headers:
* Add a `Content-Type` header with the value `application/json`
3. Set Request Body:
* Select the `raw` option and set the type to `JSON`
  Enter the JSON data in the request body, for example:
```json
{
  "title": "Sample Line Chart",
  "data": [
    {"x": 1, "y": 10},
    {"x": 2, "y": 20},
    {"x": 3, "y": 30}
  ]
}
```
4. Send Request:
* Click the Send button to send the request

### 柱状图
1. Set Request Type and URL:
* Request Type: `POST`
* URL: `http://localhost:8080/chart/bar`
2. Set Request Headers:
* Add a `Content-Type` header with the value `application/json`
3. Set Request Body:
* Select the `raw` option and set the type to `JSON`
  Enter the JSON data in the request body, for example:
```json
{
  "title": "Sample Bar Chart",
  "data": [
    {"label": "A", "value": 10},
    {"label": "B", "value": 20},
    {"label": "C", "value": 30}
  ]
}
```
4. Send Request:
* Click the Send button to send the request

### 折线柱状混合图
1. Set Request Type and URL:
* Request Type: `POST`
* URL: `http://localhost:8080/chart/linebarmixed`
2. Set Request Headers:
* Add a `Content-Type` header with the value `application/json`
3. Set Request Body:
* Select the `raw` option and set the type to `JSON`
  Enter the JSON data in the request body, for example:
```json
{
  "title": "Sample Mixed Chart",
  "line_data": [
    {"x": "A", "y": 10},
    {"x": "B", "y": 20},
    {"x": "C", "y": 30}
  ],
  "bar_data": [
    {"label": "A", "value": 10},
    {"label": "B", "value": 20},
    {"label": "C", "value": 30}
  ]
}
```
4. Send Request:
* Click the Send button to send the request
