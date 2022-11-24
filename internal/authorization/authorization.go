package authorization

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

var mySigningKey = []byte("grhetjetfb")

type Funds struct {
	Id               int             `json:"id"`
	Name             string          `json:"name"`
	Ticker           string          `json:"ticker"`
	Amount           int64           `json:"amount"`
	PricePerItem     decimal.Decimal `json:"priceperitem"`
	PurchasePrice    decimal.Decimal `json:"purchaseprice"`
	PriceCurrent     decimal.Decimal `json:"pricecurrent"`
	PercentChanges   decimal.Decimal `json:"percentchanges"`
	YearlyInvestment decimal.Decimal `json:"yearlyinvestment"`
	ClearMoney       decimal.Decimal `json:"clearmoney"`
	DatePurchase     time.Time       `json:"datepurchase"`
	DateLastUpdate   time.Time       `json:"datelastupdate"`
	Type             string          `json:"type"`
}

type User struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var u User
	json.NewDecoder(r.Body).Decode(&u)

	if CheckLogin(u) == "403 Forbidden" {
		w.WriteHeader(403)
		name := r.URL.Query().Get("name")
		io.WriteString(w, fmt.Sprintf("You don't have permission to access / on this server. %s", name))
	}
}

func CheckLoginDb(email string, pwd string) bool {
	db, err := sql.Open("sqlite3", "./database/DB_Golang")
	checkError(err)
	rows, err := db.Query("SELECT email, password FROM users WHERE email = $1 AND password = $2", email, pwd)
	checkError(err)

	var user []User
	for rows.Next() {
		var email string
		var password string
		err = rows.Scan(&email, &password)
		user = append(user, User{Email: email, Pwd: password})

		checkError(err)

		if user != nil {
			return true
		} else {
			return false
		}

	}

	return false
}

func CheckLogin(u User) string {

	Pwd := u.Pwd + "Питер"
	pwd := base64.StdEncoding.EncodeToString([]byte(Pwd))
	if !CheckLoginDb(u.Email, pwd) {
		err := "403 Forbidden"
		return err
	}

	validToken, err := GenerateJWT()
	fmt.Println(validToken)

	if err != nil {
		fmt.Println(err)
	}

	return validToken
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "Connor Kenway"
	// claims["exp"] = time.Now().Add(time.Hour * 1000).Unix()
	claims["exp"] = time.Now().Add(1 * time.Minute)

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong: %s", err.Error())
	}

	return tokenString, nil
}

func CheckAuth(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Connection", "close")
		defer r.Body.Close()
		now := time.Now()
		var access bool

		if r.Header["Token"] != nil {
			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(r.Header["Token"][0], claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				timeToken := claims["exp"]
				strTime := fmt.Sprintf("%v", timeToken)
				time2, _ := time.Parse(time.RFC3339, strTime)
				access = (now.Before(time2))
				return mySigningKey, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				return
			}

			if !access {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, "Not Authorized")
			} else {
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}
