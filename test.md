## 🧪 TESTING STRATEGY

### Test Pyramid
```
        ┌─────────────┐
        │   E2E Tests │  (Az sayıda, kritik flow'lar)
        └─────────────┘
      ┌─────────────────┐
      │ Integration Tests│  (Orta sayıda, API endpoints)
      └─────────────────┘
    ┌─────────────────────┐
    │    Unit Tests       │  (Çok sayıda, business logic)
    └─────────────────────┘
```

### Test Structure
```
internal/modules/{module}/
├── domain/
│   └── {entity}_test.go          # Domain logic tests
├── repository/
│   ├── repository_test.go        # Repository interface tests
│   └── postgres_test.go          # PostgreSQL implementation tests
├── service/
│   └── service_test.go           # Service layer tests (MOST IMPORTANT)
├── http/
│   └── handler_test.go           # HTTP handler tests
└── integration_test.go           # End-to-end integration tests
```

### Required Test Coverage

- **Service Layer**: Minimum %80 coverage (CRITICAL!)
- **Repository Layer**: Minimum %70 coverage
- **Handler Layer**: Minimum %60 coverage
- **Domain Layer**: %100 coverage (basit olduğu için)

### Testing Tools
```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
    "github.com/DATA-DOG/go-sqlmock"
)
```

---

## 📋 1) UNIT TESTS (Service Layer)

