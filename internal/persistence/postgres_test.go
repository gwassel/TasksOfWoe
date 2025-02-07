package persistence

//func TestAddTask(t *testing.T) {
//	db, err := NewDB("user=youruser dbname=task_tracker_test sslmode=disable password=yourpassword")
//	assert.NoError(t, err)
//
//	// Clean up the database before testing
//	_, err = db.Exec("DELETE FROM tasks")
//	assert.NoError(t, err)
//
//	// Test adding a task
//	err = db.AddTask(123, "Test task")
//	assert.NoError(t, err)
//
//	// Verify the task was added
//	var tasks []domain.Task
//	err = db.Select(&tasks, "SELECT * FROM tasks WHERE user_id = $1", 123)
//	assert.NoError(t, err)
//	assert.Equal(t, 1, len(tasks))
//	assert.Equal(t, "Test task", tasks[0].Task)
//}
//
//func TestListTasks(t *testing.T) {
//	db, err := NewDB("user=youruser dbname=task_tracker_test sslmode=disable password=yourpassword")
//	assert.NoError(t, err)
//
//	// Clean up the database before testing
//	_, err = db.Exec("DELETE FROM tasks")
//	assert.NoError(t, err)
//
//	// Add test tasks
//	_, err = db.Exec("INSERT INTO tasks (user_id, task, created_at) VALUES ($1, $2, $3)", 123, "Task 1", time.Now())
//	assert.NoError(t, err)
//	_, err = db.Exec("INSERT INTO tasks (user_id, task, created_at) VALUES ($1, $2, $3)", 123, "Task 2", time.Now())
//	assert.NoError(t, err)
//
//	// Test listing tasks
//	tasks, err := db.ListTasks(123)
//	assert.NoError(t, err)
//	assert.Equal(t, 2, len(tasks))
//}
//
//func TestCompleteTask(t *testing.T) {
//	db, err := NewDB("user=youruser dbname=task_tracker_test sslmode=disable password=yourpassword")
//	assert.NoError(t, err)
//
//	// Clean up the database before testing
//	_, err = db.Exec("DELETE FROM tasks")
//	assert.NoError(t, err)
//
//	// Add a test task
//	_, err = db.Exec("INSERT INTO tasks (user_id, task, created_at) VALUES ($1, $2, $3)", 123, "Test task", time.Now())
//	assert.NoError(t, err)
//
//	// Get the task ID
//	var task domain.Task
//	err = db.Get(&task, "SELECT id FROM tasks WHERE user_id = $1", 123)
//	assert.NoError(t, err)
//
//	// Test completing the task
//	err = db.CompleteTask(123, task.ID)
//	assert.NoError(t, err)
//
//	// Verify the task was completed
//	var completed bool
//	err = db.Get(&completed, "SELECT completed FROM tasks WHERE id = $1", task.ID)
//	assert.NoError(t, err)
//	assert.True(t, completed)
//}
