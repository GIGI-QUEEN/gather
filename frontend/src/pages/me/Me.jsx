import { UserProfile } from "../../components/UserProfile"
import { useGet } from "../../hooks/useGet"

const Me = () => {
  const { data: me, setHit } = useGet("/me")

  return <UserProfile user={me} hit={setHit} />
}

export default Me
