package handlers

import (
	"github.com/kubektl/v0-blog-backend/services"
)

type Handlers struct {
	Auth    *AuthHandler
	User    *UserHandler
	Article *ArticleHandler
	Social  *SocialHandler
	Team    *TeamHandler
	Admin   *AdminHandler
}

func NewHandlers(
	auth *services.AuthService,
	users *services.UserService,
	articles *services.ArticleService,
	social *services.SocialService,
	teams *services.TeamService,
) *Handlers {
	return &Handlers{
		Auth:    NewAuthHandler(auth, users),
		User:    NewUserHandler(users, social, articles),
		Article: NewArticleHandler(articles),
		Social:  NewSocialHandler(social),
		Team:    NewTeamHandler(teams),
		Admin:   NewAdminHandler(teams),
	}
}
