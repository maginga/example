package filewatcher

import (
	"context"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

var triggerMd = trigger.NewMetadata(&HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type Trigger struct {
	handlers []trigger.Handler
	logger   log.Logger
}

type Factory struct {
}

func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	return &Trigger{}, nil
}

func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	t.handlers = ctx.GetHandlers()
	t.logger = ctx.Logger()
	return nil
}

func (t *Trigger) Start() error {
	t.logger.Debug("Start")
	handlers := t.handlers
	t.logger.Debug("Processing handlers")

	for _, handler := range handlers {
		fmt.Println("Starting File watching process")

		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			t.logger.Error("Error metadata: ", err.Error())
		}

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			t.logger.Error(err)
		}
		defer watcher.Close()

		done := make(chan bool)
		go func() {
			for {
				select {
				case event := <-watcher.Events:

					if event.Op&fsnotify.Write == fsnotify.Write {
						trgData := make(map[string]interface{})
						trgData["fileName"] = event.Name
						response, err := handler.Handle(context.Background(), trgData)

						fmt.Println("modified file:", event.Name)
						if err != nil {
							t.logger.Error("Error starting action: ", err.Error())
						} else {
							t.logger.Debugf("Action call successful: %v", response)
						}
					}
				case err := <-watcher.Errors:
					fmt.Println("error:", err)
				}
			}
		}()

		err = watcher.Add(s.DirName)
		if err != nil {
			t.logger.Error(err)
		}
		<-done

	}
	return nil
}

func (t *Trigger) Stop() error {
	return nil
}
