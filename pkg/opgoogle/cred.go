package opgoogle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	_ "embed"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	//go:embed credentials.json
	credJSON []byte
	//go:embed token.json
	tokenJSON []byte
	//go:embed secret.json
	secretJSON []byte
)

var (
	ctx      = context.Background()
	docSrv   *docs.Service
	sheetSrv *sheets.Service
)

// GetDocsService 新建谷歌文档服务
func GetDocsService() (*docs.Service, error) {
	var err error

	if docSrv == nil {
		docSrv, err = newDocsService()
		if err != nil {
			return nil, err
		}
	}

	return docSrv, nil
}

// GetSheetService 新建谷歌报表服务
func GetSheetService() (*sheets.Service, error) {
	var err error

	if sheetSrv == nil {
		sheetSrv, err = newSheetServiceByOAuth()
		if err != nil {
			return nil, err
		}
	}

	return sheetSrv, nil
}

func newDocsService() (*docs.Service, error) {
	creds, err := google.CredentialsFromJSON(ctx, secretJSON, docs.DocumentsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("get credentials error:%w", err)
	}

	srv, err := docs.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("new docs service error:%w", err)
	}

	return srv, nil
}

func newSheetService() (*sheets.Service, error) {
	creds, err := google.CredentialsFromJSON(ctx, secretJSON, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("get credentials error:%w", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("new sheet service error:%w", err)
	}

	return srv, nil
}

func newSheetServiceByOAuth() (*sheets.Service, error) {
	config, err := google.ConfigFromJSON(credJSON, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %w", err)
	}

	tok := &oauth2.Token{}
	err = json.NewDecoder(bytes.NewReader(tokenJSON)).Decode(tok)
	if err != nil {
		return nil, fmt.Errorf("get token error:%w", err)
	}
	client := config.Client(context.Background(), tok)

	return sheets.NewService(ctx, option.WithHTTPClient(client))
}
