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

##### Struct JSON

| Name | Type | Description |
|---|---|---|
| name | string | Name of the song |
| artist | string | Artist name of the song |
| album | string | Album that contains the song |
| writer | string | songwriter |
| genre | string | Genre of the song |
| released | date | Year of recording the song |
| recorded | date | Record date of the song |
| length | time | Time in minutes of the duration of the song |
| track_number | int | Track number on the disc |

##### Example

  ```json
  {
      "name": "Oceans (Where Feet May Fail)",
      "artist": "Hillsong UNITED",
      "album": "Zion (Deluxe Edition)",
      "writer": "Joel Houston, Matt Crocker & Salomon Ligthelm",
      "genre": "Christian & Gospel",
      "released": "2013",
      "recorded": "2012",
      "length": "8:56",
      "track_number": 4
  }
  ```
 
#### HTTP Status Code

      400 NOT FOUND

##### Struct JSON

| Name | Type | Description |
|---|---|---|
| status_message | string | Error message |
| status_code | int | HTTP Status Code |

##### Example
  ```json
  { 
    "status_messaje" : "The required information could not be obtained...",
    "status_code" : 400
  }
  ```

### Nota
Para realizar un test puede utilizar la herramienta [**Postman**](https://www.getpostman.com/).
