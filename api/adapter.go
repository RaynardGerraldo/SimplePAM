package main 

import (
    "github.com/gorilla/websocket"
)

type WebSocketRW struct {
    Conn *websocket.Conn
    readBuf []byte
}

func (w *WebSocketRW) Write(p []byte) (int, error) {
    err := w.Conn.WriteMessage(websocket.TextMessage, p)
    if err != nil {
        return 0, err
    }
    return len(p), nil
}

func (w *WebSocketRW) Read(p []byte) (int, error) {
    if len(w.readBuf) > 0 {
        n := copy(p, w.readBuf)
        w.readBuf = w.readBuf[n:]
        return n, nil
    }

    _, message, err := w.Conn.ReadMessage()
    if err != nil {
        return 0, err // Connection closed or error
    }

    n := copy(p, message)

    if n < len(message) {
        w.readBuf = message[n:]
    }

    return n, nil
}

