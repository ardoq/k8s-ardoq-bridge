/*
Copyright © 2020 alexsimonjones@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"context"
	"flag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"K8SArdoqBridge/app/controllers"
	"K8SArdoqBridge/app/lib/runtime"
	"K8SArdoqBridge/app/lib/subscription"
	"K8SArdoqBridge/app/lib/watcher"
	"K8SArdoqBridge/app/subscriptions"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var (
	masterURL          string
	kubeconfig         string
	addr               = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	leaseLockName      string
	leaseLockNamespace string
	id                 string
)

func main() {

	klog.InitFlags(nil)
	flag.Parse()

	if leaseLockName == "" {
		klog.Fatal("unable to get lease lock resource name (missing lease-lock-name flag).")
	}
	if leaseLockNamespace == "" {
		klog.Fatal("unable to get lease lock resource namespace (missing lease-lock-namespace flag).")
	}

	start := time.Now()
	klog.Infof("Starting @ %s", start.String())

	go func() {
		http.Handle("/metrics", promhttp.Handler())

		log.Fatal(http.ListenAndServe(*addr, nil))
	}()
	go func() {
		klog.Error(http.ListenAndServe(":7777", http.DefaultServeMux))
	}()

	klog.Info("Got watcher client...")

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	klog.Info("Built config from flags...")

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building watcher clientset: %s", err.Error())
	}

	klog.Info("Created new KubeConfig")

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
		case <-ctx.Done():
			cancel()
		}
	}()

	//initialise cluster
	if os.Getenv("ARDOQ_CLUSTER") == "" {
		klog.Fatalf("ARDOQ_CLUSTER is a required environment variable")
	}
	klog.Info("Initialising cluster in Ardoq")
	controllers.GenericUpsert("Cluster", os.Getenv("ARDOQ_CLUSTER"))

	//start Resource Consumers
	go controllers.ResourceUpsertConsumer()
	go controllers.ResourceDeleteConsumer()

	klog.Info("Starting event buffer...")

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      leaseLockName,
			Namespace: leaseLockNamespace,
		},
		Client: kubeClient.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock: lock,
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: true,
		LeaseDuration:   10 * time.Second,
		RenewDeadline:   7 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				// we're notified when we start - this is where you would
				// usually put your code
				/*
					This is a default template file.
					Add subscriptions and watchers to make it your own.
				*/
				err = runtime.EventBuffer(ctx, kubeClient,
					&subscription.Registry{
						Subscriptions: []subscription.ISubscription{
							subscriptions.DeploymentSubscriber{},
							subscriptions.StatefulsetSubscriber{},
							subscriptions.NodeSubscriber{},
						},
					}, []watcher.IObject{
						kubeClient.AppsV1().Deployments(""),
						kubeClient.AppsV1().StatefulSets(""),
						kubeClient.CoreV1().Nodes(),
					})
				if err != nil {
					klog.Error(err)
				}
			},
			OnStoppedLeading: func() {
				// we can do cleanup here
				klog.Infof("leader lost: %s", id)
				os.Exit(0)
			},
			OnNewLeader: func(identity string) {
				// we're notified when new leader elected
				if identity == id {
					// I just got the lock
					return
				}
				klog.Infof("new leader elected: %s", identity)
			},
		},
	})

}

func init() {
	hostname, _ := os.Hostname()
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&leaseLockName, "lease-lock-name", "k8s-ardoq-bridge", "the lease lock resource name")
	flag.StringVar(&leaseLockNamespace, "lease-lock-namespace", "default", "the lease lock resource namespace")
	flag.StringVar(&id, "id", hostname, "the holder identity name")
}
