package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type HealthPack struct {
	// Props interface{} `json:"props"`
	Data interface{} `json:"data"`
	// Date  string      `json:"date"`
}

type Props struct {
	DeviceName string    `json:"deviceName"`
	MacAddress string    `json:"macAddress"`
	DeviceIp   time.Time `json:"deviceIp"`
}

type Inertias struct {
	FAccX        string `json:"fAccX"`        // :"0.0"
	FAccY        string `json:"fAccY"`        // :"0.0"
	FAccZ        string `json:"fAccZ"`        // :"0.0"
	AccX         string `json:"accX"`         // :"944.000000"
	AccY         string `json:"accY"`         // :"-356.000000"
	AccZ         string `json:"accZ"`         // :"-16960.000000"
	GyrX         string `json:"gyrX"`         // :"-491.000000"
	GyrY         string `json:"gyrY"`         // :"15.000000"
	GyrZ         string `json:"gyrZ"`         // :"196.000000"
	ContimpactoX string `json:"contimpactoX"` // :"944.000000"
	ContimpactoY string `json:"contimpactoY"` // :"412.000000"
	ContimpactoZ string `json:"contimpactoZ"` // :"17968.000000"
	Pitch        string `json:"pitch"`        // :"-1.200725"
	Roll         string `json:"roll"`         // :"-3.185352"
	Yaw          string `json:"yaw"`          // :"0.000000"}
}

type Tracking struct {
	Latitude               string `json:"latitude"`               // :"0.000000"
	Longitude              string `json:"longitude"`              // :"0.000000"
	Tempbateriasecundaria  string `json:"tempbateriasecundaria"`  // :"24.2"
	Tempbateriaprincipal   string `json:"tempbateriaprincipal"`   // :"23.6"
	Temperaturacondensador string `json:"temperaturacondensador"` // :"28.1"
	Temperaturacuba1       string `json:"temperaturacuba1"`       // :"-0.9"
	Temperaturacuba2       string `json:"temperaturacuba2"`       // :"11.4"
	TemperaturaexternaLL   string `json:"temperaturaexternaLL"`   // :"36.3"
	TemperaturaexternaLS   string `json:"temperaturaexternaLS"`   // :"8.4"
	Temperaturaexterna     string `json:"temperaturaexterna"`     // :"28.1"
	Temperaturadissipador  string `json:"temperaturadissipador"`  // :"26.5"
	Correntebateria        string `json:"correntebateria"`        // :"0.0"
	Correntecompressor     string `json:"correntecompressor"`     // :"3.9"
	Correntepeltier        string `json:"correntepeltier"`        // :"0.0"
	Correntecooler         string `json:"correntecooler"`         // :"0.0"
	Correnteexaustor       string `json:"correnteexaustor"`       // :"0.0"
	Temperaturacompressor  string `json:"temperaturacompressor"`  // :"44.8"
	Setpoint_pid1          string `json:"setpoint_pid1"`          // :"-5.0"
	Valor_pid1_atual       string `json:"valor_pid1_atual"`       // :"-0.9"
	Esforco_pid1           string `json:"esforco_pid1"`           // :"536.0"
	Setpoint_pid2          string `json:"setpoint_pid2"`          // :"50.0"
	Valor_pid2_atual       string `json:"valor_pid2_atual"`       // :"44.8"
	Esforco_pid2           string `json:"esforco_pid2"`           // :"0.0"}
}

type Status struct {
}

type Ischemia struct {
	IdModal                 string `json:"idModal"`
	IdOperador              string `json:"IdOperador"`              // :""
	Niveldepermissao        string `json:"niveldepermissao"`        // :""
	Nome                    string `json:"nome"`                    // :""
	Numtransplante          string `json:"numtransplante"`          // :""
	Numeroempresa           string `json:"numeroempresa"`           // :""
	Orgao                   string `json:"orgao"`                   // :"tecidos_oculares"
	Tempo_total_isquemia    string `json:"tempo_total_isquemia"`    // :"00/00/00 00:00:00"
	Tempo_restante_isquemia string `json:"tempo_restante_isquemia"` // :"00d 00:00"
	Hora_isquemia           string `json:"hora_isquemia"`           // :""
	Timeinfo_sp2            string `json:"timeinfo_sp2"`            // :"24:12:09 16:18:53"}
}

