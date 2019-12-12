package libs

import (
    "../db"
    "../models"
    "crypto/rand"
    "encoding/base64"
)

var refreshTokens map[string]string

func GenerateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return nil, err
    }

    return b, nil
}

func GenerateRandomString(s int) (string, error) {
    b, err := GenerateRandomBytes(s)
    return base64.URLEncoding.EncodeToString(b), err
}

func InitDB() {
    refreshTokens = make(map[string]string)
}

func FetchUser(username string, password string) models.User {
    con := db.ConnectGORM()
    con.SingularTable(true)
    user := models.User{}
    con.First(&user, "name=? and password=?", username, password)
    return user
}

func StoreRefreshToken() (jti string, err error) {
    jti, err = GenerateRandomString(32)
    if err != nil {
        return jti, err
    }

    for refreshTokens[jti] != "" {
        jti, err = GenerateRandomString(32)
        if err != nil {
            return jti, err
        }
    }

    refreshTokens[jti] = "valid"

    return jti, err
}

func DeleteRefreshToken(jti string) {
    delete(refreshTokens, jti)
}

func CheckRefreshToken(jti string) bool {
    return refreshTokens[jti] != ""
}
