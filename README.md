## Bonus counting service «Gophermart»

---

### Main

Service is HTTP API with next business requirements:

* registration and user authorization;
* accepts of order numbers from authorized users;
* accounting and maintaining of order numbers list from registered user;
* accounting of a registered user;
* verification of order numbers using the loyalty counting system;
* accrual reward to the user's loyalty account for each suitable order number.

![image](https://pictures.s3.yandex.net:443/resources/gophermart2x_1634502166.png)

Bonus counting service «Gophermart»
Service is HTTP API with next business requirements:

  - registration and user authorization;

  - accepts of order numbers from authorized users;

  - accounting and maintaining of order numbers list from registered user;

  - accounting of a registered user;

  - verification of order numbers using the loyalty counting system;

  - accrual reward to the user's loyalty account for each suitable order number.

## API

### Sing up

`POST /api/user/register`

	POST /api/user/register HTTP/1.1
	Content-Type: application/json

	{
	"login": "<login>",
	"password": "<password>"
	}

- `200` — User registered;

- `400` — Wrong request format;

- `409` — Login already used;

- `500` — Internal error.

---

### Login

`POST /api/user/login`

	POST /api/user/login HTTP/1.1
	Content-Type: application/json

	{
	"login": "<login>",
	"password": "<password>"
	}

- `200` — User authorized successfully;

- `400` — Wrong request format;

- `401` — Wrong login or password;

- `500` — Internal error.

---

### User upload new order number

`POST /api/user/orders`

	POST /api/user/login HTTP/1.1
	Content-Type: application/json

	12345678903

- `200` — Order number already uploaded before;

- `202` — New order namber accepted;

- `400` — Wrong request format;

- `401` — Unathorized;

- `409` — Order number already uploaded by other user;

- `422` — Wrong format of order number;

- `500` — Internal error.

---

### User gets list of orders with status and balance

`GET /api/user/orders`

	GET /api/user/orders HTTP/1.1
	Content-Length: 0

- `200` — List of orders in JSON

 200 OK HTTP/1.1
 Content-Type: application/json
 [
	 {
		   "number": "9278923470",
	    	"status": "PROCESSED",
	    	"accrual": 500,
	    	"uploaded_at": "2020-12-10T15:15:45+03:00"
	    },
	    {
	    	"number": "12345678903",
		   "status": "PROCESSING",
	    	"uploaded_at": "2020-12-10T15:12:01+03:00"
	    },
	    {
	    	"number": "346436439",
	    	"status": "INVALID",
	    	"uploaded_at": "2020-12-09T16:09:53+03:00"
	    }
 ]

- `204` — No orders;

- `401` — Unathorized;

- `500` — Internal error.

---

### User gets current balance

`GET /api/user/balance`

	GET /api/user/balance HTTP/1.1
	Content-Length: 0

- `200` — Balance

		 200 OK HTTP/1.1
		 Content-Type: application/json

	    {
	    	"current": 500.5,
	    	"withdrawn": 42
	    }

- `401` — Unathorized;

- `500` — Internal error.

---

### User request to pay new order using bonus account

`POST /api/user/balance/withdraw`

	 POST /api/user/balance/withdraw HTTP/1.1
	 Content-Type: application/json

	 {
		"order": "2377225624",
	    "sum": 751
	 }

- `200` — Success;

- `401` — Unathorized;

- `402` - Not enough money;

- `422` - Wrong order number;

- `500` — Internal error.

---

### User gets list of withdrawals approved

`GET /api/user/withdrawals`

	 GET /api/user/withdrawals HTTP/1.1
	 Content-Length: 0

- `200` — Success

 200 OK HTTP/1.1
 Content-Type: application/json

 [
	    {
	        "order": "2377225624",
	        "sum": 500,
	        "processed_at": "2020-12-09T16:09:57+03:00"
	    }
 ]

- `204` - Not withdrawals yet;

- `401` — Unathorized;

- `500` — Internal error.

## Order statuses

  - `REGISTERED` — registeres;

  - `INVALID` — no bonus for this order number;

  - `PROCESSING` — calculating in progress;

  - `PROCESSED` — order accepted, balance updated;

## Revards calculation system

Bonus counting service makes order numbers validation using third party service as client. Final stages are INVAID or PROCESSED.

### Settings

Database connection link (required): environment variable 'DATABASE_URI' or flag '-d'

Third party service: environment variable 'ACCRUAL_SYSTEM_ADDRESS' or flag '-r'

Service port (default:'localhost:8080'): environment variable 'RUN_ADDRESS' or flag '-a'
