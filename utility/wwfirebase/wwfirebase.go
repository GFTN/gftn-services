// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package wwfirebase

// [START authenticate_db_imports]
import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"firebase.google.com/go/auth"

	"context"

	"github.com/gorilla/mux"

	"google.golang.org/api/option"

	"firebase.google.com/go/db"
	logging "github.com/op/go-logging"

	firebase "firebase.google.com/go"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

var LOGGER = logging.MustGetLogger("wwfirebase")

var (
	AppContext   context.Context = context.Background()
	FbRef        *db.Ref
	FbClient     *db.Client
	FbAuthClient *auth.Client
)

// [END authenticate_db_imports]

// AuthenticateForAuthService:  Returns dbclient and authclient.
func AuthenticateWithAdminPrivileges() (*db.Client, *auth.Client, error) {

	// TODO: per Chase currently using 'next-gftn', need to add a environment
	// variable to detect production vs. development and replace credentials with
	// production credentials. Also, these credentials should not be included inside
	// the actual code and should be dynamically retrieved for security purposes.

	// Development Credentials
	firebaseCredentialsBase64 := os.Getenv(global_environment.ENV_KEY_FIREBASE_CREDENTIALS)
	firebaseDbURL := os.Getenv(global_environment.ENV_KEY_FIREBASE_DB_URL)

	// decoding Base64 based credentials
	firebaseCredentials, err := base64.StdEncoding.DecodeString(firebaseCredentialsBase64)

	if err != nil {
		LOGGER.Fatal("Error decoding base64: ", err)
	}

	// firebaseCredentials := string(firebaseCredentialsByte)

	// [START authenticate_with_admin_privileges]
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: firebaseDbURL,
	}

	// Fetch the service account key JSON file contents

	//commenting file based credential parsing
	//opt := option.WithCredentialsFile(firebaseCredentials)

	//Now aws secrets should be set for FIREBASE_CREDENTIALS to json string for firebase authentication
	// opt := option.WithCredentialsJSON([]byte(firebaseCredentials))
	opt := option.WithCredentialsJSON(firebaseCredentials)

	// firebase app:
	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		LOGGER.Fatal("Error initializing app:", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		LOGGER.Fatal("Error initializing database client:", err)
	}

	clientAuth, errAuth := app.Auth(ctx)
	if errAuth != nil {
		LOGGER.Fatal("Error initializing authenticate client:", errAuth)
	}

	// ctx := AppContext
	// ref := FbRef
	// fmt.Println("Ref: ", ref)

	// Print out all logs from root ref:
	// var data interface{}
	// if err := ref.Get(ctx, &data); err != nil {
	// 	LOGGER.Fatal("Error reading from database:", err)
	// }

	// fmt.Println(data)
	// fmt.Printf("Data: %v", data)
	// [END authenticate_with_admin_privileges]

	return client, clientAuth, err
}

// Logging middleware
func WithLogging() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			fmt.Println("\nRefValue: ", FbRef)
			// var data interface{}

			// if err := ref.Get(ctx, &data); err != nil {
			// 	LOGGER.Fatal("Error reading from database:", err)
			// }
			// fmt.Println(data)
			// fmt.Printf("Hello from logging %v", data)

			LOGGER.Info("Logged connection from %s", r.RemoteAddr)
			LOGGER.Info("Logged request:", r.Method, "on", r.URL)
			h.ServeHTTP(w, r)
		})
	}

}

func PrintLog(message string) string {
	LOGGER.Info("\n\nTimeStamp:", time.Now().String(), "\n\nError Details:", message)
	return message
}

func GetRootRef() *db.Ref {

	// As an admin, the app has access to read and write all data, regradless of Security Rules
	ref := FbClient.NewRef("/")
	return ref
}

// func AuthenticateWithLimitedPrivileges() {
// 	// [START authenticate_with_limited_privileges]
// 	ctx := context.Background()
// 	// Initialize the app with a custom auth variable, limiting the server's access
// 	ao := map[string]interface{}{"uid": "my-service-worker"}
// 	conf := &firebase.Config{
// 		DatabaseURL:  "https://databaseName.firebaseio.com",
// 		AuthOverride: &ao,
// 	}

// 	// Fetch the service account key JSON file contents
// 	opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")

// 	app, err := firebase.NewApp(ctx, conf, opt)
// 	if err != nil {
// 		LOGGER.Fatal("Error initializing app:", err)
// 	}

