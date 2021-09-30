// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package database

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"time"

	"github.com/globalsign/mgo"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//struct to hold driver instance, db name and collection name
type MongoDBConnect struct {
	mongoSession *mgo.Session
	database     string
}

/*
 *Create object to establish session to MongoDB for pool of socket connections
 */
func InitializeConnection(addrs []string, timeout int, authDatabase string, username string, password string, workDatabase string) (conn MongoDBConnect, err error) {
	conn = MongoDBConnect{}

	//creating DialInfo object to establish a session to MongoDB
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    addrs,
		Timeout:  time.Duration(timeout) * time.Second,
		Database: authDatabase,
		Username: username,
		Password: password,
	}

	conn.database = workDatabase

	//Creating a session object which creates a pool of socket connections to MongoDB
	conn.mongoSession, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		LOGGER.Error("Error while creating MongoDB socket connections pool: ", err)
		return
	}

	LOGGER.Infof("\t* Dialing MongoDB session for socket connections pool for the Database: [%s] in Mongo Host: [%s]", conn.database, addrs)
	conn.mongoSession.SetMode(mgo.Eventual, true)

	LOGGER.Infof("\t* MongoDB socket connections pool for [%s] initialized successfully.", conn.database)
	return
}

/*
 * Request a socket connection from the session and retrieve collection to process query.
 */
func (mongoDBConnection MongoDBConnect) GetSocketConn() (session *mgo.Session) {
	session = mongoDBConnection.mongoSession.Copy()
	return
}

/*
 * Get the collection object from the session and collection name provided
 */
func (mongoDBConnection MongoDBConnect) GetCollection(session *mgo.Session, collectionName string) (collection *mgo.Collection) {
	collection = session.DB(mongoDBConnection.database).C(collectionName)
	return
}

/*
 * Close the session when the goroutine exits
 */
func (mongoDBConnect MongoDBConnect) CloseSession() {
	defer mongoDBConnect.mongoSession.Close()
}

/*
 * Connect to Mongo Atlas
 */

func InitializeAtlasConnection(username string, password string, id string) (*mongo.Client, error) {

	LOGGER.Infof("\t* Establishing Mongo DB connection...")
	urlEncodedPassword := url.QueryEscape(password)
	envVersion := os.Getenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION)
	ConnectionURI := "mongodb+srv://" + username + ":" + urlEncodedPassword + "@" + envVersion + "-" + id + ".mongodb.net/test?retryWrites=true"
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientAtlas, err := mongo.Connect(ctx, options.Client().ApplyURI(ConnectionURI))
	if err != nil {
		return &mongo.Client{}, err
	}

	// Check the connection
	err = clientAtlas.Ping(ctx, readpref.Primary())
	if err != nil {
		return &mongo.Client{}, err
	}

	LOGGER.Infof("\t* Mongo Atlas DB connection established!")
	return clientAtlas, nil
}

// used when return result will be array
func ParseResult(cursor *mongo.Cursor, ctx context.Context) ([]byte, error) {
	var interfaces []interface{}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var temp map[string]interface{}
		err := cursor.Decode(&temp)
		if err != nil {
			return []byte{}, errors.New("Encounter error while decoding cursor")
		}
		interfaces = append(interfaces, temp)
	}
	if err := cursor.Err(); err != nil {
		return []byte{}, errors.New("Encounter error while traversing through cursor")
	}

	bytes, err := json.Marshal(interfaces)
	if err != nil {
		return []byte{}, errors.New("Encounter error while marshaling mongo result")
	}

	return bytes, nil
}