### Service Test Template
```go
// internal/modules/task/service/service_test.go

package service

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "lifeplanner/internal/modules/task/domain"
    "lifeplanner/internal/modules/task/dto"
)

// Mock Repository
type MockTaskRepository struct {
    mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
    args := m.Called(ctx, task)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id int) (*domain.Task, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetSubtasks(ctx context.Context, parentID int) ([]*domain.Task, error) {
    args := m.Called(ctx, parentID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *domain.Task) error {
    args := m.Called(ctx, task)
    return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id int) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockTaskRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Task, error) {
    args := m.Called(ctx, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*domain.Task), args.Error(1)
}

// Test Suite
type TaskServiceTestSuite struct {
    suite.Suite
    mockRepo *MockTaskRepository
    service  TaskService
}

func (suite *TaskServiceTestSuite) SetupTest() {
    suite.mockRepo = new(MockTaskRepository)
    suite.service = NewTaskService(suite.mockRepo)
}

func (suite *TaskServiceTestSuite) TearDownTest() {
    suite.mockRepo.AssertExpectations(suite.T())
}

// Test: Create Task - Success
func (suite *TaskServiceTestSuite) TestCreateTask_Success() {
    // Arrange
    ctx := context.Background()
    userID := 1
    req := &dto.CreateTaskRequest{
        Title:       "Test Task",
        Description: "Test Description",
        Priority:    "high",
    }

    expectedTask := &domain.Task{
        ID:          123,
        UserID:      userID,
        Title:       req.Title,
        Description: req.Description,
        Priority:    req.Priority,
        IsCompleted: false,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    suite.mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Task")).
        Return(expectedTask, nil)
    suite.mockRepo.On("GetSubtasks", ctx, 123).
        Return([]*domain.Task{}, nil)

    // Act
    result, err := suite.service.CreateTask(ctx, req, userID)

    // Assert
    assert.NoError(suite.T(), err)
    assert.NotNil(suite.T(), result)
    assert.Equal(suite.T(), 123, result.ID)
    assert.Equal(suite.T(), "Test Task", result.Title)
    assert.Equal(suite.T(), 0.0, result.ProgressPercentage)
}

// Test: Create Task - Validation Error
func (suite *TaskServiceTestSuite) TestCreateTask_EmptyTitle() {
    // Arrange
    ctx := context.Background()
    userID := 1
    req := &dto.CreateTaskRequest{
        Title:    "", // Empty title
        Priority: "high",
    }

    // Act
    result, err := suite.service.CreateTask(ctx, req, userID)

    // Assert
    assert.Error(suite.T(), err)
    assert.Nil(suite.T(), result)
    assert.Contains(suite.T(), err.Error(), "title")
}

// Test: Complete Subtask - Success
func (suite *TaskServiceTestSuite) TestCompleteSubtask_Success() {
    // Arrange
    ctx := context.Background()
    userID := 1
    subtaskID := 100
    parentID := 50

    subtask := &domain.Task{
        ID:           subtaskID,
        UserID:       userID,
        ParentTaskID: &parentID,
        Title:        "Subtask 1",
        IsCompleted:  false,
    }

    suite.mockRepo.On("GetByID", ctx, subtaskID).Return(subtask, nil)
    suite.mockRepo.On("Update", ctx, mock.MatchedBy(func(t *domain.Task) bool {
        return t.ID == subtaskID && t.IsCompleted == true
    })).Return(nil)
    suite.mockRepo.On("GetSubtasks", ctx, parentID).Return([]*domain.Task{}, nil)

    // Act
    err := suite.service.CompleteSubtask(ctx, subtaskID, userID)

    // Assert
    assert.NoError(suite.T(), err)
    assert.True(suite.T(), subtask.IsCompleted)
    assert.NotNil(suite.T(), subtask.CompletedAt)
}

// Test: Complete Subtask - Unauthorized
func (suite *TaskServiceTestSuite) TestCompleteSubtask_Unauthorized() {
    // Arrange
    ctx := context.Background()
    userID := 1
    subtaskID := 100

    subtask := &domain.Task{
        ID:          subtaskID,
        UserID:      999, // Different user!
        Title:       "Subtask 1",
        IsCompleted: false,
    }

    suite.mockRepo.On("GetByID", ctx, subtaskID).Return(subtask, nil)

    // Act
    err := suite.service.CompleteSubtask(ctx, subtaskID, userID)

    // Assert
    assert.Error(suite.T(), err)
    assert.Contains(suite.T(), err.Error(), "unauthorized")
}

// Test: Complete Subtask - Auto Complete Parent
func (suite *TaskServiceTestSuite) TestCompleteSubtask_AutoCompleteParent() {
    // Arrange
    ctx := context.Background()
    userID := 1
    subtaskID := 100
    parentID := 50

    subtask := &domain.Task{
        ID:           subtaskID,
        UserID:       userID,
        ParentTaskID: &parentID,
        Title:        "Subtask 3",
        IsCompleted:  false,
    }

    // Other subtasks already completed
    otherSubtasks := []*domain.Task{
        {ID: 101, IsCompleted: true},
        {ID: 102, IsCompleted: true},
    }

    parentTask := &domain.Task{
        ID:          parentID,
        Title:       "Parent Task",
        IsCompleted: false,
    }

    suite.mockRepo.On("GetByID", ctx, subtaskID).Return(subtask, nil)
    suite.mockRepo.On("Update", ctx, mock.MatchedBy(func(t *domain.Task) bool {
        return t.ID == subtaskID && t.IsCompleted == true
    })).Return(nil)
    suite.mockRepo.On("GetSubtasks", ctx, parentID).Return(otherSubtasks, nil)
    suite.mockRepo.On("GetByID", ctx, parentID).Return(parentTask, nil)
    suite.mockRepo.On("Update", ctx, mock.MatchedBy(func(t *domain.Task) bool {
        return t.ID == parentID && t.IsCompleted == true
    })).Return(nil)

    // Act
    err := suite.service.CompleteSubtask(ctx, subtaskID, userID)

    // Assert
    assert.NoError(suite.T(), err)
    assert.True(suite.T(), subtask.IsCompleted)
}

// Test: Calculate Progress - With Subtasks
func (suite *TaskServiceTestSuite) TestCalculateProgress_WithSubtasks() {
    // Arrange
    ctx := context.Background()
    taskID := 50

    subtasks := []*domain.Task{
        {ID: 1, IsCompleted: true},
        {ID: 2, IsCompleted: true},
        {ID: 3, IsCompleted: false},
        {ID: 4, IsCompleted: false},
    }

    suite.mockRepo.On("GetSubtasks", ctx, taskID).Return(subtasks, nil)

    // Act
    progress, err := suite.service.CalculateProgress(ctx, taskID)

    // Assert
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 50.0, progress) // 2/4 = 50%
}

// Test: Calculate Progress - No Subtasks, Completed
func (suite *TaskServiceTestSuite) TestCalculateProgress_NoSubtasks_Completed() {
    // Arrange
    ctx := context.Background()
    taskID := 50

    task := &domain.Task{
        ID:          taskID,
        IsCompleted: true,
    }

    suite.mockRepo.On("GetSubtasks", ctx, taskID).Return([]*domain.Task{}, nil)
    suite.mockRepo.On("GetByID", ctx, taskID).Return(task, nil)

    // Act
    progress, err := suite.service.CalculateProgress(ctx, taskID)

    // Assert
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 100.0, progress)
}

// Run test suite
func TestTaskServiceTestSuite(t *testing.T) {
    suite.Run(t, new(TaskServiceTestSuite))
}

// Table-driven tests example
func TestTaskService_Priority_Validation(t *testing.T) {
    tests := []struct {
        name        string
        priority    string
        expectError bool
    }{
        {"Valid Low", "low", false},
        {"Valid Medium", "medium", false},
        {"Valid High", "high", false},
        {"Invalid Priority", "urgent", true},
        {"Empty Priority", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic here
        })
    }
}
```

