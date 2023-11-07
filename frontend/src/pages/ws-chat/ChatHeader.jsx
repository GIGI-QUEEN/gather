import { useNavigate } from "react-router-dom"
import { usernameOrName } from "../../components/utils/usernameOrName"
export const ChatHeader = ({ header }) => {
  const nav = useNavigate()
  const handleClick = () => {
    if (header?.username || header?.firstname) {
      nav(`/user/${header?.id}`)
    } else {
      nav(`/group/${header?.id}`)
    }
  }
  return (
    <div className="messages-container__header" onClick={handleClick}>
      {header?.username || header?.firstname
        ? usernameOrName(header)
        : header?.title}
    </div>
  )
}
