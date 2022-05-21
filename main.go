package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dhowden/tag"
	_ "github.com/joho/godotenv/autoload"
	"github.com/parnurzeal/gorequest"
)

type Song struct {
	File string `json:"file,omitempty"`
}

func main() {
	np, err := http.Get("http://localhost:" + os.Getenv("PORT") + "/NP")
	if err != nil {
		fmt.Println(err)
	}
	defer np.Body.Close()
	npBody, _ := io.ReadAll(np.Body)
	var n Song
	json.Unmarshal(npBody, &n)
	f, _ := os.Open(n.File)
	music, _ := tag.ReadFrom(f)
	fmt.Printf("%s - %s\n", music.Artist(), music.Album())
	b64, md5 := encode(music.Picture().Data)
	fmt.Println(md5)
	b64img := "data:image/jpeg;base64," + b64
	payload := fmt.Sprintf(`{"name": "%s", "image": "%s", "type": "1"}`, md5, b64img)
	req := gorequest.New()
	_, body, _ := req.Post("https://discord.com/api/v9/oauth2/applications/"+os.Getenv("APP_ID")+"/assets").
		Set("Authorization", os.Getenv("TOKEN")).
		Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36").
		Set("origin", "discordapp.com").
		Set("cache-control", "no-cache").
		Type("json").
		Send(payload).
		End()
	fmt.Println(body)
}

func encode(data []byte) (string, string) {
	b64 := base64.StdEncoding.EncodeToString(data)
	md5sum := fmt.Sprintf("%x", md5.Sum([]byte(b64)))
	return b64, md5sum
}