---

## 📊 2) REPOSITORY TESTS (Integration with DB)
```go
// internal/modules/task/repository/postgres_test.go

package repository

import (
    "context"
    "database/sql"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"

    "lifeplanner/internal/modules/task/domain"
)

func TestPostgresTaskRepository_Create(t *testing.T) {
    // Create mock DB
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewPostgresTaskRepository(db)

    // Arrange
    ctx := context.Background()
    now := time.Now()
    task := &domain.Task{
        UserID:      1,
        Title:       "Test Task",
        Description: "Test Description",
        Priority:    "high",
        IsCompleted: false,
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    // Expect query
    mock.ExpectQuery(`INSERT INTO tasks`).
        WithArgs(
            task.UserID,
            task.ParentTaskID,
            task.Title,
            task.Description,
            task.DueDate,
            task.EstimatedStart,
            task.EstimatedEnd,
            task.Priority,
            task.IsCompleted,
            sqlmock.AnyArg(),
            sqlmock.AnyArg(),
        ).
        WillReturnRows(
            sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
                AddRow(123, now, now),
        )

    // Act
    result, err := repo.Create(ctx, task)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, 123, result.ID)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresTaskRepository_GetSubtasks(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewPostgresTaskRepository(db)

    ctx := context.Background()
    parentID := 50

    rows := sqlmock.NewRows([]string{
        "id", "user_id", "parent_task_id", "title", "description",
        "due_date", "estimated_start", "estimated_end", "actual_start", "actual_end",
        "priority", "is_completed", "completed_at", "created_at", "updated_at",
    }).
        AddRow(1, 1, 50, "Subtask 1", "Desc", nil, nil, nil, nil, nil, "low", false, nil, time.Now(), time.Now()).
        AddRow(2, 1, 50, "Subtask 2", "Desc", nil, nil, nil, nil, nil, "medium", true, time.Now(), time.Now(), time.Now())

    mock.ExpectQuery(`SELECT .+ FROM tasks WHERE parent_task_id`).
        WithArgs(parentID).
        WillReturnRows(rows)

    // Act
    subtasks, err := repo.GetSubtasks(ctx, parentID)

    // Assert
    assert.NoError(t, err)
    assert.Len(t, subtasks, 2)
    assert.Equal(t, "Subtask 1", subtasks[0].Title)
    assert.Equal(t, "Subtask 2", subtasks[1].Title)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

---

## 🌐 3) HANDLER TESTS (HTTP)
```go
// internal/modules/task/http/handler_test.go

package http

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "lifeplanner/internal/modules/task/dto"
)

// Mock Service
type MockTaskService struct {
    mock.Mock
}

func (m *MockTaskService) CreateTask(ctx context.Context, req *dto.CreateTaskRequest, userID int) (*dto.TaskResponse, error) {
    args := m.Called(ctx, req, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.TaskResponse), args.Error(1)
}

// ... other methods

