# Task API

Base URL: `/api/tasks`

## Endpoints

### POST /tasks
Create a new task or subtask
- Auth: Required
- Body: `CreateTaskRequest`
- If `parent_task_id` is provided, creates a subtask

### GET /tasks
Get all tasks for current user
- Auth: Required

### GET /tasks/parent
Get only parent tasks (no subtasks)
- Auth: Required

### GET /tasks/{id}
Get task by ID with progress info
- Auth: Required
- Returns: progress_percentage, completed_subtasks, total_subtasks

### PUT /tasks/{id}
Update task
- Auth: Required

### DELETE /tasks/{id}
Delete task
- Auth: Required

### GET /tasks/{id}/subtasks
Get subtasks of a parent task
- Auth: Required

### POST /tasks/{id}/complete
Complete a subtask (auto-completes parent if all done)
- Auth: Required
- **Business Logic**:
  - Updates subtask `is_completed = true`
  - Recalculates parent progress_percentage
  - Auto-completes parent if all subtasks done

For complete API documentation, see `/api/openapi.yaml`
