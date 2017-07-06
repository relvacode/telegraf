package deluge

import (
	"github.com/relvacode/go-libdeluge"
	"testing"
)

func TestDeluge(t *testing.T) {
	d := &Deluge{c: delugeclient.New(delugeclient.Settings{
		Hostname: "localhost",
		Port:     58846,
		Login:    "localclient",
		Password: "2be16cfc0c993a1f1727101aaf47556eacb20029",
	})}
	if err := d.c.Connect(); err != nil {
		t.Fatal(err)
	}
	methods, err := d.c.MethodsList()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(methods)
}
