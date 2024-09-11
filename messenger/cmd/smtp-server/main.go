package smtp_server

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-smtp"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"
)

var dkimPrivateKey *rsa.PrivateKey

var dkimOptions = &dkim.SignOptions{
	Domain:   "example.com",
	Selector: "default",
	Signer:   dkimPrivateKey,
}

type Backend struct {
}

type Session struct {
	From string
	To   string
}

func init() {
	privateKeyPEM, err := ioutil.ReadFile("private.pem")
	if err != nil {
		log.Fatalf("failed to read private key: %v", err)
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		log.Fatal("failed to parse PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", err)
	}

	dkimPrivateKey = privateKey
}

func main() {
	s := smtp.NewServer(&Backend{})

	s.Addr = ":2525"
	s.Domain = "localhost"
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxRecipients = 50
	s.MaxMessageBytes = 1024 * 1024
	s.AllowInsecureAuth = true

	log.Println("Starting SMTP server on ", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (bck *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Println("Mail from:", from)
	s.From = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	s.To = append(s.To, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if data, err := io.ReadAll(r); err != nil {
		return err
	} else {
		fmt.Println("Message received: ", string(data))
		for _, recipient := range s.To {
			if err := sendMail(s.From, recipient, data); err != nil {
				fmt.Println("failed to send mail to %s: %v:", recipient, err)
			} else {
				fmt.Printf("Email sent successfully to %s", recipient)
			}
		}

		return nil
	}
}

func (s *Session) AuthPlain(username, password string) error {
	if username != "test" || password != "test123" {
		return errors.New("Invalid username or password")
	}
	return nil
}

func lookupMX(domain string) ([]*net.MX, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, fmt.Errorf("Error looking up MX records: %v", err)
	}

	return mxRecords, nil
}

func (s *Session) Logout() error {
	return nil
}

func sendMail(from string, to int32, msg []byte) error {
	domain := strings.Split(to, "@")[1]

	mxRecords, err := lookupMX(domain)
	if err != nil {
		return fmt.Errorf("Error looking up MX records: %v", err)
	}

	for _, mxRecord := range mxRecords {
		host := mxRecord.Host

		for _, port := range []int{25, 587, 465} {
			addr := fmt.Sprintf("%s:%d", host, port)

			var c *smtp.Client
			var err error

			switch port {
			case 465:
				tlsConfig := &tls.Config{ServerName: host}
				conn, err := tls.Dial("tcp", addr, tlsConfig)
				if err != nil {
					continue
				}

				c, err = smtp.NewClient(conn, host)
				if err != nil {
					log.Fatalf("Failed to create SMTP client: %v", err)
				}

			case 25, 587:
				c, err = smtp.Dial(addr)
				if err != nil {
					continue
				}

				if port == 587 {
					if err = c.StartTLS(&tls.Config{ServerName: host}); err != nil {
						c.Close()
						continue
					}
				}
			}

			var b bytes.Buffer
			if err := dkim.Sign(&b, bytes.NewReader{msg}, dkimOptions); err != nil {
				return fmt.Errorf("Error signing message with DKIM: %v", err)
			}
			signedData := b.Bytes()

			if err != nil {
				continue
			}

			if err := c.Mail(from); err != nil {
				c.Close()
				continue
			}

			if err = c.Rcpt(to); err != nil {
				c.Close()
				continue
			}

			w, err := c.Data()

			if err != nil {
				c.Close()
				continue
			}

			_, err = w.Write(signedData)
			if err != nil {
				c.Close()
				continue
			}

			err = w.Close()

			if err != nil {
				c.Close()
				continue
			}

			c.Quit()

			return nil
		}
	}
	return fmt.Errorf("failed to send mail to %s", to)
}
