import React from "react";
import { Component } from "react";
import Drawer from "react-drag-drawer";
import "./authPopup.scss";
import Form from "react-bootstrap/Form";
import { Button } from "react-bootstrap";

class AuthPopup extends Component {
  constructor(props) {
    super(props);
    this.handleInputChange = this.handleInputChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleInputChange(event) {
    const target = event.target;
    const value = target.type === "checkbox" ? target.checked : target.value;
    const name = target.name;
    this.setState({
      [name]: value
    });
  }

  handleSubmit(event) {
    if (this.state.AuthAnswer === "") {
      alert("please fill in the answer");
      event.preventDefault();
      return;
    }
    alert("Answer send to peer !");
    event.preventDefault();
    // send the auth answer
    this.props.sendAuthAnswer(this.state.AuthAnswer, this.props.peerID);
  }

  render() {
    return (
      <div>
        <Drawer
          open={this.props.open}
          modalElementClass="Modal"
          className="Drawer"
        >
          <div className="DrawerCard">
            <h2>
              {" "}
              <strong style={{ color: "red" }}>{this.props.peerID}</strong>{" "}
              wants to authenticate !
            </h2>
            <br />
            <Form style={{ width: "90%" }} onSubmit={this.handleSubmit}>
              <Form.Group>
                <Form.Label>SMP Authentication Question :</Form.Label>
                <p style={{ fontWeight: "bold" }}>{this.props.authQuestion}</p>
              </Form.Group>
              <Form.Group>
                <Form.Label>SMP Authentication Answer </Form.Label>
                <Form.Control
                  type="text"
                  onChange={this.handleInputChange}
                  name="AuthAnswer"
                />
              </Form.Group>
              <Form.Group>
                <Button variant="success" type="submit" block>
                  Send the Answer !
                </Button>
              </Form.Group>
            </Form>
          </div>
        </Drawer>
      </div>
    );
  }
}

export default AuthPopup;
