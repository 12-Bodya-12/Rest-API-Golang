package authorization

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

var mySigningKey = []byte("grhetjetfb")

const (
	// Initialize connection constants.
	HOST     = "localhost"
	DATABASE = "postgres"
	USER     = "postgres"
	PASSWORD = "57kq10!!"
)

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

var user = User{
	Email: "1",
	Pwd:   "1",
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Login
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	checkLogin(u)
}

func checkLogin(u User) string {

	if user.Email != u.Email || user.Pwd != u.Pwd {
		fmt.Println("NOT CORRECT")
		err := "error"
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
	claims["user"] = "Elliot Forbes"
	claims["exp"] = time.Now().Add(time.Hour * 1000).Unix()

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

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				return
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func myRUBCurrentFunds(fundType string) []Funds {
	var amountShares []Funds
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE)

	db, err := sql.Open("postgres", connectionString)
	checkError(err)

	err = db.Ping()
	checkError(err)
	// fmt.Println("Successfully created connection to database")

	rows, err := db.Query("SELECT * FROM funds WHERE type = $1 ORDER BY ticker ASC", fundType)
	checkError(err)

	for rows.Next() {
		bk := Funds{}
		err = rows.Scan(&bk.Id, &bk.Name, &bk.Ticker, &bk.Amount, &bk.PricePerItem, &bk.PurchasePrice, &bk.PriceCurrent, &bk.PercentChanges, &bk.YearlyInvestment, &bk.ClearMoney, &bk.DatePurchase, &bk.DateLastUpdate, &bk.Type)
		checkError(err)

		amountShares = append(amountShares, bk)
	}

	defer rows.Close()

	return amountShares
}

func getRUBFundsShares(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var ArrShares = myRUBCurrentFunds("share")
	json.NewEncoder(w).Encode(ArrShares)
}

func Auth() {
	fmt.Println("GO")

	r := mux.NewRouter()
	//////////////////////////////////////////////////
	////////////////////// Login /////////////////////
	//////////////////////////////////////////////////

	r.HandleFunc("/login", login).Methods("POST")

	//////////////////////////////////////////////////
	//////////////////// GET /////////////////////////
	//////////////////////////////////////////////////

	r.Handle("/", CheckAuth(getRUBFundsShares)).Methods("GET")
}
