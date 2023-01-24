package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type User struct {
	Id       int                    `json:"user_id"`
	Calendar map[time.Time][]string `json:"calendar"`
}
type Data struct {
	UserId    int    `json:"user_id"`
	Date      string `json:"date"`
	Action    string `json:"action"`
	OldAction string `json:"old_action"`
}
type CustomError struct {
	Message string `json:"error"`
}
type CustomResponse struct {
	Result interface{} `json:"result"`
}

var (
	Users map[int]User
)

const layout = "2006-Jan-02"

func main() {
	Users = make(map[int]User)
	Route()

	log.Fatal(http.ListenAndServe(":8000", Middleware(http.DefaultServeMux.ServeHTTP)))
}
func Route() {
	http.HandleFunc("/create_event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			HttpError(w, fmt.Errorf("method not allowed"),
				"", "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()

		var data Data
		if errUnmarshBody := UnmarshalBody(r.Body, &data); errUnmarshBody != nil {
			HttpError(w, errUnmarshBody, "", "bad request", http.StatusBadRequest)
			return
		}

		if !ValidateUser(data.UserId) {
			Users[data.UserId] = User{
				Calendar: map[time.Time][]string{},
			}
		}

		user, errCreate := CreateAction(data)
		if errCreate != nil {
			HttpError(w, errCreate, errCreate.Error(), "bad request", http.StatusBadRequest)
			return
		}

		Users[user.Id] = user

		result := CustomResponse{Result: Users}
		resultJ, _ := json.Marshal(result)

		w.Header().Set("Content-Type", "application/json")
		w.Write(resultJ)
	})

	http.HandleFunc("/update_event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			HttpError(w, nil, "", "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()

		var data Data
		if errUnmarshBody := UnmarshalBody(r.Body, &data); errUnmarshBody != nil {
			HttpError(w, errUnmarshBody, "", "bad request", http.StatusBadRequest)
			return
		}
		t, errParseDate := data.parseDate()
		if errParseDate != nil {
			HttpError(w, errParseDate, errParseDate.Error(), "bad request", http.StatusBadRequest)
			return
		}

		if !ValidateUser(data.UserId) {
			HttpError(w, fmt.Errorf("wrong id"), "",
				"bad request", http.StatusBadRequest)
			return
		}

		user := Users[data.UserId]

		if _, ok := user.Calendar[t]; !ok {
			HttpError(w, nil, "", "bad request", http.StatusBadRequest)
			return
		}

		if errUpdate := UpdateAction(&user, data); errUpdate != nil {
			HttpError(w, nil, "", "bad request", http.StatusBadRequest)
			return
		}

		resp := CustomResponse{
			Result: user,
		}
		respJ, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.Write(respJ)
	})

	http.HandleFunc("/delete_event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			HttpError(w, nil, "", "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()

		var data Data
		if errUnmarshalBody := UnmarshalBody(r.Body, &data); errUnmarshalBody != nil {
			HttpError(w, nil, "", "bad request", http.StatusBadRequest)
			return
		}

		t, errParseDate := data.parseDate()
		if errParseDate != nil {
			HttpError(w, errParseDate, errParseDate.Error(), "bad request", http.StatusBadRequest)
			return
		}

		if !ValidateUser(data.UserId) {
			HttpError(w, nil, "", "user not found", http.StatusBadRequest)
			return
		}

		if _, ok := Users[data.UserId].Calendar[t]; !ok {
			HttpError(w, nil, "", "date not found", http.StatusBadRequest)
			return
		}

		user, errDel := DeleteAction(data)
		if errDel != nil {
			HttpError(w, nil, "", "date not found", http.StatusBadRequest)
			return
		}
		resp := CustomResponse{
			Result: user,
		}
		respJ, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.Write(respJ)
	})

	http.HandleFunc("/event_for_day", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			HttpError(w, nil, "", "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()

		usrID, date, errParseUrl := parseUrl(r.URL.Query())
		if errParseUrl != nil {
			HttpError(w, errParseUrl, errParseUrl.Error(), "bad request", http.StatusBadRequest)
			return
		}
		if !ValidateUser(usrID) {
			HttpError(w, fmt.Errorf("user not found"), "",
				"bad request", http.StatusBadRequest)
			return
		}
		user := Users[usrID]

		actions := user.GetActionsForDay(date)

		actionsJ, _ := json.Marshal(actions)

		w.Header().Set("Content-Type", "application/json")
		w.Write(actionsJ)
	})

	http.HandleFunc("/event_for_week", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			HttpError(w, fmt.Errorf("not allowed method"), "", "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()

		usrID, date, errParseUrl := parseUrl(r.URL.Query())
		if errParseUrl != nil {
			HttpError(w, errParseUrl, errParseUrl.Error(), "bad request", http.StatusBadRequest)
			return
		}
		if !ValidateUser(usrID) {
			HttpError(w, fmt.Errorf("user not found"), "",
				"bad request", http.StatusBadRequest)
			return
		}
		user := Users[usrID]
		actions := user.GetActionForWeek(date)

		resp := CustomResponse{
			Result: actions,
		}
		respJ, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.Write(respJ)
	})

	http.HandleFunc("/event_for_month", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			HttpError(w, fmt.Errorf("not allowed method"), "", "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()

		usrID, date, errParseUrl := parseUrl(r.URL.Query())
		if errParseUrl != nil {
			HttpError(w, errParseUrl, errParseUrl.Error(), "bad request", http.StatusBadRequest)
			return
		}
		if !ValidateUser(usrID) {
			HttpError(w, fmt.Errorf("user not found"), "",
				"bad request", http.StatusBadRequest)
			return
		}
		user := Users[usrID]
		actions := user.GetActionsForMonth(date)

		resp := CustomResponse{Result: actions}
		respJ, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.Write(respJ)
	})
}

