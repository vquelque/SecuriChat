import React from "react";
import { Component } from "react";
import Drawer from "react-drag-drawer";
import "./addContact.scss";
import Form from "react-bootstrap/Form";
import { Button } from "react-bootstrap";

class AddContact extends Component {
  constructor(props) {
    super(props);
    this.state = {
      toggle: false
    };

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
    if (this.state.contactID === "") {
      alert("please fill contactID");
      event.preventDefault();
      return;
    }
    if (this.state.AuthQuestion === "" || this.state.AuthAnswer === "") {
      alert("Contact added without authentication !!");
    } else {
      alert("Contact added. Key echange and authentication will start..");
    }
    event.preventDefault();
    this.toggle();
    this.props.addContact(
      this.state.contactID,
      this.state.AuthQuestion,
      this.state.AuthAnswer
    );
  }

  toggle = () => {
    let { toggle } = this.state;
    this.setState({ toggle: !toggle });
  };
  render() {
    return (
      <div>
        <Drawer
          open={this.state.toggle}
          onRequestClose={this.toggle}
          modalElementClass="Modal"
          className="Drawer"
        >
          <div className="DrawerCard">
            <h2>Add a friend</h2>
            <Form style={{ width: "90%" }} onSubmit={this.handleSubmit}>
              <Form.Group>
                <Form.Label>Peer ID</Form.Label>
                <Form.Control
                  type="text"
                  placeholder="Enter Peer ID"
                  name="contactID"
                  onChange={this.handleInputChange}
                />
              </Form.Group>

              <Form.Group>
                <Form.Label>SMP Authentication Question </Form.Label>
                <Form.Control
                  type="text"
                  onChange={this.handleInputChange}
                  name="AuthQuestion"
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
                Add contact !
              </Button>
            </Form>
          </div>
        </Drawer>
        <button
          onClick={this.toggle}
          style={{ position: "absolute", bottom: 0 }}
        >
          Add a contact
        </button>
      </div>
    );
  }
}

export default AddContact;
