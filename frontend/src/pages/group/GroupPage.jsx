import React, { useState } from "react"
import "./group.scss"
import { useGet } from "../../hooks/useGet"
import { useParams } from "react-router-dom"
import Modal from "../../components/Modal"
import { GroupEvents, CreateEvent } from "./GroupEvents"
import { GroupPosts, CreateGroupPost } from "./GroupPosts"
import { GroupPostsAndEvents } from "./GroupPostsAndEvents"
import { GroupBanner } from "./GroupBanner"
import { FollowersToInviteList } from "./FollowersToInviteList"
const GroupPage = () => {
  const { id } = useParams()
  const { data: group, setHit } = useGet(`/group/${id}`)

  const { data: me } = useGet("/me")
  const [switched, setSwitched] = useState("posts") //state to switch between posts and events
  return (
    <div className="content-container">
      <GroupBanner group={group} setHit={setHit} me={me} />
      {group?.join_approved === true && group?.join_requested === true ? (
        <GroupPostsAndEvents switched={switched} setSwitched={setSwitched} />
      ) : null}
      {switched === "posts" ? (
        <GroupPosts group={group} setHit={setHit} />
      ) : (
        <GroupEvents
          events={group?.group_events}
          setHit={setHit}
          group_id={id}
        />
      )}
      <CreateGroupPost group_id={id} setHit={setHit} />
      <CreateEvent group_id={id} setHit={setHit} />
      <FollowersToInviteList
        followers={me?.followers}
        group_id={group?.id}
        group_title={group?.title}
      />
      <Modal />
    </div>
  )
}

export default GroupPage
