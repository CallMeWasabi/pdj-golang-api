package db

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type FirestoreCommunicator struct {
	Client *firestore.Client
	Ctx    context.Context
}

var Provider = FirestoreCommunicator{}

func InitFirestoreCommunicator() {
	Provider.Ctx = context.Background()

	opt := option.WithCredentialsFile("firebase/key/serviceAccountKey.json")
	app, err := firebase.NewApp(Provider.Ctx, nil, opt)
	if err != nil {
		panic("error initializing app: " + err.Error())
	}

	Provider.Client, err = app.Firestore(Provider.Ctx)
	if err != nil {
		panic("error creating firebase client: " + err.Error())
	}
}

func CloseFirestoreCommunicator() {
	Provider.Client.Close()
	Provider.Ctx = nil
}
