package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// type Client struct {
//     Cluster  //向集群里增加 etcd 服务端节点之类，属于管理员操作。
//     KV 		 //我们主要使用的功能，即 K-V 键值库的操作。
//     Lease    //租约相关操作，比如申请一个 TTL=10 秒的租约（应用给 key 可以实现键值的自动过期）。
//     Watcher  //观察订阅，从而监听最新的数据变化。
//     Auth     //管理 etcd 的用户和权限，属于管理员操作。
//     Maintenance   //维护 etcd ，比如主动迁移 etcd 的 leader 节点，属于管理员操作。

//	    // Username is a user name for authentication.
//	    Username string
//	    // Password is a password for authentication.
//	    Password string
//	    // contains filtered or unexported fields
//	}
var (
	etcdCli *clientv3.Client
	etcdCtx context.Context
)

func init() {
	etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
		// Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("etcd init err:", err)
	}
	// defer etcdCli.Close()
}

func EtcdInit() {

	// 我们通过方法clientv3.NewKV()来获得 KV 接口的实现（实现中内置了错误重试机制）：
	kv := clientv3.NewKV(etcdCli)
	etcdCtx = context.TODO()
	// 第一个参数是 goroutine 的上下文 Context 。后面两个参数分别是 key 和 value ，对于 etcd 来说， key=/test/key1 只是一个字符串而已，但是对我们而言却可以模拟出目录层级关系。
	putResp, _ := kv.Put(etcdCtx, "/test/key1", "Hello etcd!")
	fmt.Println(putResp.Header)
	// 在上面的例子中，我没有传递 opOption ，所以就是获取 key=/test/key1 的最新版本数据。这里 err 并不能反馈出 key 是否存在（只能反馈出本次操作因为各种原因异常了），我们需要通过 GetResponse （实际上是 pb.RangeResponse ）判断 key 是否存在
	// 除了上面例子中的三个的参数，还支持一个变长参数，可以传递一些控制项来影响 Put 的行为，例如可以携带一个 lease ID 来支持 key 过期。
	kv.Put(etcdCtx, "/test/key2", "Hello World!")
	// 再写一个同前缀的干扰项
	kv.Put(etcdCtx, "/testspam", "spam")

	getResp, _ := kv.Get(etcdCtx, "/test/key1")
	// 在上面的例子中，我没有传递 opOption ，所以就是获取 key=/test/key1 的最新版本数据。这里 err 并不能反馈出 key 是否存在（只能反馈出本次操作因为各种原因异常了），我们需要通过 GetResponse （实际上是 pb.RangeResponse ）判断 key 是否存在
	fmt.Println(getResp.Kvs)
	// RangeResponse.More和Count，当我们使用withLimit()等选项进行 Get 时会发挥作用，相当于翻页查询。
	// 接下来，我们通过给 Get 查询增加 WithPrefix 选项，获取 /test 目录下的所有子元素：

	rangeResp, _ := kv.Get(etcdCtx, "/test/", clientv3.WithPrefix())
	fmt.Println(rangeResp.Kvs, rangeResp.Count)

	for _, v := range rangeResp.Kvs {
		fmt.Println(v)
		fmt.Println(v.Lease)
	}
	/*
		clientv3.WithFromKey() 表示针对的key操作是大于等于当前给定的key
		clientv3.WithPrevKV() 表示返回的 response 中含有之前删除的值，否则
		下面的 delResp.PrevKvs 为空
	*/

	delResp, err := kv.Delete(context.TODO(), "/testspam",
		clientv3.WithFromKey(), clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	}
	// 查看被删除的 key 和 value 是什么
	if delResp.PrevKvs != nil {
		// if len(delResp.PrevKvs) != 0 {
		for _, kvpair := range delResp.PrevKvs {
			fmt.Println("已删除:", string(kvpair.Key), string(kvpair.Value))
		}
	}

}

