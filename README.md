## ping-prober-config-generator

该组件监听consul service的变化，当watch到service发生变化时，会执行模板文件重新渲染，进而生成Prometheus的配置文件。