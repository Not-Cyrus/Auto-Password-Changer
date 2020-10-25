package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	discordAPI      = "https://discordapp.com/api/v8/users/@me"
	oldToken        string
	oldPassword     string
	intervalOfSleep int
	passwordLength  int
	canRun          = true
)

type passwordInformation struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

type toSave struct {
	Email         string
	SavedPassword string `json:"NewPassword"`
	Token         string `json:"NewToken"`
}

func main() {
	fmt.Print("Enter your Password: ")
	fmt.Scan(&oldPassword)
	fmt.Print("Enter your Token: ")
	fmt.Scan(&oldToken)
	fmt.Print("Enter the amount of minutes it'll take: ")
	fmt.Scan(&intervalOfSleep)
	s := fmt.Sprintf("Okay, It will change every %d minutes.", intervalOfSleep)
	fmt.Println(s)
	fmt.Print("One last thing, Length of the password?: ")
	fmt.Scan(&passwordLength)
	for {
		if canRun {
			writePassword()
			time.Sleep(time.Duration(intervalOfSleep) * time.Minute)
		}
	}
}

func generatePassword() string {
	rand.Seed(time.Now().Unix())
	charSet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	var output strings.Builder
	length := passwordLength
	for i := 0; i < length; i++ {
		randomChar := charSet[rand.Intn(len(charSet))]
		output.WriteRune(randomChar)
	}
	newOutput := output.String()
	output.Reset()
	return newOutput
}

func changePassword() (string, string, string) {
	generatedPassword := generatePassword()
	passInfo := passwordInformation{oldPassword, generatedPassword}
	bytes := new(bytes.Buffer)
	json.NewEncoder(bytes).Encode(passInfo)
	Client := &http.Client{}
	req, _ := http.NewRequest("PATCH", discordAPI, bytes)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", oldToken)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.130 Safari/537.36")
	Res, _ := Client.Do(req)
	b, _ := ioutil.ReadAll(Res.Body)
	token, email := decode(b)
	if len(token) == 0 {
		fmt.Println("You sure you put in the correct password/Token combination? Whatever, here's a generated password: " + generatedPassword)
		canRun = false
	} else {
		return email, token, generatedPassword
	}
	return "", "", ""
}

func writePassword() {
	email, token, password := changePassword()
	if len(token) != 0 && len(password) != 0 {
		infoToSave, _ := encode(toSave{email, password, token})
		oldPassword = password
		oldToken = token
		ioutil.WriteFile("LoginInfo.json", infoToSave, os.ModePerm)
		fmt.Println("Check the LoginInfo.json file for your new login.")
	}
}

func encode(toEncode interface{}) ([]byte, error) {
	return json.Marshal(toEncode)
}

func decode(toDecode []byte) (string, string) {
	var a map[string]string
	json.Unmarshal([]byte(toDecode), &a)
	return a["token"], a["email"]
}
