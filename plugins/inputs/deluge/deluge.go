// package deluge implements a telegraf input for the deluge torrent daemon
package deluge

import (
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/relvacode/go-libdeluge"
)

type RPCMethod func(*delugeclient.Client, telegraf.Accumulator) error

type Deluge struct {
	Hostname string
	Port     uint
	Login    string
	Password string

	// methods are the list of used RPC methods to accumulate
	methods []RPCMethod
}

func (Deluge) Description() string {
	return "Read stats from the Deluge torrent daemon"
}

func (Deluge) SampleConfig() string {
	return `
  ## IP or hostname of the deluge daemon
  Hostname = "localhost"
  
  ## The port number of the deluge daemon
  Port = 58846

  ## Login and Password are optional but required for most deluge daemon configuration
  ## Authentication credentials can be found in $DELUGE_HOME/.config/deluge/auth
  # Login = "localclient"
  # Password = "password"
	`
}

func (d *Deluge) Settings() delugeclient.Settings {
	return delugeclient.Settings{
		Hostname:         d.Hostname,
		Port:             d.Port,
		Login:            d.Login,
		Password:         d.Password,
		ReadWriteTimeout: time.Second * 5,
	}
}

func (d *Deluge) Gather(acc telegraf.Accumulator) error {
	c := delugeclient.New(d.Settings())
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Close()

	for _, m := range d.methods {
		if err := m(c, acc); err != nil {
			return err
		}
	}
	return nil
}

func GetSessionStatus() RPCMethod {
	keys := []string{"upload_rate",
		"download_rate",
		"payload_upload_rate",
		"payload_download_rate",
		"dht_upload_rate",
		"dht_download_rate",
		"tracker_upload_rate",
		"tracker_download_rate",
		"total_redundant_bytes",
		"total_failed_bytes",
		"total_download",
		"total_upload",
		"num_peers",
		"up_bandwidth_queue",
		"down_bandwidth_queue",
		"dht_nodes",
	}
	fields := map[string]interface{}{}
	tags := map[string]string{}
	return func(c *delugeclient.Client, acc telegraf.Accumulator) error {
		s, err := c.SessionStatus(keys)
		if err != nil {
			return err
		}
		for k, v := range s {
			fields[k] = v
		}
		acc.AddFields("deluge", fields, tags)
		return nil
	}
}

func GetFreeSpace() RPCMethod {
	fields := map[string]interface{}{}
	return func(c *delugeclient.Client, acc telegraf.Accumulator) error {
		free, err := c.FreeSpace()
		if err != nil {
			return err
		}
		fields["free_space"] = free
		acc.AddFields("deluge", fields, nil)
		return nil
	}
}

func init() {
	inputs.Add("deluge", func() telegraf.Input {
		return &Deluge{
			methods: []RPCMethod{
				GetSessionStatus(),
				GetFreeSpace(),
			},
		}
	})
}
