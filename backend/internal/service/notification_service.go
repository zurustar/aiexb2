// backend/internal/service/notification_service.go
package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/queue"
	"github.com/your-org/esms/internal/repository"
)

var (
	ErrTemplateNotFound = errors.New("notification template not found")
	ErrInvalidEmail     = errors.New("invalid email address")
)

// NotificationType は通知の種類を表します
type NotificationType string

const (
	NotificationTypeReservationCreated  NotificationType = "reservation_created"
	NotificationTypeReservationApproved NotificationType = "reservation_approved"
	NotificationTypeReservationRejected NotificationType = "reservation_rejected"
	NotificationTypeReservationCanceled NotificationType = "reservation_canceled"
	NotificationTypeReservationReminder NotificationType = "reservation_reminder"
)

// EmailSender はメール送信インターフェース
type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}

// NotificationService は通知に関するビジネスロジックを提供します
type NotificationService struct {
	userRepo    repository.UserRepository
	jobQueue    queue.JobQueue
	emailSender EmailSender
	templates   map[NotificationType]*template.Template
	sentCache   map[string]time.Time // 重複送信防止用キャッシュ (簡易実装)
}

// NewNotificationService は新しいNotificationServiceを作成します
func NewNotificationService(
	userRepo repository.UserRepository,
	jobQueue queue.JobQueue,
	emailSender EmailSender,
) *NotificationService {
	svc := &NotificationService{
		userRepo:    userRepo,
		jobQueue:    jobQueue,
		emailSender: emailSender,
		templates:   make(map[NotificationType]*template.Template),
		sentCache:   make(map[string]time.Time),
	}

	// テンプレート初期化
	svc.initTemplates()

	return svc
}

// initTemplates はメールテンプレートを初期化します
func (s *NotificationService) initTemplates() {
	// 予約作成通知テンプレート
	s.templates[NotificationTypeReservationCreated] = template.Must(template.New("reservation_created").Parse(`
予約が作成されました

タイトル: {{.Title}}
開始時刻: {{.StartAt}}
終了時刻: {{.EndAt}}
主催者: {{.OrganizerName}}

詳細はシステムでご確認ください。
`))

	// 予約承認通知テンプレート
	s.templates[NotificationTypeReservationApproved] = template.Must(template.New("reservation_approved").Parse(`
予約が承認されました

タイトル: {{.Title}}
開始時刻: {{.StartAt}}
終了時刻: {{.EndAt}}

ご利用をお待ちしております。
`))

	// 予約却下通知テンプレート
	s.templates[NotificationTypeReservationRejected] = template.Must(template.New("reservation_rejected").Parse(`
予約が却下されました

タイトル: {{.Title}}
開始時刻: {{.StartAt}}
終了時刻: {{.EndAt}}
理由: {{.Reason}}

別の日時でのご予約をご検討ください。
`))

	// 予約キャンセル通知テンプレート
	s.templates[NotificationTypeReservationCanceled] = template.Must(template.New("reservation_canceled").Parse(`
予約がキャンセルされました

タイトル: {{.Title}}
開始時刻: {{.StartAt}}
終了時刻: {{.EndAt}}

キャンセルされた予約の詳細はシステムでご確認ください。
`))

	// リマインダー通知テンプレート
	s.templates[NotificationTypeReservationReminder] = template.Must(template.New("reservation_reminder").Parse(`
予約のリマインダー

タイトル: {{.Title}}
開始時刻: {{.StartAt}}
終了時刻: {{.EndAt}}
場所: {{.Location}}

まもなく予約時刻です。ご準備ください。
`))
}

