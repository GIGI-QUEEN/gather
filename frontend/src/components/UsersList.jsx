import React from "react"
import { useNavigate } from "react-router-dom"
import Image from "./Image"
import "../styles/follow-list.scss"
import { usernameOrName } from "../components/utils/usernameOrName"
const UsersList = ({ users, title, pageOwner }) => {
  const nav = useNavigate()

  return (
    <div className="content-container">
      <div className="users-container">
        <h1>
          <span onClick={() => nav(`/user/${pageOwner.id}`)}>
            {typeof pageOwner === "string"
              ? pageOwner
              : pageOwner?.username + "'s"}
          </span>{" "}
          {title}
        </h1>
        {users?.map((user) => (
          <div
            className="user-container"
            key={user.id}
            onClick={() => nav(`/user/${user.id}`)}
          >
            <Image image={user?.avatar} className={"list-avatar"} />
            {usernameOrName(user)}
          </div>
        ))}
      </div>
    </div>
  )
}

export default UsersList
