# Plex Input plugin

This input plugin will measures session statistics from a Plex Media Server

### Configuration:

For remote access enabled servers it is usually required to connect over https.
InecureSkipVerify is usually needed as Plex signs certificates valid for their own hostname i.e `https://10-42-1-1.7j38cgj2euhr2fh6f14d30e048d7gh4.plex.direct:32400`.

```
[[inputs.plex]]
  ## URL is the http or https address of the Plex media server
  URL = "http://localhost:32400"

  ## X-Plex-Token
  ## See here for how to get one:
  ## https://support.plex.tv/hc/en-us/articles/204059436-Finding-an-authentication-token-X-Plex-Token
  Token = ""

  ## Optional SSL config
  # InsecureSkipVerify = true
```

### Measurements & Fields:

- bandwidth (from `Session.Bandwidth`)
- transcode_speed (from `TranscodeSession.Speed`)
- video_sessions (length of `Video` sessions)
- track_sessions (length of `Track` sessions)

### Tags:

- platform (from `Player.Platform`)
- device  (from `Player.Device`)
- user (from `User.Title`)
- video_codec (from `Media.VideoCodec`)
- audio_codec (from `Media.AudioCodec`)
- media_type (from `Type`)
- resolution (from `Media.VideoResolution`)

### Example Output:

```
$ ./telegraf --config telegraf.conf --input-filter ping --test
* Plugin: plex, Collection 1
plex,device=Windows,user=User1,video_codec=h264,audio_codec=aac,media_type=episode,resolution=480p,host=hostname.local,platform=Chrome bandwidth=1981i,transcode_speed=0 1499447110000000000
plex,media_type=episode,resolution=480p,host=hostname.local,platform=Chromecast,device=Chromecast,user=User2,video_codec=h264,audio_codec=aac bandwidth=1298i,transcode_speed=0 1499447110000000000
plex,host=hostname.local video_sessions=2i,track_sessions=0i 1499447110000000000
```