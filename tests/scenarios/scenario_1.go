package scenarios

import (
	"fmt"
	"time"

	"github.com/nanoservices/tests/api"
	"github.com/nanoservices/tests/models"
)

func Scenario1(user models.User) {
	token, err := api.LoginUser(user)
	if err != nil {
		fmt.Println("[ ] Login failed:", err)
		return
	}
	fmt.Println("[*] Login successful")

	post := models.CreatePostRequest{
		Title:       fmt.Sprintf("Post by %s", user.Username),
		Description: fmt.Sprintf("Content created at %s", time.Now().Format(time.RFC3339)),
		Tags:        []string{"test", "example"},
	}
	postID, err := api.CreatePost(token, post)
	if err != nil {
		fmt.Println("[ ] Create post failed:", err)
		return
	}
	fmt.Printf("[*] Post created: ID=%s\n", postID)

	fmt.Println("[...] Waiting for stats processing (3s)...")
	time.Sleep(3 * time.Second)

	stats, err := api.GetPostStats(token, postID)
	if err != nil {
		fmt.Println("[ ] Get stats failed:", err)
		return
	}
	fmt.Printf("[*] Post stats: Views=%d, Likes=%d\n", stats.Views, stats.Likes)
}
