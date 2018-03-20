package wsmap

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hdsbrute"
)

const setMapCmd = "setmap"

// SetMapCommand ...
var SetMapCommand = hdsbrute.Command{
	Cmd: setMapCmd,
	HelpFunc: func(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, getHelpMessage(b))
	},
	Exec: setMapHandleFunc,
}

// handlerFunc answers calls to map and map [coord|color] message
func setMapHandleFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	if m.Attachments == nil || len(m.Attachments) == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No picture found! Please attach one to this message.\n%s", getHelpMessage(b)))
		return
	}

	if query == nil || len(query) == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Don't know what to do with the attachment!\n\n%s", getHelpMessage(b)))
		return
	}

	var picType = query[0]

	if picType != "labels" && picType != "screenshot" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unknown picture type: %s! Use *labels* or *screenshot*!", picType))
		return
	}

	err := sendPictureToBackend(picType, m.Attachments[0].URL)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to update picture! %v", err))
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("STUB: picture set as [%s]", query[0]))
}

func sendPictureToBackend(picType, picUrl string) error {
	// download the pic from discord's cdn
	resp, err := http.Get(picUrl)
	if err != nil {
		return fmt.Errorf("failed to download picture. %v", err)
	}
	defer resp.Body.Close()

	// copy the content to the request to the backend
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	defer writer.Close()

	part, err := writer.CreateFormFile(picType, picType)
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, resp.Body); err != nil {
		return fmt.Errorf("failed to add picture to backend request. %v", err)
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/map?secret=%s", backendURL, backendSecret)

	if resp, err := http.Post(url, writer.FormDataContentType(), body); err != nil || resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send picture to backend. %s %v", resp.Status, err)
	}

	return nil
}

func getHelpMessage(b *hdsbrute.Brute) string {
	var msg = []string{
		fmt.Sprintf("**Description**:"),
		fmt.Sprintf("Updates the `labels` or map `screenshot` layer for the `%s%s` command.", b.Prefix, mapCmd),
		fmt.Sprintf(""),
		fmt.Sprintf("**Usage**:"),
		fmt.Sprintf("Attach a picture to a message and set its text as follows:"),
		fmt.Sprintf(""),
		fmt.Sprintf("`%s%s [labels|screenshot]`", b.Prefix, mapCmd),
		fmt.Sprintf(""),
		fmt.Sprintf("**labels** is the layer with planet names and levels. It needs to be a transparent `.png` file"),
		fmt.Sprintf("**screenshot** is the screenshot of the WS cropped to the edges of the outer hexes. Format must be `.jpg`"),
	}

	return strings.Join(msg, "\n")
}
