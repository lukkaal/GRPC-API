package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/resolver"
)

// single addr
type ServerNode struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`
	Version string `json:"version"`
	Weight  string `json:"weight"`
}

// exp: /user/v1/ 或者 /user/
func BuildPrefix(server ServerNode) string {
	if server.Version == "" {
		return fmt.Sprintf("/%s/", server.Name)
	}
	return fmt.Sprintf("/%s/%s/", server.Name, server.Version)
}

// exp: /user/v1/127.0.0.1:8080 或者 /user/127.0.0.1:8080
func BuildRegisterPath(server ServerNode) string {
	return fmt.Sprintf("%s%s", BuildPrefix(server), server.Addr)
}

// parse from JSON (from etcd) to get serverinfo
func ParseValue(value []byte) (ServerNode, error) {
	server := ServerNode{}
	if err := json.Unmarshal(value, &server); err != nil {
		return server, err
	}
	return server, nil
}

// etcd key string -> ServerNode
func SplitPath(path string) (ServerNode, error) {
	server := &ServerNode{}
	strs := strings.Split(path, "/")
	if len(strs) == 0 {
		return *server, errors.New("invalid path")
	}

	server.Addr = strs[len(strs)-1]
	return *server, nil
}

func Exist(l []resolver.Address,
	addr resolver.Address) bool {
	for i := range l {
		if l[i].Addr == addr.Addr {
			return true
		}
	}
	return false
}

func Remove(s []resolver.Address,
	addr resolver.Address) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr.Addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

// define own resolver (exp: etcd:///user-service)
func BuildResolverUrl(app string) string {
	return schema + ":///" + app
}
