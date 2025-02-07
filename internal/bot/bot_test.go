package bot

//
//import (
//	"testing"
//
//	"github.com/gwassel/TasksOfWoe/internal/persistence/mock"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//)
//
//type MockBotAPI struct {
//	mock.MockDB
//}
//
//func (m *MockBotAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
//	args := m.Called(c)
//	return args.Get(0).(tgbotapi.Message), args.Error(1)
//}
//
//func TestHandleMessage_AddTask(t *testing.T) {
//	// Create mock instances
//	mockAPI := new(MockBotAPI)
//	mockDB := new(mock.MockDB)
//
//	// Create the bot with mock dependencies
//	bot := NewBot(mockAPI, mockDB)
//
//	// Mock the database call
//	mockDB.On("AddTask", int64(123), "Test task").Return(nil)
//
//	// Mock the bot API call
//	mockAPI.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil)
//
//	// Create a test message
//	message := &tgbotapi.Message{
//		From: &tgbotapi.User{ID: 123},
//		Text: "/add Test task",
//	}
//
//	// Call the handler
//	bot.HandleMessage(message)
//
//	// Verify the mocks were called as expected
//	mockDB.AssertCalled(t, "AddTask", int64(123), "Test task")
//	mockAPI.AssertCalled(t, "Send", mock.Anything)
//}
//
//func TestHandleMessage_ListTasks(t *testing.T) {
//	mockAPI := new(MockBotAPI)
//	mockDB := new(mock.MockDB)
//
//	bot := NewBot(mockAPI, mockDB)
//
//	// Mock the database call
//	mockDB.On("ListTasks", int64(123)).Return([]database.Task{
//		{ID: 1, UserID: 123, Task: "Task 1", Completed: false},
//		{ID: 2, UserID: 123, Task: "Task 2", Completed: true},
//	}, nil)
//
//	// Mock the bot API call
//	mockAPI.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil)
//
//	// Create a test message
//	message := &tgbotapi.Message{
//		From: &tgbotapi.User{ID: 123},
//		Text: "/list",
//	}
//
//	// Call the handler
//	bot.HandleMessage(message)
//
//	// Verify the mocks were called as expected
//	mockDB.AssertCalled(t, "ListTasks", int64(123))
//	mockAPI.AssertCalled(t, "Send", mock.Anything)
//}
//
//func TestHandleMessage_CompleteTask(t *testing.T) {
//	mockAPI := new(MockBotAPI)
//	mockDB := new(mock.MockDB)
//
//	bot := NewBot(mockAPI, mockDB)
//
//	// Mock the database call
//	mockDB.On("CompleteTask", int64(123), "1").Return(nil)
//
//	// Mock the bot API call
//	mockAPI.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil)
//
//	// Create a test message
//	message := &tgbotapi.Message{
//		From: &tgbotapi.User{ID: 123},
//		Text: "/complete 1",
//	}
//
//	// Call the handler
//	bot.HandleMessage(message)
//
//	// Verify the mocks were called as expected
//	mockDB.AssertCalled(t, "CompleteTask", int64(123), "1")
//	mockAPI.AssertCalled(t, "Send", mock.Anything)
//}
