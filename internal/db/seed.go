package db

import (
	"context"
	"database/sql"
	"log"

	"math/rand"

	"github.com/ErickLeal/gopher/internal/store"
	"github.com/brianvoe/gofakeit/v7"
)

func Seed(store store.Storage, db *sql.DB) {
	gofakeit.Seed(0)

	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.UserModel {
	users := make([]*store.UserModel, num)

	for i := 0; i < num; i++ {
		users[i] = &store.UserModel{
			Username: gofakeit.Name(),
			Email:    gofakeit.Email(),
		}
	}

	return users
}

func generatePosts(num int, users []*store.UserModel) []*store.PostModel {
	posts := make([]*store.PostModel, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.PostModel{
			UserID:  user.ID,
			Title:   gofakeit.JobTitle(),
			Content: gofakeit.Company(),
			Tags: []string{
				gofakeit.AppName(),
				gofakeit.Language(),
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.UserModel, posts []*store.PostModel) []*store.CommentModel {
	cms := make([]*store.CommentModel, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.CommentModel{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: gofakeit.Phrase(),
		}
	}
	return cms
}
