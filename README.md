# k8s-test
这篇文章通过微服务（或者说多进程服务），把docker、k8s和istio串了起来。涉及相应的开发、部署、升级，本文把主干引了出来，但每一步更多的枝叶还需要取完善（可参考结尾链接）。实验环境是一台Macbook Pro上，安装Docker Desktop。

#### 1，创建两个Http服务，用来做示例。
access-app：做为集群对外的接入层，所有外部进来的请求都先到这里来
logic-app：做为集群内逻辑层，提供某一服务。
具体来说，`外部请求 <-> access <-> logic`，实际应用中，服务模块（logic层）可能会很多，但大致流程如此，这里为了演示，做到最简化。
两个应用都是用Golang来构建的http server，用其它语言也是无妨。

#### 2，用docker跑起来
Golang是需要编译的，我开发环境是Mac，所以不能直接用mac上的环境来编译生成。需要用跟运行环境一样的go来编译。运行环境是alpine，所以编译环境用了golang:alpine，
```
# 文件 Makefile
APP_NAME=access-http
SRC=$(wildcard *.go)
EXE=app
$(EXE):$(SRC)
    # 先在golang:alpine环境下编译生成可运行文件
    docker run --rm -it --name go-compile -v ${PWD}:/go/src --workdir /go/src golang:alpine go build -o $(EXE)
    # 再把可运行文件打包到image中
    docker build -t k8s-test/$(APP_NAME) .
```

```
# 文件Dockerfile
FROM alpine:latest
# 安装了curl，方便进入container做测试
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --update curl && rm -rf /var/cache/apk/*
WORKDIR /app/
COPY ./app .
EXPOSE 8080
ENTRYPOINT ["./app"]
```
之后，用docker images就可以看到本地有`k8s-test/access-http`镜像了

#### 3，部署到kubenate
k8s实现了多机docker上container的编排，基本上单机docker就用docker-compose，多机就用k8s。

先在mac上的docker desktop里面安装k8s。本来应该是一件很容易的事情，但是一旦东西拉不下来就变得很麻烦。参考https://github.com/AliyunContainerService/k8s-for-docker-desktop
操作到dashboard那一步就可以，后面ingress不是必需的，也一直跑不过去。

部署：
配置文件在`deploy-k8s`里面了，大致上是为access-http和logic-http分别创建pod和service，这里pod不是直接创建的，而是放在一组deployment里面，depolyment会自动创建repliset并把pod放里面，实现副本和各种策略。access-service是要对外入口，用NodePort或者LoadBalance类型。logic-service则不用，用默认的ClusterIP即可。
到此k8s的部署就算完成了。

#### 4，部署到Istio
有了k8s，为什么还需Istio呢？k8s只是对多机上的docker进程做管理，而istio是在k8s的基础上，提供了一套完善的service mesh，功能上提供了更灵活的路由和流量控制、熔断、超时、服务重试、故障注入等各种高级微服务功能，具体实现上，istio通过注入k8s（或者插入了一个插件），实现了sidecar（边车）模型，对代码是完全无侵入的，任何能跑在docker上的程序都可以跑在istio上。与之对应的是spring cloud微服务框架，只能跑在java程序上，且对代码是有侵入的。

istio是个大东西，学习起来不会是几个小时就能掌握的。安装和学习资料参考官方文档：
https://preliminary.istio.io/zh/docs/setup/getting-started/
安装：
```
curl -L https://istio.io/downloadIstio | sh -
```
当前是1.4.4版本，安装完成后，就可以部署access-http和logic-http，配置文件在deploy-istio里面。
首先是`access-istio.yaml`和`logic-istio.yaml`，跟在k8s里面的配置基本一样的，除了access-service不用对外暴露端口；
接着是`gateway.yaml`，里面创建了Gateway和VirtualService，这两个都是istio引入的新对象。其中gateway用了自己的`istio: ingressgateway`。
再接下来，演示了灰度升级，之前都是access-http访问logic-http，升级之后就变成access-http要访问logic-http-v1或logic-http-v2，要做到不中断服务，同时支持方便可配的灰度升级。先在`destination.yaml`里面新建了两个destination：v1和v2（它们都对应在Deployment中有自己的配置），然后新增一个VirtualService对它们做路由：
```
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: logic
  namespace: istio-test
spec:
  hosts:
    - logic
  http:
    - route:
        - destination:
            host: logic
            subset: v2
          weight: 10
        - destination:
            host: logic
            subset: v1
          weight: 90
```
上面例子表示90%的流量到`subset: v1`，10%的流量到`subset: v2`。

istio还有很多更高级的演示，可以看下官网上的[BookInfo示例](https://preliminary.istio.io/zh/docs/examples/bookinfo/)

参考：
https://preliminary.istio.io/zh/docs/concepts/what-is-istio/
https://github.com/AliyunContainerService/k8s-for-docker-desktop
https://sanyuesha.com/2019/05/17/kubernetes-tutorial-for-beginner/
