# Login

Digunakan untuk memperoleh token authentikasi dan authorisasi.

**URL** : `/api/login/`

**Method** : `POST`

**Auth required** : NO

**Data constraints**

```json
{
  "email": "[alamat email pengguna dalam format yang valid]",
  "password": "[password pengguna]"
}
```

**Data example**

```json
{
  "email": "example@email.com",
  "password": "examplepassword"
}
```

## Success Response

**Code** : `200 OK`

**Content example**

```json
{
  "status": "Success",
  "response_code": 200,
  "message": "Login Success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBhc3N3b3JkMSIsImV4cCI6MTc3NDY4MTQ2OSwicHViX2lkIjoiNDRjZGNkNWEtNjg4ZS00Y2FmLWI4MTUtOTZkZWVmODZlYzhhIiwicm9sZSI6InVzZXIiLCJ1c2VyX2lkIjo0fQ.RE5-MGE1GltnBlo5ODA29sQ2BsF-F8UjhrU-wFUIpWY",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzQ3NjQyNjksInVzZXJfaWQiOjR9.QlTZ6zLbqVoZ_VIWr4g2RbrRpdrryKbhe8EsKb2LqAw",
    "user": {
      "public_id": "44cdcd5a-688e-4caf-b815-96deef86ec8a",
      "name": "fakhri",
      "email": "password1",
      "role": "user",
      "created_at": "2026-03-02T13:34:51.507941+08:00",
      "updated_at": "2026-03-03T22:06:34.280412+08:00"
    }
  }
}
```

## Error Response

**Condition** : Jika 'username' atau 'password' salah.

**Code** : `400 BAD REQUEST`

**Content** :

```json
{
  "status": "Error bad request",
  "response_code": 400,
  "message": "Login gagal, sepertinya email atau passwordnya bermasalah",
  "error": "Invalid credentials"
}
```
