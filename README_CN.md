# XRTools Web

XRTools Web 是一个使用 Go 语言编写的 Web 应用程序，支持 WebSocket 连接、消息广播和文件上传功能。

方便不同设备之间的消息同步。比如 Pico / visionPro 和 PC 之间同步消息。mac 和 pc 之间使用也方便，不用别的软件。

此外，项目还内置了 WebXR 360 度全景视频播放器，支持在 Vision Pro、Quest、Pico 等 VR 设备中播放沉浸式全景视频，支持 180°/360° 空间视频格式和立体视频模式。

* PC 上运行服务端

![](./doc/serve.png)

* PC 浏览器访问发送消息

![](./doc/pc.png)

* visionPro 接收消息并复制到剪贴板

![](./doc/visionPro.png)

* pico 接受消息并复制到剪贴板

![](./doc/pico.png)


## 功能

- 建立 WebSocket 连接
- 接收和发送消息
- 广播消息给所有连接的客户端
- 复制消息到剪贴板
- 上传文件
- WebXR 360 全景视频播放器(支持visionPro,quest,pico)

## 技术栈

- Go
- WebSocket
- HTML/CSS/JavaScript

## 安装

1. 克隆仓库：

   ```bash
   git clone https://github.com/jjunk1989/xrtools_web
   ```

2. 进入项目目录：

   ```bash
   cd xrtools_web
   ```

3. 生成SSL证书和私钥（仅用于开发和测试目的）：

   ```bash
    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
   ```

4. 运行项目：

   ```bash
   go run main.go -port=8443
   ```

## 使用

### 消息同步功能

1. 打开浏览器，访问 https://localhost:8443。

2. 在输入框中输入消息，点击"发送"按钮发送消息。

3. 点击"复制"按钮复制消息到剪贴板

### 全景视频播放器

1. 访问 https://localhost:8443/xrvideo_player.html

2. 将360度全景视频文件放入 `uploads/` 目录

3. 在播放器页面选择要播放的视频

4. 支持以下功能：
   - 180°/360° 空间视频切换
   - 支持单目和立体视频格式
   - 在VR设备（Vision Pro、Quest、Pico）中进入沉浸式VR模式
   - 使用鼠标或触摸拖拽来旋转视角

5. VR模式使用：
   - 在支持WebXR的VR设备浏览器中打开播放器
   - 点击"Enter VR"按钮进入沉浸式体验
   - 在VR环境中360度观看视频

## 命令行参数

- port: 指定服务器监听的端口号，默认为443
