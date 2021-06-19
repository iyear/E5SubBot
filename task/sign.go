package task

import (
	"fmt"
	"github.com/iyear/E5SubBot/config"
	"github.com/iyear/E5SubBot/model"
	"go.uber.org/zap"
)

func Sign(clients []*model.Client) []*model.ErrClient {
	var errClients []*model.ErrClient

	done := make(chan struct{})
	in := make(chan *model.ErrClient, 5)
	out := make(chan *model.ErrClient, 5)

	go func() {
		for _, client := range clients {
			in <- &model.ErrClient{
				Client: client,
				Err:    nil,
			}
		}
		close(in)
	}()
	for i := 0; i < config.MaxGoroutines; i++ {
		go func() {
			for {
				select {
				case errCli, f := <-in:
					if !f {
						continue
					}

					err := errCli.GetOutlookMails()
					errCli.Err = err
					out <- errCli
				case <-done:
					return
				}
			}
		}()
	}
	for i := 0; i < len(clients); i++ {
		errClient := <-out
		if errClient.Err == nil {
			fmt.Printf("%s OK\n", errClient.MsId)
		} else {
			zap.S().Errorw("failed to sign",
				"error", errClient.Err,
				"id", errClient.ID,
			)
			//fmt.Printf("%s %s\n",errClient.MsId,errClient.Err)
		}
		errClients = append(errClients, errClient)
	}
	close(done)
	return errClients
}
