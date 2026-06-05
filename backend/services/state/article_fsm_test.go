package state_test

import (
	"testing"

	"github.com/kubektl/v0-blog-backend/apperrors"
	"github.com/kubektl/v0-blog-backend/models"
	"github.com/kubektl/v0-blog-backend/services/state"
)

func TestArticleTransitions(t *testing.T) {
	cases := []struct {
		from, to string
		ok       bool
	}{
		{models.ArticleStatusDraft, models.ArticleStatusPublished, true},
		{models.ArticleStatusDraft, models.ArticleStatusArchived, true},
		{models.ArticleStatusPublished, models.ArticleStatusDraft, true},
		{models.ArticleStatusPublished, models.ArticleStatusArchived, true},
		{models.ArticleStatusArchived, models.ArticleStatusDraft, true},
		{models.ArticleStatusDraft, models.ArticleStatusDraft, true},
		{models.ArticleStatusPublished, models.ArticleStatusPublished, true},
		{models.ArticleStatusArchived, models.ArticleStatusPublished, false},
		{models.ArticleStatusDraft, "invalid", false},
	}

	for _, tc := range cases {
		err := state.ValidateArticleTransition(tc.from, tc.to)
		if tc.ok && err != nil {
			t.Fatalf("expected %s -> %s allowed, got %v", tc.from, tc.to, err)
		}
		if !tc.ok && err != apperrors.ErrInvalidTransition {
			t.Fatalf("expected %s -> %s rejected, got %v", tc.from, tc.to, err)
		}
	}
}

func TestVerificationTransitions(t *testing.T) {
	cases := []struct {
		from, to string
		ok       bool
	}{
		{"", models.VerificationStatusPending, true},
		{models.VerificationStatusPending, models.VerificationStatusApproved, true},
		{models.VerificationStatusPending, models.VerificationStatusRejected, true},
		{models.VerificationStatusRejected, models.VerificationStatusPending, true},
		{models.VerificationStatusApproved, models.VerificationStatusPending, false},
		{models.VerificationStatusPending, models.VerificationStatusPending, true},
	}

	for _, tc := range cases {
		err := state.ValidateVerificationTransition(tc.from, tc.to)
		if tc.ok && err != nil {
			t.Fatalf("expected %s -> %s allowed, got %v", tc.from, tc.to, err)
		}
		if !tc.ok && err != apperrors.ErrInvalidTransition {
			t.Fatalf("expected %s -> %s rejected, got %v", tc.from, tc.to, err)
		}
	}
}
