var socket = new WebSocket("ws://localhost:8080/ws");

let connect = callback => {
  console.log("Attempting Connection...");

  socket.onopen = () => {
    console.log("Successfully Connected");
  };

  socket.onmessage = msg => {
    var json = JSON.parse(msg.data);
    console.log(json);
    var origin = json.Origin;
    var message = json.Text;
    var room = json.Room;
    var authenticated = json.Authenticated;
    var authQuestion = json.AuthQuestion;
    callback(origin, message, room, authenticated, authQuestion);
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };
};

let sendMsg = msg => {
  console.log("sending msg: ", msg);
  socket.send(msg);
};

let init = init_state => {
  fetch("http://localhost:8080/init")
    .then(res => res.json())
    .then(
      result => {
        init_state(result);
      },
      // Remarque : il est important de traiter les erreurs ici
      // au lieu d'utiliser un bloc catch(), pour ne pas passer à la trappe
      // des exceptions provenant de réels bugs du composant.
      error => {
        console.log("error initializing app" + error);
      }
    );
};

export { connect, sendMsg, init };
