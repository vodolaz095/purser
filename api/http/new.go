package purser_client

import (
	"context"
	"fmt"
	"net/http"
)

// New создаёт новый HTTP клиент для доступа к API
func New(apiURL, token string) (Client, error) {
	client := Client{
		Server: apiURL,
		Client: http.DefaultClient,
		RequestEditors: []RequestEditorFn{
			func(ctx context.Context, req *http.Request) error {
				req.Header.Add("Authorization", "Bearer "+token)
				req.Header.Add("User-Agent", "purser-http-client")
				return nil
			},
		},
	}
	resp, err := client.GetPing(context.Background())
	if err != nil {
		return Client{}, err
	}
	if resp.StatusCode == http.StatusNoContent {
		return client, nil
	}
	return Client{}, fmt.Errorf("wrong response code %s", resp.Status)
}