// 	client, err := app.Database(ctx)
// 	if err != nil {
// 		LOGGER.Fatal("Error initializing database client:", err)
// 	}

// 	// The app only has access as defined in the Security Rules
// 	ref := client.NewRef("/some_resource")
// 	var data map[string]interface{}
// 	if err := ref.Get(ctx, &data); err != nil {
// 		LOGGER.Fatal("Error reading from database:", err)
// 	}
// 	fmt.Println(data)
// 	// [END authenticate_with_limited_privileges]
// }

// func AuthenticateWithGuestPrivileges() {
// 	// [START authenticate_with_guest_privileges]
// 	ctx := context.Background()
// 	// Initialize the app with a nil auth variable, limiting the server's access
// 	var nilMap map[string]interface{}
// 	conf := &firebase.Config{
// 		DatabaseURL:  "https://databaseName.firebaseio.com",
// 		AuthOverride: &nilMap,
// 	}

// 	// Fetch the service account key JSON file contents
// 	opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")

// 	app, err := firebase.NewApp(ctx, conf, opt)
// 	if err != nil {
// 		LOGGER.Fatal("Error initializing app:", err)
// 	}

// 	client, err := app.Database(ctx)
// 	if err != nil {
// 		LOGGER.Fatal("Error initializing database client:", err)
// 	}

// 	// The app only has access to public data as defined in the Security Rules
// 	ref := client.NewRef("/some_resource")
// 	var data map[string]interface{}
// 	if err := ref.Get(ctx, &data); err != nil {
// 		LOGGER.Fatal("Error reading from database:", err)
// 	}
// 	fmt.Println(data)
// 	// [END authenticate_with_guest_privileges]
// }

func getReference(ctx context.Context, app *firebase.App) {
	// [START get_reference]
	// Create a database client from App.
	client, err := app.Database(ctx)
	if err != nil {
		LOGGER.Fatal("Error initializing database client:", err)
	}

	// Get a database reference to our blog.
	ref := client.NewRef("server/saving-data/fireblog")
	// [END get_reference]
	fmt.Println(ref.Path)
}

// [START user_type]

