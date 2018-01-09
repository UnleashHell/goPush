package push

import (
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"goPush/lib/config"
	"goPush/lib/log"
	"runtime"
)

type Ios struct {
	Queue      chan *Message
	ClientPool chan *apns2.Client
}

var IosInstance *Ios
var topic string

func (this *Ios) InitClient() {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	channelNum := config.Instance.GetInt("worker", "iosChannelNum")
	this.Queue = make(chan *Message, channelNum)
	this.ClientPool = make(chan *apns2.Client, cores)
	pemPath := config.Instance.Get("app", "pemPath")
	pemPass := config.Instance.Get("app", "pemPass")
	topic = config.Instance.Get("app", "package")
	cert, err := certificate.FromPemFile(pemPath, pemPass)
	if err != nil {
		panic(err)
	}
	//每个cpu核心创建一个client
	for i := 0; i < cores; i++ {
		client := apns2.NewClient(cert).Production()
		this.ClientPool <- client
	}
	workerNum := config.Instance.GetInt("worker", "iosWorkerNum")
	for i := 0; i < workerNum; i++ {
		go this.worker()
	}
}

func (this *Ios) Push(model *Message) bool {
	select {
	case this.Queue <- model:
		return true
	default:
		return false
	}
}

func (this *Ios) worker() {
	for message := range this.Queue {
		client := <-this.ClientPool
		data := payload.NewPayload().Alert(message.Alert).Badge(message.Badge).
			Sound(message.Sound)
		notification := &apns2.Notification{}
		notification.DeviceToken = message.Token
		notification.Payload = data
		var res *apns2.Response
		var err error
		notification.Topic = topic
		res, err = client.Push(notification)
		if err != nil {
			log.Instance.Error("send error %v", err)
		} else {
			if res.Sent() {
				log.Instance.Info("message:%s token:%s ",
					message.Alert, message.Token)
			} else {
				log.Instance.Error("message:%s token:%s statusCode:%d reason:%s",
					message.Alert, message.Token, res.StatusCode, res.Reason)
			}
		}
		this.ClientPool <- client
	}
}
