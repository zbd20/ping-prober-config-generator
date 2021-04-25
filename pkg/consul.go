package pkg

import (
	"log"
	"strings"

	"github.com/hashicorp/consul/api"
)

func getConsulClient() (*api.Client, error) {
	cfg := GetConfig()
	cCfg := &api.Config{
		Address: cfg.Consul.Address,
		Scheme:  cfg.Consul.Scheme,
	}

	clt, err := api.NewClient(cCfg)
	if err != nil {
		return nil, err
	}

	return clt, nil
}

func SetConsulKV(key string, val []byte) error {
	cdb, err := getConsulClient()
	if err != nil {
		log.Println("get consul client error:", err)
		return err
	}
	kv, _, err := cdb.KV().Get(key, &api.QueryOptions{})
	if err != nil {
		log.Println("get key error:", err)
		return err
	}

	if kv != nil {
		if strings.Compare(string(kv.Value), string(val)) == 0 {
			log.Println(key, "values do not change")
			return nil
		}
	}

	log.Println("KV pair changed, start reset the value")
	if _, err := cdb.KV().Put(&api.KVPair{
		Key:   key,
		Value: val,
	}, &api.WriteOptions{}); err != nil {
		log.Println("set value error:", err)
		return err
	}

	return nil
}
