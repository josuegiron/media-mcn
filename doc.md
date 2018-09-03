# Search songs

Returns information of the songs hosted on the indicated provider.
### Endpoint

  <http://localhost:3006/search/song/{provider}/{q}>
     
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
| songs | Array | Contains an arrangement of songs |
| songs.title | string | Name of the song |
| songs.artist | string | Artist name of the song |
| songs.album | string | Album that contains the song |
| songs.rank | string | Position of the song in the provider |
| songs.genre | string | Genre of the song |
| songs.released | fieldDataType | Date of recording the song |
| songs.explicit | bool | Record date of the song |
| songs.length_ms | fieldDataType | Time in milliseconds of the duration of the song |
| songs.track_number | int | Track number on the disc |
| songs.url | string | URL of the track file |
| next_page | string | Following results |
| previous_page | string | Previous results |
| total | int | Total results found |

##### Example

  ```json
  {
    "songs": [
        {
            "title": "Oceans (Where Feet May Fail) [Live at Red Rocks] - Live at Red Rocks",
            "artist": "Hillsong United",
            "album": "Oceans EP",
            "rank": 53,
            "genre": "",
            "released": "2013-10-10",
            "explicit": false,
            "length_ms": 589200,
            "track_number": 4,
            "url": "https://open.spotify.com/track/7IDiC6vVdVWMIwQPpf4SWi"
        }, 
        ...
        ],
    "next_page": "https://api.spotify.com/v1/...",
    "previous_page": "https://api.spotify.com/v1/...",
    "total": 31
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

#### HTTP Status Code

      200 OK

##### Struct JSON

| Name | Type | Description |
|---|---|---|
| status_message | string | The song was not found |
| status_code | int | HTTP Status Code |

##### Example
  ```json
  { 
    "status_messaje" : "The song was not found...",
    "status_code" : 200
  }
  ```
### Note
To perform a test you can use the [**Postman**](https://www.getpostman.com/) tool.
