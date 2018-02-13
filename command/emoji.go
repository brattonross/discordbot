package command

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/disintegration/imaging"
	log "github.com/sirupsen/logrus"
)

// Merges two images; an emoji and a base image.
// x and y represent the point at which the emoji should be drawn.
// w and h represent the width and height that the emoji should be resized to.
func mergeEmojiImage(emoji, baseImage io.Reader, x, y, w, h int) (io.Reader, error) {
	ei, err := png.Decode(emoji)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to decode custom emoji image")
		return nil, err
	}
	ei = imaging.Resize(ei, w, h, imaging.Lanczos)
	bi, err := png.Decode(baseImage)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to decode image")
		return nil, err
	}
	bi = imaging.FlipH(bi)
	rgba := image.NewRGBA(image.Rectangle{
		image.ZP,
		image.Point{bi.Bounds().Dx(), bi.Bounds().Dy()},
	})
	draw.Draw(rgba, bi.Bounds(), bi, image.ZP, draw.Src)
	draw.Draw(rgba, ei.Bounds().Add(image.Pt(x, y)), ei, image.ZP, draw.Over)
	buf := &bytes.Buffer{}
	err = png.Encode(buf, rgba)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to encode image")
		return nil, err
	}
	return buf, nil
}

// Retrieves a custom emoji full size image from Discord CDN.
func getCustomEmoji(eID string) (io.Reader, error) {
	url := discordgo.EndpointEmoji(eID)
	log.Info("performing GET on url: ", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Info("GET successful for ID: ", eID)
	return bytes.NewReader(b), nil
}
