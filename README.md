SaleStock Backend Dev Test
==========================

Introduction
------------
This is a submission for Salestock Backend Dev test built using Golang (Toko Ijah inventory).
The application is a webserver acting as endpoints for REST API.

Offline installation
====================
Prequisites:
------------
1. Golang
2. Sqlite3
3. dep (install from https://github.com/golang/dep)

Installation steps
------------------
1. Download from this repo, extract to your local machine
2. If required, change the port no of the web server by modifying "port" value in http config file (located at `repository/inventory/server/config/http/httpConfig.json`)
3. run `install.sh` script

Note: 
- The default port no of the http server is 8123.
- The "main" package file to run the http server is `main.go` located at `repository/inventory/server/http/main`
- The `install.sh` script restores a sqlite database (from file `ijahDump.sql`) to `/tmp`
- The http server needs access to sqlite database file `ijah.db` (location defaults to `/tmp`). the path to the file can be changed in config file (entry "filePath" under "database" in config file `repository/inventory/server/config/http/httpConfig.json`
- The http server logs access by writing to a file. It's possible to change the access log file location prior to running the http server. The config entry is "path" (under "combinedLog") is located at `repository/inventory/server/config/http/httpConfig.json`
- Errors during the http server execution are logged to a file and can also be modified. The config entry is "path" (under "combinedLog") is located at `repository/inventory/server/config/http/httpConfig.json`
 
Running Unit Test
-----------------
run command: `go test /path/to/ijah-inventory/repository/inventory/domain/inventory/service`

API Documentation
=================
List of services
----------------
The following services are provided:
1. **Get SKU Info** (for getting info about an item/SKU)
2. **Add SKU** (for adding new item type to stock)
3. **Update SKU** (for updating item/SKU properties, e.g. quantity, buying and selling price)
4. **Create Sale** (for creating a new sale)
5. **Update Sale Status** (for updating a sale status)
6. **Get All Stock Value** (for getting valuation of all SKU/items in stock)
7. **Get All Sales Value** (for getting valuation of all sales data)


API Format
----------
Note: the following assumes the http server is running on http://127.0.0.1:8123

### 1. Create Order

URL: `http://127.0.0.1/itemInfo?sku=<skuCode>`

METHOD: `HTTP GET`
Query string variables:
- sku: the sku of the item to view

Sample response:
```javascript
{
	"code": "S",
	"message": "Inquiry successful",
	"data": {
		"Sku": "SSI-D00864612-LL-NAV",
		"Name": "Deklia Plain Casual Blouse (L,Navy)",
		"Quantity": 85,
		"BuyPrice": 55000,
		"SellPrice": 60000
	}
}
````

### 2. Add SKU

URL: `http://127.0.0.1/addSKU`

METHOD: `HTTP POST`

Post Variables:
+ **orderId** : id of the order to add the product into. Use the order id returned from 'Create Order' service.
+ **itemId** : id of product to add. Refer to the 'product_id' in table 'm_product' in the database. Some values to try are: 'MBP01' and 'MBA01'.
+ **quantity** : quantity of item to add

Sample response:
```javascript
{
	"code": "S",
	"message": "Addition successful",
	"data": null
}
````

### 3. Update SKU

URL: `http://127.0.0.1/updateSKU`

METHOD: `HTTP POST`

Post Variables:
+ **sku** : the sku of the item to update.
+ **quantity** : item quantity.
+ **buyPrice** : item buying price
+ **sellPrice** : item selling price

Sample response:
```javascript
{
	"code": "S",
	"message": "Update successful",
	"data": null
}
````

### 4. Create Sale

URL: `http://127.0.0.1:8123/createSale`

METHOD: `HTTP POST`

Post Variables:
+ **invoiceNo** : the id of the sale invoice.
+ **note** : note of the sale.
+ **sku[x]** : sku of item in the sale.
+ **quantity[x]** : quantity of item in the sale.

Note: 
- replace 'x' with a number 
- Every sku[x] and quantity[x] with the same number is considered a pair, e.g. quantity[1] value is the quantity of item which sku is in sku[1] 

Sample HTTP request:
```
POST /createSale
content-type: application/x-www-form-urlencoded
user-agent: PostmanRuntime/7.1.1
accept: */*
host: localhost:8123
content-length: 139

invoiceId=invABC&note=invoice baru&sku[1]=SSI-D00864612-LL-NAV&quantity[1]=3&sku[2]=SSI-D01322234-LL-WHI&quantity[2]=5
````

Sample response:
```javascript
{
	"code": "S",
	"message": "Sale created successfully",
	"data": null
}
````

### 5. Update Sale Status

URL: `http://127.0.0.1:8123/updateSale`

METHOD: `HTTP POST`

Post Variables:
+ **invoiceNo** : the invoice id of the sale to update
+ **status** : the status of the sale


Sample response:
```javascript
{
	"code": "S",
	"message": "Update successful",
	"data": null
}
````

### 6. Get All Stock Value

URL: `http://127.0.0.1:8123/getStockValue`

METHOD: `HTTP GET`

Query string variables: None

Sample response:
```javascript
{
	"code": "S",
	"message": "Inquiry successful",
	"data": {
		"date": "2018-01-22T00:54:40.4121035+07:00",
		"totalQuantity": 600,
		"totalAmount": 40568000,
		"totalItemKind": 5,
		"items": {
			"SSI-D00791015-LL-BWH": {
				"sku": "SSI-D00791015-LL-BWH",
				"quantity": 154,
				"buyPrice": 62000,
				"totalAmount": 9548000
			},
			"SSI-D00864612-LL-NAV": {
				"sku": "SSI-D00864612-LL-NAV",
				"quantity": 85,
				"buyPrice": 55000,
				"totalAmount": 4675000
			},
			"SSI-D01037807-X3-BWH": {
				"sku": "SSI-D01037807-X3-BWH",
				"quantity": 74,
				"buyPrice": 85000,
				"totalAmount": 6290000
			},
			"SSI-D01220307-XL-SAL": {
				"sku": "SSI-D01220307-XL-SAL",
				"quantity": 182,
				"buyPrice": 75000,
				"totalAmount": 13650000
			},
			"SSI-D01322234-LL-WHI": {
				"sku": "SSI-D01322234-LL-WHI",
				"quantity": 105,
				"buyPrice": 61000,
				"totalAmount": 6405000
			}
		}
	}
}
````

### 7. Get All Sales Value

URL: `http://127.0.0.1:8123/getSalesValue`

METHOD: `HTTP GET`

Query String variables:
+ **startTime** : the start date of sales period to summarize (use format: YYYY-MM-DD, e.g. 2017-11-30)
+ **endTime** : the end date of sales peiod to summarize (use format: YYYY-MM-DD, e.g. 2017-12-31).

Sample response:
```javascript
{
	"code": "S",
	"message": "Inquiry successful",
	"data": {
		"startDate": "2016-12-31T00:00:00Z",
		"endDate": "2017-12-31T00:00:00Z",
		"totalQuantity": 54,
		"totalItemKind": 4,
		"saleCount": 3,
		"omzet": 4074200,
		"totalProfit": 308200,
		"items": [{
				"sku": "SSI-D00791015-LL-BWH",
				"quantity": 2,
				"buyPrice": 50000,
				"sellPrice": 60000,
				"profit": 20000
			}, {
				"sku": "SSI-D00864612-LL-NAV",
				"quantity": 5,
				"buyPrice": 68000,
				"sellPrice": 70000,
				"profit": 10000
			}, {
				"sku": "SSI-D01037807-X3-BWH",
				"quantity": 21,
				"buyPrice": 75000,
				"sellPrice": 80000,
				"profit": 105000
			}, {
				"sku": "SSI-D00791015-LL-BWH",
				"quantity": 7,
				"buyPrice": 52000,
				"sellPrice": 61000,
				"profit": 63000
			}, {
				"sku": "SSI-D01322234-LL-WHI",
				"quantity": 19,
				"buyPrice": 73000,
				"sellPrice": 78800,
				"profit": 110200
			}
		]
	}
}
````
