package repository_test

import (
	"context"
	"flag"
	"go-clean-architecture/todo/models"
	"go-clean-architecture/todo/repository"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	// All tests that use mtest.Setup() are expected to be integration tests, so skip them when the
	// -short flag is included in the "go test" command. Also, we have to parse flags here to use
	// testing.Short() because flags aren't parsed before TestMain() is called.
	flag.Parse()
	if testing.Short() {
		log.Print("skipping mtest integration test in short mode")
		return
	}

	if err := mtest.Setup(); err != nil {
		log.Fatal(err)
	}
	defer os.Exit(m.Run())
	if err := mtest.Teardown(); err != nil {
		log.Fatal(err)
	}
}

func TestTodoFindAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mtest.ClusterURI()))
	assert.NoError(t, err)
	defer client.Disconnect(ctx)

	repo := repository.NewMongoTodoRepository(client)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("when success", func(mt *mtest.T) {
		bsonData, err := bson.Marshal(&models.Todo{})
		assert.NoError(mt, err)

		var bsonD bson.D
		err = bson.Unmarshal(bsonData, &bsonD)
		assert.NoError(mt, err)

		find := mtest.CreateCursorResponse(
			1,
			"todo",
			mtest.FirstBatch,
			bsonD)
		getMore := mtest.CreateCursorResponse(
			1,
			"todo",
			mtest.NextBatch,
			bsonD,
		)
		killCursors := mtest.CreateCursorResponse(
			0,
			"todo",
			mtest.NextBatch,
		)
		mt.AddMockResponses(find, getMore, killCursors)

		repo.FindAll("", 10, 0)
	})
}
