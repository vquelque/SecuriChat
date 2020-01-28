import React, { Component } from "react";
import "./App.css";
import { connect, sendMsg, init } from "./api";
import Header from "./components/header/Header";
import ChatBox from "./components/chatBox/ChatBox";
import Input from "./components/input/input";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      messages: [],
      peerId: ""
    };
  }

  componentDidMount() {
    init(data => {
      this.setState(prevstate => ({
        peerId: data.PeerId,
        PubRSAKey: data.PubRSAKey
      }));
    });

    connect((origin, text) => {
      var msg = {
        origin: origin,
        text: text
      };
      this.setState(prevState => ({
        messages: [...this.state.messages, msg]
      }));
    });
  }

  send = text => {
    var message = JSON.stringify({
      destination: "",
      text: text
    });
    sendMsg(message);
  };

  render() {
    return (
      <div className="App">
        <aside className="sidebar left-sidebar">
          <div className="user-profile">
            <span className="username">{this.state.peerId}</span>
            <span className="user-id">{this.state.PubRSAKey}</span>
          </div>
        </aside>
        <section className="chat-screen">
          <Header className="chat-header" peerId={this.state.peerId} />
          <ChatBox messages={this.state.messages} id={this.state.peerId} />
          <footer className="chat-footer">
            <Input send={this.send} />
          </footer>
        </section>
      </div>
    );
  }
}

export default App;
