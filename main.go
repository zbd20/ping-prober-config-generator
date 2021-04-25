package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"

	"gitlab.ushareit.me/sgt/hawkeye/ping-monitor-config-generator/pkg"
)

var (
	cfgFile = flag.String("config", "/etc/config/config.prod.yaml", "config path")
	pingKey = "network-monitor/prometheus/ping-probe.yml"
)

type probeInfo struct {
	Address string
	Port    int
	Tags    map[string]string
}

type renderInfo struct {
	Name    string
	Probes  []probe
	Address string
}

type probe struct {
	Targets []string
	Labels  map[string]string
}

func main() {
	flag.Parse()

	if err := pkg.InitConfig(*cfgFile); err != nil {
		log.Println("get config error:", err)
		os.Exit(1)
	}

	watchService()
}

func watchService() {
	wc := make(map[string]interface{})
	wc["type"] = "service"
	wc["service"] = "ping-probe"
	wp, err := watch.Parse(wc)
	if err != nil {
		log.Println(err)
	}
	wp.Handler = func(idx uint64, result interface{}) {
		services := result.([]*api.ServiceEntry)
		if len(services) != 0 {
			var allProbe []probeInfo
			for _, svc := range services {
				var tmpProbe probeInfo
				tmpProbe.Address = svc.Service.Address
				tmpProbe.Port = svc.Service.Port
				tags := make(map[string]string)
				for _, tagStr := range svc.Service.Tags {
					kv := strings.Split(tagStr, "=")
					tags[kv[0]] = kv[1]
				}
				tmpProbe.Tags = tags
				allProbe = append(allProbe, tmpProbe)
			}

			var renders []renderInfo
			for idx1, sourceProbe := range allProbe {
				var tmpRender renderInfo
				tmpRender.Name = fmt.Sprintf("%v_%v_%v",
					sourceProbe.Tags["account"],
					sourceProbe.Tags["region"],
					sourceProbe.Tags["vpc"],
				)
				tmpRender.Address = fmt.Sprintf("%v:%v", sourceProbe.Address, sourceProbe.Port)
				for idx2, targetProbe := range allProbe {
					if idx1 != idx2 {
						var tmpProbe probe
						tmpProbe.Targets = []string{targetProbe.Address}
						tmpProbe.Labels = map[string]string{
							"source_region": sourceProbe.Tags["region"],
							"source_vpc":    sourceProbe.Tags["vpc"],
							"target_region": targetProbe.Tags["region"],
							"target_vpc":    targetProbe.Tags["vpc"],
						}

						tmpRender.Probes = append(tmpRender.Probes, tmpProbe)
					}
				}

				renders = append(renders, tmpRender)
			}

			result, err := pkg.ProberRender(renders)
			if err != nil {
				log.Println(err)
			}

			if err := pkg.SetConsulKV(pingKey, []byte(result)); err != nil {
				log.Println(err)
			}
		}

	}

	if err := wp.Run(pkg.GetConfig().Consul.Address); err != nil {
		log.Println(err)
	}
}
