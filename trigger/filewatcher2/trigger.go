package filewatcher2

import (
	"context"
	"strconv"

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
	t.logger.Info("handlers: " + strconv.Itoa(len(t.handlers)))
	return nil
}

func (t *Trigger) Start() error {
	t.logger.Info("Processing handlers.")

	done := make(chan bool)
	for _, handler := range t.handlers {
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

		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
						trgData := make(map[string]interface{})
						trgData["fileName"] = event.Name
						trgData["assetName"] = s.AssetName
						trgData["sensorName"] = s.SensorName
						trgData["sensorType"] = s.SensorType

						t.logger.Infof("[%s] File change event was triggered. Modified file: %s", s.AssetName, event.Name)
						response, err := handler.Handle(context.Background(), trgData)
						if err != nil {
							t.logger.Error("Error starting action: ", err.Error())
						} else {
							t.logger.Info("Action was successfully called.")
							t.logger.Debugf("Action Response: %v", response)
						}
					}
				case err := <-watcher.Errors:
					t.logger.Error("Error: ", err.Error())
				}
			}
		}()

		err = watcher.Add(s.DirName)
		if err != nil {
			t.logger.Error(err)
		}
		t.logger.Info("Watching : ", s.DirName)
	}
	<-done

	return nil
}

func (t *Trigger) Stop() error {
	return nil
}
