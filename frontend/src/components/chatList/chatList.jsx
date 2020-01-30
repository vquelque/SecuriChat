import React from "react";

const ChatList = props => {
  const { rooms, currentRoom, connectToRoom } = props;

  const roomList = rooms.map(room => {
    const roomIcon = room.authenticated === "AUTHENTICATION_OK" ? "🔒" : "🌐";
    const isRoomActive = room.id === currentRoom ? "active" : "";
    return (
      <li
        className={isRoomActive}
        key={room.id}
        onClick={() => connectToRoom(room)}
      >
        <span className="room-icon">{roomIcon}</span>
        <span className="room-name">{room.id}</span>
      </li>
    );
  });
  return (
    <div className="rooms">
      <ul className="chat-rooms">{roomList}</ul>
    </div>
  );
};

export default ChatList;
