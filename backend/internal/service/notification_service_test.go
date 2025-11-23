// backend/internal/service/notification_service_test.go
package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/queue"
	"github.com/your-org/esms/internal/service"
)

// Mock EmailSender
type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) Send(ctx context.Context, to, subject, body string) error {
	args := m.Called(ctx, to, subject, body)
	return args.Error(0)
}

// Mock JobQueue
type MockJobQueue struct {
	mock.Mock
}

func (m *MockJobQueue) Enqueue(ctx context.Context, jobType string, payload map[string]interface{}) (string, error) {
	args := m.Called(ctx, jobType, payload)
	return args.String(0), args.Error(1)
}

func (m *MockJobQueue) Dequeue(ctx context.Context) (*queue.Job, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*queue.Job), args.Error(1)
}

func TestNotificationService_NotifyReservationCreated(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockJobQueue := new(MockJobQueue)
	mockEmailSender := new(MockEmailSender)

	svc := service.NewNotificationService(mockUserRepo, mockJobQueue, mockEmailSender)

	ctx := context.Background()
	reservation := &domain.Reservation{
		ID:      uuid.New(),
		Title:   "Test Meeting",
		StartAt: time.Now().Add(24 * time.Hour),
		EndAt:   time.Now().Add(25 * time.Hour),
	}

	organizer := &domain.User{
		ID:    uuid.New(),
		Email: "organizer@example.com",
		Name:  "Test Organizer",
	}

	mockJobQueue.On("Enqueue", ctx, "send_email", mock.MatchedBy(func(payload map[string]interface{}) bool {
		return payload["to"] == organizer.Email && payload["subject"] == "予約が作成されました"
	})).Return("job-id", nil)

	err := svc.NotifyReservationCreated(ctx, reservation, organizer)

	assert.NoError(t, err)
	mockJobQueue.AssertExpectations(t)
}

func TestNotificationService_NotifyReservationCreated_DuplicatePrevention(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockJobQueue := new(MockJobQueue)
	mockEmailSender := new(MockEmailSender)

	svc := service.NewNotificationService(mockUserRepo, mockJobQueue, mockEmailSender)

	ctx := context.Background()
	reservation := &domain.Reservation{
		ID:      uuid.New(),
		Title:   "Test Meeting",
		StartAt: time.Now().Add(24 * time.Hour),
		EndAt:   time.Now().Add(25 * time.Hour),
	}

	organizer := &domain.User{
		ID:    uuid.New(),
		Email: "organizer@example.com",
		Name:  "Test Organizer",
	}

	// 1回目の送信
	mockJobQueue.On("Enqueue", ctx, "send_email", mock.Anything).Return("job-id", nil).Once()
	err := svc.NotifyReservationCreated(ctx, reservation, organizer)
	assert.NoError(t, err)

	// 2回目の送信（重複）- Enqueueは呼ばれない
	err = svc.NotifyReservationCreated(ctx, reservation, organizer)
	assert.NoError(t, err)

	// Enqueueが1回だけ呼ばれたことを確認
	mockJobQueue.AssertNumberOfCalls(t, "Enqueue", 1)
}

func TestNotificationService_NotifyReservationApproved(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockJobQueue := new(MockJobQueue)
	mockEmailSender := new(MockEmailSender)

	svc := service.NewNotificationService(mockUserRepo, mockJobQueue, mockEmailSender)

	ctx := context.Background()
	reservation := &domain.Reservation{
		ID:      uuid.New(),
		Title:   "Test Meeting",
		StartAt: time.Now().Add(24 * time.Hour),
		EndAt:   time.Now().Add(25 * time.Hour),
	}

	organizer := &domain.User{
		ID:    uuid.New(),
		Email: "organizer@example.com",
		Name:  "Test Organizer",
	}

	mockJobQueue.On("Enqueue", ctx, "send_email", mock.MatchedBy(func(payload map[string]interface{}) bool {
		return payload["to"] == organizer.Email && payload["subject"] == "予約が承認されました"
	})).Return("job-id", nil)

	err := svc.NotifyReservationApproved(ctx, reservation, organizer)

	assert.NoError(t, err)
	mockJobQueue.AssertExpectations(t)
}

func TestNotificationService_NotifyReservationRejected(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockJobQueue := new(MockJobQueue)
	mockEmailSender := new(MockEmailSender)

	svc := service.NewNotificationService(mockUserRepo, mockJobQueue, mockEmailSender)

	ctx := context.Background()
	reservation := &domain.Reservation{
		ID:      uuid.New(),
		Title:   "Test Meeting",
		StartAt: time.Now().Add(24 * time.Hour),
		EndAt:   time.Now().Add(25 * time.Hour),
	}

	organizer := &domain.User{
		ID:    uuid.New(),
		Email: "organizer@example.com",
		Name:  "Test Organizer",
	}

	reason := "リソースが不足しています"

	mockJobQueue.On("Enqueue", ctx, "send_email", mock.MatchedBy(func(payload map[string]interface{}) bool {
		return payload["to"] == organizer.Email && payload["subject"] == "予約が却下されました"
	})).Return("job-id", nil)

	err := svc.NotifyReservationRejected(ctx, reservation, organizer, reason)

	assert.NoError(t, err)
	mockJobQueue.AssertExpectations(t)
}