// NotifyReservationCreated は予約作成通知を送信します
func (s *NotificationService) NotifyReservationCreated(ctx context.Context, reservation *domain.Reservation, organizer *domain.User) error {
	// 重複送信チェック
	cacheKey := fmt.Sprintf("created_%s", reservation.ID.String())
	if s.isDuplicate(cacheKey) {
		return nil // 重複送信を防止
	}

	// テンプレートデータ準備
	data := map[string]interface{}{
		"Title":         reservation.Title,
		"StartAt":       reservation.StartAt.Format("2006-01-02 15:04"),
		"EndAt":         reservation.EndAt.Format("2006-01-02 15:04"),
		"OrganizerName": organizer.Name,
	}

	// メール本文生成
	body, err := s.renderTemplate(NotificationTypeReservationCreated, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// 非同期でメール送信ジョブをキューに追加
	payload := map[string]interface{}{
		"to":      organizer.Email,
		"subject": "予約が作成されました",
		"body":    body,
	}

	_, err = s.jobQueue.Enqueue(ctx, "send_email", payload)
	if err != nil {
		return fmt.Errorf("failed to enqueue email job: %w", err)
	}

	// 送信済みキャッシュに追加
	s.markAsSent(cacheKey)

	return nil
}

// NotifyReservationApproved は予約承認通知を送信します
func (s *NotificationService) NotifyReservationApproved(ctx context.Context, reservation *domain.Reservation, organizer *domain.User) error {
	cacheKey := fmt.Sprintf("approved_%s", reservation.ID.String())
	if s.isDuplicate(cacheKey) {
		return nil
	}

	data := map[string]interface{}{
		"Title":   reservation.Title,
		"StartAt": reservation.StartAt.Format("2006-01-02 15:04"),
		"EndAt":   reservation.EndAt.Format("2006-01-02 15:04"),
	}

	body, err := s.renderTemplate(NotificationTypeReservationApproved, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	payload := map[string]interface{}{
		"to":      organizer.Email,
		"subject": "予約が承認されました",
		"body":    body,
	}

	_, err = s.jobQueue.Enqueue(ctx, "send_email", payload)
	if err != nil {
		return fmt.Errorf("failed to enqueue email job: %w", err)
	}

	s.markAsSent(cacheKey)
	return nil
}

// NotifyReservationRejected は予約却下通知を送信します
func (s *NotificationService) NotifyReservationRejected(ctx context.Context, reservation *domain.Reservation, organizer *domain.User, reason string) error {
	cacheKey := fmt.Sprintf("rejected_%s", reservation.ID.String())
	if s.isDuplicate(cacheKey) {
		return nil
	}

	data := map[string]interface{}{
		"Title":   reservation.Title,
		"StartAt": reservation.StartAt.Format("2006-01-02 15:04"),
		"EndAt":   reservation.EndAt.Format("2006-01-02 15:04"),
		"Reason":  reason,
	}

	body, err := s.renderTemplate(NotificationTypeReservationRejected, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	payload := map[string]interface{}{
		"to":      organizer.Email,
		"subject": "予約が却下されました",
		"body":    body,
	}

	_, err = s.jobQueue.Enqueue(ctx, "send_email", payload)
	if err != nil {
		return fmt.Errorf("failed to enqueue email job: %w", err)
	}

	s.markAsSent(cacheKey)
	return nil
}

// renderTemplate はテンプレートをレンダリングします
func (s *NotificationService) renderTemplate(notifType NotificationType, data map[string]interface{}) (string, error) {
	tmpl, ok := s.templates[notifType]
	if !ok {
		return "", ErrTemplateNotFound
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// isDuplicate は重複送信かチェックします（簡易実装）
func (s *NotificationService) isDuplicate(cacheKey string) bool {
	sentAt, exists := s.sentCache[cacheKey]
	if !exists {
		return false
	}

	// 1時間以内の送信は重複とみなす
	return time.Since(sentAt) < 1*time.Hour
}

// markAsSent は送信済みとしてマークします
func (s *NotificationService) markAsSent(cacheKey string) {
	s.sentCache[cacheKey] = time.Now()

	// キャッシュサイズ制限（簡易実装）
	if len(s.sentCache) > 10000 {
		// 古いエントリを削除
		for k, v := range s.sentCache {
			if time.Since(v) > 24*time.Hour {
				delete(s.sentCache, k)
			}
		}
	}
}
