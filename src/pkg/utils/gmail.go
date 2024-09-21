package utils

// gmail送信用のコード
// 事前にtoken.jsonを取得しておく必要がある
// https://zenn.dev/happy663/articles/36a21fa960a0f8

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GetClient - Gmailのクライアントを取得
func getClient(config *oauth2.Config, path_token string) *http.Client {
	// tokFile := "../../token.json"
	tok, err := tokenFromFile(path_token)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(path_token, tok)
	}
	return config.Client(context.Background(), tok)
}

// getTokenFromWeb - Webからトークンを取得
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// tokenFromFile - ファイルからトークンを取得
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken - トークンを保存
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func Send_Gmail(recipient string, subject string, body string, path_credentials string, path_token string) {
	// Gmailの送信
	ctx := context.Background()
	// credentials.jsonを読み込む
	b, err := os.ReadFile(path_credentials)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// Gmailの設定を取得
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	// Gmailのクライアントを取得
	client := getClient(config, path_token)

	// Gmailのサービスを取得
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	fmt.Println("Created Gmail service", srv)
	//追記
	msgStr := "From: 'me'\r\n" +
		"reply-to: hoge@gmail.com\r\n" + //送信元
		"To: " + recipient + "\r\n" + //送信先
		"Subject:" + subject + "\r\n" +
		"\r\n" + body
	// 文字化け対策
	reader := strings.NewReader(msgStr)
	transformer := japanese.ISO2022JP.NewEncoder()
	msgISO2022JP, err := ioutil.ReadAll(transform.NewReader(reader, transformer))
	if err != nil {
		log.Fatalf("Unable to convert to ISO2022JP: %v", err)
	}
	// メール送信
	msg := []byte(msgISO2022JP)
	message := gmail.Message{}
	message.Raw = base64.StdEncoding.EncodeToString(msg)
	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		fmt.Printf("%v", err)
	}
}
