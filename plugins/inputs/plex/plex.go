// package plex implements a telegraf input for the plex media server
package plex

import (
	"crypto/tls"
	"net/http"
	"strconv"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/jrudio/go-plex-client"
)

type PlexMediaServer struct {
	URL                string
	Token              string
	InsecureSkipVerify bool
}

func (PlexMediaServer) Description() string {
	return "Reads stats from a Plex media server"
}

func (PlexMediaServer) SampleConfig() string {
	return `
  ## URL is the http or https address of the Plex media server
  URL = "http://localhost:32400"

  ## X-Plex-Token
  ## See here for how to get one:
  ## https://support.plex.tv/hc/en-us/articles/204059436-Finding-an-authentication-token-X-Plex-Token
  Token = ""

  ## Optional SSL config
  # InsecureSkipVerify = true
	`
}

func (pms *PlexMediaServer) transport() http.RoundTripper {
	return &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: pms.InsecureSkipVerify,
		},
	}
}

func (pms *PlexMediaServer) Gather(acc telegraf.Accumulator) error {
	c, err := plex.New(pms.URL, pms.Token)
	if err != nil {
		return err
	}
	c.HTTPClient.Transport = pms.transport()

	sessions, err := c.GetSessions()
	if err != nil {
		return err
	}
	for _, stream := range sessions.Video {
		tags := map[string]string{
			"platform":    stream.Player.Platform,
			"device":      stream.Player.Device,
			"user":        stream.User.Title,
			"video_codec": stream.Media.VideoCodec,
			"audio_codec": stream.Media.AudioCodec,
			"media_type":  stream.Type,
			"resolution":  stream.Media.VideoResolution,
		}
		fields := map[string]interface{}{
			"active": 1,
		}

		i, err := strconv.Atoi(stream.Session.Bandwidth)
		if err == nil {
			fields["bandwidth"] = i
		}
		if stream.TranscodeSession.Key != "" {
			f, err := strconv.ParseFloat(stream.TranscodeSession.Speed, 64)
			if err != nil {
				return err
			}
			fields["transcode_speed"] = f
		}
		acc.AddFields("plex", fields, tags)
	}
	return nil
}

func init() {
	inputs.Add("plex", func() telegraf.Input {
		return new(PlexMediaServer)
	})
}
