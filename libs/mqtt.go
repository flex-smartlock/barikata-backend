package libs

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

func MQTT_Test() {
	mqtt_client := generateMqttClient()

	sub(mqtt_client)
	publish(mqtt_client)

	mqtt_client.Disconnect(250)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func publish(client mqtt.Client) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func sub(client mqtt.Client) {
	topic := "topic/test"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}

func generateMqttClient() mqtt.Client {
	var broker = os.Getenv("SERVER_MQTT_BROKER")
	var clientID = os.Getenv("SERVER_MQTT_CLIENT_ID") + fmt.Sprintf("%d", rand.Intn(10000000000000))
	var username = os.Getenv("SERVER_MQTT_USERNAME")
	var password = os.Getenv("SERVER_MQTT_PASSWORD")
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	mqtt_client := mqtt.NewClient(opts)
	if token := mqtt_client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return mqtt_client
}

func init() {
	godotenv.Load()
}