type Alarm struct {
	Alarms []string `json:"alarms"`
}

func connLostHandler(c MQTT.Client, err error) {
	fmt.Printf("Connection lost, reason: %v\n", err)
	os.Exit(1)
}

func main() {
	id := uuid.New().String()
	var sbMqttSubClientId strings.Builder
	var sbMqttPubClientId strings.Builder
	var sbPubTopic strings.Builder
	sbMqttSubClientId.WriteString("mqtt-topic-rewrite-lns-chirpstackv4-")
	sbMqttSubClientId.WriteString(id)
	sbMqttPubClientId.WriteString("mqtt-topic-rewrite-lns-chirpstackv4-")
	sbMqttPubClientId.WriteString(id)

	mqttSubBroker := "mqtt://smartcampus.maua.br:1883"
	mqttSubClientId := sbMqttSubClientId.String()
	mqttSubUser := "PUBLIC"
	mqttSubPassword := "public"
	mqttSubQos := 0

	mqttSubOpts := MQTT.NewClientOptions()
	mqttSubOpts.AddBroker(mqttSubBroker)
	mqttSubOpts.SetClientID(mqttSubClientId)
	mqttSubOpts.SetUsername(mqttSubUser)
	mqttSubOpts.SetPassword(mqttSubPassword)
	mqttSubOpts.SetConnectionLostHandler(connLostHandler)

	mqttSubTopics := map[string]byte{
		"IMT/SaoRafael/+/+": byte(mqttSubQos),
	}

	mqttPubBroker := "mqtt://mqtt.maua.br:1883"
	mqttPubClientId := sbMqttPubClientId.String()
	mqttPubUser := "public"
	mqttPubPassword := "public"
	mqttPubQos := 0

	mqttPubOpts := MQTT.NewClientOptions()
	mqttPubOpts.AddBroker(mqttPubBroker)
	mqttPubOpts.SetClientID(mqttPubClientId)
	mqttPubOpts.SetUsername(mqttPubUser)
	mqttPubOpts.SetPassword(mqttPubPassword)

	c := make(chan [2]string)

	mqttSubOpts.SetDefaultPublishHandler(func(mqttSubClient MQTT.Client, msg MQTT.Message) {
		c <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	mqttSubClient := MQTT.NewClient(mqttSubOpts)
	if token := mqttSubClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", mqttSubBroker)
	}

	pClient := MQTT.NewClient(mqttPubOpts)
	if token := pClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", mqttPubBroker)
	}

	if token := mqttSubClient.SubscribeMultiple(mqttSubTopics, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for {
		incoming := <-c
		s := strings.Split(incoming[0], "/")
		// IMT/SaoRafael/DEVICE_ID/origin
		var measurement string
		var etc string
		deviceId := s[2]

		// Unmarshal Json to HealthPack struct
		var healthPack HealthPack
		json.Unmarshal([]byte(incoming[1]), &healthPack)

		// Marshal healthPack.Data to Json string
		dataMessage, err := json.Marshal(healthPack.Data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Unmarshal Json to HealthPack struct to dataMap
		var dataMap map[string]interface{}
		json.Unmarshal([]byte(dataMessage), &dataMap)

		if _, exists := dataMap["accX"]; exists {
			measurement = "Inertias"
		}
		if _, exists := dataMap["latitude"]; exists {
			measurement = "Tracking"
		}
		if _, exists := dataMap["vbateriaprincipal"]; exists {
			measurement = "Status"
		}
		if _, exists := dataMap["hora_isquemia"]; exists {
			measurement = "Ischemia"
		}
		if _, exists := dataMap["alarms"]; exists {
			measurement = "Alarms"
		}

		if s[3] == "rx" {
			etc = "wifi"
		} else if s[3] == "rx_gsm" {
			etc = "gsm"
		}

		sbPubTopic.Reset()
		sbPubTopic.WriteString("OpenDataTelemetry/IMT/HealthPack/")
		sbPubTopic.WriteString(measurement)
		sbPubTopic.WriteString("/")
		sbPubTopic.WriteString(deviceId)
		sbPubTopic.WriteString("/up/")
		sbPubTopic.WriteString(etc)

		// fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		token := pClient.Publish(sbPubTopic.String(), byte(mqttPubQos), false, incoming[1])
		token.Wait()

	}
}
