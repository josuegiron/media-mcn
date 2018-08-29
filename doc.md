# Search songs

Returns information of the songs hosted on the indicated provider.
### Endpoint

  <http://localhost:3006/search/song>
     
### HTTP Method:

 	 GET
    
### Query String
 
| Parameter	 | Required | Valid Options | Description |
|---|---|---|---|
| q | Yes | Any words | Keyword for the search |
| provider | Yes | spotify, deezer | The content provider |

  
### Responses

#### HTTP Status Code

      200 OK

##### Schema

| Name | Type | Description |
|---|---|---|
| Fecha | date | Fecha en que se toma el tipo de cambio del dolar. |
| Referencia | double | Valor actual del tipo de cambio del dolar. |

##### Example

  ```json
  {
      "Fecha": "13/08/2018",
      "Referencia": "7.48156"
  }
  ```
 
#### HTTP Status Code

      404 NOT FOUND

##### Schema

| Name | Type | Description |
|---|---|---|
| status_message | string | Error message |
| status_code | int | HTTP Status Code |

##### Example
  ```json
  { 
    "status_messaje" : "No se pudo obtener el valor actual del dolar...",
    "status_code" : 404
  }
  ```

### Nota
Para realizar un test puede utilizar la herramienta [**Postman**](https://www.getpostman.com/).
