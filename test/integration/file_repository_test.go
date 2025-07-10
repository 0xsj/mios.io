// test/integration/file_repository_test.go
package integration

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/0xsj/mios.io/config"
	db "github.com/0xsj/mios.io/db/sqlc"
	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FileRepositoryTestSuite struct {
	suite.Suite
	pool       *pgxpool.Pool
	queries    *db.Queries
	userRepo   repository.UserRepository
	logger     log.Logger
	cleanup    func()
}

func (suite *FileRepositoryTestSuite) SetupSuite() {
	// Load test configuration
	cfg := config.LoadConfig("test", "../..")
	
	// Create test logger
	suite.logger = log.Development().WithLayer("FileRepoTest")
	
	// Connect to test database
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName+"_test")
	
	pool, err := pgxpool.Connect(context.Background(), dbURL)
	require.NoError(suite.T(), err)
	
	suite.pool = pool
	suite.queries = db.New(pool)
	suite.userRepo = repository.NewUserRepository(suite.queries, suite.logger)
	
	// Setup cleanup function
	suite.cleanup = func() {
		pool.Close()
	}
}

func (suite *FileRepositoryTestSuite) TearDownSuite() {
	if suite.cleanup != nil {
		suite.cleanup()
	}
}

func (suite *FileRepositoryTestSuite) SetupTest() {
	// Clean database before each test
	suite.cleanDatabase()
}

func (suite *FileRepositoryTestSuite) cleanDatabase() {
	ctx := context.Background()
	
	// Clean in correct order due to foreign keys
	_, err := suite.pool.Exec(ctx, "DELETE FROM analytics")
	require.NoError(suite.T(), err)
	
	_, err = suite.pool.Exec(ctx, "DELETE FROM content_items")
	require.NoError(suite.T(), err)
	
	_, err = suite.pool.Exec(ctx, "DELETE FROM auth")
	require.NoError(suite.T(), err)
	
	_, err = suite.pool.Exec(ctx, "DELETE FROM users")
	require.NoError(suite.T(), err)
}

func (suite *FileRepositoryTestSuite) createTestUser() *db.User {
	ctx := context.Background()
	
	params := repository.CreateUserParams{
		Username: "testuser",
		Handle:   "testuser",
		Email:    "test@example.com",
	}
	
	user, err := suite.userRepo.CreateUser(ctx, params)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), user)
	
	return user
}

func TestFileRepositoryTestSuite(t *testing.T) {
	// Skip integration tests if no database
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}
	
	suite.Run(t, new(FileRepositoryTestSuite))
}