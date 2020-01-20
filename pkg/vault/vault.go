package vault

import (
	path_helper "path"
	"strings"

	"github.com/pkg/errors"

	"github.com/hashicorp/vault/api"
)

func Walk(client *api.Client, path string) ([]string, error) {

	s, err := client.Logical().List(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error listing path : %s", path)
	}
	key, found := s.Data["keys"]

	if !found {
		return nil, errors.New("no keys found")
	}

	keyArr, ok := key.([]interface{})
	if !ok {
		return nil, errors.New("error - the 'key' isn't an []interface{}")
	}

	ret := []string{}

	for _, p := range keyArr {

		pStr := p.(string)

		fullPath := path_helper.Join(path, pStr)

		if strings.HasSuffix(pStr, "/") {

			r, err := Walk(client, fullPath)

			if err != nil {
				return nil, errors.Wrap(err, "error walking path")
			}

			ret = append(ret, r...)

		} else {

			ret = append(ret, fullPath)
		}

	}
	return ret, nil

}
