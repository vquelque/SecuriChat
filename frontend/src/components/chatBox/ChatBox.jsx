import React, { Component } from "react";
import "./ChatBox.scss";

class ChatBox extends Component {
  render() {
    const { messages, currentRoom } = this.props;
    const roomMsg = messages.filter(m => m.room === currentRoom);
    return (
      <ul className="Messages-list">
        {roomMsg.map((msg, index) => this.renderMessage(msg, index))}
      </ul>
    );
  }

  renderMessage(message, index) {
    const { origin, text } = message;
    const { id } = this.props;
    const messageFromMe = id === origin;
    const cssClass = messageFromMe
      ? "Messages-message currentMember"
      : "Messages-message";
    return (
      <li className={cssClass} key={index}>
        <div className="Message-content">
          <div className="username">{origin}</div>
          <div className="text">{text}</div>
        </div>
      </li>
    );
  }
}

export default ChatBox;
