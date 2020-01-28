import React, { Component } from "react";
import "./App.css";
import { connect, sendMsg } from "./api";
import Header from "./components/header/Header";
import ChatBox from "./components/chatBox/ChatBox";
import Input from "./components/input/input";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      messages: [],
      id: "wlRvQdUpIJ"
    };
  }

  componentDidMount() {
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
        <Header />
        <ChatBox messages={this.state.messages} id={this.state.id} />
        <Input send={this.send} />
      </div>
    );
  }
}

export default App;
