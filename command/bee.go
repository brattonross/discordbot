package command

import (
	"bytes"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
)

const baseBee = "/home/ross/Pictures/headlessbee.png"

// Bee command
type Bee struct{}

// Execute ...
func (b *Bee) Execute(ctx *Context) {
	if len(ctx.Args) < 1 {
		return
	}

	emote := ctx.Args[0]
	i := strings.LastIndex(emote, ":")
	if i == -1 {
		log.WithFields(log.Fields{
			"argument": emote,
		}).Warn("command called with invalid argument")
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Sorry, I can currently only use bee command with a custom emoji!")
		return
	}
	eID := emote[i+1 : len(emote)-1]
	er, err := getCustomEmoji(eID)
	if err != nil {
		log.WithFields(log.Fields{
			"emojiID": eID,
			"error":   err,
		}).Error("failed to get custom emoji")
		return
	}
	bs, err := ioutil.ReadFile(baseBee)
	if err != nil {
		log.WithFields(log.Fields{
			"path":  baseBee,
			"error": err,
		}).Error("failed to read file")
	}
	br := bytes.NewReader(bs)
	// Flip the base bee image if the command was called with a backslash.
	if ctx.Prefix == "\\" {

	}
	img, err := mergeEmojiImage(er, br, 260, 25, 512, 512)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"emojiID": eID,
			"image":   baseBee,
		}).Error("failed to merge emoji with image")
		return
	}
	_, err = ctx.Session.ChannelFileSend(ctx.Message.ChannelID, "bee.png", img)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to send message")
		return
	}
}
