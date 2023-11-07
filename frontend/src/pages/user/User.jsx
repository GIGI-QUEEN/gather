import { useParams } from "react-router-dom"
import "../../styles/user-profile.scss"
import { UserProfile } from "../../components/UserProfile"
import { useGet } from "../../hooks/useGet"

const UserPage = () => {
  const { id } = useParams()

  const { data: user, setHit, setFollow } = useGet(`/user/${id}`)

  return <UserProfile user={user} id={id} follow={setFollow} hit={setHit} />
}

export { UserPage }
