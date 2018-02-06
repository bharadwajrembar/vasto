package shell

import (
	"errors"
	"github.com/chrislusf/vasto/client"
	"io"
)

type CommandEnv struct {
	env           map[string]string
	keyspace      string
	dataCenter    string
	clusterClient *client.ClusterClient
}

type Command interface {
	Name() string
	Help() string
	Do(*client.VastoClient, []string, *CommandEnv, io.Writer) error
}

var commands = []Command{}

var InvalidArguments = errors.New("invalid arguments")
