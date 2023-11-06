package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type ResponseBody struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/token", handler)
    http.HandleFunc("/verify", verify)
	log.Println("Server started and listening on http://localhost:3000")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Println("Error starting server", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := ResponseBody{
        Message: "Data",
    }

	token, err := generateJWT("Data from BE")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	res.Token = token

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func verify(w http.ResponseWriter, r *http.Request) {
    tokenString := r.Header.Get("Token")
    if tokenString == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("No token found in header"))
        return
    }

    token, err := verifyJWT(tokenString)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if ok {
        fmt.Printf("Claims: %v\n", claims)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Valid Token"))
}

func generateJWT(data string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = now.Add(10 * time.Minute).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["authorized"] = true
	claims["user"] = "username"
	claims["data"] = data

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return token, nil
	}

	return nil, fmt.Errorf("invalid token")
}
