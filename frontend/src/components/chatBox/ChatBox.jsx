import React, { Component } from "react";
import "./ChatBox.scss";

class ChatBox extends Component {
  render() {
    const { messages } = this.props;
    return (
      <ul className="Messages-list">
        {messages.map(msg => this.renderMessage(msg))}
      </ul>
    );
  }

  renderMessage(message) {
    const { origin, text } = message;
    const { id } = this.props;
    const messageFromMe = id === origin;
    const cssClass = messageFromMe
      ? "Messages-message currentMember"
      : "Messages-message";
    return (
      <li className={cssClass} key="">
        <div className="Message-content">
          <div className="username">{origin}</div>
          <div className="text">{text}</div>
        </div>
      </li>
    );
  }
}

export default ChatBox;