func EtcdLease() {
	// type Lease interface {
	// 	// Grant 创建一个新租约
	// 	Grant(ctx context.Context, ttl int64) (*LeaseGrantResponse, error)
	// 	// Revoke 销毁给定租约ID的租约
	// 	Revoke(ctx context.Context, id LeaseID) (*LeaseRevokeResponse, error)
	// 	// TimeToLive retrieves the lease information of the given lease ID.
	// 	TimeToLive(ctx context.Context, id LeaseID, opts ...LeaseOption) (*LeaseTimeToLiveResponse, error)
	// 	// Leases retrieves all leases.
	// 	Leases(ctx context.Context) (*LeaseLeasesResponse, error)
	// 	// KeepAlive keeps the given lease alive forever.
	// 	KeepAlive(ctx context.Context, id LeaseID) (<-chan *LeaseKeepAliveResponse, error)
	// 	// KeepAliveOnce renews the lease once. In most of the cases, KeepAlive
	// 	// should be used instead of KeepAliveOnce.
	// 	KeepAliveOnce(ctx context.Context, id LeaseID) (*LeaseKeepAliveResponse, error)
	// 	// Close releases all resources Lease keeps for efficient communication
	// 	// with the etcd server.
	// 	Close() error
	// }
	lease := clientv3.NewLease(etcdCli)
	kv := clientv3.NewKV(etcdCli)
	// 要想实现 key 自动过期，首先得创建一个租约，下面的代码创建一个 TTL 为 10 秒的租约：
	etcdCtx := context.TODO()
	grantResp, _ := lease.Grant(etcdCtx, 10)
	// 这里特别需要注意，有一种情况是在 Put 之前 Lease 已经过期了，那么这个 Put 操作会返回 error ，此时你需要重新分配 Lease 。
	kv.Put(etcdCtx, "/test/vanish", "vanish in 10s", clientv3.WithLease(grantResp.ID))

	// 当我们实现服务注册时，需要主动给 Lease 进行续约，
	// 通常是以小于 TTL 的间隔循环调用 Lease 的 KeepAliveOnce() 方法对租约进行续期，
	// 一旦某个服务节点出错无法完成租约的续期，等 key 过期后客户端即无法在查询服务时获得对应节点的服务，
	// 这样就通过租约到期实现了服务的错误隔离。
	// 类似看门狗机制?需要一直喂食，停止喂食代表出现错误
	// 或者使用KeepAlive()方法，其会返回<-chan *LeaseKeepAliveResponse只读通道，每次自动续租成功后会向通道中发送信号。
	// keepResp, _ := lease.KeepAliveOnce(etcdCtx, grantResp.ID)

	// 一般都用KeepAlive()方法，
	// KeepAlive 和 Put 一样，如果在执行之前 Lease 就已经过期了，那么需要重新分配 Lease 。
	//  etcd 并没有提供 API 来实现原子的 Put with Lease ，需要我们自己判断 err 重新分配 Lease 。

	leaseId := grantResp.ID
	// 自动永久续租
	keepRespChan, err := lease.KeepAlive(etcdCtx, leaseId)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 处理续约应答的协程
	go func() {
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约已经失效了")
					goto END
				} else { // 每秒会续租一次, 所以就会受到一次应答
					fmt.Println("收到自动续租应答:", keepResp.ID)
				}
			}
		}
	END:
	}()

	// 获得kv API子集
	// kv := clientv3.NewKV(client)

	// Put一个KV, 让它与租约关联起来, 从而实现10秒后自动过期
	putResp, err := kv.Put(context.TODO(), "/demo/A/B1", "hello", clientv3.WithLease(leaseId))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入成功:", putResp.Header.Revision)

	// 定时的看一下key过期了没有
	for {
		getResp, err := kv.Get(context.TODO(), "/demo/A/B1")
		if err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没过期:", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
}

// Op 字面意思就是”操作”， Get 和 Put 都属于 Op ，只是为了简化用户开发而开放的特殊 API 。

// KV 对象有一个 Do 方法接受一个 Op ：
// 其参数 Op 是一个抽象的操作，可以是 Put/Get/Delete… ；而 OpResponse 是一个抽象的结果，可以是 PutResponse/GetResponse…
func EtcdOP() {
	defer etcdCli.Close()
	ops := []clientv3.Op{
		clientv3.OpPut("put-key", "123"),
		clientv3.OpGet("put-key"),
		clientv3.OpPut("put-key", "456"),
	}
	for _, op := range ops {
		if ksv, err := etcdCli.Do(context.TODO(), op); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(ksv)
		}
	}

}

