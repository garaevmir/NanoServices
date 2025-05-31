package scenarios

import (
	"fmt"
	"time"

	"math/rand"

	"github.com/nanoservices/tests/api"
	"github.com/nanoservices/tests/models"
)

func Scenario2(user models.User) {
	token, err := api.LoginUser(user)
	if err != nil {
		fmt.Println("[ ] Login failed:", err)
		return
	}
	fmt.Println("[*] Login successful")

	posts, total, err := api.ListPosts(token)
	if err != nil {
		fmt.Println("[ ] List posts failed:", err)
		return
	}

	fmt.Printf("[*] Found %d posts (total: %d)\n", len(posts), total)

	if len(posts) == 0 {
		fmt.Println("No posts found, creating sample post...")
		post := models.CreatePostRequest{
			Title:       "Sample Post",
			Description: "Default content for testing",
			Tags:        []string{"sample"},
		}
		_, err := api.CreatePost(token, post)
		if err != nil {
			fmt.Println("[ ] Create sample post failed:", err)
			return
		}
		posts, total, err = api.ListPosts(token)
		if err != nil {
			fmt.Println("[ ] List posts failed after creation:", err)
			return
		}
	}

	rand.Seed(time.Now().UnixNano())
	randomPost := posts[rand.Intn(len(posts))]
	fmt.Printf("Selected random post: ID=%s, Title=%s\n", randomPost.ID, randomPost.Title)

	stats, err := api.GetPostStats(token, randomPost.ID)
	if err != nil {
		fmt.Println("[ ] Get stats failed:", err)
		return
	}
	fmt.Printf("[*] Initial stats: Views=%d\n", stats.Views)

	fmt.Println("[...] Simulating post view...")
	if err := api.ViewPost(token, randomPost.ID); err != nil {
		fmt.Println("[ ] View post failed:", err)
		return
	}
	fmt.Println("[*] Post viewed")

	fmt.Println("[...] Waiting for stats update (2s)...")
	time.Sleep(2 * time.Second)

	updatedStats, err := api.GetPostStats(token, randomPost.ID)
	if err != nil {
		fmt.Println("[ ] Get updated stats failed:", err)
		return
	}
	fmt.Printf("[*] Updated stats: Views=%d (delta: +%d)\n",
		updatedStats.Views, updatedStats.Views-stats.Views)

	fmt.Println("[...] Requesting views trend...")
	trendItems, err := api.GetViewsTrend(token, randomPost.ID, "7d")
	if err != nil {
		fmt.Println("[ ] Get views trend failed:", err)
		return
	}

	fmt.Println("[*] Views trend (7 days):")
	totalViews := 0
	for _, item := range trendItems {
		fmt.Printf("  - %s: %d views\n", item.Date, item.Count)
		totalViews += item.Count
	}
	fmt.Printf("Total views in period: %d\n", totalViews)
}
