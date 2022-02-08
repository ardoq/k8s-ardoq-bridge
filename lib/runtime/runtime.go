package runtime

import (
	"K8SArdoqBridge/app/lib/metrics"
	"K8SArdoqBridge/app/lib/subscription"
	"K8SArdoqBridge/app/lib/watcher"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	k "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"sync"
	"time"
)

var (
	minWatchTimeout = (60 * 24 * 365) * time.Minute //should run for atleast 1 year without hanging up
	timeoutSeconds  = int64(minWatchTimeout.Seconds())
)

func EventBuffer(context context.Context, client k.Interface,
	registry *subscription.Registry, obj []watcher.IObject) error {

	if len(obj) == 0 {
		return errors.New("no watchers selected, exiting")
	}
	var watchers []<-chan watch.Event
	for _, o := range obj {
		funcObj := o
		w, err := funcObj.Watch(context, metav1.ListOptions{
			TimeoutSeconds:      &timeoutSeconds,
			Watch:               true,
			AllowWatchBookmarks: true})
		defer w.Stop()
		if err != nil {
			switch {
			case err == io.EOF:
				// watch closed normally
				klog.Infof("closed with EOF")
			case err == io.ErrUnexpectedEOF:
				klog.Infof("closed with unexpected EOF")
			}
			klog.Error(err)
		}
		watchers = append(watchers, w.ResultChan())
	}
	log.Debugf("%+v", watchers)
	var wg sync.WaitGroup
	wg.Add(len(watchers))
	for x, o := range watchers {
		x := x
		o := o
		go func() {
			err := func(t int, c <-chan watch.Event) error {
				defer wg.Done()
				counter := 0
				for {
					select {
					case update, hasUpdate := <-c:
						if hasUpdate {
							err := registry.OnEvent(subscription.Message{
								Event:  update,
								Client: client,
							})
							if err != nil {
								klog.Error(err)
								return err
							}
							metrics.TotalEventOps.Inc()
						}
						if !hasUpdate {
							// the channel got closed, so we need to restart
							klog.Error("Kubernetes hung up on us, restarting event watcher")
						}
						//case <-time.After(30 * time.Minute):
						//	// deal with the issue where we get no events
						//	klog.Fatal("Timeout, restarting event watcher")
					}
					counter++
				}
			}(x, o)
			if err != nil {
				klog.Error(err)
			}
		}()
	}
	wg.Wait()
	return nil
}