func (u *User) GetActionsForMonth(date time.Time) map[int][]string {
	actions := make(map[int][]string)
	nextMonth := date.Month() + 1

	for date.Month() != nextMonth {
		v, ok := u.Calendar[date]
		if ok {
			actions[date.Day()] = v
		}

		date = date.Add(24 * time.Hour)
	}

	return actions
}

func (u *User) GetActionForWeek(date time.Time) map[time.Weekday][]string {
	actions := make(map[time.Weekday][]string)
	for date.Weekday() != time.Sunday {
		v, ok := u.Calendar[date]
		if !ok {
			date = date.Add(24 * time.Hour)
			continue
		}

		actions[date.Weekday()] = v
		date = date.Add(24 * time.Hour)

	}

	actions[date.Weekday()] = u.Calendar[date]
	return actions
}

func (u *User) GetActionsForDay(date time.Time) []string {
	return u.Calendar[date]
}

func CreateAction(data Data) (User, error) {
	t, errParse := time.Parse(layout, data.Date)
	if errParse != nil {
		return User{}, errParse
	}

	Users[data.UserId].Calendar[t] = append(Users[data.UserId].Calendar[t], data.Action)
	user := User{
		Id:       data.UserId,
		Calendar: Users[data.UserId].Calendar,
	}

	return user, nil
}

func UpdateAction(user *User, data Data) error {
	t, errParse := time.Parse(layout, data.Date)
	if errParse != nil {
		return errParse
	}

	for i, v := range user.Calendar[t] {
		if v == data.OldAction {
			user.Calendar[t][i] = data.Action
			return nil
		}
	}

	return fmt.Errorf("action not found")
}

func DeleteAction(data Data) (User, error) {
	t, errParse := time.Parse(layout, data.Date)
	if errParse != nil {
		return User{}, errParse
	}

	user := Users[data.UserId]
	for i, v := range user.Calendar[t] {
		if v == data.Action {
			user.Calendar[t] = append(user.Calendar[t][:i], user.Calendar[t][i+1:]...)
			break
		}
	}

	return user, nil
}

func ValidateUser(id int) bool {
	_, ok := Users[id]
	return ok
}

func (d *Data) parseDate() (time.Time, error) {
	t, errParse := time.Parse(layout, d.Date)
	if errParse != nil {
		return time.Time{}, errParse
	}

	return t, nil
}

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("From: %v, Method: %v", r.RemoteAddr, r.Method)
		next(w, r)
	}
}

func parseUrl(value url.Values) (int, time.Time, error) {
	usr_id := value.Get("user_id")
	dateStr := value.Get("date")

	userID, errParseID := strconv.Atoi(usr_id)
	if errParseID != nil {
		return -1, time.Time{}, errParseID
	}

	t, errParseDate := time.Parse(layout, dateStr)
	if errParseDate != nil {
		return -1, time.Time{}, errParseID
	}

	return userID, t, nil
}

func UnmarshalBody(r io.Reader, v interface{}) error {
	resp, errResp := ioutil.ReadAll(r)
	if errResp != nil {
		//ErrorPorcessing.HttpError(w, errResp, "failed to get body", "Bad Request", http.StatusBadRequest)
		return fmt.Errorf("server error: %w", errResp)
	}

	if errUnmarshalJson := json.Unmarshal(resp, &v); errUnmarshalJson != nil {
		//ErrorPorcessing.HttpError(w, errUnmarshalJson, "failed to get Json in Authorization", "Server Error", http.StatusInternalServerError)
		return fmt.Errorf("server error: %w", errUnmarshalJson)
	}

	return nil
}

func HttpError(w http.ResponseWriter, err error, msgForLogger string, msgForResponse string, code int) {
	w.Header().Set("Content-Type", "application/json")
	ce := CustomError{
		Message: msgForResponse,
	}

	res, errGetJson := json.Marshal(ce)
	if errGetJson != nil {
		fmt.Println(errGetJson)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	fmt.Println(msgForLogger + ": " + err.Error())
	w.WriteHeader(code)
	w.Write(res)
}
