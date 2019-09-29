package database

import (
	"chapter-4/shared/config"
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
)

// Booking is the model for the booking object
// The JSON tags determine the mapping
type Booking struct {
	Name      string `json: name`
	Email     string `json: email`
	ID        string `json: id`
	Processed bool   `json: processed`
}

func (b *Booking) Exists() bool {
	return b.ID != ""
}

func NewClient() (*DB, error) {
	ctx := context.Background()
	firestore, err := firestore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, err
	}
	return &DB{firestore}, nil
}

type DB struct {
	firestore *firestore.Client
}

func (db *DB) SaveBooking(
	ctx context.Context,
	b *Booking) error {

	collection := db.firestore.Collection("bookings")

	var docRef *firestore.DocumentRef
	if b.Exists() {
		docRef = collection.Doc(b.ID)
	} else {
		docRef = collection.NewDoc()
	}
	_, err := docRef.Set(ctx, b)

	if err != nil {
		return fmt.Errorf("while saving b: %v", err)
	}
	b.ID = docRef.ID
	log.Printf("Saved b [%v]", b.ID)

	return nil
}

func (db *DB) GetBooking(ctx context.Context, ID string) (*Booking, error) {
	docRef, err := db.firestore.Collection("bookings").Doc(ID).Get(ctx)
	if docRef != nil && !docRef.Exists() {
		return &Booking{ID: ID}, nil
	}
	if err != nil {
		return nil, err
	}
	booking := Booking{}
	err = docRef.DataTo(&booking)
	if err != nil {
		return nil, err
	}
	booking.ID = ID
	return &booking, nil
}
