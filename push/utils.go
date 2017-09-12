package push

import (
	"bytes"
	"encoding/json"
	"errors"
	"html"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

//Request struct
type Request struct {
	AppID            string
	PlayerID         string
	Title            string
	Message          string
	ErrorCode        float64
	ErrorDescription string
}

//NewPush create class
func NewPush(playerID string, title string, message string) *Request {
	return &Request{
		PlayerID: playerID,
		Title:    title,
		Message:  message,
	}
}

func (r *Request) Send() (bool, error) {

	AppID := viper.GetString("services.push.oneSignal.appId")
	APIKey := viper.GetString("services.push.oneSignal.APIKey")
	sendURL := viper.GetString("services.push.oneSignal.send")

	//https://documentation.onesignal.com/reference#create-notification
	var jsonStr = []byte(`{
		"app_id": "` + AppID + `",
		"headings": ` + r.Title + `,
		"contents": ` + r.Message + `,
		"include_player_ids": ["` + r.PlayerID + `"]
	}`)

	req, err := http.NewRequest("POST", sendURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Basic "+APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result interface{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return false, err
	}

	resultList := result.(map[string]interface{})
	if resultList["errors"] != nil {
		return false, errors.New(resultList["errors"].(string))
	}

	return true, nil
}

//ParseTemplate func
//OneSignal требует отправлять объект с обязательно структурой
// {"en": "text", "..": ""}
func (r *Request) ParseTemplate(name string, tpl string, data interface{}) error {
	t, err := template.New(name).Parse(tpl)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.Message = html.UnescapeString(buf.String())
	return nil
}
