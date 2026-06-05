package state

import (
	"github.com/kubektl/v0-blog-backend/apperrors"
	"github.com/kubektl/v0-blog-backend/models"
)

var articleTransitions = map[string]map[string]bool{
	models.ArticleStatusDraft: {
		models.ArticleStatusPublished: true,
		models.ArticleStatusArchived:  true,
	},
	models.ArticleStatusPublished: {
		models.ArticleStatusDraft:    true,
		models.ArticleStatusArchived: true,
	},
	models.ArticleStatusArchived: {
		models.ArticleStatusDraft: true,
	},
}

func CanTransitionArticle(from, to string) bool {
	if from == to {
		return true
	}
	targets, ok := articleTransitions[from]
	if !ok {
		return false
	}
	return targets[to]
}

func ValidateArticleTransition(from, to string) error {
	if CanTransitionArticle(from, to) {
		return nil
	}
	return apperrors.ErrInvalidTransition
}

var verificationTransitions = map[string]map[string]bool{
	"": {
		models.VerificationStatusPending: true,
	},
	models.VerificationStatusPending: {
		models.VerificationStatusApproved: true,
		models.VerificationStatusRejected: true,
	},
	models.VerificationStatusRejected: {
		models.VerificationStatusPending: true,
	},
}

func CanTransitionVerification(from, to string) bool {
	if from == to {
		return true
	}
	targets, ok := verificationTransitions[from]
	if !ok {
		return false
	}
	return targets[to]
}

func ValidateVerificationTransition(from, to string) error {
	if CanTransitionVerification(from, to) {
		return nil
	}
	return apperrors.ErrInvalidTransition
}
