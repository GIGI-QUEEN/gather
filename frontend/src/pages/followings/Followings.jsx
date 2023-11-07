import UsersList from "../../components/UsersList"
import { useParams } from "react-router-dom"
import { useGet } from "../../hooks/useGet"

const Followings = () => {
  const { id } = useParams()
  const { data: user } = useGet(`/user/${id}`)
  return (
    <div>
      <UsersList
        users={user?.followings}
        title={"Followings"}
        pageOwner={user}
      />
    </div>
  )
}

export default Followings
