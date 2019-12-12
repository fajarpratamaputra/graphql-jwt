package middleware

import (
    "time"
    "../libs"
    jwt "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo"
    "fmt"
)

type MyClaim struct {
    UserId int64
    IsAdmin bool
    RefreshJti string
    jwt.StandardClaims
}

func createRefreshTokenString(userid int64) (refreshTokenString string, err error) {
    refreshJti, err := libs.StoreRefreshToken()
    if err != nil {
        return "", err
    }

    if userid != 0 {
        // Create token
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyClaim{
            UserId: userid,
            IsAdmin: false,
            RefreshJti: refreshJti,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
            }})

        // Generate encoded token and send it as response.
        t, err := token.SignedString([]byte("secret2"))
        if err != nil {
            return "", err
        }
        return t, err
    }
    return "", echo.ErrUnauthorized
}

func createAuthTokenString(userid int64) (authTokenString string, err error) {
    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyClaim{
        UserId: userid,
        IsAdmin: true,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
        }})

    // Generate encoded token and send it as response.
    t, err := token.SignedString([]byte("secret1"))
    if err != nil {
        return "", err
    }
    return t, err
}

func CreateNewTokens(username string, password string) (authTokenString string, refreshTokenString string, err error) {
    user := libs.FetchUser(username, password)

    refreshTokenString, err = createRefreshTokenString(user.Id)

    if err != nil {
        return "", "", err
    }

    authTokenString, err = createAuthTokenString(user.Id)

    if err != nil {
        return "", "", err
    }

    return
}

func UpdateRefreshTokenExp(myClaim *MyClaim, oldTokenString string) (newTokenString, newRefreshTokenString string, err error) {
    myClaim2 := MyClaim{}
    _, err = jwt.ParseWithClaims(oldTokenString, &myClaim2, func(token *jwt.Token) (interface{}, error) {
        return []byte("secret1"), nil
    })

    if err != nil {
        return "", "", err
    }

    if !libs.CheckRefreshToken(myClaim.RefreshJti) || myClaim.UserId != myClaim2.UserId {
        return "", "", fmt.Errorf("error: %s", "old token is invalid")
    }

    libs.DeleteRefreshToken(myClaim2.RefreshJti)

    newRefreshTokenString, err = createRefreshTokenString(myClaim2.UserId)

    if err != nil {
        return "", "", err
    }

    newTokenString, err = createAuthTokenString(myClaim2.UserId)

    if err != nil {
        return "", "", err
    }

    return
}
