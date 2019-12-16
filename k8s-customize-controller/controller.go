package main

import (
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1beta1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	bolingcavalryv1 "k8s-customize-controller/pkg/apis/bolingcavalry/v1"
	clientset "k8s-customize-controller/pkg/client/clientset/versioned"
	//programmerscheme "k8s-customize-controller/pkg/client/clientset/versioned/scheme"
	informers "k8s-customize-controller/pkg/client/informers/externalversions/bolingcavalry/v1"
	listers "k8s-customize-controller/pkg/client/listers/bolingcavalry/v1"
)

const controllerAgentName = "programmer-controller"

const (
	SuccessSynced = "Synced"
	ErrResourceExists = "ErrResourceExists"
	MessageResourceExists = "Resource %q already exists and is not managed by Programmer"
	MessageResourceSynced = "Programmer synced successfully"
)

type Controller struct {
	kubeclientset kubernetes.Interface
	programmerclientset clientset.Interface

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced

	programmersLister listers.ProgrammerLister
	programmersSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface

	recorder record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	programmerclientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	programmerInformer informers.ProgrammerInformer) *Controller {

	//utilruntime.Must(programmerscheme.AddToScheme(scheme.Scheme))
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:    kubeclientset,
		programmerclientset: programmerclientset,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		programmersLister:   programmerInformer.Lister(),
		programmersSynced:   programmerInformer.Informer().HasSynced,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Programmers"),
		recorder:         recorder,
	}

	glog.Info("Setting up event handlers")
	programmerInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueProgrammer,
		UpdateFunc: func(old, new interface{}) {
			oldProgrammer := old.(*bolingcavalryv1.Programmer)
			newProgrammer := new.(*bolingcavalryv1.Programmer)
			if oldProgrammer.ResourceVersion == newProgrammer.ResourceVersion {
				return
			}
			controller.enqueueProgrammer(new)
		},
		DeleteFunc: controller.enqueueProgrammerForDelete,
	})

	return controller
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	glog.Info("开始controller业务，开始一次缓存数据同步")
	if ok := cache.WaitForCacheSync(stopCh, c.programmersSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("worker启动")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("worker已经启动")
	<-stopCh
	glog.Info("worker已经结束")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// 取数据处理
func (c *Controller) processNextWorkItem() bool {

	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// 在syncHandler中处理业务
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}

		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// 处理
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// 从缓存中取对象
	programmer, err := c.programmersLister.Programmers(namespace).Get(name)
	if err != nil {
		// 如果Student对象被删除了，就会走到这里，所以应该在这里加入执行
		if errors.IsNotFound(err) {
			glog.Infof("Student对象被删除，请在这里执行实际的删除业务: %s/%s ...", namespace, name)

			return nil
		}

		runtime.HandleError(fmt.Errorf("failed to list student by: %s/%s", namespace, name))

		return err
	}

	glog.Infof("这里是student对象的期望状态: %#v ...", programmer)
	glog.Infof("实际状态是从业务层面得到的，此处应该去的实际状态，与期望状态做对比，并根据差异做出响应(新增或者删除)")

	deploymentName := programmer.Spec.DeploymentName
	if deploymentName == "" {
		utilruntime.HandleError(fmt.Errorf("%s: deployment name must be specified", key))
		return nil
	}

	deployment, err := c.deploymentsLister.Deployments(programmer.Namespace).Get(deploymentName)
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1beta1().Deployments(programmer.Namespace).Create(newDeployment(programmer))
	}

	if err != nil {
		return err
	}

	if !metav1.IsControlledBy(deployment, programmer) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
		c.recorder.Event(programmer, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	if programmer.Spec.Replicas != nil && *programmer.Spec.Replicas != *deployment.Spec.Replicas {
		glog.V(4).Infof("Programmer %s replicas: %d, deployment replicas: %d", name, *programmer.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = c.kubeclientset.AppsV1beta1().Deployments(programmer.Namespace).Update(newDeployment(programmer))
	}

	if err != nil {
		return err
	}

	err = c.updateProgrammerStatus(programmer, deployment)
	if err != nil {
		return err
	}

	c.recorder.Event(programmer, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) updateProgrammerStatus(programmer *bolingcavalryv1.Programmer, deployment *appsv1.Deployment) error {
	programmerCopy := programmer.DeepCopy()
	programmerCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	_, err := c.programmerclientset.BolingcavalryV1().Programmers(programmer.Namespace).Update(programmerCopy)
	return err
}

// 数据先放入缓存，再入队列
func (c *Controller) enqueueProgrammer(obj interface{}) {
	var key string
	var err error
	// 将对象放入缓存
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}

	// 将key放入队列
	c.workqueue.AddRateLimited(key)
}

// 删除操作
func (c *Controller) enqueueProgrammerForDelete(obj interface{}) {
	var key string
	var err error
	// 从缓存中删除指定对象
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	//再将key放入队列
	c.workqueue.AddRateLimited(key)
}

func newDeployment(programmer *bolingcavalryv1.Programmer) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "nginx",
		"controller": programmer.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      programmer.Spec.DeploymentName,
			Namespace: programmer.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(programmer, bolingcavalryv1.SchemeGroupVersion.WithKind("Programmer")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: programmer.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
		},
	}
}
