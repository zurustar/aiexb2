// backend/internal/domain/user_test.go
package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/your-org/esms/internal/domain"
)

func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		name string
		user domain.User
		want bool
	}{
		{
			name: "Valid active user",
			user: domain.User{IsActive: true, DeletedAt: nil},
			want: true,
		},
		{
			name: "Inactive user",
			user: domain.User{IsActive: false, DeletedAt: nil},
			want: false,
		},
		{
			name: "Deleted user",
			user: domain.User{IsActive: true, DeletedAt: &time.Time{}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.user.IsValid())
		})
	}
}

func TestUser_HasRole(t *testing.T) {
	user := domain.User{Role: domain.RoleAdmin}
	assert.True(t, user.HasRole(domain.RoleAdmin))
	assert.False(t, user.HasRole(domain.RoleGeneral))
}

func TestUser_CanManage(t *testing.T) {
	admin := domain.User{ID: uuid.New(), Role: domain.RoleAdmin}
	manager := domain.User{ID: uuid.New(), Role: domain.RoleManager}
	general := domain.User{ID: uuid.New(), Role: domain.RoleGeneral, ManagerID: &manager.ID}
	other := domain.User{ID: uuid.New(), Role: domain.RoleGeneral}

	tests := []struct {
		name   string
		user   domain.User
		target domain.User
		want   bool
	}{
		{"Admin can manage anyone", admin, general, true},
		{"Manager can manage subordinate", manager, general, true},
		{"Manager cannot manage non-subordinate", manager, other, false},
		{"General cannot manage anyone", general, other, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.user.CanManage(&tt.target))
		})
	}
}

func TestUser_CanAccessResource(t *testing.T) {
	admin := domain.User{Role: domain.RoleAdmin}
	manager := domain.User{Role: domain.RoleManager}
	general := domain.User{Role: domain.RoleGeneral}
	auditor := domain.User{Role: domain.RoleAuditor}

	roleAdmin := domain.RoleAdmin
	roleManager := domain.RoleManager
	roleGeneral := domain.RoleGeneral

	tests := []struct {
		name         string
		user         domain.User
		requiredRole *domain.Role
		want         bool
	}{
		{"No required role (nil)", general, nil, true},
		{"Admin accessing Admin resource", admin, &roleAdmin, true},
		{"Manager accessing Admin resource", manager, &roleAdmin, false},
		{"Manager accessing Manager resource", manager, &roleManager, true},
		{"General accessing Manager resource", general, &roleManager, false},
		{"Admin accessing General resource", admin, &roleGeneral, true},
		{"Auditor accessing General resource", auditor, &roleGeneral, false}, // Auditor is special
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.user.CanAccessResource(tt.requiredRole))
		})
	}
}

func TestUser_Penalty(t *testing.T) {
	future := time.Now().Add(1 * time.Hour)
	past := time.Now().Add(-1 * time.Hour)

	tests := []struct {
		name string
		user domain.User
		want bool
	}{
		{"No penalty score", domain.User{PenaltyScore: 0}, false},
		{"Penalty with future expiration", domain.User{PenaltyScore: 1, PenaltyScoreExpireAt: &future}, true},
		{"Penalty with past expiration", domain.User{PenaltyScore: 1, PenaltyScoreExpireAt: &past}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.user.HasActivePenalty())
		})
	}
}
