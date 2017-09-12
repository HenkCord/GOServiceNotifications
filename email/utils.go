package email

import (
	"bytes"
	"html"
	"html/template"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"

	"github.com/spf13/viper"
)

// Regular expression from WebCore's HTML5 email input: http://goo.gl/7SZbzj
var emailRegexp = regexp.MustCompile("(?i)" + // case insensitive
	"^[a-z0-9!#$%&'*+/=?^_`{|}~.-]+" + // local part
	"@" +
	"[a-z0-9-]+(\\.[a-z0-9-]+)*$") // domain part

// IsValidEmail returns true if the given string is a valid email address.
//
// It uses a simple regular expression to check the address validity.
func IsValidEmail(email string) bool {
	if len(email) > 254 {
		return false
	}
	return emailRegexp.MatchString(email)
}

// splitEmail splits email address into local and domain parts.
// The last returned value is false if splitting fails.
func splitEmail(email string) (local string, domain string, ok bool) {
	parts := strings.Split(email, "@")
	if len(parts) < 2 {
		return
	}
	local = parts[0]
	domain = parts[1]
	// Check that the parts contain enough characters.
	if len(local) < 1 {
		return
	}
	if len(domain) < len("x.xx") {
		return
	}
	return local, domain, true
}

// NormalizeEmail returns a normalized email address.
// It returns an empty string if the email is not valid.
func NormalizeEmail(email string) string {
	// Trim whitespace.
	email = strings.TrimSpace(email)
	// Make sure it is valid.
	if !IsValidEmail(email) {
		return ""
	}
	// Split email into parts.
	local, domain, ok := splitEmail(email)
	if !ok {
		return ""
	}
	// Remove trailing dot from domain.
	domain = strings.TrimRight(domain, ".")
	// Convert domain to lower case.
	domain = strings.ToLower(domain)
	// Combine and return the result.
	return local + "@" + domain
}

//Request struct
type Request struct {
	name    string
	from    string
	to      string
	subject string
	body    string
}

//NewEmail create class
func NewEmail(to string, subject string, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) Send() (bool, error) {
	host := viper.GetString("services.email.smtp.host")
	password := viper.GetString("services.email.smtp.password")
	username := viper.GetString("services.email.smtp.username")
	portS := viper.GetString("services.email.smtp.port")
	port, _ := strconv.Atoi(portS)

	m := gomail.NewMessage()
	m.SetHeader("From", username, "Bitkom24")
	m.SetHeader("To", r.to)
	//m.SetAddressHeader("Cc", r.from, "Dan")
	m.SetHeader("Subject", r.subject)
	m.SetBody("text/html", r.body)

	d := gomail.NewDialer(host, port, username, password)
	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return false, err
	}
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
	r.body = html.UnescapeString(buf.String())
	return nil
}
