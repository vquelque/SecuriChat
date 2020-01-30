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
            <h2>{this.props.peerID} wants to authenticate !</h2>
            <Form style={{ width: "90%" }} onSubmit={this.handleSubmit}>
              <Form.Group>
                <Form.Label>SMP Authentication Question : </Form.Label>
                <input
                  type="text"
                  readonly
                  class="form-control-plaintext"
                  id="authQuestion"
                  value={this.props.authQuestion}
                />
              </Form.Group>
              <Form.Group>
                <Form.Label>SMP Authentication Answer </Form.Label>
                <Form.Control
                  type="text"
                  onChange={this.handleInputChange}
                  name="AuthAnswer"
                />
              </Form.Group>
              <Button
                variant="primary"
                type="submit"
                style={{
                  display: "flex",
                  justifyContent: "center",
                  alignItems: "center"
                }}
              >
                Send the Answer !
              </Button>
            </Form>
          </div>
        </Drawer>
      </div>
    );
  }
}

export default AuthPopup;
