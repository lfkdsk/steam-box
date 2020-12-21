package main

import (
	// "bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/YouEclipse/steam-box/pkg/steambox"
	"github.com/google/go-github/github"
)

func main() {
	var err error
	steamAPIKey := os.Getenv("STEAM_API_KEY")
	steamID, _ := strconv.ParseUint(os.Getenv("STEAM_ID"), 10, 64)
	appIDs := os.Getenv("APP_ID")
	appIDList := make([]uint32, 0)

	for _, appID := range strings.Split(appIDs, ",") {
		appid, err := strconv.ParseUint(appID, 10, 32)
		if err != nil {
			continue
		}
		appIDList = append(appIDList, uint32(appid))
	}

	ghToken := os.Getenv("GH_TOKEN")
	ghUsername := os.Getenv("GH_USER")
	gistID := os.Getenv("GIST_ID")

	steamOption := "ALLTIME" // options for types of games to list: RECENT (recently played games), ALLTIME <default> (playtime of games in descending order)
	if os.Getenv("STEAM_OPTION") != "" {
		steamOption = os.Getenv("STEAM_OPTION")
	}

	multiLined := false // boolean for whether hours should have their own line - YES = true, NO = false
	if os.Getenv("MULTILINE") != "" {
		lineOption := os.Getenv("MULTILINE")
		if lineOption == "YES" {
			multiLined = true
		}
	}
	
	// updateOption := os.Getenv("UPDATE_OPTION") // options for update: GIST (Gist only), MARKDOWN (README only), GIST_AND_MARKDOWN (Gist and README)
	// markdownFile := os.Getenv("MARKDOWN_FILE") // the markdown filename (e.g. MYFILE.md)

	var updateGist, updateMarkdown bool
	if updateOption == "MARKDOWN" {
		updateMarkdown = true
	} else if updateOption == "GIST_AND_MARKDOWN" {
		updateGist = true
		updateMarkdown = true
	} else {
		updateGist = true
	}

	box := steambox.NewBox(steamAPIKey, ghUsername, ghToken)

	fmt.Println(steamAPIKey)
	fmt.Println(ghUsername)
	fmt.Println(ghToken)
	fmt.Println(steamID)	
	
	ctx := context.Background()

	var (
		filename string
		lines []string
	)

	fmt.Println(steamOption)
	
	if steamOption == "ALLTIME" {
		filename = "🎮 Steam playtime leaderboard"
		lines, err = box.GetPlayTime(ctx, steamID, multiLined, appIDList...)
		if err != nil {
			panic("GetPlayTime err:" + err.Error())
		}
		fmt.Println(lines)
	} else if steamOption == "RECENT" {
		filename = "🎮 Recently played Steam games"
		lines, err = box.GetRecentGames(ctx, steamID, multiLined)
		if err != nil {
			panic("GetRecentGames err:" + err.Error())
		}
		fmt.Println(lines)
	}
	
	if updateGist {
		gist, err := box.GetGist(ctx, gistID)
		if err != nil {
			panic("GetGist err:" + err.Error())
		}

		f := gist.Files[github.GistFilename(filename)]

		f.Content = github.String(strings.Join(lines, "\n"))

		err = box.UpdateGist(ctx, gistID, gist)
		
		if err != nil {
			panic("UpdateGist err:" + err.Error())
		}
	}
}
