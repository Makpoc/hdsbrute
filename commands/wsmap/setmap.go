package wsmap

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hdsbrute"
)

const (
	officerRoleId = "351841217153204245"
	setMapCmd     = "setmap"
)

var adminRoles []string

// SetMapCommand ...
var SetMapCommand = hdsbrute.Command{
	Cmd: []string{setMapCmd},
	Init: func(brute *hdsbrute.Brute) error {
		envRoles, ok := os.LookupEnv("ADMIN_ROLES")
		if ok {
			for _, role := range strings.Split(envRoles, ",") {
				adminRoles = append(adminRoles, role)
			}
		}

		return nil
	},
	HelpFunc: func(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, getHelpMessage(b))
	},
	Exec: setMapHandleFunc,
}

// handlerFunc answers calls to map and map [coord|color] message
func setMapHandleFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	hasPermission, err := hasPermissions(s, m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to query your permissions. %v", err))
		return
	}

	if !hasPermission {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You do not have permissions to use this command. Please contact an <@&%s>", officerRoleId))
		return
	}

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

	if err := sendPictureToBackend(picType, m.Attachments[0].URL); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to update picture! %v", err))
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Picture set as [%s]", query[0]))
}

func sendPictureToBackend(picType, picUrl string) error {
	// download the pic from discord's cdn
	discordResp, err := http.Get(picUrl)
	if err != nil {
		return fmt.Errorf("failed to download picture. %v", err)
	}
	defer discordResp.Body.Close()

	// copy the content to the request to the backend
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	defer writer.Close()

	part, err := writer.CreateFormFile(picType, picType)
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, discordResp.Body); err != nil {
		return fmt.Errorf("failed to add picture to backend request. %v", err)
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/map?secret=%s", backendURL, backendSecret)

	backendResp, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		return fmt.Errorf("failed to send picture to backend. %v", err)
	}
	if backendResp != nil && backendResp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send picture to backend. %s", backendResp.StatusCode)
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
		fmt.Sprintf("`%s%s [labels|screenshot]`", b.Prefix, setMapCmd),
		fmt.Sprintf(""),
		fmt.Sprintf("**labels** is the layer with planet names and levels. It needs to be a transparent `.png` file"),
		fmt.Sprintf("**screenshot** is the screenshot of the WS cropped to the edges of the outer hexes. Format must be `.jpg`"),
	}

	return strings.Join(msg, "\n")
}

func hasPermissions(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	officers, err := hdsbrute.GetMembersByRole(s, m.ChannelID, adminRoles)
	if err != nil {
		return false, err
	}

	for _, officer := range officers {
		if m.Author.Username == officer.User.Username {
			return true, nil
		}
	}

	return false, nil
}
