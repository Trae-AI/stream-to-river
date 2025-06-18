package coze

import (
	"net/http"
	"sync"
	"time"

	"github.com/coze-dev/coze-go"
)

var (
	cozeClient coze.CozeAPI
	once       sync.Once
)

func GetCozeClient() coze.CozeAPI {
	once.Do(func() {
		var auth coze.Auth
		if CozeConf.Auth == "JWT" {
			client, e := coze.NewJWTOAuthClient(coze.NewJWTOAuthClientParam{
				ClientID:      CozeConf.ClientID,
				PublicKey:     CozeConf.PublishKey,
				PrivateKeyPEM: CozeConf.PrivateKey,
			}, coze.WithAuthBaseURL(CozeConf.BaseURL), coze.WithAuthHttpClient(&http.Client{Timeout: time.Second * 5}))
			if e != nil {
				panic(e)
			}

			auth = coze.NewJWTAuth(client, nil)
		} else {
			auth = coze.NewTokenAuth(CozeConf.Token)
		}

		hc := &http.Client{
			Timeout: time.Second * 30,
		}
		cozeClient = coze.NewCozeAPI(auth, coze.WithBaseURL(CozeConf.BaseURL), coze.WithHttpClient(hc))
	})
	return cozeClient
}
