package gin_frame

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huangzhilin/gin-frame/core"

	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Path       string                     //配置类型：file->本地配置文件;consul->consul配置中心
	Type       string                     //配置路径: 如./config/config_develop.yaml 或 http://192.168.3.91:8500
	Token      string                     //consul的token，type=file时为空
	KvPath     string                     //consul的远程key/value路径
	RemoteType string                     //consul远程配置类型，如：yaml、json、hcl
	Router     *gin.Engine                //gin路由
	httpServer *http.Server               //http服务
	Validators map[string]core.Validators //gin 数据验证
}

//Init 初始化配置,包括：初始化gorm、redis、zap日志等等初始化行为
func (t ServerConfig) Init() ServerConfig {
	core.InitConfig(t.Path, t.Type, t.Token, t.KvPath, t.RemoteType)

	//gin的mod设置 支持 release 和 debug 两种模式
	gin.SetMode(viper.GetString("appMode"))

	//gin 验证器 中文翻译
	if err := core.InitTrans("zh", t.Validators); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
	}

	core.InitZap() //初始化zap日志

	return t
}

//Gorm 初始化gorm
func (t ServerConfig) Gorm() ServerConfig {
	DB, DBList = core.InitGorm()
	return t
}

//Redis 初始化redis
func (t ServerConfig) Redis() ServerConfig {
	Redis, _ = core.InitRedisClient()
	return t
}

//RunHttp 开启并运行http服务
func (t ServerConfig) RunHttp() {
	//可以添加一些公共部分的路由
	t.Router.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	//第一步、安装swag： go get -u github.com/swaggo/swag/cmd/swag
	//第二步、在项目根目录下使用swag工具生成接口文档数据，多出docs文件夹：swag init
	//第三步、swagger生成文档，访问http://127.0.0.1:8887/swagger/index.html

	errChan := make(chan error) //如果程序退出让服务注销掉

	t.httpServer = &http.Server{
		Addr:    ":" + viper.GetString("post"),
		Handler: t.Router, // 初始化路由，加载APP目录中的路由配置
	}

	go (func() {
		fmt.Println("Listening and serving HTTP on :" + viper.GetString("post"))
		if err := t.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	})()

	go (func() {
		sigC := make(chan os.Signal)
		signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM) //监听kill  或 ctrl+c 等退出程序
		errChan <- fmt.Errorf("%s", <-sigC)
	})()

	getErr := <-errChan //如果没有错误，则一直阻塞，不会进行下面的操作
	log.Println(getErr)

	//实现优雅关闭服务
	t.httpStop()
}

//httpStop 实现优雅关闭服务
func (t ServerConfig) httpStop() {
	log.Println("准备关闭服务...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := t.httpServer.Shutdown(ctx); err != nil {
		log.Fatal("关闭服务错误: ", err)
	}
	log.Println("服务已关闭！")
}
