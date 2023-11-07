import React from "react"
import UsersList from "../../components/UsersList"
import { useGet } from "../../hooks/useGet"
const Users = () => {
  const { data: users } = useGet("/user/all")
  return (
    <div>
      <UsersList users={users} pageOwner={"All Gather users"} />
    </div>
  )
}
export default Users
