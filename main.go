package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// 消去を行う際に必要な項目です。適切に埋めてください。
const (
	UserID             = "NEKO"
	TimelineChannel    = "home"
	APIToken           = "4YlfkvEIObF2ftHv8hVY2k5Ltv6bYhSE"
	MisskeyHost        = "msk.ilnk.info"
	RequestIntervalSec = 15 //DOSまがいの負荷をかけないように遅延秒数を指定する。
)

func main() {
	for {
		postIds := getPostIDs()
		if len(postIds) == 0 {
			break
		}

		for _, id := range postIds {
			deletePost(id)
			time.Sleep(time.Second * RequestIntervalSec)
		}
	}

	fmt.Println("処理が終了しました。")
}

func postHTTPRequest(endpoint string, data interface{}) ([]byte, error) {
	url := "https://" + MisskeyHost + endpoint
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		fmt.Println("何らかの通信エラーが発生しました。")
		fmt.Printf("HTTPステータス: %d\n", resp.StatusCode)
		fmt.Printf("%s\n", respBody)
		return nil, fmt.Errorf("HTTP error")
	}

	return respBody, nil
}

func deletePost(id string) {
	data := map[string]interface{}{
		"i":      APIToken,
		"noteId": id,
	}
	_, err := postHTTPRequest("/api/notes/delete", data)
	if err != nil {
		fmt.Printf("投稿の消去時に問題が発生しました %s: %v\n", id, err)
		return
	}

	timestamp := time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf("deleted %s -- %s\n", id, timestamp)
}

// nsfwな投稿を残す場合はexcludeNsfwをfalseにする
func getPostIDs() []string {
	data := map[string]interface{}{
		"i":           TimelineChannel,
		"excludeNsfw": true,
		"limit":       100,
		"userId":      UserID,
	}
	respBody, err := postHTTPRequest("/api/users/notes", data)
	if err != nil {
		fmt.Println("投稿IDを取得できませんでした。IDs:", err)
		return nil
	}

	var respData []struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		fmt.Println("投稿IDをパース出来ませんでした。 IDs:", err)
		return nil
	}

	count := len(respData)
	fmt.Printf("消去された投稿数[ %d ] = >>\n", count)

	var postIDs []string
	for _, item := range respData {
		postIDs = append(postIDs, item.ID)
	}

	return postIDs
}
