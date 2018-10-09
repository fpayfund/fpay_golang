/* The MIT License (MIT)
Copyright © 2018 by Atlas Lee(atlas@fpay.io)

Permission is hereby granted, free of charge, to any person obtaining a
copy of this software and associated documentation files (the “Software”),
to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"errors"
	"fpay/account"
	"fpay/node"
	"time"
	"zlog"
)

// 1. 公共库
// 1.1. 数据结构
// 1.1.1. 账户

func TestCreateAccount() (err error) {

	_, err = account.Create()
	if err == nil {
		zlog.Debugln("Succeed!")
	} else {
		zlog.Fatalln("Failed!")
	}
	return
}

// 1.1.3. 通用请求
type BaseAO struct {
}

// 1.1.5. 区块

// 1.2. 全局数据

// 2. 命令行
// 2.1. Option命令行选项
// 3. 测试进程
// 4. 服务进程
// 4.1. 管理进程
// 4.2. 公共服务进程
// 4.2.1. 最新区块缓存（双向链表）
// 4.2.2. 地址缓存（KV哈希缓存）
// 4.2.4. 地址持久化（KV哈希持久化）
// 4.2.4. 全量区块持久化
// 4.4. 父节点进程
// 4.4.1. 根节点状态
// 4.4.2. 评委状态
// 4.4.4. 子节点状态
// 4.4.4. 叶子节点状态
// 4.4.5. 轻节点状态
// 4.4. 子节点进程
// 4.4.1. 基础服务
// 4.4.1.1. 创建账户

// 4.4.1. 根节点状态
// 4.4.2. 评委状态
// 4.4.4. 子节点状态
// 4.4.4. 叶子节点状态
// 4.4.5. 轻节点状态
// 4.5. 动态优化进程
// 4.5.1. 性能监控
// 5. 节点服务
// 5.1. Node: 节点
// 5.1.1. Root: 根节点
//

type FounderSettings struct { /* 创世块设定 */
}

type Settings struct { /* 设置，可从命令行解释 */
}

const (
	INITIALIZE = iota
	BOOKKEEPER
	REVIEWER
	TRANSFERRER
	COLLECTOR
	SUBMITTER
	SHUTTER
)

var StateNames = [7]string{
	"INITIALIZE",
	"BOOKKEEPER",
	"REVIEWER",
	"TRANSFERRER",
	"COLLECTOR",
	"SUBMITTER",
	"SHUTTER"}

type Service struct { /* 节点上下文。此对象线程不安全 */
	state       uint8
	parents     []node.Node /* 保存连接的父节点 */
	nodes       []node.Node
	connections uint16 /* 连接数 */
}

type InitializeContext struct {
	timeout time.Time
	state   uint8
}

func Create(settings *Settings) (service *Service) { /* 根据设置创建服务 */
	return nil
}

// initialize
// 初始化状态，启动时执行
// 主要任务是寻找父节点，确定自己的状态
// 如果找不到，则自己就是根节点
func (this *Service) initialize(context interface{}) (nextState uint8, nextContext *InitializeContext, err error) {
	/*
		1. 查找Parent
		1.1. 如果有指定节点，尝试连接指定节点
		1.2. 如果没有指定节点，尝试连接默认节点
		2. 找不到，自己作为Root
		3. 找到，请求与Parent建立连接
	*/
	if context == nil {
		nextContext = new(InitializeContext)
		nextContext.timeout = time.Unix(time.Now().Unix()+30, 0) /* 30秒超时 */
	} else {
		var ok bool
		nextContext, ok = context.(*InitializeContext)
		if !ok {
			return SHUTTER, nil, errors.New("InitializeContext initialization failed.")
		}
	}

	/* 是否timeout */
	if time.Now().After(nextContext.timeout) && this.connections == 0 {
		/* 如果没找到其它FPAY节点，则本节点为出块节点 */
		return BOOKKEEPER, nil, nil
	}

	return INITIALIZE, nextContext, nil
}

func (this *Service) bookkeeper(context interface{}) (nextState uint8, nextContext interface{}, err error) { /* 出块者状态 */

	return INITIALIZE, nil, nil
}

func (this *Service) reviewer(context interface{}) (nextState uint8, nextContext interface{}, err error) { /* 评委状态 */
	return INITIALIZE, nil, nil
}

func (this *Service) transferrer(context interface{}) (nextState uint8, nextContext interface{}, err error) { /* 转送者状态 */
	return INITIALIZE, nil, nil
}

func (this *Service) collector(context interface{}) (nextState uint8, nextContext interface{}, err error) { /* 收集者状态 */
	return INITIALIZE, nil, nil
}

func (this *Service) submitter(context interface{}) (nextState uint8, nextContext interface{}, err error) { /* 提交者 */
	return INITIALIZE, nil, nil
}

func (this *Service) shutter(context interface{}) (nextState uint8, nextContext interface{}, err error) { /* 关闭者 */
	return INITIALIZE, nil, nil
}

func (this *Service) Run() (err error) {
	var (
		lastState uint8       = INITIALIZE
		context   interface{} = nil
		nextState uint8
	)

	for {
		switch this.state {
		case INITIALIZE:
			if lastState == this.state {
				nextState, context, err = this.initialize(context)
			} else {
				nextState, context, err = this.initialize(nil)
			}

		case BOOKKEEPER:
			if lastState == this.state {
				nextState, context, err = this.bookkeeper(context)
			} else {
				nextState, context, err = this.bookkeeper(nil)
			}

		case REVIEWER:
			if lastState == this.state {
				nextState, context, err = this.reviewer(context)
			} else {
				nextState, context, err = this.reviewer(nil)
			}

		case TRANSFERRER:
			if lastState == this.state {
				nextState, context, err = this.transferrer(context)
			} else {
				nextState, context, err = this.transferrer(nil)
			}

		case COLLECTOR:
			if lastState == this.state {
				nextState, context, err = this.collector(context)
			} else {
				nextState, context, err = this.collector(nil)
			}

		case SUBMITTER:
			if lastState == this.state {
				nextState, context, err = this.submitter(context)
			} else {
				nextState, context, err = this.submitter(nil)
			}

		case SHUTTER:
			if lastState == this.state {
				nextState, context, err = this.shutter(context)
			} else {
				nextState, context, err = this.shutter(nil)
			}
		}

		lastState = nextState
	}
}

func TestFPAY() {
	var (
		settings *Settings
		service  *Service
	)

	settings = new(Settings)
	service = Create(settings)
	service.Run()
}

// 6. 测试主函数
func TestAll() (err error) {
	//err = TestCreateAccount()

	return nil
}

func main() {
	zlog.Infoln("Build Succeed!")

	err := TestAll()
	if err == nil {
		zlog.Infoln("TestAll Succeed!")
	} else {
		zlog.Infoln("TestAll Failed!")
	}
}

// 6. 主函数
