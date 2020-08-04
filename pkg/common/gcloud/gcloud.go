package gcloud

import (
	"context"
	"io/ioutil"

	"golang.org/x/oauth2/google"

	cm "o.o/backend/pkg/common"
)

func LoadCredentialsFromFile(ctx context.Context, path string, scopes []string) (*google.Credentials, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, cm.Wrap(err)
	}
	return google.CredentialsFromJSON(ctx, data, scopes...)
}
