package database

import (
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

var (
	Cluster *gocb.Cluster
	Bucket  *gocb.Bucket
)

// Connect initializes and connects to the database
func Connect(url, username, password, bucketName string, timeout time.Duration) error {
	log.Println("Connecting to database...")

	// Create cluster options
	options := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	}

	// Create cluster connection
	cluster, err := gocb.Connect(url, options)
	if err != nil {
		return err
	}
	Cluster = cluster

	// Get bucket
	bucket := cluster.Bucket(bucketName)
	err = bucket.WaitUntilReady(timeout, nil)
	if err != nil {
		log.Printf("Failed to get bucket: %v", err)
		return err
	}
	Bucket = bucket

	log.Println("Database connected successfully")

	return nil
}

// Close closes the database connection
func Close() {
	log.Println("Closing database connection...")
	if Cluster != nil {
		Cluster.Close(nil)
	}
}
