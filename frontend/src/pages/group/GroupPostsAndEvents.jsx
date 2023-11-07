import createIcon from "../../icons/create-icon.svg"
import eventsIcon from "../../icons/new-event-icon.svg"
import { handleModal } from "../../components/utils/handleModal"
export const GroupPostsAndEvents = ({ switched, setSwitched }) => {
  return (
    <div className="posts-and-events">
      <Switcher
        setSwitched={setSwitched}
        icon={createIcon}
        createButtonValue={"create"}
        text={"posts"}
        classNameValue={
          switched === "posts"
            ? "container_1 switcher switched"
            : "container_1 switcher"
        }
      />
      <hr />
      <Switcher
        setSwitched={setSwitched}
        icon={eventsIcon}
        createButtonValue={"create-event"}
        text={"events"}
        classNameValue={
          switched === "events"
            ? "container_1 switcher switched"
            : "container_1 switcher"
        }
      />
    </div>
  )
}

//createButtonValue could be "create" or "create-event"
//text is used to display button switcher button text and as value for setSwitched
//clasNameValue is used to determine which switcher is active and apply switched class to it
const Switcher = ({
  setSwitched,
  icon,
  createButtonValue,
  text,
  classNameValue,
}) => {
  return (
    <div className={classNameValue}>
      <p onClick={() => setSwitched(text)}>{text}</p>
      <div
        className="container_1 icon"
        onClick={() => {
          handleModal(createButtonValue)
        }}
      >
        <img src={icon} alt="" />
      </div>
    </div>
  )
}
