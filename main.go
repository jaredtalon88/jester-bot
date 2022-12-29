package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

const StableDiffURL = "http://127.0.0.1:7860"

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

type Gopher struct {
	Name string `json: "name"`
}

type StableDiff struct {
	Images     []string `json:"images"`
	Parameters struct {
		EnableHr                          bool        `json:"enable_hr"`
		DenoisingStrength                 int         `json:"denoising_strength"`
		FirstphaseWidth                   int         `json:"firstphase_width"`
		FirstphaseHeight                  int         `json:"firstphase_height"`
		Prompt                            string      `json:"prompt"`
		Styles                            interface{} `json:"styles"`
		Seed                              int         `json:"seed"`
		Subseed                           int         `json:"subseed"`
		SubseedStrength                   int         `json:"subseed_strength"`
		SeedResizeFromH                   int         `json:"seed_resize_from_h"`
		SeedResizeFromW                   int         `json:"seed_resize_from_w"`
		SamplerName                       interface{} `json:"sampler_name"`
		BatchSize                         int         `json:"batch_size"`
		NIter                             int         `json:"n_iter"`
		Steps                             int         `json:"steps"`
		CfgScale                          float64     `json:"cfg_scale"`
		Width                             int         `json:"width"`
		Height                            int         `json:"height"`
		RestoreFaces                      bool        `json:"restore_faces"`
		Tiling                            bool        `json:"tiling"`
		NegativePrompt                    interface{} `json:"negative_prompt"`
		Eta                               interface{} `json:"eta"`
		SChurn                            float64     `json:"s_churn"`
		STmax                             interface{} `json:"s_tmax"`
		STmin                             float64     `json:"s_tmin"`
		SNoise                            float64     `json:"s_noise"`
		OverrideSettings                  interface{} `json:"override_settings"`
		OverrideSettingsRestoreAfterwards bool        `json:"override_settings_restore_afterwards"`
		SamplerIndex                      string      `json:"sampler_index"`
	} `json:"parameters"`
	Info string `json:"info"`
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!hello" {

		// check on da bot

		_, err := s.ChannelMessageSend(m.ChannelID, "i am alive")
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.HasPrefix(m.Content, "!img") {

		// parrot back our command

		prompt := strings.TrimPrefix(m.Content, "!img ")

		values := map[string]string{"prompt": prompt, "steps": "30"}

		json_data, err := json.Marshal(values)

		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Post(StableDiffURL+"/sdapi/v1/txt2img", "application/json",
			bytes.NewBuffer(json_data))

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				log.Fatal(err)
			}

			var result StableDiff
			if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
				fmt.Println("Can not unmarshal JSON")
			}

			raw_imgdata := fmt.Sprintf(result.Images[0])

			// reader := strings.NewReader(imgdata)

			decoded_img_data := base64.NewDecoder(base64.StdEncoding, strings.NewReader(raw_imgdata))

			_, err = s.ChannelFileSend(m.ChannelID, "image.png", decoded_img_data)
			if err != nil {
				fmt.Println(err)
			}

		} else {
			fmt.Println("Error: Something went wrong")
		}
	}

}
