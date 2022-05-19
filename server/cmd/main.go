package main

import (
	"bufio"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	topic = "quickstart"
)

type KafkaProducer struct {
	*kafka.Producer
	errors chan error
	done   chan struct{}
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	errs := make(chan error, 1)

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})

	if err != nil {
		log.Fatalf("coundnt create producer: %v", err)
	}

	kp := KafkaProducer{p, errs, make(chan struct{})}

	kp.run()

	select {
	case err := <-errs:
		if err != nil {
			log.Printf("server error: %s", err)
		} else {
			log.Println("success")
		}

	case s := <-sig:
		log.Println("signal ", s)

		if err := kp.shutDown(); err != nil {
			log.Println("could not stop server gracefully: %w", err)
		}
	}

	time.Sleep(1 * time.Second)
	fmt.Println("done")
}

func run2() {

	if err := syscall.SetNonblock(0, true); err != nil {
		panic(err)
	}
	f := os.NewFile(0, "stdin")
	go func() {
		time.Sleep(time.Millisecond)
		fmt.Println("setting deadline")
		if err := f.SetDeadline(time.Now().Add(time.Second)); err != nil {
			panic(err)
		}
	}()
	var buf [1]byte
	_, err := f.Read(buf[:])
	if err == nil {
		panic("Read succeeded")
	}
	fmt.Println("read failed as expected:", err)
	if err := syscall.SetNonblock(0, false); err != nil {
		panic(err)
	}

}

func run() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	done := make(chan struct{}, 1)

	go func() {
		defer fmt.Println("Scanner done")

		i := 0
		sc := bufio.NewReader(os.Stdin)

		for {
			select {
			case <-done:
				return
			default:
				input, _, err := sc.ReadLine()
				if err != nil {
					return
				}
				fmt.Printf("msg[%d]: \"%s\"", i, input)
				i++
			}
		}
	}()

	fmt.Println("wait")
	<-sig
	fmt.Println("got sig")
	done <- struct{}{}
	os.Stdin.Close()
	//n, err := io.Copy(os.Stdin, strings.NewReader("close"))
	//fmt.Println(n, err)
	time.Sleep(2 * time.Second)
}
func (kp *KafkaProducer) run() {

	deliveryCh := make(chan kafka.Event, 1)

	// process user input
	go func() {
		defer func() {
			close(deliveryCh)
			log.Println("end process user input")
		}()

		log.Println("ready to read")
		sc := bufio.NewScanner(os.Stdin)
		i := 0

		for {
			select {
			case <-kp.done:
				kp.errors <- nil
				return
			default:
				if sc.Scan() {
					if err := kp.produceMsg(deliveryCh, topic, fmt.Sprintf("msg[%d]: \"%s\"", i, sc.Text()), -1); err != nil {
						kp.errors <- err
						return
					}
					i++
				} else {
					kp.errors <- io.EOF
					return
				}
			}
		}

	}()

	// process errors channel
	go func() {
		for e := range deliveryCh {
			m := e.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				log.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		}
		log.Println("end process delivery channel")
	}()
}

func (kp *KafkaProducer) shutDown() error {
	kp.done <- struct{}{}
	kp.Flush(5000)
	kp.Close()
	log.Println("kp flushed and closed")
	return nil
}

func (kp *KafkaProducer) produceMsg(deliveryCh chan kafka.Event, topic string, msg string, partition int32) error {

	if partition == -1 {
		partition = kafka.PartitionAny
	}

	err := kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          []byte(msg),
		Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, deliveryCh)

	return err
}
