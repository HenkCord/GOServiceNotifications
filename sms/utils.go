package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"html"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
	"github.com/ttacon/libphonenumber"
)

//PhoneNormalize change phone number have format +7 (999) 999-99-99, ...
//in +79999999999
func PhoneNormalize(phone string) (err error) {
	num, err := libphonenumber.Parse(phone, "RU")
	if err != nil {
		return err
	}
	phone = libphonenumber.Format(num, libphonenumber.E164)
	return
}

//Request struct
type Request struct {
	Phone     string
	Message   string
	SMSID     float64
	ErrorCode float64
}

//NewSMS create class
func NewSMS(phone string) *Request {
	return &Request{
		Phone: phone,
	}
}

func (r *Request) Send() (bool, error) {

	sendURL := viper.GetString("services.sms.send")
	charset := viper.GetString("services.sms.charset")
	format := viper.GetString("services.sms.fmt")

	resp, err := http.Get(sendURL +
		"&phones=" + r.Phone +
		"&mes=" + r.Message +
		"&charset=" + charset +
		"&fmt=" + format)
	if err != nil {
		return false, err
	}
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	resp.Body.Close()

	//Special for http://smsc.ru/api/
	var result interface{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return false, err
	}
	//Errors http://smsc.ru/api/
	resultList := result.(map[string]interface{})
	if resultList["error_code"] != nil {
		r.ErrorCode = resultList["error_code"].(float64)
		if r.ErrorCode == 7 || r.ErrorCode == 8 {
			return true, errors.New(resultList["error"].(string))
		}
		return false, errors.New(resultList["error"].(string))
	}

	if resultList["id"] == nil {
		return true, errors.New("invalid_response_id")
	}

	r.SMSID = resultList["id"].(float64)
	return true, nil
}

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
