package repository

// please notice the limitations
// - transaction is not supported, (if transaction is needed, please refer to github.com/strikesecurity/strikememongo)
// - ephemeralForTest is used for storage engine which is different from inMemory
// - only supports UNIX systems. CI will run on MacOS, Ubuntu Xenial, Ubuntu Trusty, and Ubuntu Precise.
//   Other flavors of Linux may or may not work.

import (
	"github.com/stretchr/testify/assert"
	"github.com/tryvium-travels/memongo"
	"log"
	"strings"
	"testing"
)

var (
	databaseName = ""
	mongoURI     = ""
)

func TestMain(m *testing.M) {
	// tested with 4.2.17, 4.4.10, 5.0.3
	mongoServer, err := memongo.StartWithOptions(&memongo.Options{MongoVersion: "5.0.3"})
	if err != nil {
		log.Fatal(err)
	}
	defer mongoServer.Stop()

	mongoURI = mongoServer.URIWithRandomDB()
	splitedDatabaseName := strings.Split(mongoURI, "/")
	databaseName = splitedDatabaseName[len(splitedDatabaseName)-1]

	setup()
	m.Run()
}

func TestAddClient(t *testing.T) {
	client := Client{
		Name: "abc",
		Email: "abc1@gmail.com",
		Surname: "somebody",
	}

	err := AddClient(&client)
	assert.Equal(t, nil, err)

	fetched, err := FindFirstByName("abc")
	assert.Equal(t, nil, err)
	assert.Equal(t, client.Email, fetched.Email)
	assert.Equal(t, client.Surname, fetched.Surname)

	fetched, err = FindFirstByName("xyz")
	assert.Error(t, err)
	assert.Equal(t, (*Client)(nil), fetched)
}

func setup() {
	err := InitMongo(databaseName, mongoURI)
	if err != nil {
		panic("cannot connect to DB")
	}
}
