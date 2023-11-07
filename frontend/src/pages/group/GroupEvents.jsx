import { useState, useContext } from 'react';
import { longDate } from '../../components/utils/dateFormatter';
import { errorToast } from '../../components/utils/toast/errorToast';
import minDate from '../../components/utils/minDate';
import closeIcon from '../../icons/close-icon.svg';
import { handleModal } from '../../components/utils/handleModal';
import { toast } from 'react-toastify';
import { myAxios } from '../../api/axios';
import { WebSocketContext } from '../../components/utils/WebSocketContext';
import '../../styles/index.scss';
import { UserContext } from '../../components/utils/UserContext';
import { clearInput } from '../../components/utils/clearInput';

export const GroupEvents = ({ events, group_id, setHit }) => {
  return (
    <div className="events-container">
      {events?.map((event, index) => (
        <Event key={index} event={event} group_id={group_id} setHit={setHit} />
      ))}
    </div>
  );
};

const Event = ({ event, group_id, setHit }) => {
  const handleGoing = async () => {
    await myAxios
      .post(`/group/${group_id}/event/${event.event_id}/event-accept`)
      .then(() => {
        setHit(true);
      });
  };
  const handleNotGoing = async () => {
    await myAxios
      .post(`/group/${group_id}/event/${event.event_id}/event-reject`)
      .then(() => {
        setHit(true);
      })
      .catch((error) => {
        errorToast(error);
        return;
      });
  };

  return (
    <div className="event-container">
      <h2 className="event-title">{event.title}</h2>
      <p className="event-info">Date: {longDate(event.event_date)}</p>
      <div className="description-and-buttons">
        <p>{event.description}</p>
        <div className="decision-buttons">
          <button
            onClick={handleGoing}
            disabled={event.going_decision === 1 ? 'picked' : ''}
          >
            {event.going_decision === 1 ? "you're going" : 'going'}
          </button>
          <button
            onClick={handleNotGoing}
            disabled={event.going_decision === 2 ? 'picked' : ''}
          >
            {event.going_decision === 2 ? "you're not going" : 'not going'}
          </button>
        </div>
      </div>
      <p className="event-info">People going: {event.members_going}</p>
    </div>
  );
};

export const CreateEvent = ({ group_id, setHit }) => {
  const { webSocket } = useContext(WebSocketContext);
  const { user } = useContext(UserContext);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [eventDate, setEventDate] = useState('');

  const minimumDate = minDate(new Date());
  const checkDate = Date.parse(eventDate) < Date.parse(minimumDate);

  const handleCreate = async () => {
    if (!checkDate && eventDate !== '') {
      await myAxios
        .post(
          `/group/${group_id}/create-event`,
          JSON.stringify({
            title,
            description,
            event_date: Date.parse(eventDate) / 1000,
          })
        )
        .then(() => {
          // check for group_id should be int
          const groupEventCreatedEvent = {
            event_type: 'ws_group_event_created_event',
            group_id: parseInt(group_id),
            group_event_title: title,
            event_creator: user.id,
          };

          webSocket?.send(JSON.stringify(groupEventCreatedEvent));

          handleModal('create-event');
          setHit(true);
          clearInput(
            ['input', 'textarea'],
            [setTitle, setDescription, setEventDate]
          );
        })
        .catch((error) => {
          errorToast(error);
          return;
        });
    } else {
      toast.error('Wrong date inserted');
      return;
    }
  };

  return (
    <div className="create-modal-container create-event" id="create-event">
      <div className="title-and-close">
        <input
          type="text"
          placeholder="title"
          onChange={(e) => setTitle(e.target.value)}
        />
        <span
          onClick={() => {
            handleModal('create-event');
          }}
        >
          <img src={closeIcon} alt="" />
        </span>
      </div>
      <div
        className={
          checkDate ? 'event-date-input wrong-date' : 'event-date-input'
        }
      >
        <input
          className="event-date"
          type="datetime-local"
          min={minimumDate}
          onChange={(e) => setEventDate(e.target.value)}
        />
        {checkDate ? <p>wrong date</p> : null}
      </div>
      <textarea
        maxLength={100}
        onChange={(e) => setDescription(e.target.value)}
      ></textarea>
      <div className="upload-and-create group-create">
        <button onClick={handleCreate}>create</button>
      </div>
    </div>
  );
};
