**Football-Data-BackEnd**
----
  <This API provides data of current football competitions as well as standings inside them_>

* **URL**

  </api/competitiors>

* **Method:**

  `GET` | 

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `{ data:competitions  }`
 
* **Error Response:**

  * **Code:** 422 UNPROCESSABLE ENTRY <br />
    **Content:** `{ error : "No such endpoint" }`

* **Sample Call:**

  <axious.get("http://localhost:8000/api/competitions")> 
  
  
  * **URL**

  </api/competitiors/{id}/standings>

* **Method:**

  `GET` 

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `{data: res.data.standings[0}`
 
* **Error Response:**

  * **Code:** 422 UNPROCESSABLE ENTRY <br />
    **Content:** `{ error : "No such endpoint" }`

* **Sample Call:**

  <axious.get("http://localhost:8000/api/competitions/{id}/standings")> 
  
* **Notes:**

  <For caching purposes it uses Redis. Download and lunch a server before using an app> 
