import React from "react";
import { Component } from "react";
import "./Header.scss";

class Header extends Component {
  render() {
    const { peerId } = this.props;
    return (
      <div className="header">
        <h2>SecuriChat - Try it, use it.</h2>
        <h3>Peer Id : {peerId} </h3>
      </div>
    );
  }
}

export default Header;