func TestTaskHandler_CreateTask_Success(t *testing.T) {
    // Arrange
    mockService := new(MockTaskService)
    handler := NewTaskHandler(mockService, nil, nil)

    reqBody := dto.CreateTaskRequest{
        Title:    "Test Task",
        Priority: "high",
    }
    body, _ := json.Marshal(reqBody)

    expectedResponse := &dto.TaskResponse{
        ID:                 123,
        Title:              "Test Task",
        Priority:           "high",
        ProgressPercentage: 0,
    }

    mockService.On("CreateTask", mock.Anything, &reqBody, 1).
        Return(expectedResponse, nil)

    req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req = req.WithContext(context.WithValue(req.Context(), "userID", 1))

    w := httptest.NewRecorder()

    // Act
    handler.CreateTask(w, req)

    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)

    var response map[string]interface{}
    json.NewDecoder(w.Body).Decode(&response)

    assert.True(t, response["success"].(bool))
    assert.NotNil(t, response["data"])

    mockService.AssertExpectations(t)
}

func TestTaskHandler_CreateTask_InvalidJSON(t *testing.T) {
    handler := NewTaskHandler(nil, nil, nil)

    req := httptest.NewRequest("POST", "/tasks", bytes.NewReader([]byte("invalid json")))
    req.Header.Set("Content-Type", "application/json")
    req = req.WithContext(context.WithValue(req.Context(), "userID", 1))

    w := httptest.NewRecorder()

    // Act
    handler.CreateTask(w, req)

    // Assert
    assert.Equal(t, http.StatusBadRequest, w.Code)

    var response map[string]interface{}
    json.NewDecoder(w.Body).Decode(&response)

    assert.False(t, response["success"].(bool))
    assert.Contains(t, response["error"], "invalid request body")
}
```

---

## 🔗 4) INTEGRATION TESTS
```go
// internal/modules/task/integration_test.go

// +build integration

package task_test

import (
    "context"
    "database/sql"
    "testing"

    _ "github.com/lib/pq"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"

    "lifeplanner/internal/modules/task/dto"
    "lifeplanner/internal/modules/task/repository"
    "lifeplanner/internal/modules/task/service"
)

type TaskIntegrationTestSuite struct {
    suite.Suite
    db      *sql.DB
    service service.TaskService
}

func (suite *TaskIntegrationTestSuite) SetupSuite() {
    // Connect to test database
    db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/lifeplanner_test?sslmode=disable")
    assert.NoError(suite.T(), err)

    suite.db = db

    // Setup service
    repo := repository.NewPostgresTaskRepository(db)
    suite.service = service.NewTaskService(repo)
}

func (suite *TaskIntegrationTestSuite) TearDownSuite() {
    suite.db.Close()
}

func (suite *TaskIntegrationTestSuite) SetupTest() {
    // Clean database before each test
    suite.db.Exec("TRUNCATE tasks CASCADE")
}

func (suite *TaskIntegrationTestSuite) TestCreateAndCompleteSubtasks() {
    ctx := context.Background()
    userID := 1

    // Create parent task
    parentReq := &dto.CreateTaskRequest{
        Title:    "Parent Task",
        Priority: "high",
    }
    parent, err := suite.service.CreateTask(ctx, parentReq, userID)
    assert.NoError(suite.T(), err)
    assert.NotNil(suite.T(), parent)

    // Create subtasks
    subtask1Req := &dto.CreateTaskRequest{
        ParentTaskID: &parent.ID,
        Title:        "Subtask 1",
        Priority:     "medium",
    }
    subtask1, err := suite.service.CreateTask(ctx, subtask1Req, userID)
    assert.NoError(suite.T(), err)

    subtask2Req := &dto.CreateTaskRequest{
        ParentTaskID: &parent.ID,
        Title:        "Subtask 2",
        Priority:     "low",
    }
    subtask2, err := suite.service.CreateTask(ctx, subtask2Req, userID)
    assert.NoError(suite.T(), err)

    // Complete first subtask
    err = suite.service.CompleteSubtask(ctx, subtask1.ID, userID)
    assert.NoError(suite.T(), err)

    // Check progress
    progress, err := suite.service.CalculateProgress(ctx, parent.ID)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 50.0, progress)

    // Complete second subtask
    err = suite.service.CompleteSubtask(ctx, subtask2.ID, userID)
    assert.NoError(suite.T(), err)

    // Check parent auto-completed
    updatedParent, err := suite.service.GetTask(ctx, parent.ID, userID)
    assert.NoError(suite.T(), err)
    assert.True(suite.T(), updatedParent.IsCompleted)
    assert.Equal(suite.T(), 100.0, updatedParent.ProgressPercentage)
}

