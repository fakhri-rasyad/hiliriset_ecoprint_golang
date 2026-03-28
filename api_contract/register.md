# Register

Digunakan untuk mendaftarkan user.

**URL** : `/api/register/`

**Method** : `POST`

**Auth required** : NO

**Data constraints**

```json
{
  "username": "[nama pengguna]",
  "email" : "[alamat email pengguna dalam format yang valid]"
  "password": "[password pengguna]"
}
```

**Data example**

```json
{
  "username" : "example_username",
  "email" : "example@email.com"
  "password": "example_password"
}
```

## Success Response

**Code** : `200 OK`

**Content example**

```json
{
  "status": "Success",
  "response_code": 200,
  "message": "Registrasi Success"
}
```

## Error Response

**Condition** : Jika 'email' berada dalam format tidak valid atau telah digunakan.

**Code** : `400 BAD REQUEST`

**Content** :

```json
{
  "status": "Error bad request",
  "response_code": 400,
  "message": "Registrasi gagal, sepertinya email perlu diubah",
  "error": "Invalid credentials"
}
```
