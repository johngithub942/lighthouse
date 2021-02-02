#Search API

The search api is an interface to query a back searchable database and return
specifically formatted results. In addition it allows for caching. 

#### Format

`https://localhost:50005/api/search?s=3ds`

#### Result

```json
[
  {
    "claimId": "4de4c7bf6e00b99b50adbc63a22e09e5bae8d25b",
    "name": "ASPHALT-3D---GAMEPLAY---3DS-CITRA-EMULATOR"
  },
  {
    "claimId": "15e171df65d745eec53a0877febb3b99611db502",
    "name": "camera-setup-auto-vert-tilt-in-3ds-max"
  },
  {
    "claimId": "46d26d2abc48b19c6e1749457b03c98895daadb1",
    "name": "DEAD-OR-ALIVE-DIMENSIONS---GAMEPLAY---3DS-CITRA-EMULATOR---4K"
  },
  {
    "claimId": "2fc05e61ccdf667a75ec9750c2a5cef0738df56a",
    "name": "POOCHY-AND-YOSHI-WOOLY-WORLD---GAMEPLAY---3DS-CITRA-EMULATOR---4K"
  },
  {
    "claimId": "f3b0ad3ba1978db264788fbc930c26ea58d385df",
    "name": "3ds-max-octane-vs-arnold-skin-with-node"
  },
  {
    "claimId": "c13b7fb88863302fa6e25065ea83889bec4d42a3",
    "name": "nintendo-3ds-games-removed-without"
  },
  {
    "claimId": "3a9ccecdc1b7608551496cef4d81916fa8e0d026",
    "name": "3D-Classics--Urban-Champion-(3DS)-Playthrough---NintendoComplete"
  },
  {
    "claimId": "6466f4ce206186ec23ca9c256ce45c58977c8f90",
    "name": "PAXEast2012StreetPassNetwork"
  },
  {
    "claimId": "83c35305c826aa6193e59daec0731af2b7ca9802",
    "name": "mario-maker-for-nintendo-3ds-100-mario"
  },
  {
    "claimId": "338c637c7697a1d03ed6275af0e845b6067e4357",
    "name": "3ds-Max-繃布作法"
  }
]
```

##Parameters

### s [`string`] <span style="color:red">REQUIRED</span>

This is the text being searched via the API. It must be at least 3 characters long.

### size [`int`] Optional

This dictates the size of the results returned. It defaults to 20 if not passed.

### from [`int`]  Optional

This dictates the starting point of the results. The max results are 10,000, therefore
pagination's highest values are `size=20` and `from=9980`, or a similar combination.

### channel [`string`] Optional

This will ensure that the search result filters directly to claims by the specific
 channel mentioned by its name field.

### channel_id [`string`] Optional

This will ensure that the search result filters directly to claims by the specific
channel claim id passed. 

### related_to [`string`] Optional

When a claim id is passed, the search query is changed dynamically to find results
 that are related to the claim not just the text being searched. It broadens the scope 
 with specific changes to the overall search query. 
 
### sort_by [`string`] Optional

This allows for specific sorting of claim fields populated in the back end. For example
you could sort by `release_date`. In addition you can control the sort via `^` the
carrot such as `^release_date` which will put it in ascending sort. 

### include [`string`] Optional 

passing in a claim field in the document will ensure the results has the field included
in addition the base response above. 

call `http://localhost:50005/search?s=3ds&include=channel`

```json
{
  "channel": "@1Player",
  "claimId": "4de4c7bf6e00b99b50adbc63a22e09e5bae8d25b",
  "name": "ASPHALT-3D---GAMEPLAY---3DS-CITRA-EMULATOR"
}
```

### contentType [`string`] Optional

content types like `video/mp4` can be passed to ensure the results from the api
are only for content of that type. 

### mediaType [`string`] Optional
Allowed values: "audio", "video", "text", "application", "image"

This allows for more generic use of the content type field of a claim. So instead
of searching for `video/mp4` you can search all video. 

### claimTyoe [`string`] Optional
Allowed values: "file", "channel"

This allows filtering to channel claims or content claims in search results

### nsfw [`bool`] Optional

if passed it will determine the type of results you get. When true you will get
only NSFW results, if false, you will not get any. If not passed, it is not explicitly
filtered from the results. 

