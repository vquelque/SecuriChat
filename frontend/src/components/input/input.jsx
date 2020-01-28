import { Component } from "react";
import React from "react";
import "./input.scss";

class Input extends Component {
  state = {
    text: ""
  };

  render() {
    return (
      <div className="message-form">
        <form onSubmit={e => this.onSubmit(e)}>
          <input
            onChange={e => this.onChange(e)}
            value={this.state.text}
            type="text"
            placeholder="Enter your message and press ENTER"
            autofocus="true"
          />
          <button>Send</button>
        </form>
      </div>
    );
  }

  onChange(e) {
    this.setState({ text: e.target.value });
  }

  onSubmit(e) {
    e.preventDefault();
    this.props.send(this.state.text);
    this.setState({ text: "" });
  }
}

export default Input;
