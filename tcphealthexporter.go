package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xeonel2/tcp-shaker"
	"gopkg.in/yaml.v2"
)

type conf struct {
	Services []Service `yaml:"services"`
}

//Service is the strucure of a config entry for a service
type Service struct {
	ServiceName string `yaml:"servicename"`
	ServiceHost string `yaml:"host"`
	ServicePort string `yaml:"port"`
	MetricName  string `yaml:"metricname"`
	Help        string `yaml:"help"`
}

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("tcpservicenames.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		panic(err)
	}
	return c
}

func heartBeat(servicename string, host string, port string) {
	hostname, err := os.Hostname()
	c := tcp.NewChecker(true)
	if err := c.InitChecker(); err != nil {
		log.Fatal("Checker init failed:", err)
	}

	timeout := time.Second * 2
	err = c.CheckAddr(host+":"+port, timeout)
	switch err {
	case tcp.ErrTimeout:
		//timeout
		GaugeMap[servicename].WithLabelValues(hostname).Set(1)

	case nil:
		//success
		GaugeMap[servicename].WithLabelValues(hostname).Set(0)
	default:
		if _, ok := err.(*tcp.ErrConnect); ok {
			// fmt.Println("Connect to "+host+" failed:", e)
			//connect to host failed
			GaugeMap[servicename].WithLabelValues(hostname).Set(1)
		} else {
			//error connecting
			GaugeMap[servicename].WithLabelValues(hostname).Set(1)
			// fmt.Println("Error occurred while connecting:", err)
		}
	}
	time.Sleep(time.Second * 5)
	heartBeat(servicename, host, port)
}

//Configuration object
var con *conf

//GaugeMap is a map of Prometheus Gauges
var GaugeMap = make(map[string]*prometheus.GaugeVec)

func main() {
	con = new(conf)
	con.getConf()
	for _, element := range con.Services {
		GaugeMap[element.ServiceName] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: element.MetricName, Help: element.Help}, []string{"hostname"})
		prometheus.MustRegister(GaugeMap[element.ServiceName])
		go heartBeat(element.ServiceName, element.ServiceHost, element.ServicePort)
	}

	fmt.Println("Starting Http server and listening on port 9112...")
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("starting heartbeat")
	if err := http.ListenAndServe("0.0.0.0:9112", nil); err != nil {
		fmt.Println("Failed to make connection" + err.Error())
	}
}
