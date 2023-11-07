import { Message } from "./Message"
import { useEffect, useRef } from "react"

export const Messages = ({ messages, user, typeOfChat }) => {
  const messagesEndRef = useRef(null)
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [messages])

  return (
    <div className="messages-container">
      {messages?.map((message, index) => (
        <Message
          key={index}
          message={message}
          user={user}
          typeOfChat={typeOfChat}
        />
      ))}
      <div ref={messagesEndRef} />
    </div>
  )
}
