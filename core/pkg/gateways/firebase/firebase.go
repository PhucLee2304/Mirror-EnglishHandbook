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
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("error creating firebase app: %v\n", err)
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Printf("error getting auth client: %v\n", err)
		return nil, err
	}

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("error getting messaging client: %v\n", err)
		return nil, err
	}

	return &Client{app, authClient, messagingClient}, nil
}

func (c *Client) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := c.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (c *Client) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := c.auth.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}
