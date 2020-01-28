import React, { Component } from "react";
import "./App.css";
import { connect, sendMsg, init } from "./api";
import Header from "./components/header/Header";
import ChatBox from "./components/chatBox/ChatBox";
import Input from "./components/input/input";
import ChatList from "./components/chatList/chatList";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      messages: [],
      peerId: "",
      roomList: ["a", "b", "c"]
    };
  }

  componentDidMount() {
    init(data => {
      this.setState(prevstate => ({
        peerId: data.PeerId,
        PubRSAKey: data.PubRSAKey
      }));
    });

    connect((origin, text, room) => {
      var msg = {
        origin: origin,
        text: text,
        room: room
      };

      this.setState(prevState => ({
        messages: [...this.state.messages, msg]
      }));
    });
  }

  send = text => {
    var message = JSON.stringify({
      room: this.state.currentRoom,
      destination: this.state.currentRoom,
      text: text
    });
    sendMsg(message);
  };

  joinChat = room => {
    this.setState(prevState => ({
      currentRoom: room
    }));
  };

  render() {
    return (
      <div className="App">
        <aside className="sidebar left-sidebar">
          <div className="user-profile">
            <span className="username">{this.state.peerId}</span>
            <span className="user-id">{this.state.PubRSAKey}</span>
          </div>
          <ChatList
            rooms={this.state.roomList}
            currentRoom={this.state.currentRoom}
            connectToRoom={this.joinChat}
          />
        </aside>
        <section className="chat-screen">
          <Header className="chat-header" peerID={this.state.currentRoom} />
          <ChatBox
            messages={this.state.messages}
            id={this.state.peerId}
            currentRoom={this.state.currentRoom}
          />
          <footer className="chat-footer">
            <Input send={this.send} />
          </footer>
        </section>
      </div>
    );
  }
}

export default App;
