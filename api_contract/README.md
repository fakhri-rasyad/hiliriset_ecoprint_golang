# Ecoprint RESTAPIDocs

## Open Endpoints

Endpoint api yang tidak memerlukan authentikasi

- [Login](login.md) : `POST /api/login/`
- [Register](register.md) : `POST /api/register/`

## Endpoint yang memerlukan autentikasi pengguna

Endpoint yang memerlukan Token JWT yang valid pada header dari request.
Token dapat diperoleh setelah melakukan login.

### Endpoint user

Endpoint yang memanipulasi maupun menampilkan informasi yang berhubungan dengan
seorang atau beberapa user. Memerlukan bearer token pada headernya:

- [Menampilkan daftar user](user/get.md) : `GET /api/users/`
- [Menampilkan info user](user/uuid/get.md) : `GET /api/users/:uuid`
- [Menghapus user](user/uuid/delete.md) : `DELETE /api/users/:uuid/

### Endpoint kompor

Endpoint yang memanipulasi maupun menampilkan informasi yang berhubungan
dengan kompor. Memerlukan bearer token pada headernya:

- [Menampilkan daftar kompor](kompor/get.md) : `GET /api/kompors/`
- [Menambahkan kompor](kompor/post.md) : `POST /api/kompors/`
- [Menampilkan informasi kompor](kompors/uuid/get.md) `GET /api/kompors/:uuid`
- [Menghapus kompor](kompors/uuid/get.md) `DELETE /api/kompors/:uuid`

### Endpoint esps

Endpoint yang memanipulasi maupun menampilkan informasi yang berhubungan
dengan esps. Memerlukan bearer token pada headernya:

- [Menampilkan daftar esps](esp/get.md) : `GET /api/esps/`
- [Menambahkan esp](esp/post.md) : `POST /api/esps/`
- [Menampilkan informasi esp](esps/uuid/get.md) `GET /api/esps/:uuid`
- [Menghapus esp](esps/uuid/get.md) `DELETE /api/esps/:uuid`

### Endpoint sessions

Endpoint yang memanipulasi maupun menampilkan informasi yang berhubungan
dengan sessions. Memerlukan bearer token pada headernya:

- [Menampilkan daftar sessions](session/get.md) : `GET /api/sessions/`
- [Menambahkan session](session/post.md) : `POST /api/sessions/`
- [Menampilkan informasi session](sessions/uuid/get.md) `GET /api/sessions/:uuid`
- [Menghapus session](sessions/uuid/get.md) `DELETE /api/sessions/:uuid`

### Endpoint records

Endpoint yang memanipulasi maupun menampilkan informasi yang berhubungan
dengan records. Memerlukan bearer token pada headernya:

- [Menampilkan informasi record](records/uuid/get.md) `GET /api/records/:uuid`
- [Menghapus record](records/uuid/get.md) `DELETE /api/records/:uuid`
