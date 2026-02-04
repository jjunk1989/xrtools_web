# XRTools Web

XRTools Web is a web application written in Go that supports WebSocket connections, message broadcasting, and file upload functionality.

It facilitates message synchronization between different devices, such as between Pico/Vision Pro and PC. It's also convenient for Mac and PC communication without needing additional software.

Additionally, the project includes a built-in WebXR 360-degree panoramic video player that supports immersive panoramic video playback on VR devices such as Vision Pro, Quest, and Pico, with support for 180°/360° spatial video formats and stereoscopic video modes.

* Run server on PC

![](./doc/serve.png)

* Send messages from PC browser

![](./doc/pc.png)

* Vision Pro receives messages and copies to clipboard

![](./doc/visionPro.png)

* Pico receives messages and copies to clipboard

![](./doc/pico.png)


## Features

- Establish WebSocket connections
- Receive and send messages
- Broadcast messages to all connected clients
- Copy messages to clipboard
- File upload
- WebXR 360 panoramic video player (supports Vision Pro, Quest, Pico)

## Tech Stack

- Go
- WebSocket
- HTML/CSS/JavaScript

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/jjunk1989/xrtools_web
   ```

2. Navigate to the project directory:

   ```bash
   cd xrtools_web
   ```

3. Generate SSL certificate and private key (for development and testing purposes only):

   ```bash
   openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
   ```

4. Run the project:

   ```bash
   go run main.go -port=8443
   ```

## Usage

### Message Synchronization

1. Open your browser and visit https://localhost:8443.

2. Enter a message in the input field and click the "Send" button to send the message.

3. Click the "Copy" button to copy the message to clipboard.

### Panoramic Video Player

1. Visit https://localhost:8443/xrvideo_player.html

2. Place your 360-degree panoramic video files in the `uploads/` directory

3. Select the video you want to play from the player page

4. Features:
   - Toggle between 180°/360° spatial video modes
   - Support for monoscopic and stereoscopic video formats
   - Enter immersive VR mode on VR devices (Vision Pro, Quest, Pico)
   - Use mouse or touch to drag and rotate the view

5. VR Mode:
   - Open the player in a WebXR-compatible VR device browser
   - Click the "Enter VR" button to enter immersive experience
   - Watch videos in 360 degrees within the VR environment

## Command Line Arguments

- port: Specify the port number for the server to listen on, default is 443
