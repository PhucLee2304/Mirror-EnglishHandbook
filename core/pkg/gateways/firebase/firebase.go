package firebase

import (
	"context"
	"core/config"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Client struct {
	app       *firebase.App
	auth      *auth.Client
	messaging *messaging.Client
}

func NewClient(cfg *config.Config) (*Client, error) {
	ctx := context.Background()

	opt := option.WithCredentialsFile(cfg.FirebaseCredentialsPath)
	log.Println(cfg.FirebaseCredentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if app == nil {
		log.Println("app init failed")
	}
	if err != nil {
		log.Printf("error creating firebase app: %v\n", err)
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if authClient == nil {
		log.Println("auth init failed")
	}
	if err != nil {
		log.Printf("error getting auth client: %v\n", err)
		return nil, err
	}

	messagingClient, err := app.Messaging(ctx)
	if messagingClient == nil {
		log.Println("messaging init failed")
	}
	if err != nil {
		log.Printf("error getting messaging client: %v\n", err)
		return nil, err
	}

	return &Client{app, authClient, messagingClient}, nil
}

func (c *Client) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	log.Println("idtoken: " + idToken)
	token, err := c.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		if auth.IsIDTokenExpired(err) {
			log.Printf("id token expired: %v\n", err)
		} else if auth.IsIDTokenInvalid(err) {
			log.Printf("id token invalid: %v\n", err)
		} else if auth.IsIDTokenRevoked(err) {
			log.Printf("id token revoked: %v\n", err)
		} else {
			log.Printf("error verifying id token: %v\n", err)
		}
		return nil, err
	}
	log.Printf("token.Audience: %v\n", token.Audience)
	log.Printf("token.Issuer: %v\n", token.Issuer)

	return token, nil
}

func (c *Client) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := c.auth.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}
