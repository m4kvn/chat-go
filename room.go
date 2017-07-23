package main

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join: // 参加
			r.clients[client] = true
		case client := <-r.leave: // 退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// 全てのクライアントにメッセージを送信
			for client := range r.clients {
				select {
				case client.send <- msg: // メッセージを送信
				default: // 送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
