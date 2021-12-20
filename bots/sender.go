package bots

import (
	"fmt"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"sync"
)

type Sender struct {
	done chan struct{}
	in   chan *Msg
	wg   sync.WaitGroup
}
type Msg struct {
	To      tb.Recipient
	What    interface{}
	Options []interface{}
}

func NewSender() *Sender {
	return &Sender{}
}
func (s *Sender) Init(goroutine int) {

	s.done = make(chan struct{})
	s.in = make(chan *Msg)

	for i := 0; i < goroutine; i++ {
		go func() {
			s.sender()
		}()
	}
}

func (s *Sender) Stop() {
	s.wg.Wait()
	close(s.done)
	fmt.Println("task finished")
}

func (s *Sender) SendMessageByID(ID int64, what interface{}, options ...interface{}) {
	s.wg.Add(1)
	go func() {
		chat, err := bot.ChatByID(strconv.FormatInt(ID, 10))
		if err != nil {
			zap.S().Errorw("failed to get chat",
				"error", err,
				"id", ID,
			)
			s.wg.Done()
			return
		}
		s.SendMessage(chat, what, options...)
	}()
}
func (s *Sender) SendMessage(to tb.Recipient, what interface{}, options ...interface{}) {
	s.in <- &Msg{
		To:      to,
		What:    what,
		Options: options,
	}
}
func (s *Sender) sender() {
	for {
		select {
		case msg, f := <-s.in:
			if !f {
				continue
			}
			fmt.Println("send to " + msg.To.Recipient())
			if _, err := bot.Send(msg.To, msg.What, msg.Options...); err != nil {
				zap.S().Errorw("failed to send msg",
					"error", err,
					"id", msg.To.Recipient(),
				)
			}
			s.wg.Done()
		case <-s.done:
			return
		}
	}
}