### free_only [`bool`] Optional

this will ensure the results you get are claims with no fee value recorded. 

### resolve [`bool`] Optional

This will add specific fields to the result, that would normally be provided from
a resolve call.

## language [`bool`] Optional

This will filter content results to the language passed. 

```json
{
  "channel": "@1Player",
  "channel_claim_id": "8698817e551b9db8020e8abf05006a3f4202f98b",
  "claimId": "4de4c7bf6e00b99b50adbc63a22e09e5bae8d25b",
  "duration": 224,
  "fee": 0,
  "name": "ASPHALT-3D---GAMEPLAY---3DS-CITRA-EMULATOR",
  "release_time": "2021-01-13T14:46:39Z",
  "thumbnail_url": "https://spee.ch/8/3599a04391347a87.png",
  "title": "ASPHALT 3D - GAMEPLAY - 3DS CITRA EMULATOR - 4K"
}
```

## Debug Paramters

These paramters are specific to implementation of the API and used for debug purposes
of that specific implementation. These are the parameters for this implementation. 

### claim_id [`string`] Optional

this will filter the results to this specific claim id. This is commonly used to
determine if a claim is "searchable" vs just not in the top n results. 

### source [`boolean`] Optional

passing this will return the entire claim document that is stored in the backend for 
review. See example, for this implementation below.

```json
{
    "bid_state": "Controlling",
    "cert_valid": true,
    "certificate_amount": 18600000000,
    "channel": "@1Player",
    "channel_claim_id": "8698817e551b9db8020e8abf05006a3f4202f98b",
    "claimId": "4de4c7bf6e00b99b50adbc63a22e09e5bae8d25b",
    "claim_cnt": 199,
    "claim_type": "stream",
    "content_type": "video/mp4",
    "description": "Testeando el juego Asphalt 3D en el emulador de 3DS Citra a 4K\n\nAsphalt 3D (lanzado en Japón como Asphalt 3D: Nitro Racing (アスファルト 3ディー: ニトロ レーシング?)) es un juego de carreras publicado por Konami en Japón y Ubisoft en todo el mundo y desarrollado por Gameloft para la consola portátil Nintendo 3DS. Fue lanzado en Japón el 10 de marzo de 2011, en Europa el 25 de marzo de 2011, en América del Norte el 27 de marzo de 2011 y en Australia el 31 de marzo de 2011. Es parte de Asphalt y fue uno de los ocho 3DS títulos de lanzamiento publicados por Ubisoft. Fue revelado en la Electronic Entertainment Expo 2010 (E3 2010).\n\nUna conversión directa de un juego Apple iOS, Asphalt 6: Adrenaline, Asphalt 3D incluye 17 pistas basadas en ubicaciones de la vida real y 42 vehículos deportivos con licencia. Cuenta con varios modos de juego que incluyen multijugador para hasta seis jugadores que utilizan el juego inalámbrico local. Asphalt 3D recibió críticas negativas, con gran parte de las críticas dirigidas a controles deficientes, numerosos errores, efectos visuales deficientes y una velocidad de fotogramas entrecortada. Ha recibido puntajes de compilación del 43% y 47% en Metacritic y GameRankings, respectivamente. ",
    "duration": 224,
    "effective_amount": 313769902,
    "fee": 0,
    "frame_height": 2160,
    "frame_width": 3840,
    "id": 9186084,
    "name": "ASPHALT-3D---GAMEPLAY---3DS-CITRA-EMULATOR",
    "release_time": "2021-01-13T14:46:39Z",
    "stripped_name": "ASPHALT3DGAMEPLAY3DSCITRAEMULATOR",
    "sub_cnt": 0,
    "tags": [
      "asphalt 3d",
      "citra",
      "espanol",
      "Gameplay",
      "videojuegos"
    ],
    "thumbnail_url": "https://spee.ch/8/3599a04391347a87.png",
    "title": "ASPHALT 3D - GAMEPLAY - 3DS CITRA EMULATOR - 4K",
    "transaction_time": "2021-01-14T10:16:37Z",
    "view_cnt": 12
}
``` 

### debug [`boolean`] Optional

This paramter is used to debug the choices made by the specific backend. In this case
elastic search is being used so it shows the entire calculation used for the top results
returned by the api. It is the elastic search explanation. 