// User is a json-serializable type.
type User struct {
	DateOfBirth string `json:"date_of_birth,omitempty"`
	FullName    string `json:"full_name,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
}

// [END user_type]

func setValue(ctx context.Context, ref *db.Ref) {
	// [START set_value]
	usersRef := ref.Child("users")
	err := usersRef.Set(ctx, map[string]*User{
		"alanisawesome": {
			DateOfBirth: "June 23, 1912",
			FullName:    "Alan Turing",
		},
		"gracehop": {
			DateOfBirth: "December 9, 1906",
			FullName:    "Grace Hopper",
		},
	})
	if err != nil {
		LOGGER.Fatal("Error setting value:", err)
	}
	// [END set_value]
}

func setChildValue(ctx context.Context, usersRef *db.Ref) {
	// [START set_child_value]
	if err := usersRef.Child("alanisawesome").Set(ctx, &User{
		DateOfBirth: "June 23, 1912",
		FullName:    "Alan Turing",
	}); err != nil {
		LOGGER.Fatal("Error setting value:", err)
	}

	if err := usersRef.Child("gracehop").Set(ctx, &User{
		DateOfBirth: "December 9, 1906",
		FullName:    "Grace Hopper",
	}); err != nil {
		LOGGER.Fatal("Error setting value:", err)
	}
	// [END set_child_value]
}

func updateChild(ctx context.Context, usersRef *db.Ref) {
	// [START update_child]
	hopperRef := usersRef.Child("gracehop")
	if err := hopperRef.Update(ctx, map[string]interface{}{
		"nickname": "Amazing Grace",
	}); err != nil {
		LOGGER.Fatal("Error updating child:", err)
	}
	// [END update_child]
}

func updateChildren(ctx context.Context, usersRef *db.Ref) {
	// [START update_children]
	if err := usersRef.Update(ctx, map[string]interface{}{
		"alanisawesome/nickname": "Alan The Machine",
		"gracehop/nickname":      "Amazing Grace",
	}); err != nil {
		LOGGER.Fatal("Error updating children:", err)
	}
	// [END update_children]
}

func overwriteValue(ctx context.Context, usersRef *db.Ref) {
	// [START overwrite_value]
	if err := usersRef.Update(ctx, map[string]interface{}{
		"alanisawesome": &User{Nickname: "Alan The Machine"},
		"gracehop":      &User{Nickname: "Amazing Grace"},
	}); err != nil {
		LOGGER.Fatal("Error updating children:", err)
	}
	// [END overwrite_value]
}

// [START post_type]

// Post is a json-serializable type.
type Post struct {
	Author string `json:"author,omitempty"`
	Title  string `json:"title,omitempty"`
}

// [END post_type]

func pushValue(ctx context.Context, ref *db.Ref) {
	// [START push_value]
	postsRef := ref.Child("posts")

	newPostRef, err := postsRef.Push(ctx, nil)
	if err != nil {
		LOGGER.Fatal("Error pushing child node:", err)
	}

	if err := newPostRef.Set(ctx, &Post{
		Author: "gracehop",
		Title:  "Announcing COBOL, a New Programming Language",
	}); err != nil {
		LOGGER.Fatal("Error setting value:", err)
	}

	// We can also chain the two calls together
	if _, err := postsRef.Push(ctx, &Post{
		Author: "alanisawesome",
		Title:  "The Turing Machine",
	}); err != nil {
		LOGGER.Fatal("Error pushing child node:", err)
	}
	// [END push_value]
}

func pushAndSetValue(ctx context.Context, postsRef *db.Ref) {
	// [START push_and_set_value]
	if _, err := postsRef.Push(ctx, &Post{
		Author: "gracehop",
		Title:  "Announcing COBOL, a New Programming Language",
	}); err != nil {
		LOGGER.Fatal("Error pushing child node:", err)
	}
	// [END push_and_set_value]
}

func pushKey(ctx context.Context, postsRef *db.Ref) {
	// [START push_key]
	// Generate a reference to a new location and add some data using Push()
	newPostRef, err := postsRef.Push(ctx, nil)
	if err != nil {
		LOGGER.Fatal("Error pushing child node:", err)
	}

	// Get the unique key generated by Push()
	postID := newPostRef.Key
	// [END push_key]
	fmt.Println(postID)
}

func transaction(ctx context.Context, client *db.Client) {
	// [START transaction]
	fn := func(t db.TransactionNode) (interface{}, error) {
		var currentValue int
		if err := t.Unmarshal(&currentValue); err != nil {
			return nil, err
		}
		return currentValue + 1, nil
	}

	ref := client.NewRef("server/saving-data/fireblog/posts/-JRHTHaIs-jNPLXOQivY/upvotes")
	if err := ref.Transaction(ctx, fn); err != nil {
		LOGGER.Fatal("Transaction failed to commit:", err)
	}
	// [END transaction]
}

func readValue(ctx context.Context, app *firebase.App) {
	// [START read_value]
	// Create a database client from App.
	client, err := app.Database(ctx)
	if err != nil {
		LOGGER.Fatal("Error initializing database client:", err)
	}

	// Get a database reference to our posts
	ref := client.NewRef("server/saving-data/fireblog/posts")

	// Read the data at the posts reference (this is a blocking operation)
	var post Post
	if err := ref.Get(ctx, &post); err != nil {
		LOGGER.Fatal("Error reading value:", err)
	}
	// [END read_value]
	fmt.Println(ref.Path)
}

// [START dinosaur_type]

// Dinosaur is a json-serializable type.
type Dinosaur struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

// [END dinosaur_type]

func orderByChild(ctx context.Context, client *db.Client) {
	// [START order_by_child]
	ref := client.NewRef("dinosaurs")

	results, err := ref.OrderByChild("height").GetOrdered(ctx)
	if err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}
	for _, r := range results {
		var d Dinosaur
		if err := r.Unmarshal(&d); err != nil {
			LOGGER.Fatal("Error unmarshaling result:", err)
		}
		fmt.Printf("%s was %d meteres tall", r.Key(), d.Height)
	}
	// [END order_by_child]
}

func orderByNestedChild(ctx context.Context, client *db.Client) {
	// [START order_by_nested_child]
	ref := client.NewRef("dinosaurs")

	results, err := ref.OrderByChild("dimensions/height").GetOrdered(ctx)
	if err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}
	for _, r := range results {
		var d Dinosaur
		if err := r.Unmarshal(&d); err != nil {
			LOGGER.Fatal("Error unmarshaling result:", err)
		}
		fmt.Printf("%s was %d meteres tall", r.Key(), d.Height)
	}
	// [END order_by_nested_child]
}

func orderByKey(ctx context.Context, client *db.Client) {
	// [START order_by_key]
	ref := client.NewRef("dinosaurs")

	results, err := ref.OrderByKey().GetOrdered(ctx)
	if err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}
	snapshot := make([]Dinosaur, len(results))
	for i, r := range results {
		var d Dinosaur
		if err := r.Unmarshal(&d); err != nil {
			LOGGER.Fatal("Error unmarshaling result:", err)
		}
		snapshot[i] = d
	}
	fmt.Println(snapshot)
	// [END order_by_key]
}

func orderByValue(ctx context.Context, client *db.Client) {
	// [START order_by_value]
	ref := client.NewRef("scores")

	results, err := ref.OrderByValue().GetOrdered(ctx)
	if err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}
	for _, r := range results {
		var score int
		if err := r.Unmarshal(&score); err != nil {
			LOGGER.Fatal("Error unmarshaling result:", err)
		}
		fmt.Printf("The %s dinosaur's score is %d\n", r.Key(), score)
	}
	// [END order_by_value]
}

func limitToLast(ctx context.Context, client *db.Client) {
	// [START limit_query_1]
	ref := client.NewRef("dinosaurs")

	results, err := ref.OrderByChild("weight").LimitToLast(2).GetOrdered(ctx)
	if err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}
	for _, r := range results {
		fmt.Println(r.Key())
	}
	// [END limit_query_1]
}

func limitToFirst(ctx context.Context, client *db.Client) {
	// [START limit_query_2]
	ref := client.NewRef("dinosaurs")

	results, err := ref.OrderByChild("height").LimitToFirst(2).GetOrdered(ctx)
	if err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}
	for _, r := range results {
		fmt.Println(r.Key())
	}
	// [END limit_query_2]
}

func limitWithValueOrder(ctx context.Context, client *db.Client) {
	// [START limit_query_3]
	ref := client.NewRef("scores")

	results, err := ref.OrderByValue().LimitToLast(3).GetOrdered(ctx)
	if err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}
	for _, r := range results {
		var score int
		if err := r.Unmarshal(&score); err != nil {
			LOGGER.Fatal("Error unmarshaling result:", err)
		}
		fmt.Printf("The %s dinosaur's score is %d\n", r.Key(), score)
	}
	// [END limit_query_3]
}

func startAt(ctx context.Context, client *db.Client) {
	// [START range_query_1]
	ref := client.NewRef("dinosaurs")

	results, err := ref.OrderByChild("height").StartAt(3).GetOrdered(ctx)
	if err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}
	for _, r := range results {
		fmt.Println(r.Key())
	}
	// [END range_query_1]
}

func endAt(ctx context.Context, client *db.Client) {
	// [START range_query_2]
	ref := client.NewRef("dinosaurs")

	results, err := ref.OrderByKey().EndAt("pterodactyl").GetOrdered(ctx)
	if err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}
	for _, r := range results {
		fmt.Println(r.Key())
	}
	// [END range_query_2]
}

func startAndEndAt(ctx context.Context, client *db.Client) {
	// [START range_query_3]
	ref := client.NewRef("dinosaurs")

	results, err := ref.OrderByKey().StartAt("b").EndAt("b\uf8ff").GetOrdered(ctx)
	if err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}
	for _, r := range results {
		fmt.Println(r.Key())
	}
	// [END range_query_3]
}

func equalTo(ctx context.Context, client *db.Client) {
	// [START range_query_4]
	ref := client.NewRef("dinosaurs")

	results, err := ref.OrderByChild("height").EqualTo(25).GetOrdered(ctx)
	if err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}
	for _, r := range results {
		fmt.Println(r.Key())
	}
	// [END range_query_4]
}

func complexQuery(ctx context.Context, client *db.Client) {
	// [START complex_query]
	ref := client.NewRef("dinosaurs")

	var favDinoHeight int
	if err := ref.Child("stegosaurus").Child("height").Get(ctx, &favDinoHeight); err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}

	query := ref.OrderByChild("height").EndAt(favDinoHeight).LimitToLast(2)
	results, err := query.GetOrdered(ctx)
	if err != nil {
		LOGGER.Fatal("Error querying database:", err)
	}
	if len(results) == 2 {
		// Data is ordered by increasing height, so we want the first entry.
		// Second entry is stegosarus.
		fmt.Printf("The dinosaur just shorter than the stegosaurus is %s\n", results[0].Key())
	} else {
		fmt.Println("The stegosaurus is the shortest dino")
	}
	// [END complex_query]
}