// etcd 中事务是原子执行的，只支持 if … then … else … 这种表达。首先来看一下 Txn 中定义的方法：\
// Txn 必须是这样使用的：If(满足条件) Then(执行若干Op) Else(执行若干Op)。
// If 中支持传入多个 Cmp 比较条件，如果所有条件满足，则执行 Then 中的 Op （上一节介绍过Op），否则执行 Else中 的 Op 。
func EtcdTxn() {
	client := etcdCli
	// lease实现锁自动过期:
	// op操作
	// txn事务: if else then

	// 1, 上锁 (创建租约, 自动续租, 拿着租约去抢占一个key)
	lease := clientv3.NewLease(client)

	// 申请一个5秒的租约
	leaseGrantResp, err := lease.Grant(context.TODO(), 5)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 拿到租约的ID
	leaseId := leaseGrantResp.ID

	// 准备一个用于取消自动续租的context
	ctx, cancelFunc := context.WithCancel(context.TODO())

	// 确保函数退出后, 自动续租会停止
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	// 5秒后会取消自动续租
	keepRespChan, err := lease.KeepAlive(ctx, leaseId)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 处理续约应答的协程
	go func() {
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约已经失效了")
					goto END
				} else { // 每秒会续租一次, 所以就会受到一次应答
					fmt.Println("收到自动续租应答:", keepResp.ID)
				}
			}
		}
	END:
	}()

	//  if 不存在key， then 设置它, else 抢锁失败
	kv := clientv3.NewKV(client)

	// 创建事务
	txn := kv.Txn(context.TODO())

	// 定义事务

	// 如果key不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision("/demo/A/B1"), "=", 0)).
		Then(clientv3.OpPut("/demo/A/B1", "xxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/demo/A/B1")) // 否则抢锁失败

	// 提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		fmt.Println(err)
		return // 没有问题
	}

	// 判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用:", string(
			txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	// 2, 处理业务

	fmt.Println("处理任务")
	time.Sleep(5 * time.Second)

	// 3, 释放锁(取消自动续租, 释放租约)
	// defer 会把租约释放掉, 关联的KV就被删除了
}

// Watch 用于监听某个键的变化, Watch调用后返回一个WatchChan，
// 当监听的 key 有变化后会向WatchChan发送WatchResponse。

// Watch 的典型应用场景是应用于系统配置的热加载，我们可以在系统读取到存储在 etcd key 中的配置后，用 Watch 监听 key 的变化。
// 在单独的 goroutine 中接收 WatchChan 发送过来的数据，并将更新应用到系统设置的配置变量中，
// 比如像下面这样在 goroutine 中更新变量 appConfig ，这样系统就实现了配置变量的热加载。
type AppConfig struct {
	config1 string
	config2 string
}

var appConfig AppConfig

func watchConfig(clt *clientv3.Client, key string, ss interface{}) {
	watchCh := clt.Watch(context.TODO(), key)
	go func() {
		for res := range watchCh {
			value := res.Events[0].Kv.Value
			if err := json.Unmarshal(value, ss); err != nil {
				fmt.Println("now", time.Now(), "watchConfig err", err)
				continue
			}
			fmt.Println("now", time.Now(), "watchConfig", ss)
		}
	}()
}
func EtcdWatch() {
	client := etcdCli

	// 获得kv API子集
	kv := clientv3.NewKV(client)

	// 模拟etcd中KV的变化
	go func() {
		for {
			kv.Put(context.TODO(), "/demo/A/B1", "i am B1")

			kv.Delete(context.TODO(), "/demo/A/B1")

			time.Sleep(1 * time.Second)
		}
	}()

	// 先GET到当前的值，并监听后续变化
	getResp, err := kv.Get(context.TODO(), "/demo/A/B1")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 现在key是存在的
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值:", string(getResp.Kvs[0].Value))
	}

	// 当前etcd集群事务ID, 单调递增的
	watchStartRevision := getResp.Header.Revision + 1

	// 创建一个watcher
	watcher := clientv3.NewWatcher(client)

	// 启动监听
	fmt.Println("从该版本向后监听:", watchStartRevision)

	// 创建一个 5s 后取消的上下文
	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})

	// 该监听动作在 5s 后取消
	watchRespChan := watcher.Watch(ctx, "/demo/A/B1", clientv3.WithRev(watchStartRevision))

	// 处理kv变化事件
	for watchResp := range watchRespChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			// mvccpb.PUT:
			case 0:
				fmt.Println("修改为:", string(event.Kv.Value), "Revision:",
					event.Kv.CreateRevision, event.Kv.ModRevision)
				// mvccpb.Delete:
			case 1:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}

}
