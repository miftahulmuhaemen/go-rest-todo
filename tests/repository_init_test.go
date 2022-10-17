package tests

import (
	"context"
	"fmt"
	"go-rest-todo/business"
	core "go-rest-todo/core/user"
	repo "go-rest-todo/modules/repository/user"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var DB *mongo.Database

type Config struct {
	DatabaseRepo  string
	DatabaseTag   string
	DatabaseName  string
	ConnectionURI string
	Port          string
	EnvVars       []string
}

func TestMain(m *testing.M) {

	var (
		username = "root"
		password = "password"
		dbname   = "ews"
	)

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "latest",
		Env: []string{
			fmt.Sprintf("%s=%s", "MONGO_INITDB_ROOT_USERNAME", username),
			fmt.Sprintf("%s=%s", "MONGO_INITDB_ROOT_PASSWORD", password),
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	err = pool.Retry(func() error {
		var err error
		MongoClient, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@localhost:%s", username, password, resource.GetPort("27017/tcp"))))
		if err != nil {
			return err
		}
		DB = MongoClient.Database(dbname)
		return MongoClient.Ping(context.TODO(), nil)
	})
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	exitCode := m.Run()
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(exitCode)
}

func TestRepositoryUserCreate(t *testing.T) {

	clock := clockwork.NewFakeClock()
	repo := repo.NewMongoDBRepository(DB)
	tests := []struct {
		name           string
		input          core.User
		expectedReturn core.User
		errExpectation error
	}{
		{
			name:           "Input empty;",
			errExpectation: business.ErrInvalidID,
		},
		{
			name: "RoleID is not ObjectID;",
			input: core.User{
				RoleID:    "3f82fa1f-42a4-401b-8607-0674b94b6dab",
				Username:  "user",
				Password:  "password",
				Name:      "name",
				CreatedAt: clock.Now(),
				Updatet:   clock.Now(),
			},
			errExpectation: business.ErrInvalidID,
		},
		{
			name: "Success;",
			input: core.User{
				RoleID:    "634ce6d867fb6b4853384815",
				Username:  "user",
				Password:  "password",
				Name:      "name",
				CreatedAt: clock.Now(),
				Updatet:   clock.Now(),
			},
			expectedReturn: core.User{
				RoleID:    "634ce6d867fb6b4853384815",
				Username:  "user",
				Password:  "password",
				Name:      "name",
				CreatedAt: clock.Now(),
				Updatet:   clock.Now(),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {

			t.Parallel()

			ret, err := repo.Create(test.input)
			if err != nil {
				log.Println(err)
				assert.True(t, strings.Contains(err.Error(), test.errExpectation.Error()))
			} else {
				assert.NotEmpty(t, ret.ID)
				ret.ID = test.expectedReturn.ID
				assert.Equal(t, test.expectedReturn, ret)
			}
		})
	}
}
