import React, { Component } from "react";
import "./App.css";
import { connect, sendMsg, init } from "./api";
import Header from "./components/header/Header";
import ChatBox from "./components/chatBox/ChatBox";
import Input from "./components/input/input";
import ChatList from "./components/chatList/chatList";
import AddContact from "./components/addContact/addContat";
import AuthPopup from "./components/authPopup";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      messages: [],
      peerId: "",
      roomList: [],
      authPopup: false,
      authQuestion: "",
      origin: ""
    };
  }

  componentDidMount() {
    init(data => {
      this.setState(prevstate => ({
        peerId: data.PeerId,
        PubRSAKey: data.PubRSAKey
      }));
    });

    connect(this.messageHandler);
  }

  send = text => {
    if (this.state.currentRoom == null) {
      alert("Please select a room first !");
      return;
    }
    var message = JSON.stringify({
      room: this.state.currentRoom,
      destination: this.state.currentRoom,
      text: text
    });
    sendMsg(message);
  };

  addRoom = (id, authenticated) => {
    var room = {
      id: id,
      authenticated: authenticated
    };
    console.log("adding room " + id);
    if (!roomAlreadyPresent(room, this.state.roomList)) {
      this.setState(() => ({
        roomList: [...this.state.roomList, room]
      }));
    }
  };

  joinChat = room => {
    this.setState(prevState => ({
      currentRoom: room.id
    }));
  };

  addContact = (contactID, AuthQuestion, AuthAnswer) => {
    var message = JSON.stringify({
      Room: contactID,
      AuthQuestion: AuthQuestion,
      AuthAnswer: AuthAnswer
    });
    sendMsg(message);
  };

  messageHandler = (origin, text, room, authenticated, authQuestion) => {
    if (authQuestion !== "") {
      //Handle auth question popup
      console.log("AuthQuestion received");
      this.openAuthPopup(authQuestion, origin);
    } else {
      var msg = {
        origin: origin,
        text: text,
        room: room,
        authenticated: authenticated
      };
      this.setState(prevState => ({
        messages: [...this.state.messages, msg]
      }));
      this.addRoom(room, authenticated);
    }
  };

  openAuthPopup = (authQuestion, origin) => {
    this.setState({
      authPopup: true,
      authQuestion: authQuestion,
      origin: origin
    });
  };

  closeAuthPopup = () => {
    this.setState({
      authPopup: false,
      authQuestion: "",
      origin: ""
    });
  };

  sendAuthAnswer = (authAnswer, peerId) => {
    var message = JSON.stringify({
      Room: peerId,
      AuthAnswer: authAnswer
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
          <ChatList
            rooms={this.state.roomList}
            currentRoom={this.state.currentRoom}
            connectToRoom={this.joinChat}
          />
          <AddContact addContact={this.addContact} />
        </aside>
        <section className="chat-screen">
          <Header className="chat-header" peerID={this.state.currentRoom} />
          <ChatBox
            messages={this.state.messages}
            id={this.state.peerId}
            currentRoom={this.state.currentRoom}
          />
          <AuthPopup
            peerID={this.state.origin}
            authQuestion={this.state.authQuestion}
            open={this.state.authPopup}
            sendAuthAnswer={this.sendAuthAnswer}
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

function roomAlreadyPresent(room, roomList) {
  for (var i = 0; i < roomList.length; i++) {
    if (roomList[i].id === room.id) {
      return true;
    }
  }
}
