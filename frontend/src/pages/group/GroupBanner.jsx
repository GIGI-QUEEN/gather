import { JoinLeaveCancelButtons } from './GroupBannerButtons';
import { InviteButton } from './GroupBannerButtons';
export const GroupBanner = ({ group, setHit }) => {
  const posts = `posts: ${group?.group_posts ? group?.group_posts?.length : 0}`;
  const members = `members: ${group?.group_members}`;
  return (
    <div className="banner ">
      <div className="first-section">
        <GroupTitle title={group?.title} />
        <div className="stats-container">
          <GroupStats stat={posts} />
          <GroupStats stat={members} />
        </div>
      </div>
      <div className="second-section">
        <GroupDescription description={group?.description} />
        <div className="banner-buttons">
          <InviteButton group={group} />

          <JoinLeaveCancelButtons group={group} setHit={setHit} />
        </div>
      </div>
      <div className="third-section"></div>
    </div>
  );
};

const GroupTitle = ({ title }) => {
  return <div className="container_1 group-title">{title}</div>;
};

const GroupStats = ({ stat }) => {
  return <div className="container_1 group-stats">{stat}</div>;
};

const GroupDescription = ({ description }) => {
  return (
    <div className="container_1 description">
      <p>{description}</p>
    </div>
  );
};