func TestTaskIntegrationTestSuite(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests")
    }
    suite.Run(t, new(TaskIntegrationTestSuite))
}
```

---

## 🏃 RUNNING TESTS

### Run Commands
```bash
# Run all unit tests
go test ./internal/modules/task/... -v

# Run with coverage
go test ./internal/modules/task/... -cover -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out

# Run only service tests
go test ./internal/modules/task/service/... -v

# Run integration tests (requires DB)
go test ./internal/modules/task/... -tags=integration -v

# Skip integration tests
go test ./internal/modules/task/... -short

# Run specific test
go test ./internal/modules/task/service/... -run TestCompleteSubtask_Success -v

# Run tests with race detector
go test ./internal/modules/task/... -race

# Generate coverage for all modules
go test ./internal/modules/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

---

## 📋 TEST CHECKLIST (Per Module)

Before marking a module as "done", ensure:

### Service Layer
- [ ] All public methods have tests
- [ ] Success cases covered
- [ ] Error cases covered
- [ ] Edge cases covered (nil, empty, boundary values)
- [ ] Authorization checks tested
- [ ] Business logic validated
- [ ] Mock dependencies properly
- [ ] Coverage > 80%

### Repository Layer
- [ ] CRUD operations tested
- [ ] SQL queries validated (with sqlmock)
- [ ] Error handling tested
- [ ] Edge cases covered
- [ ] Coverage > 70%

### Handler Layer
- [ ] All endpoints tested
- [ ] Valid request/response tested
- [ ] Invalid JSON tested
- [ ] Validation errors tested
- [ ] Authorization tested
- [ ] Error responses validated
- [ ] Coverage > 60%

### Integration Tests
- [ ] Full flow tested (create → update → delete)
- [ ] Complex scenarios tested (subtask auto-completion)
- [ ] Database constraints validated
- [ ] Transaction rollback tested

### Domain Layer
- [ ] All business methods tested
- [ ] Calculations validated
- [ ] Coverage = 100% (should be easy)

---

## 🎯 TEST NAMING CONVENTIONS
```go
// Pattern: Test{MethodName}_{Scenario}_{ExpectedResult}

func TestCreateTask_ValidInput_Success(t *testing.T) {}
func TestCreateTask_EmptyTitle_ReturnsError(t *testing.T) {}
func TestCompleteSubtask_AllSubtasksDone_AutoCompletesParent(t *testing.T) {}
func TestCompleteSubtask_DifferentUser_ReturnsUnauthorized(t *testing.T) {}
func TestGetTask_NotFound_ReturnsError(t *testing.T) {}
```

---

## 🚨 CRITICAL TESTING RULES

1. **Service layer MUST be thoroughly tested** - This is where business logic lives
2. **Mock ALL external dependencies** - Database, HTTP clients, etc.
3. **Test error paths** - Not just happy paths
4. **Use table-driven tests** - For multiple similar scenarios
5. **Integration tests require real DB** - Use Docker for test DB
6. **Never skip tests** - All modules must have tests before PR
7. **Keep tests fast** - Unit tests should run in < 1 second
8. **Clean test data** - Always cleanup after tests
9. **Test isolation** - Tests should not depend on each other
10. **Descriptive test names** - Should explain what's being tested

---

## 🐳 TEST DATABASE SETUP
```bash
# docker-compose.test.yml
version: '3.8'
services:
  postgres-test:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
      POSTGRES_DB: lifeplanner_test
    ports:
      - "5433:5432"
    tmpfs:
      - /var/lib/postgresql/data  # In-memory for speed

# Run test DB
docker-compose -f docker-compose.test.yml up -d

# Run migrations on test DB
migrate -path internal/infrastructure/database/migrations \
        -database "postgres://test:test@localhost:5433/lifeplanner_test?sslmode=disable" up

# Run integration tests
go test ./internal/modules/... -tags=integration
```

---

## 📊 CI/CD PIPELINE TEST STAGES
```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run unit tests
        run: go test ./internal/modules/... -short -cover

  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v3
      - name: Run migrations
        run: migrate -path migrations -database $DB_URL up
      - name: Run integration tests
        run: go test ./internal/modules/... -tags=integration
```

---

**REMEMBER**:
- Tests are documentation
- Tests catch bugs early
- Tests enable refactoring
- Tests give confidence
- **Write tests BEFORE marking module as done!**
